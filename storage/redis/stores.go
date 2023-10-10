package redis

import (
	"context"
	"github.com/zerok-ai/zk-utils-go/ds"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/storage/redis/clientDBNames"
	"github.com/zerok-ai/zk-utils-go/storage/redis/config"
	"time"
)

type StoreFactory struct {
	redisConfig config.RedisConfig
	ctx         context.Context
	mapOfStores map[string]interface{}
}

var storeFactory *StoreFactory

func GetStoreFactory(redisConfig config.RedisConfig, ctx context.Context) *StoreFactory {
	if storeFactory == nil {
		storeFactory = &StoreFactory{
			redisConfig: redisConfig,
			ctx:         ctx,
			mapOfStores: make(map[string]interface{}),
		}
	}
	return storeFactory
}

// GetExecutorAttrStore returns the store. If the store has already been created, it returns the same store.
func (sf *StoreFactory) GetExecutorAttrStore() *cache.AttributeCache {

	if localCache, ok := sf.mapOfStores[clientDBNames.ExecutorAttrDBName]; ok {
		return localCache.(*cache.AttributeCache)
	}

	dbName := clientDBNames.ExecutorAttrDBName
	noExpiryCache := ds.GetCacheWithExpiry[map[string]string](ds.NoExpiry)
	redisClient := config.GetRedisConnection(dbName, sf.redisConfig)

	localCache := cache.GetAttributeCache(redisClient, noExpiryCache, nil, sf.ctx)
	sf.mapOfStores[clientDBNames.ExecutorAttrDBName] = localCache

	return localCache
}

// GetPodDetailsStore returns the store. If the store has already been created, it returns the same store.
func (sf *StoreFactory) GetPodDetailsStore() *LocalCacheHSetStore {

	if localCache, ok := sf.mapOfStores[clientDBNames.PodDetailsDBName]; ok {
		return localCache.(*LocalCacheHSetStore)
	}

	dbName := clientDBNames.PodDetailsDBName
	expiry := int64(5 * time.Minute)
	expiryCache := ds.GetCacheWithExpiry[map[string]string](expiry)
	redisClient := config.GetRedisConnection(dbName, sf.redisConfig)
	localCache := GetLocalCacheHSetStore(redisClient, expiryCache, nil, sf.ctx)
	sf.mapOfStores[clientDBNames.PodDetailsDBName] = localCache

	return localCache
}
