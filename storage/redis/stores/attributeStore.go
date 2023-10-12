package stores

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zerok-ai/zk-utils-go/ds"
	zklogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	zkRedis "github.com/zerok-ai/zk-utils-go/storage/redis"
	"sort"
)

const LoggerTag = "attribute-cache"

type ProtocolKeyMap map[string][]cache.AttribStoreKey
type NameMapExecutorProtocol map[string]ProtocolKeyMap

type ExecutorAttrStore struct {
	localCacheHSetStore     LocalCacheHSetStore
	nameMapExecutorProtocol *NameMapExecutorProtocol
}

func GetExecutorAttrStore(rc *redis.Client, localCache ds.Cache[map[string]string], csh zkRedis.CacheStoreHook[map[string]string], ctx context.Context) *ExecutorAttrStore {
	attributeCache := (&ExecutorAttrStore{
		localCacheHSetStore: *GetLocalCacheHSetStore(rc, localCache, csh, ctx),
	}).initialize()

	return attributeCache
}

func (attributeCache *ExecutorAttrStore) SetCache(cache ds.Cache[map[string]string]) {
	attributeCache.localCacheHSetStore.SetCache(cache)
}

func (attributeCache *ExecutorAttrStore) initialize() *ExecutorAttrStore {
	attributeCache.nameMapExecutorProtocol = attributeCache.populateAttributeDatasetsFromRedis()
	return attributeCache
}

func (attributeCache *ExecutorAttrStore) Close() {
	attributeCache.localCacheHSetStore.Close()
}

func (attributeCache *ExecutorAttrStore) PutInLocalCache(key string, value *map[string]string) {
	attributeCache.localCacheHSetStore.PutInLocalCache(key, value)
}

func (attributeCache *ExecutorAttrStore) GetFromLocalCache(key string) (*map[string]string, bool) {
	return attributeCache.localCacheHSetStore.Get(key)
}

func (attributeCache *ExecutorAttrStore) GetFromRedis(key string) (*map[string]string, error) {
	value, err := attributeCache.localCacheHSetStore.GetFromRedis(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Get returns the value of the attribute for the given executor, attributeVersion, protocol and attributeName
/* 	It first checks for the protocol, if not found, it checks for the `GENERAL` protocol
   	While searching for the protocol or `GENERAL`, it finds the closest key that is smaller or equal to the input key
	If no key is found, it returns nil.

	For example: Let's say that redis has the following keys:
		OTEL_1.17.0_GENERAL
		OTEL_1.7.0_HTTP
		OTEL_1.7.0_GENERAL
		OTEL_1.21.0_HTTP
		EBPF_0.1.0-alpha_HTTP

	If the input key is `OTEL_1.21.2_HTTP` then the attribute will be searched in `OTEL_1.21.0_HTTP`
	and then in `OTEL_1.17.0_GENERAL`

*/

func (attributeCache *ExecutorAttrStore) GetAttributeFromStore(key cache.AttribStoreKey, attributeName string) (string, bool) {
	return attributeCache.Get(model.ExecutorName(key.Executor), key.Version, model.ProtocolName(key.Protocol), attributeName)
}

func (attributeCache *ExecutorAttrStore) Get(executor model.ExecutorName, attributeVersion string, protocol model.ProtocolName, attributeName string) (string, bool) {

	defer func() {
		if r := recover(); r != nil {
			zklogger.ErrorF(LoggerTag, "In ExecutorAttrStore.Get %s\n: Recovered from panic: %v", executor)
		}
	}()

	protocols := []model.ProtocolName{protocol, model.ProtocolGeneral}
	for _, proto := range protocols {

		// 1. get the closest key
		closestProtocolKey := attributeCache.getClosestKey(string(executor), attributeVersion, proto)

		// 2. get data for closest key from local cache
		dataFromLocalCache, _ := attributeCache.localCacheHSetStore.Get(closestProtocolKey.Value)
		if dataFromLocalCache != nil {
			zklogger.DebugF(LoggerTag, "looking for attribute: %s in key %s", attributeName, closestProtocolKey.Value)
			returnVal, gotVal := (*dataFromLocalCache)[attributeName]
			if gotVal {
				return returnVal, true
			}
		}

	}
	return "", false
}

var BlankKey = &cache.AttribStoreKey{Value: ""}

func (attributeCache *ExecutorAttrStore) getClosestKey(executor string, attributeVersion string, protocol model.ProtocolName) *cache.AttribStoreKey {

	inputKey, err := cache.ParseKey(fmt.Sprintf("%s_%s_%s", executor, attributeVersion, protocol))
	if err != nil {
		return BlankKey
	}

	protocolData, ok := (*attributeCache.nameMapExecutorProtocol)[executor]
	if !ok {
		return BlankKey
	}

	keys, ok := protocolData[string(protocol)]
	if !ok || len(keys) == 0 {
		return BlankKey
	}

	// find the closest key smaller than the input key
	index := 0
	for index = 0; index < len(keys); index++ {
		if keys[index].IsGreaterThan(inputKey) {
			break
		}
	}
	if index == 0 {
		return BlankKey
	}
	return &keys[index-1]
}

func (attributeCache *ExecutorAttrStore) populateAttributeDatasetsFromRedis() *NameMapExecutorProtocol {

	//1. fetch data for the `protocol` and `GENERAL` protocol
	strKeys, err := attributeCache.localCacheHSetStore.GetAllKeysFromRedis("*")
	if err != nil {
		executorData := make(NameMapExecutorProtocol)
		return &executorData
	}

	//2. load data into `NameMapExecutorProtocol` object
	executorData := PopulateExecutorData(strKeys)

	return executorData
}

func PopulateExecutorData(strKeys *[]string) *NameMapExecutorProtocol {

	if strKeys == nil {
		executorData := make(NameMapExecutorProtocol)
		return &executorData
	}

	nameMapExecutorProtocol := make(NameMapExecutorProtocol)

	//2. load data into `NameMapExecutorProtocol` object
	for _, key := range *strKeys {
		parsedKey, err1 := cache.ParseKey(key)
		if err1 != nil {
			continue
		}

		// get the protocol
		protocolKeys, ok := nameMapExecutorProtocol[parsedKey.Executor]
		if !ok {
			protocolKeys = make(ProtocolKeyMap)
			nameMapExecutorProtocol[parsedKey.Executor] = protocolKeys
		}

		// get the keys
		keys, ok := protocolKeys[parsedKey.Protocol]
		if !ok {
			keys = make([]cache.AttribStoreKey, 0)
		}

		keys = append(keys, parsedKey)
		protocolKeys[parsedKey.Protocol] = keys
	}

	// sort the version list
	for _, protocol := range nameMapExecutorProtocol {
		for _, keys := range protocol {
			sort.Sort(cache.ByVersion(keys))
		}
	}
	return &nameMapExecutorProtocol
}
