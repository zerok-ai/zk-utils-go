package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type RedisConfig struct {
	Host        string         `yaml:"host"`
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
	password := os.Getenv("ZK_REDIS_PASSWORD")
	return redis.NewClient(&redis.Options{
		Addr:        fmt.Sprint(redisConfig.Host, ":", redisConfig.Port),
		Password:    password,
		DB:          redisConfig.DBs[dbName],
		ReadTimeout: readTimeout,
	})
}
