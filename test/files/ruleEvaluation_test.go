package files

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	zklogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
	"github.com/zerok-ai/zk-utils-go/test/files/helpers"
	"testing"
)

func TestRuleEvaluation(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("ruleEvaluation/othertests/schema.json", &w, "ruleEvaluation/othertests/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestRuleEvaluationBooleanAnd(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("ruleEvaluation/bool/and/schema.json", &w, "ruleEvaluation/bool/and/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestRuleEvaluationBooleanOR(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("ruleEvaluation/bool/or/schema.json", &w, "ruleEvaluation/bool/or/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, false)
}

func TestRuleEvaluationBooleanNotEqual(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("ruleEvaluation/bool/not_equal/schema.json", &w, "ruleEvaluation/bool/not_equal/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestRuleEvaluationString(t *testing.T) {

	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("ruleEvaluation/string/schema.json", &w, "ruleEvaluation/string/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestRuleEvaluationInteger(t *testing.T) {
	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("./ruleEvaluation/integer/schema.json", &w, "./ruleEvaluation/integer/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestRuleEvaluationFloat(t *testing.T) {
	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("./ruleEvaluation/float/schema.json", &w, "./ruleEvaluation/float/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestRuleEvaluationID(t *testing.T) {
	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("./ruleEvaluation/id/schema.json", &w, "./ruleEvaluation/id/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestRuleEvaluationExistsOperator(t *testing.T) {
	var dataStore map[string]interface{}
	var w model.Workload
	err := helpers.LoadObjects("./ruleEvaluation/exists/schema.json", &w, "./ruleEvaluation/exists/data.json", &dataStore)
	assert.NoError(t, err)

	helpers.Validate(t, w, dataStore, true)
}

func TestFunction(t *testing.T) {
	input := "path123.#func1(param1).#func2().#func3(param31, param32)"

	configPath := "config/config.yaml"
	sf := helpers.GetStoreFactory(configPath)
	ff := functions.NewFunctionFactory(sf.GetPodDetailsStore(), sf.GetExecutorAttrStore())

	key, err := cache.ParseKey("OTEL_1.7.0_HTTP")
	assert.NoError(t, err)

	functionArr := ff.GetPathAndFunctions(input, &key)

	assert.Greater(t, len(functionArr), 0)
}

func TestValueFromStore(t *testing.T) {

	configPath := "config/config.yaml"
	sf := helpers.GetStoreFactory(configPath)
	ff := functions.NewFunctionFactory(sf.GetPodDetailsStore(), sf.GetExecutorAttrStore())

	dataPath := "ruleEvaluation/ip/data.json"
	var dataObject *map[string]interface{}

	// load the data
	err := helpers.LoadFile(dataPath, &dataObject, false)
	assert.NoError(t, err)

	key, err := cache.ParseKey("OTEL_1.7.0_HTTP")
	assert.NoError(t, err)

	input := "dest_service"
	value, ok := ff.EvaluateString(input, *dataObject, &key)
	zklogger.InfoF("result", fmt.Sprintf("%v", value))

	assert.Equal(t, true, ok)

}
