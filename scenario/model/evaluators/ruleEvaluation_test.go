package evaluators

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"testing"
)

func TestSort(t *testing.T) {
	workloadJS := string(common.GetBytesFromFile("test/files/unsortedWorkloadJs.json"))
	print("workloadJS: ", workloadJS, "\n")

	var w model.Workload
	err := json.Unmarshal([]byte(workloadJS), &w)
	assert.NoError(t, err)

	w.Rule.Rules.Sort()
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
