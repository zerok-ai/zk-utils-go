package stores

import (
	"context"
	"github.com/zerok-ai/zk-utils-go/ds"
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
		var mapOfStores = make(map[string]interface{})
		storeFactory = &StoreFactory{
			redisConfig: redisConfig,
			ctx:         ctx,
			mapOfStores: mapOfStores,
		}
	}
	return storeFactory
}

// GetExecutorAttrStore returns the store. If the store has already been created, it returns the same store.
func (sf *StoreFactory) GetExecutorAttrStore() *ExecutorAttrStore {

	if localCache, ok := sf.mapOfStores[clientDBNames.ExecutorAttrDBName]; ok {
		return localCache.(*ExecutorAttrStore)
	}

	//create new store
	dbName := clientDBNames.ExecutorAttrDBName
	noExpiryCache := ds.GetCacheWithExpiry[map[string]string](ds.NoExpiry)
	redisClient := config.GetRedisConnection(dbName, sf.redisConfig)

	// save and return
	executorAttrStore := GetExecutorAttrStore(redisClient, noExpiryCache, nil, sf.ctx)
	if sf.mapOfStores == nil {
		sf.mapOfStores = make(map[string]interface{})
	}
	sf.mapOfStores[clientDBNames.ExecutorAttrDBName] = executorAttrStore

	return executorAttrStore
}

// GetPodDetailsStore returns the store. If the store has already been created, it returns the same store.
func (sf *StoreFactory) GetPodDetailsStore() *LocalCacheHSetStore {

	if localCache, ok := sf.mapOfStores[clientDBNames.PodDetailsDBName]; ok {
		return localCache.(*LocalCacheHSetStore)
	}

	//create new store
	dbName := clientDBNames.PodDetailsDBName
	expiry := int64(5 * time.Minute)
	expiryCache := ds.GetCacheWithExpiry[map[string]string](expiry)
	redisClient := config.GetRedisConnection(dbName, sf.redisConfig)

	//save and return
	localCache := GetLocalCacheHSetStore(redisClient, expiryCache, nil, sf.ctx)
	if sf.mapOfStores == nil {
		sf.mapOfStores = make(map[string]interface{})
	}
	sf.mapOfStores[clientDBNames.PodDetailsDBName] = localCache

	return localCache
}
