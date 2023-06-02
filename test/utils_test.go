package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/crypto"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"sort"
	"testing"
)

func TestCalculateHash(t *testing.T) {
	input := "example string"
	expected := "0d355fb4-a093-520b-9d7e-f1f615e3fa10"
	result := crypto.CalculateHash(input)

	if result.String() != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestSortedScenariosHash(t *testing.T) {
	unsortedWorkloadJS := string(GetBytesFromFile("files/unsortedWorkloadJs.json"))
	var wUnsorted model.Workload
	errUnsorted := json.Unmarshal([]byte(unsortedWorkloadJS), &wUnsorted)
	assert.NoError(t, errUnsorted)
	wUnsorted.Rule.Rules.Sort()

	x, _ := json.Marshal(wUnsorted)
	fmt.Print(x)

	sortedWorkloadJS := string(GetBytesFromFile("files/sortedWorkloadJs.json"))
	var wSorted model.Workload
	errSorted := json.Unmarshal([]byte(sortedWorkloadJS), &wSorted)
	assert.NoError(t, errSorted)
	sort.Sort(wSorted.Rule.Rules)
	a, b := model.WorkLoadUUID(wUnsorted), model.WorkLoadUUID(wSorted)

	fmt.Print(a, b)
	if a == b {
		fmt.Println("equal")
	} else {
		fmt.Println(a.String() + ":::" + b.String())
	}
	fmt.Println(a.String() + "pppp" + b.String())

	assert.Equal(t, a, b)
}
