package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zerok-ai/zk-utils-go/ds"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
)

// LocalCacheHSetStore is a cache store that uses LRU cache for local caching
type LocalCacheHSetStore struct {
	redisClient    *redis.Client
	localCache     ds.Cache[map[string]string]
	cacheStoreHook CacheStoreHook[map[string]string]
	context        context.Context
}

func GetLocalCacheHSetStore(rc *redis.Client, localCache ds.Cache[map[string]string], csh CacheStoreHook[map[string]string], ctx context.Context) *LocalCacheHSetStore {
	localCacheHSetStore := (&LocalCacheHSetStore{
		redisClient:    rc,
		localCache:     localCache,
		cacheStoreHook: csh,
		context:        ctx,
	}).initialize()

	return localCacheHSetStore

}

func (localCacheHSetStore *LocalCacheHSetStore) SetCache(cache ds.Cache[map[string]string]) {
	localCacheHSetStore.localCache = cache
}

func (localCacheHSetStore *LocalCacheHSetStore) initialize() *LocalCacheHSetStore {
	return localCacheHSetStore
}

func (localCacheHSetStore *LocalCacheHSetStore) Close() {
	err := localCacheHSetStore.redisClient.Close()
	if err != nil {
		return
	}
}

func (localCacheHSetStore *LocalCacheHSetStore) PutInLocalCache(key string, value *map[string]string) {
	localCacheHSetStore.localCache.Put(key, value)
}

func (localCacheHSetStore *LocalCacheHSetStore) GetFromLocalCache(key string) (*map[string]string, bool) {
	return localCacheHSetStore.localCache.Get(key)
}

func (localCacheHSetStore *LocalCacheHSetStore) GetFromRedis(key string) (*map[string]string, error) {
	value, err := localCacheHSetStore.redisClient.HGetAll(localCacheHSetStore.context, key).Result()
	if err != nil {
		return nil, err
	}

	return &value, nil
}

// Get returns the value for the given key. If the value is not present in the cache, it is fetched from the DB and stored in the cache
// The function returns the value and a boolean indicating if the value was fetched from the cache
func (localCacheHSetStore *LocalCacheHSetStore) Get(key string) (*map[string]string, bool) {
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
func (localCacheHSetStore *LocalCacheHSetStore) saveLocally(key string, value *map[string]string) {
	var err *zkErrors.ZkError = nil
	if localCacheHSetStore.cacheStoreHook != nil {
		err = localCacheHSetStore.cacheStoreHook.PreCacheSaveHookAsync(key, value)
	}

	if err == nil {
		localCacheHSetStore.PutInLocalCache(key, value)
	}
}
