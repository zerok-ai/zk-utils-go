package files

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	zklogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
	"github.com/zerok-ai/zk-utils-go/test/files/helpers"
	"testing"
)

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

func TestNonFunction(t *testing.T) {
	input := "alpha.beta.gamma"

	configPath := "config/config.yaml"
	sf := helpers.GetStoreFactory(configPath)
	ff := functions.NewFunctionFactory(sf.GetPodDetailsStore(), sf.GetExecutorAttrStore())
	key, err := cache.ParseKey("OTEL_1.7.0_HTTP")
	assert.NoError(t, err)

	functionArr := ff.GetPathAndFunctions(input, &key)

	assert.Equal(t, 0, len(functionArr))
}

func TestValueFromStore(t *testing.T) {

	configPath := "config/config.yaml"
	sf := helpers.GetStoreFactory(configPath)
	ff := functions.NewFunctionFactory(sf.GetPodDetailsStore(), sf.GetExecutorAttrStore())

	dataPath := "ruleEvaluation/ip/data.json"
	var dataObject *map[string]interface{}

	// load the data
	err := helpers.LoadFile(dataPath, &dataObject, true)
	assert.NoError(t, err)

	key, err := cache.ParseKey("OTEL_1.7.0_HTTP")
	assert.NoError(t, err)

	input := "dest_service"
	value, ok := ff.EvaluateString(input, *dataObject, &key)
	zklogger.InfoF("result", fmt.Sprintf("%v", value))

	assert.Equal(t, true, ok)

}
