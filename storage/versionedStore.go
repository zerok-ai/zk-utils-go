package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zerok-ai/zk-utils-go/interfaces"
	zkTime "github.com/zerok-ai/zk-utils-go/ticker"
	"time"
)

var LATEST = fmt.Errorf("version passed is already latest")

type RedisConfig struct {
	Host        string `yaml:"host" env:"REDIS_HOST" env-description:"Database host"`
	Port        string `yaml:"port" env:"REDIS_PORT" env-description:"Database port"`
	DB          int    `yaml:"db" env:"REDIS_DB" env-description:"Database to load"`
	ReadTimeout int    `yaml:"readTimeout"`
}

type VersionedStoreConfig struct {
	RefreshTimeSec int `yaml:"RefreshTimeSec" env:"REFRESH_TIME_SEC" env-description:"Database host"`
}

type VersionedStore[T interfaces.ZKComparable] struct {
	redisClient        *redis.Client
	versionHashSetName string
	versions           map[string]string
	localKeyValueCache map[string]*T
}

type Version struct {
	key     string
	version int
}

func GetVersionedStore[T interfaces.ZKComparable](redisConfig *RedisConfig, model T) (*VersionedStore[T], error) {

	if redisConfig == nil {
		return nil, fmt.Errorf("redis config not found")
	}
	readTimeout := time.Duration(redisConfig.ReadTimeout) * time.Second
	_redisClient := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprint(redisConfig.Host, ":", redisConfig.Port),
		Password:    "",
		DB:          redisConfig.DB,
		ReadTimeout: readTimeout,
	})

	versionStore := (&VersionedStore[T]{
		redisClient:        _redisClient,
		versionHashSetName: "zk_value_version",
	}).initialize()

	return versionStore, nil
}

var filterPullTickInterval time.Duration = 10 * time.Second

func (versionStore VersionedStore[T]) initialize() *VersionedStore[T] {

	// trigger recurring filter pull
	tickerFilterPull := time.NewTicker(filterPullTickInterval)

	task := func() {
		err := versionStore.RefreshLocalCache()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	zkTime.RunTaskOnTicks(tickerFilterPull, task)
	task()

	return &versionStore
}

func (versionStore *VersionedStore[T]) Value(key string) (*T, error) {

	// get the value from local store
	localVal := versionStore.localKeyValueCache[key]
	if localVal != nil {
		return localVal, nil
	}

	valueFromStore, err := versionStore.valueFromStore(key)
	if err == nil {
		versionStore.localKeyValueCache[key] = valueFromStore
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
			versionStore.localKeyValueCache[key] = &value
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
	versionStore.localKeyValueCache[key] = &value
	var newVersion string
	newVersion, err = versionStore.version(key)
	if err == nil {
		versionStore.versions[key] = newVersion
	} else {
		//	not sure what to do here
	}

	return nil
}

func (versionStore *VersionedStore[T]) getAllVersions() (map[string]string, error) {
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

func (versionStore *VersionedStore[T]) valuesForKeysFromStore(keys []string) ([]*T, error) {
	rdb := versionStore.redisClient

	// get the values
	opt, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}

	sliceToReturn := make([]*T, len(opt))
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

	// 1. get the new filter versions for each filter
	newFilterVersions, err := versionStore.getAllVersions()
	if err != nil {
		return fmt.Errorf("error in getting versions for values: %v", err)
	}

	// 2. collect the filters which have the same version in a new map
	newFilters := make(map[string]*T)
	var missingOrOldFilters []string
	for key, newVersion := range newFilterVersions {
		oldVersion, ok := versionStore.versions[key]
		if ok {
			if oldVersion == newVersion {
				newFilters[key] = versionStore.localKeyValueCache[key]
				continue
			}
		}
		missingOrOldFilters = append(missingOrOldFilters, key)
	}
	// 2.1 nothing new, go home
	if len(missingOrOldFilters) == 0 {
		return nil
	}

	// 3. get the filters which don't exist
	newRawFilters, err := versionStore.valuesForKeysFromStore(missingOrOldFilters)
	if err != nil {
		return fmt.Errorf("error in fetching new data for cache: %v", err)
	}
	// populate new cache
	for i, v := range newRawFilters {
		newFilters[missingOrOldFilters[i]] = v
	}

	// 4. assign the new objects to filter processors
	versionStore.localKeyValueCache = newFilters
	versionStore.versions = newFilterVersions
	return nil
}
