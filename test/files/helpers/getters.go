package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
	"github.com/zerok-ai/zk-utils-go/test/files/config"
)

func GetRuleEvaluator() *evaluators.RuleEvaluator {
	executor := "OTEL"

	configPath := "config/config.yaml"
	sf := GetStoreFactory(configPath)
	executorAttrDB := sf.GetExecutorAttrStore()
	podDetailsStore := sf.GetPodDetailsStore()

	return evaluators.NewRuleEvaluator(model.ExecutorName(executor), executorAttrDB, podDetailsStore)
}

func GetStoreFactory(configPath string) *stores.StoreFactory {
	var cfg config.AppConfigs
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	redisConfig := cfg.Redis
	ctx := context.Background()
	return stores.GetStoreFactory(redisConfig, ctx)
}

func LoadFile(path string, dataObj interface{}, printObject bool) error {
	if dataObj == nil {
		return fmt.Errorf("dataObj is nil")
	}
	jsonString := string(common.GetBytesFromFile(path))
	if printObject {
		print("jsonString: ", jsonString, "\n")
	}
	return json.Unmarshal([]byte(jsonString), &dataObj)
}
