# Redis

## RedisConfig
This is a `struct` which has all the following configurations for redis:
- `Host`: The host of the redis server.
- `Port`: The port of the redis server.
- `Password`: The password of the redis server.
- `DBs`: Map of usable dbs. It has the following structure: `[DB name]`:`[DB number]`
- `ReadTimeout`: The maximum amount of time for a read.


## Versioned Store
Versioned Store is a local cache built over remote redis server. It has the following features:
- A version is maintained against value for each key.
- Any change in value of a key will increment the version of that key.
- The cache will be updated on initialization and at a regular fixed time interval. 
- This refresh interval can be configured during initialization.
- On refresh, it will fetch the version of each key from redis server and store them in local cache.
- It will also refresh the values from redis server where there is a change in version of a value and store them in local cache.

### Initialization

Call the following function to initialize a versioned store:

```go
func GetVersionedStore[T interfaces.ZKComparable](redisConfig *config.RedisConfig, dbName string, syncTimeInterval time.Duration) (*VersionedStore[T], error)
```

- `T`: The type of the value to be stored in the cache.
- `redisConfig`: The redis config to be used.
- `dbName`: The name of the db to be used to get the DB number from redis config.
- `syncTimeInterval`: The time interval after which the cache will be refreshed.
