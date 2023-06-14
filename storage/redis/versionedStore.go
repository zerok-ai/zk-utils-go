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
	AutoSync           bool

	tickerTask *ticker.TickerTask
	mutex      sync.Mutex
}

func (versionStore *VersionedStore[T]) safeAddToLocalVersionMap(key string, value string) {
	versionStore.mutex.Lock()
	versionStore.mutex.Unlock()
	versionStore.localVersions[key] = value
}

func (versionStore *VersionedStore[T]) safeAddToLocalKeyValueCacheMap(key string, value *T) {
	versionStore.mutex.Lock()
	versionStore.mutex.Unlock()
	versionStore.localKeyValueCache[key] = value
}

type Version struct {
	key     string
	version int
}

func GetVersionedStore[T interfaces.ZKComparable](redisConfig *config.RedisConfig, dbName string, autoSync bool, model T) (*VersionedStore[T], error) {

	if redisConfig == nil {
		return nil, fmt.Errorf("redis config not found")
	}
	readTimeout := time.Duration(redisConfig.ReadTimeout) * time.Second
	_redisClient := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprint(redisConfig.Host, ":", redisConfig.Port),
		Password:    "",
		DB:          redisConfig.DBs[dbName],
		ReadTimeout: readTimeout,
	})

	versionStore := (&VersionedStore[T]{
		redisClient:        _redisClient,
		versionHashSetName: "zk_value_version",
		localVersions:      map[string]string{},
		localKeyValueCache: map[string]*T{},
		AutoSync:           autoSync,
	}).initialize()

	return versionStore, nil
}

var filterPullTickInterval time.Duration = 2 * time.Minute

func (versionStore *VersionedStore[T]) initialize() *VersionedStore[T] {

	// trigger recurring filter pull
	task := func() {
		err := versionStore.RefreshLocalCache()
		if err != nil {
			zkLogger.Error(LogTag, err.Error())
		}
	}

	if versionStore.AutoSync {
		ticker.GetNewTickerTask("version ticker", filterPullTickInterval, task)
	}
	task()

	return versionStore
}

func (versionStore *VersionedStore[T]) Value(key string) (*T, error) {

	// get the value from local store
	localVal := versionStore.localKeyValueCache[key]
	if localVal != nil {
		return localVal, nil
	}

	var valueFromStore *T
	var err error
	// get the version and value from remote store
	versionFromStore, err := versionStore.getVersionFromDB(key)
	if err == nil {
		versionStore.safeAddToLocalVersionMap(key, versionFromStore)
		valueFromStore, err = versionStore.valueFromStore(key)
		if err == nil {
			versionStore.safeAddToLocalKeyValueCacheMap(key, valueFromStore)
		}
	}

	return valueFromStore, err
}

func (versionStore *VersionedStore[T]) SetValue(key string, value T) error {
	rdb := versionStore.redisClient

	// 1. check if the local value is different from the new value
	localVal := versionStore.localKeyValueCache[key]
	if localVal != nil && (*localVal).Equals(value) {
		return LATEST
	}

	// 2. get the value from remote and return if it matches the new value
	remoteVal := rdb.Get(context.Background(), key)
	if err := remoteVal.Err(); err == nil {
		var remoteT *T
		if err := json.Unmarshal([]byte(remoteVal.Val()), &remoteT); err != nil {
			return err
		}

		if (*remoteT).Equals(value) {
			return LATEST
		}
	} else if err != redis.Nil {
		return err
	}

	// 3. set value in remote store
	// a. create a Redis transaction: this doesn't support rollback
	ctx := context.Background()
	tx := rdb.TxPipeline()

	// b. run set command
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

	//4. update value in local cache
	versionStore.safeAddToLocalKeyValueCacheMap(key, &value)
	var newVersion string
	newVersion, err = versionStore.version(key)
	if err == nil {
		versionStore.safeAddToLocalVersionMap(key, newVersion)
	} else {
		//	not sure what to do here
	}

	return nil
}

func (versionStore *VersionedStore[T]) getVersionFromDB(key string) (string, error) {
	rdb := versionStore.redisClient

	// get the old value
	result := rdb.HGet(context.Background(), versionStore.versionHashSetName, key)
	return result.Val(), result.Err()
}

func (versionStore *VersionedStore[T]) getAllVersionsFromDB() (map[string]string, error) {
	rdb := versionStore.redisClient

	// get the old value
	versions := rdb.HGetAll(context.Background(), versionStore.versionHashSetName)
	return versions.Val(), versions.Err()
}

// value Get Redis `GET key` command. It returns storage.LATEST error when oldVersion==<version in the store>.
func (versionStore *VersionedStore[T]) valueFromStore(key string) (*T, error) {
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

func (versionStore *VersionedStore[T]) valuesForKeysFromDB(keys []string) ([]*T, error) {
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

func (versionStore *VersionedStore[T]) version(key string) (string, error) {
	rdb := versionStore.redisClient

	// get the old value
	version := rdb.HGet(context.Background(), versionStore.versionHashSetName, key)
	return version.Val(), version.Err()
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
	return nil
}

func (versionStore *VersionedStore[T]) Length() (int64, error) {
	// get the number of hash key-value pairs
	return versionStore.redisClient.HLen(context.Background(), versionStore.versionHashSetName).Result()
}

func (versionStore *VersionedStore[T]) RefreshLocalCache() error {

	// 1. get the new localVersions for all the keys
	versionsFromDB, err := versionStore.getAllVersionsFromDB()
	if err != nil {
		return fmt.Errorf("error in getting localVersions for values: %v", err)
	}

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
	// 2.1 nothing new, go home
	if len(missingOrOldDataKeys) == 0 {
		return nil
	}

	// 3. get the values which are not present locally or are old
	newRawDataPair, err := versionStore.valuesForKeysFromDB(missingOrOldDataKeys)
	if err != nil {
		return fmt.Errorf("error in fetching new data for cache: %v", err)
	}
	// populate new cache
	for i, v := range newRawDataPair {
		newDataPair[missingOrOldDataKeys[i]] = v
	}

	// 4. assign the new objects to filter processors
	versionStore.localKeyValueCache = newDataPair
	versionStore.localVersions = versionsFromDB
	return nil
}
