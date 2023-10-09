package scenario

import (
	"context"
	"github.com/zerok-ai/zk-utils-go/ds"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/storage/redis/clientDBNames"
	"github.com/zerok-ai/zk-utils-go/storage/redis/config"
)

func GetAttributeNamesStore(redisConfig config.RedisConfig, ctx context.Context) *cache.AttributeCache {

	dbName := clientDBNames.ExecutorAttrDBName
	noExpiryCache := ds.GetCacheWithExpiry[map[string]string](ds.NoExpiry)
	redisClient := config.GetRedisConnection(dbName, redisConfig)

	localCache := cache.GetAttributeCache(redisClient, noExpiryCache, nil, ctx)
	return localCache
}
