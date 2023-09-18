package files

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators"
	"testing"
)

func TestRuleEvaluation(t *testing.T) {

	var dataStore evaluators.DataStore
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/bool/equal/schema.json", &w, "test/files/ruleEvaluation/bool/equal/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationBooleanAnd(t *testing.T) {
	var dataStore evaluators.DataStore
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/bool/and/schema.json", &w, "test/files/ruleEvaluation/bool/and/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationBooleanOR(t *testing.T) {
	var dataStore evaluators.DataStore
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/bool/or/schema.json", &w, "test/files/ruleEvaluation/bool/or/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, false)
}

func TestRuleEvaluationBooleanNotEqual(t *testing.T) {
	var dataStore evaluators.DataStore
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/bool/not_equal/schema.json", &w, "test/files/ruleEvaluation/bool/not_equal/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationString(t *testing.T) {
	var dataStore evaluators.DataStore
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/string/schema.json", &w, "test/files/ruleEvaluation/string/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationInteger(t *testing.T) {
	var dataStore evaluators.DataStore
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/integer/schema.json", &w, "test/files/ruleEvaluation/integer/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func TestRuleEvaluationFloat(t *testing.T) {
	var dataStore evaluators.DataStore
	var w model.Workload
	err := loadObjects("test/files/ruleEvaluation/float/schema.json", &w, "test/files/ruleEvaluation/float/data.json", &dataStore)
	assert.NoError(t, err)

	validate(t, w, dataStore, true)
}

func validate(t *testing.T, w model.Workload, dataStore evaluators.DataStore, expected bool) {
	var result bool
	ruleEvaluator := evaluators.NewRuleEvaluator()
	result, err := ruleEvaluator.EvalRule(w.Rule, dataStore)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func loadObjects[T comparable](schemaPath string, schemaObj *T, dataPath string, dataObject *evaluators.DataStore) error {
	err := loadFile(schemaPath, &schemaObj, false)
	if err != nil {
		return err
	}

	err = loadFile(dataPath, &dataObject, true)
	if err != nil {
		return err
	}
	return nil
}

func loadFile(path string, dataObj interface{}, printObject bool) error {
	if dataObj == nil {
		return fmt.Errorf("dataObj is nil")
	}
	jsonString := string(common.GetBytesFromFile(path))
	if printObject {
		print("jsonString: ", jsonString, "\n")
	}
	return json.Unmarshal([]byte(jsonString), &dataObj)
}
