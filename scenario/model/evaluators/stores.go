package evaluators

import (
	"context"
	"github.com/zerok-ai/zk-utils-go/ds"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	zkRedis "github.com/zerok-ai/zk-utils-go/storage/redis"
	"github.com/zerok-ai/zk-utils-go/storage/redis/clientDBNames"
	"github.com/zerok-ai/zk-utils-go/storage/redis/config"
	"time"
)

func GetAttributeNamesStore(redisConfig config.RedisConfig, ctx context.Context) *cache.AttributeCache {

	dbName := clientDBNames.ExecutorAttrDBName
	noExpiryCache := ds.GetCacheWithExpiry[map[string]string](ds.NoExpiry)
	redisClient := config.GetRedisConnection(dbName, redisConfig)

	localCache := cache.GetAttributeCache(redisClient, noExpiryCache, nil, ctx)
	return localCache
}

func GetExpiryBasedCacheStore(redisConfig config.RedisConfig, ctx context.Context) *zkRedis.LocalCacheHSetStore {

	dbName := clientDBNames.PodDetailsDBName
	expiry := int64(5 * time.Minute)
	expiryCache := ds.GetCacheWithExpiry[map[string]string](expiry)
	redisClient := config.GetRedisConnection(dbName, redisConfig)
	localCache := zkRedis.GetLocalCacheHSetStore(redisClient, expiryCache, nil, ctx)

	return localCache
}
