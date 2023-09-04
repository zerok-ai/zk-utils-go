package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zerok-ai/zk-utils-go/interfaces"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/redis/config"
	ticker "github.com/zerok-ai/zk-utils-go/ticker"
	"sync"
	"time"
)

var LATEST = fmt.Errorf("version passed is already latest")
var LogTag = "redis_versionedStore"

type VersionedStoreConfig struct {
	RefreshTimeSec int `yaml:"RefreshTimeSec" env:"REFRESH_TIME_SEC" env-description:"Database host"`
}

type VersionedStore[T interfaces.ZKComparable] struct {
	redisClient        *redis.Client
	versionHashSetName string
	localVersions      map[string]string
	localKeyValueCache map[string]*T

	tickerTask *ticker.TickerTask
	mutex      sync.Mutex
}

func GetVersionedStore[T interfaces.ZKComparable](redisConfig *config.RedisConfig, dbName string, syncTimeInterval time.Duration) (*VersionedStore[T], error) {
	if redisConfig == nil {
		return nil, fmt.Errorf("redis config not found")
	}

	versionStore := (&VersionedStore[T]{
		redisClient:        config.GetRedisConnection(dbName, *redisConfig),
		versionHashSetName: "zk_value_version",
		localVersions:      map[string]string{},
		localKeyValueCache: map[string]*T{},
	}).initialize(dbName, syncTimeInterval)

	return versionStore, nil
}

func (versionStore *VersionedStore[T]) initialize(tickerName string, syncTimeInterval time.Duration) *VersionedStore[T] {

	// trigger recurring filter pull
	task := func() {
		err := versionStore.refreshLocalCache()
		if err != nil {
			zkLogger.Error(LogTag, err)
		}
	}
	versionStore.tickerTask = ticker.GetNewTickerTask(tickerName, syncTimeInterval, task).Start()

	return versionStore
}

func (versionStore *VersionedStore[T]) Close() {
	versionStore.tickerTask.Stop()
	err := versionStore.redisClient.Close()
	if err != nil {
		return
	}
}

func (versionStore *VersionedStore[T]) safeAddToLocalVersionMap(key string, value string) {
	versionStore.mutex.Lock()
	defer versionStore.mutex.Unlock()
	versionStore.localVersions[key] = value
}

func (versionStore *VersionedStore[T]) safeAddToLocalKeyValueCacheMap(key string, value *T) {
	versionStore.mutex.Lock()
	defer versionStore.mutex.Unlock()
	versionStore.localKeyValueCache[key] = value
}

// setToLocalCache refreshes the local cache from the remote store using the variables passed to the
// method. If the variables are nil, it will fetch the values from the remote store before setting them
// in the local cache.
func (versionStore *VersionedStore[T]) setToLocalCache(key string, value *T, version *string) (*T, *string, error) {

	if version == nil {
		// get the version from remote store
		versionFromStore, err := versionStore.getVersionFromDB(key)
		if err != nil {
			return value, version, err
		}
		version = &versionFromStore
	}

	// set the version in local store
	versionStore.safeAddToLocalVersionMap(key, *version)

	if value == nil {
		// get the value from remote store
		valueFromStore, err := versionStore.getValueFromDB(key)
		if err != nil {
			return value, version, err
		}
		value = valueFromStore
	}

	// set the value in local store
	versionStore.safeAddToLocalKeyValueCacheMap(key, value)
	return value, version, nil
}

func (versionStore *VersionedStore[T]) GetAllValues() map[string]*T {
	return versionStore.localKeyValueCache
}

func (versionStore *VersionedStore[T]) GetValue(key string) (*T, error) {

	// get the value from local store
	localVal := versionStore.localKeyValueCache[key]
	if localVal != nil {
		return localVal, nil
	}

	valueFromStore, _, err := versionStore.setToLocalCache(key, nil, nil)

	return valueFromStore, err
}

// value Get Redis `GET key` command. It returns storage.LATEST error when oldVersion==<version in the store>.
func (versionStore *VersionedStore[T]) getValueFromDB(key string) (*T, error) {
	rdb := versionStore.redisClient

	// get the value
	opt := rdb.Get(context.Background(), key)
	err := opt.Err()

	var value *T
	if err == nil {
		err = json.Unmarshal([]byte(opt.Val()), &value)
	}
	return value, err
}

func (versionStore *VersionedStore[T]) getMultipleValuesFromDB(keys []string) ([]*T, error) {
	rdb := versionStore.redisClient

	// get the values
	opt, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}

	var sliceToReturn []*T
	for i := 0; i < len(opt); i++ {
		var valToAppend *T
		switch value := opt[i].(type) {
		case string:
			if err := json.Unmarshal([]byte(value), &valToAppend); err != nil {
				valToAppend = nil
			}
		case []byte:
			if err := json.Unmarshal(value, &valToAppend); err != nil {
				valToAppend = nil
			}
		default:
			valToAppend = nil
		}
		sliceToReturn = append(sliceToReturn, valToAppend)
	}
	return sliceToReturn, err
}

func (versionStore *VersionedStore[T]) SetValue(key string, value T) error {

	// 1. check if the previous value is different from the new value
	localVal, _ := versionStore.GetValue(key)
	if localVal != nil && (*localVal).Equals(value) {
		return LATEST
	}

	// 2. set value in remote store
	return versionStore.setValueForced(key, value)
}

func (versionStore *VersionedStore[T]) setValueForced(key string, value T) error {

	rdb := versionStore.redisClient

	// a. create a Redis transaction: this doesn't support rollback
	ctx := context.Background()
	tx := rdb.TxPipeline()

	// b. run set command for value and version
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	valueToWrite := string(bytes)
	tx.Set(ctx, key, valueToWrite, 0)
	tx.HIncrBy(ctx, versionStore.versionHashSetName, key, 1)

	// c. Execute the transaction
	if _, err := tx.Exec(ctx); err != nil {
		return err
	}

	// d. reset value in local cache
	_, _, err = versionStore.setToLocalCache(key, &value, nil)
	return err
}

func (versionStore *VersionedStore[T]) getVersionFromDB(key string) (string, error) {
	rdb := versionStore.redisClient

	// get the old value
	version := rdb.HGet(context.Background(), versionStore.versionHashSetName, key)
	return version.Val(), version.Err()
}

func (versionStore *VersionedStore[T]) getAllVersionsFromDB() (map[string]string, error) {
	rdb := versionStore.redisClient

	// get the old value
	versions := rdb.HGetAll(context.Background(), versionStore.versionHashSetName)
	return versions.Val(), versions.Err()
}

func (versionStore *VersionedStore[T]) Delete(key string) error {
	rdb := versionStore.redisClient

	// create a transaction
	ctx := context.Background()
	tx := rdb.TxPipeline()

	// delete version
	tx.HDel(context.Background(), versionStore.versionHashSetName, key)
	tx.Del(context.Background(), key)

	// Execute the transaction
	if _, err := tx.Exec(ctx); err != nil {
		return err
	}

	versionStore.mutex.Lock()
	defer versionStore.mutex.Unlock()
	delete(versionStore.localVersions, key)
	delete(versionStore.localKeyValueCache, key)

	return nil
}

func (versionStore *VersionedStore[T]) Length() (int64, error) {
	// get the number of hash key-value pairs
	return versionStore.redisClient.HLen(context.Background(), versionStore.versionHashSetName).Result()
}

func (versionStore *VersionedStore[T]) refreshLocalCache() error {

	zkLogger.Debug(LogTag, "Triggered refreshLocalCache.")

	// 1. get the new localVersions for all the keys
	versionsFromDB, err := versionStore.getAllVersionsFromDB()
	if err != nil {
		return fmt.Errorf("error in getting localVersions for values: %v", err)
	}

	zkLogger.Debug(LogTag, "VersionFromDB ", versionsFromDB)

	// 2. collect the data points which have the same versionFromDb in a new map
	newDataPair := make(map[string]*T)
	var missingOrOldDataKeys []string
	for key, versionFromDb := range versionsFromDB {
		oldVersion, ok := versionStore.localVersions[key]
		if ok {
			if oldVersion == versionFromDb {
				newDataPair[key] = versionStore.localKeyValueCache[key]
				continue
			}
		}
		missingOrOldDataKeys = append(missingOrOldDataKeys, key)
	}

	zkLogger.Debug(LogTag, "MissingOrOldKeys ", missingOrOldDataKeys)

	// 2.1 nothing new, go home
	if len(missingOrOldDataKeys) == 0 {
		return nil
	}

	// 3. get the values which are not present locally or are old
	newRawDataPair, err := versionStore.getMultipleValuesFromDB(missingOrOldDataKeys)
	if err != nil {
		return fmt.Errorf("error in fetching new data for cache: %v", err)
	}

	zkLogger.Debug(LogTag, "newRawDataPair ", newRawDataPair)

	// populate new cache
	for i, v := range newRawDataPair {
		newDataPair[missingOrOldDataKeys[i]] = v
	}

	zkLogger.Debug(LogTag, "newDataPair ", newDataPair)

	// 4. assign the new objects to filter processors
	versionStore.mutex.Lock()
	defer versionStore.mutex.Unlock()
	versionStore.localKeyValueCache = newDataPair
	versionStore.localVersions = versionsFromDB
	return nil
}
