package model

import (
	"encoding/json"
	"fmt"
	"github.com/zerok-ai/zk-utils-go/utils"
	"sort"
)

type DataTypes string
type InputTypes = string
type OperatorTypes = string
type ValueTypes interface {
}
type Rule struct {
	ID       string        `json:"id,omitempty"`
	Field    string        `json:"field,omitempty"`
	Type     DataTypes     `json:"type,omitempty"`
	Input    InputTypes    `json:"input,omitempty"`
	Operator OperatorTypes `json:"operator,omitempty"`
	Key      string        `json:"key,omitempty"`
	Value    ValueTypes    `json:"value,omitempty"`
}

type WorkloadRule struct {
	Service         *string         `json:"service,omitempty"`
	TraceRole       *string         `json:"trace_role,omitempty"`
	Protocol        *string         `json:"protocol,omitempty"`
	ConditionalRule ConditionalRule `json:"conditional_rule,omitempty"`
}

type ConditionalRule struct {
	Condition *string   `json:"condition,omitempty"`
	RuleSet   []RuleSet `json:"rules,omitempty"`
}

type RuleSet struct {
	Rule
	ConditionalRule
}

type FilterRule struct {
	Version   int                     `json:"version"`
	Workloads map[string]WorkloadRule `json:"workloads"`
	FilterId  string                  `json:"filter_id"`
	Filters   []Filters               `json:"filters"`
}

type Filters struct {
	Type        string    `json:"type"`
	Condition   string    `json:"condition"`
	Filters     []Filters `json:"filters"`
	WorkloadSet []string  `json:"workload_id_set"`
}

type Rules []RuleSet

func (a Rules) Len() int      { return len(a) }
func (a Rules) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Rules) Less(i, j int) bool {
	if a[i].Condition == nil && a[j].Condition != nil {
		sort.Sort(Rules(a[j].RuleSet))
		return true
	} else if a[i].Condition != nil && a[j].Condition == nil {
		sort.Sort(Rules(a[i].RuleSet))
		return false
	} else if a[i].Condition != nil && a[j].Condition != nil {
		var ret *bool
		if *a[i].Condition != *a[j].Condition {
			ret = utils.ToPtr(*a[i].Condition < *a[j].Condition)
		}

		sort.Sort(Rules(a[i].RuleSet))
		sort.Sort(Rules(a[j].RuleSet))

		str1, err1 := json.Marshal(a[i])
		str2, err2 := json.Marshal(a[j])

		if err1 != nil {
			fmt.Println("err1: ", err1)
			panic(err1)
		}

		if err1 != nil {
			fmt.Println("err2: ", err2)
			panic(err1)
		}

		if ret == nil {
			ret = utils.ToPtr(string(str1) < string(str2))
		}

		return *ret
	} else {

		if a[i].ID != a[j].ID {
			return a[i].ID < a[j].ID
		}

		if a[i].Field != a[j].Field {
			return a[i].Field < a[j].Field
		}

		if a[i].Type != a[j].Type {
			return a[i].Type < a[j].Type
		}

		if a[i].Input != a[j].Input {
			return a[i].Input < a[j].Input
		}

		if a[i].Operator != a[j].Operator {
			return a[i].Operator < a[j].Operator
		}

		if a[i].Key != a[j].Key {
			return a[i].Key < a[j].Key
		}

		x := a[i].Value.(string)
		y := a[j].Value.(string)
		if x != y {
			return x < y
		}

		return true
	}
}
