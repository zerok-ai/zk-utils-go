package helpers

import (
	"context"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
	"github.com/zerok-ai/zk-utils-go/test/files/config"
)

func GetRuleEvaluator() *evaluators.RuleEvaluator {
	executor := "OTEL"

	sf := GetStoreFactory()
	executorAttrDB := sf.GetExecutorAttrStore()
	podDetailsStore := sf.GetPodDetailsStore()

	return evaluators.NewRuleEvaluator(model.ExecutorName(executor), executorAttrDB, podDetailsStore)
}

func GetStoreFactory() *stores.StoreFactory {
	configPath := "./config/config.yaml"
	var cfg config.AppConfigs
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	redisConfig := cfg.Redis
	ctx := context.Background()
	return stores.GetStoreFactory(redisConfig, ctx)
}
