package helpers

import (
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"testing"
)

func Validate(t *testing.T, w model.Workload, dataStore map[string]interface{}, expected bool) {
	var result bool

	ruleEvaluator := GetRuleEvaluator()
	key, err := cache.ParseKey("OTEL_1.21.0_HTTP")
	result, err = ruleEvaluator.EvalRule(w.Rule, key, dataStore)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func LoadObjects[T comparable](schemaPath string, schemaObj *T, dataPath string, dataObject *map[string]interface{}) error {
	err := LoadFile(schemaPath, &schemaObj, false)
	if err != nil {
		return err
	}

	err = LoadFile(dataPath, &dataObject, false)
	if err != nil {
		return err
	}
	return nil
}
