package files

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
	"testing"
)

func TestFunction(t *testing.T) {
	input := "path123.#func1(param1).#func2().#func3(param31, param32)"
	ff := functions.NewFunctionFactory(nil)
	path, functionArr := ff.GetPathAndFunctions(input)

	assert.Greater(t, len(path), 0)
	assert.Greater(t, len(functionArr), 0)

	assert.Equal(t, "path123", path)
}

func TestNonFunction(t *testing.T) {
	input := "alpha.beta.gamma"
	ff := functions.NewFunctionFactory(nil)
	path, functionArr := ff.GetPathAndFunctions(input)

	assert.Greater(t, len(path), 0)
	assert.Equal(t, 0, len(functionArr))
}
