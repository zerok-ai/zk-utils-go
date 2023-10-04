package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zerok-ai/zk-utils-go/ds"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	zkRedis "github.com/zerok-ai/zk-utils-go/storage/redis"
	"sort"
)

type Protocol map[string][]Key
type ExecutorProtocol map[string]Protocol
type ExecutorDataPopulator func()

type AttributeCache struct {
	localCacheHSetStore zkRedis.LocalCacheHSetStore
	executors           *ExecutorProtocol
}

func GetAttributeCache(rc *redis.Client, localCache ds.Cache[map[string]string], csh zkRedis.CacheStoreHook[map[string]string], ctx context.Context) *AttributeCache {
	attributeCache := (&AttributeCache{
		localCacheHSetStore: *zkRedis.GetLocalCacheHSetStore(rc, localCache, csh, ctx),
	}).initialize()

	return attributeCache
}

func (attributeCache *AttributeCache) SetCache(cache ds.Cache[map[string]string]) {
	attributeCache.localCacheHSetStore.SetCache(cache)
}

func (attributeCache *AttributeCache) initialize() *AttributeCache {
	attributeCache.executors = attributeCache.populateExecutorDataFromRedis()
	return attributeCache
}

func (attributeCache *AttributeCache) Close() {
	attributeCache.localCacheHSetStore.Close()
}

func (attributeCache *AttributeCache) PutInLocalCache(key string, value *map[string]string) {
	attributeCache.localCacheHSetStore.PutInLocalCache(key, value)
}

func (attributeCache *AttributeCache) GetFromLocalCache(key string) (*map[string]string, bool) {
	return attributeCache.localCacheHSetStore.Get(key)
}

func (attributeCache *AttributeCache) GetFromRedis(key string) (*map[string]string, error) {
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
func (attributeCache *AttributeCache) Get(executor, attributeVersion string, protocol model.Protocol, attributeName string) *string {

	protocols := []model.Protocol{protocol, model.ProtocolGeneral}
	for _, proto := range protocols {

		// 1. get the closest key
		closestProtocolKey := attributeCache.getClosestKey(executor, attributeVersion, proto)

		// 2. get data for closest key from local cache
		dataFromLocalCache, _ := attributeCache.localCacheHSetStore.Get(closestProtocolKey.Value)
		if dataFromLocalCache != nil {
			returnVal := (*dataFromLocalCache)[attributeName]
			if returnVal != "" {
				return &returnVal
			}
		}

	}
	return nil
}

var BlankKey = &Key{Value: ""}

func (attributeCache *AttributeCache) getClosestKey(executor string, attributeVersion string, protocol model.Protocol) *Key {

	inputKey, err := ParseKey(fmt.Sprintf("%s_%s_%s", executor, attributeVersion, protocol))
	if err != nil {
		return BlankKey
	}

	protocolData, ok := (*attributeCache.executors)[executor]
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
		if !keys[index].IsLessThan(inputKey) {
			break
		}
	}
	if index == 0 {
		return BlankKey
	}
	return &keys[index-1]
}

func (attributeCache *AttributeCache) populateExecutorDataFromRedis() *ExecutorProtocol {

	//1. fetch data for the `protocol` and `GENERAL` protocol
	strKeys, err := attributeCache.localCacheHSetStore.GetAllKeysFromRedis("*")
	if err != nil {
		executorData := make(ExecutorProtocol)
		return &executorData
	}

	//2. load data into `ExecutorProtocol` object
	executorData := PopulateExecutorData(strKeys)

	return executorData
}

func PopulateExecutorData(strKeys *[]string) *ExecutorProtocol {

	if strKeys != nil {
		executorData := make(ExecutorProtocol)
		return &executorData
	}

	executorData := make(ExecutorProtocol)

	//2. load data into `ExecutorProtocol` object
	for _, key := range *strKeys {
		parsedKey, err1 := ParseKey(key)
		if err1 != nil {
			continue
		}

		// get the protocol
		protocolKeys, ok := executorData[parsedKey.Executor]
		if !ok {
			protocolKeys = make(Protocol)
			executorData[parsedKey.Executor] = protocolKeys
		}

		// get the keys
		keys, ok := protocolKeys[parsedKey.Protocol]
		if !ok {
			keys = make([]Key, 0)
			protocolKeys[parsedKey.Protocol] = keys
		}

		keys = append(keys, parsedKey)
	}

	// sort the version list
	for _, protocol := range executorData {
		for _, keys := range protocol {
			sort.Sort(ByVersion(keys))
		}
	}
	return &executorData
}
