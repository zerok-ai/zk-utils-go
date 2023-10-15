package config

import "github.com/zerok-ai/zk-utils-go/storage/redis/config"

// AppConfigs is an application configuration structure
type AppConfigs struct {
	Redis config.RedisConfig `yaml:"redis"`
}
