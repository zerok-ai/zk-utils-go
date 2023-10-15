package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const LoggerTag = "redis_config"

type RedisConfig struct {
	Host        string         `yaml:"host" env:"ZK_REDIS_HOST" env-description:"Redis HOST"`
	Port        string         `yaml:"port"`
	DBs         map[string]int `yaml:"dbs"`
	ReadTimeout int            `yaml:"readTimeout"`
	Password    string         `yaml:"password" env:"ZK_REDIS_PASSWORD" env-description:"Redis password"`
}

type DB struct {
	Name string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
}

func GetRedisConnection(dbName string, redisConfig RedisConfig) *redis.Client {
	readTimeout := time.Duration(redisConfig.ReadTimeout) * time.Second
	return redis.NewClient(&redis.Options{
		Addr:        fmt.Sprint(redisConfig.Host, ":", redisConfig.Port),
		Password:    redisConfig.Password,
		DB:          redisConfig.DBs[dbName],
		ReadTimeout: readTimeout,
	})
}
