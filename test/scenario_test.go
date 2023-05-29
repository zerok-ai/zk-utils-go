package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/rules/model"
	"sort"
	"testing"
)

func TestScenarioMarshalUnMarshalSuccess(t *testing.T) {
	var s model.Scenario
	validScenarioJsonString := string(GetBytes("files/validScenarioJsonString.json"))

	err := json.Unmarshal([]byte(validScenarioJsonString), &s)
	assert.NoError(t, err)
	sJsonStr, err := json.Marshal(s)
	assert.NoError(t, err)
	dst := &bytes.Buffer{}
	err = json.Compact(dst, []byte(validScenarioJsonString))
	assert.NoError(t, err)
	assert.Equal(t, len(dst.String()), len(sJsonStr))

	nonScenarioJsonString := string(GetBytes("files/nonScenarioJsonString.json"))
	emptyScenarioJsonString := string(GetBytes("files/emptyScenarioJsonString.json"))

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
	workloadJS := string(GetBytes("files/unsortedWorkloadJs.json"))

	var w model.WorkloadRule
	err := json.Unmarshal([]byte(workloadJS), &w)
	sort.Sort(model.Rules(w.Rule.Rules))

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
