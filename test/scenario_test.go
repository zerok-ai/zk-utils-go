package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"log"
	"testing"
)

func TestScenarioMarshalUnMarshalSuccess(t *testing.T) {
	var s model.Scenario
	validScenarioJsonString := string(zkcommon.GetBytesFromFile("files/validScenarioJsonString.json"))

	err := json.Unmarshal([]byte(validScenarioJsonString), &s)
	assert.NoError(t, err)
	sJsonStr, err := json.Marshal(s)
	assert.NoError(t, err)
	dst := &bytes.Buffer{}
	err = json.Compact(dst, []byte(validScenarioJsonString))
	assert.NoError(t, err)
	assert.Equal(t, len(dst.String()), len(sJsonStr))

	nonScenarioJsonString := string(zkcommon.GetBytesFromFile("files/nonScenarioJsonString.json"))
	emptyScenarioJsonString := string(zkcommon.GetBytesFromFile("files/emptyScenarioJsonString.json"))

	var s2 model.Scenario
	_ = json.Unmarshal([]byte(nonScenarioJsonString), &s2)
	sJsonStr, err = json.Marshal(s2)
	assert.NoError(t, err)
	dst2 := &bytes.Buffer{}
	err = json.Compact(dst2, []byte(emptyScenarioJsonString))
	assert.NoError(t, err)
	assert.Equal(t, len(dst2.String()), len(sJsonStr))

}

func TestSort(t *testing.T) {
	workloadJS := string(zkcommon.GetBytesFromFile("files/unsortedWorkloadJs.json"))

	var w model.Workload
	err := json.Unmarshal([]byte(workloadJS), &w)
	w.Rule.Rules.Sort()

	assert.NoError(t, err)
	assert.Equal(t, string(*w.Rule.RuleGroup.Condition), "AND")
	assert.Equal(t, len(w.Rule.RuleGroup.Rules), 4)
	assert.Equal(t, w.Rule.RuleGroup.Rules[0].Type, "rule")
	assert.Equal(t, w.Rule.RuleGroup.Rules[1].Type, "rule")
	assert.Equal(t, w.Rule.RuleGroup.Rules[2].Type, "rule")
	assert.Equal(t, w.Rule.RuleGroup.Rules[3].Type, "rule_group")

	assert.Equal(t, *w.Rule.RuleGroup.Rules[0].ID, "id_place_1")
	assert.Equal(t, *w.Rule.RuleGroup.Rules[1].ID, "req_path_place_2")

	assert.Equal(t, string(*w.Rule.RuleGroup.Rules[3].RuleGroup.Condition), "OR")
	assert.Equal(t, len(w.Rule.RuleGroup.Rules[3].RuleGroup.Rules), 2)
	assert.Equal(t, *w.Rule.RuleGroup.Rules[3].RuleGroup.Rules[0].RuleLeaf.ID, "req_method_place_1")
	assert.Equal(t, *w.Rule.RuleGroup.Rules[3].RuleGroup.Rules[1].RuleLeaf.ID, "req_path_place_2")
}

func TestScenarioEqualitySuccess(t *testing.T) {
	var scenario1 model.Scenario
	validScenarioJsonString := string(zkcommon.GetBytesFromFile("files/validScenarioJsonString.json"))
	err := json.Unmarshal([]byte(validScenarioJsonString), &scenario1)
	assert.NoError(t, err)

	var scenario2 model.Scenario
	validScenarioJsonString = string(zkcommon.GetBytesFromFile("files/validScenarioJsonString1.json"))
	err = json.Unmarshal([]byte(validScenarioJsonString), &scenario2)
	assert.NoError(t, err)

	log.Default().Println("Checking equality using equals = ", scenario1.Equals(scenario2))
	log.Default().Println("Calling assert ")
	assert.Equal(t, scenario1, scenario2)
}
