package files

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
	"github.com/zerok-ai/zk-utils-go/test/files/helpers"
	"testing"
)

func TestFunction(t *testing.T) {
	input := "path123.#func1(param1).#func2().#func3(param31, param32)"

	sf := helpers.GetStoreFactory()
	ff := functions.NewFunctionFactory(sf.GetPodDetailsStore(), sf.GetExecutorAttrStore())

	key, err := cache.ParseKey("OTEL_1.7.0_HTTP")
	assert.NoError(t, err)

	functionArr := ff.GetPathAndFunctions(input, &key)

	assert.Greater(t, len(functionArr), 0)
}

func TestNonFunction(t *testing.T) {
	input := "alpha.beta.gamma"

	sf := helpers.GetStoreFactory()
	ff := functions.NewFunctionFactory(sf.GetPodDetailsStore(), sf.GetExecutorAttrStore())
	key, err := cache.ParseKey("OTEL_1.7.0_HTTP")
	assert.NoError(t, err)

	functionArr := ff.GetPathAndFunctions(input, &key)

	assert.Equal(t, 0, len(functionArr))
}
