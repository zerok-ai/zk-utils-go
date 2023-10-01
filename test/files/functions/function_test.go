package files

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators"
	"testing"
)

func TestFunction(t *testing.T) {
	input := "path123.#func1(param1).#func2().#func3(param31, param32)"
	path, functions := evaluators.GetPathAndFunctions(input)

	assert.Greater(t, len(path), 0)
	assert.Greater(t, len(functions), 0)

	assert.Equal(t, "path123", path)
}

func TestNonFunction(t *testing.T) {
	input := "alpha.beta.gamma"
	path, functions := evaluators.GetPathAndFunctions(input)

	assert.Greater(t, len(path), 0)
	assert.Equal(t, 0, len(functions))
}
