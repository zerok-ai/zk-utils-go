package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zerok-ai/zk-utils-go/ds"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"log"
)

type LocalCacheHSetStore interface {
	Close()
	SetCache(cache ds.Cache[map[string]string])
	PutInLocalCache(key string, value *map[string]string)
	Get(key string) (*map[string]string, bool)
	GetFromLocalCache(key string) (*map[string]string, bool)
	GetFromRedis(key string) (*map[string]string, error)
	GetAllKeysFromRedis(pattern string) (*[]string, error)
}

// LocalCacheHSetStoreInternal is a cache store that uses LRU cache for local caching
type LocalCacheHSetStoreInternal struct {
	redisClient    *redis.Client
	localCache     ds.Cache[map[string]string]
	cacheStoreHook CacheStoreHook[map[string]string]
	context        context.Context
}

func GetLocalCacheHSetStore(rc *redis.Client, localCache ds.Cache[map[string]string], csh CacheStoreHook[map[string]string], ctx context.Context) *LocalCacheHSetStore {
	internal := (&LocalCacheHSetStoreInternal{
		redisClient:    rc,
		localCache:     localCache,
		cacheStoreHook: csh,
		context:        ctx,
	}).initialize()

	var localCacheHSetStore LocalCacheHSetStore = &internal
	return &localCacheHSetStore

}

func (localCacheHSetStore *LocalCacheHSetStoreInternal) SetCache(cache ds.Cache[map[string]string]) {
	localCacheHSetStore.localCache = cache
}

func (localCacheHSetStore *LocalCacheHSetStoreInternal) initialize() LocalCacheHSetStoreInternal {
	return *localCacheHSetStore
}

func (localCacheHSetStore *LocalCacheHSetStoreInternal) Close() {
	err := localCacheHSetStore.redisClient.Close()
	if err != nil {
		return
	}
}

func (localCacheHSetStore *LocalCacheHSetStoreInternal) PutInLocalCache(key string, value *map[string]string) {
	localCacheHSetStore.localCache.Put(key, value)
}

func (localCacheHSetStore *LocalCacheHSetStoreInternal) GetFromLocalCache(key string) (*map[string]string, bool) {
	return localCacheHSetStore.localCache.Get(key)
}

func (localCacheHSetStore *LocalCacheHSetStoreInternal) GetFromRedis(key string) (*map[string]string, error) {
	value, err := localCacheHSetStore.redisClient.HGetAll(localCacheHSetStore.context, key).Result()
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (localCacheHSetStore *LocalCacheHSetStoreInternal) GetAllKeysFromRedis(pattern string) (*[]string, error) {
	// Use the client to retrieve all keys using the KEYS command (not recommended for large databases)
	keys, err := localCacheHSetStore.redisClient.Keys(localCacheHSetStore.context, pattern).Result()
	if err != nil {
		log.Fatalf("Error getting keys: %v", err)
	}

	return &keys, nil
}

// Get returns the value for the given key. If the value is not present in the cache, it is fetched from the DB and stored in the cache
// The function returns the value and a boolean indicating if the value was fetched from the cache
func (localCacheHSetStore *LocalCacheHSetStoreInternal) Get(key string) (*map[string]string, bool) {
	value, fromCache := localCacheHSetStore.GetFromLocalCache(key)
	if value == nil {
		fromCache = false
		var err error
		value, err = localCacheHSetStore.GetFromRedis(key)
		if err != nil {
			return nil, fromCache
		}
		defer localCacheHSetStore.saveLocally(key, value)
	}
	return value, fromCache
}

// Put puts the given key-value pair in the cache and DB
func (localCacheHSetStore *LocalCacheHSetStoreInternal) saveLocally(key string, value *map[string]string) {
	var err *zkErrors.ZkError = nil
	if localCacheHSetStore.cacheStoreHook != nil {
		err = localCacheHSetStore.cacheStoreHook.PreCacheSaveHookAsync(key, value)
	}

	if err == nil {
		localCacheHSetStore.PutInLocalCache(key, value)
	}
}
