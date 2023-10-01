package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zerok-ai/zk-utils-go/ds"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
)

type CacheStore interface {
	Close()
	SetCache(cache ds.Cache[string])
	PutInLocalCache(key string, value *string)
	GetFromLocalCache(key string) (*string, bool)
	GetFromRedis(keys []string) ([]*string, error)
}

type CacheStoreHook[T any] interface {
	PreCacheSaveHookAsync(key string, value *T) *zkErrors.ZkError
}

/*-------------- Implementations of LocalCache -----------------*/

// LocalCacheKVStore is a cache store that uses LRU cache for local caching
type LocalCacheKVStore[T any] struct {
	redisClient    *redis.Client
	localCache     ds.Cache[T]
	cacheStoreHook CacheStoreHook[T]
	context        context.Context
}

func GetLocalCacheStore[T any](rc *redis.Client, localCache ds.Cache[T], csh CacheStoreHook[T], ctx context.Context) *LocalCacheKVStore[T] {
	localCacheStore := (&LocalCacheKVStore[T]{
		redisClient:    rc,
		localCache:     localCache,
		cacheStoreHook: csh,
		context:        ctx,
	}).initialize()

	return localCacheStore

}

func (localCacheKVStore *LocalCacheKVStore[T]) SetCache(cache ds.Cache[T]) {
	localCacheKVStore.localCache = cache
}

func (localCacheKVStore *LocalCacheKVStore[T]) initialize() *LocalCacheKVStore[T] {
	return localCacheKVStore
}

func (localCacheKVStore *LocalCacheKVStore[T]) Close() {
	err := localCacheKVStore.redisClient.Close()
	if err != nil {
		return
	}
}

func (localCacheKVStore *LocalCacheKVStore[T]) PutInLocalCache(key string, value *T) {
	localCacheKVStore.localCache.Put(key, value)
}

func (localCacheKVStore *LocalCacheKVStore[T]) GetFromLocalCache(key string) (*T, bool) {
	return localCacheKVStore.localCache.Get(key)
}

func (localCacheKVStore *LocalCacheKVStore[T]) GetFromRedis(keys []string) ([]*T, error) {
	values, err := localCacheKVStore.redisClient.MGet(localCacheKVStore.context, keys...).Result()
	if err != nil {
		return nil, err
	}

	// Process the retrieved values
	responseArray := make([]*T, len(values))
	for i, value := range values {
		// Check if the value can be typecast to T
		if typeCastedValue, ok := value.(T); ok {
			responseArray[i] = &typeCastedValue
		} else {
			responseArray[i] = nil
		}
	}
	return responseArray, nil
}

// Get returns the value for the given key. If the value is not present in the cache, it is fetched from the DB and stored in the cache
// The function returns the value and a boolean indicating if the value was fetched from the cache
func (localCacheKVStore *LocalCacheKVStore[T]) Get(key string) (*T, bool) {
	value, fromCache := localCacheKVStore.GetFromLocalCache(key)
	if value == nil {
		fromCache = false
		valueFromDB, err := localCacheKVStore.GetFromRedis([]string{key})
		if err != nil {
			return nil, fromCache
		}
		value = valueFromDB[0]
		defer localCacheKVStore.saveLocally(key, value)
	}
	return value, fromCache
}

// Put puts the given key-value pair in the cache and DB
func (localCacheKVStore *LocalCacheKVStore[T]) saveLocally(key string, value *T) {
	var err *zkErrors.ZkError = nil
	if localCacheKVStore.cacheStoreHook != nil {
		err = localCacheKVStore.cacheStoreHook.PreCacheSaveHookAsync(key, value)
	}

	if err == nil {
		localCacheKVStore.PutInLocalCache(key, value)
	}
}
