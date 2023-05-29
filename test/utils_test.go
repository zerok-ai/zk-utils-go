package test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/rules/model"
	"github.com/zerok-ai/zk-utils-go/utils"
	"sort"
	"testing"
)

func TestCalculateHash(t *testing.T) {
	input := "example string"
	expected := "0d355fb4-a093-520b-9d7e-f1f615e3fa10"
	result := utils.CalculateHash(input)

	if result.String() != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestSortedScenariosHash(t *testing.T) {
	unsortedWorkloadJS := string(GetBytes("files/unsortedWorkloadJs.json"))

	var wUnsorted model.WorkloadRule
	errUnsorted := json.Unmarshal([]byte(unsortedWorkloadJS), &wUnsorted)
	assert.NoError(t, errUnsorted)
	sort.Sort(model.Rules(wUnsorted.Rule.Rules))

	sortedWorkloadJS := string(GetBytes("files/sortedWorkloadJs.json"))

	var wSorted model.WorkloadRule
	errSorted := json.Unmarshal([]byte(sortedWorkloadJS), &wSorted)
	assert.NoError(t, errSorted)
	sort.Sort(model.Rules(wSorted.Rule.Rules))

	assert.Equal(t, model.WorkLoadUUID(wUnsorted), model.WorkLoadUUID(wSorted))
}
