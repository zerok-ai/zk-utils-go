package files

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/test/files/helpers"
	"testing"
)

func TestRuleEvaluation(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/bool/equal/schema.json", &w, "test/files/ruleEvaluation/bool/equal/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationBooleanAnd(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("ruleEvaluation/bool/and/schema.json", &w, "ruleEvaluation/bool/and/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationBooleanOR(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("ruleEvaluation/bool/or/schema.json", &w, "ruleEvaluation/bool/or/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, false)
}

func TestRuleEvaluationBooleanNotEqual(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("ruleEvaluation/bool/not_equal/schema.json", &w, "ruleEvaluation/bool/not_equal/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationString(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("ruleEvaluation/string/schema.json", &w, "ruleEvaluation/string/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationInteger(t *testing.T) {
	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("./ruleEvaluation/integer/schema.json", &w, "./ruleEvaluation/integer/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationFloat(t *testing.T) {
	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("./ruleEvaluation/float/schema.json", &w, "./ruleEvaluation/float/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

//func TestMatchingVersionKey(t *testing.T) {
//
//	keys := []string{
//		"OTEL_1.21.0_GENERAL",
//		"OTEL_1.7.0_HTTP",
//		"OTEL_1.7.0_GENERAL",
//		"OTEL_1.21.0_HTTP",
//		"EBPF_0.1.0-alpha_HTTP",
//	}
//
//	//protocol := "HTTP"
//	//executor := "OTEL"
//
//	//attributeCache := cache.AttributeCache{}
//	//attributeCache.Executors = PopulateExecutorDataFromRedis(config.RedisConfig{}, context.Background())
//
//	//var parsedKeys []cache.Key
//	//for _, key := range keys {
//	//	parsedKey, err := cache.ParseKey(key)
//	//	if err != nil {
//	//		fmt.Printf("Error parsing key: %v\n", err)
//	//		continue
//	//	}
//	//	parsedKeys = append(parsedKeys, parsedKey)
//	//}
//	//
//	//// Sort the keys using the custom sorting criteria.
//	//sort.Sort(cache.ByVersion(parsedKeys))
//	//
//	//// Print the sorted keys.
//	//for _, parsedKey := range parsedKeys {
//	//	fmt.Println(parsedKey.Value)
//	//}
//
//}

func TestRuleEvaluationID(t *testing.T) {
	var dataStore map[string]interface{}
	var w model.Workload
	err := loadObjects("./ruleEvaluation/id/schema.json", &w, "./ruleEvaluation/id/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func validate(t *testing.T, w model.Workload, dataStore map[string]interface{}, expected bool) {
	var result bool

	ruleEvaluator := helpers.GetRuleEvaluator()
	key, err := cache.ParseKey("OTEL_1.21.0_GENERAL")
	result, err = ruleEvaluator.EvalRule(w.Rule, key, "HTTP", dataStore)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func loadObjects[T comparable](schemaPath string, schemaObj *T, dataPath string, dataObject *map[string]interface{}) error {
	err := helpers.LoadFile(schemaPath, &schemaObj, false)
	if err != nil {
		return err
	}

	err = helpers.LoadFile(dataPath, &dataObject, true)
	if err != nil {
		return err
	}
	return nil
}
