package model

import (
	"encoding/json"
	"fmt"
	"github.com/zerok-ai/zk-utils-go/utils"
	"sort"
)

type Scenario struct {
	Version    string                  `json:"version"`
	ScenarioId string                  `json:"scenario_id"`
	Workloads  map[string]WorkloadRule `json:"workloads"`
	Filter     Filter                  `json:"filter"`
}

type WorkloadRule struct {
	Service   string    `json:"service,omitempty"`
	TraceRole TraceRole `json:"trace_role,omitempty"`
	Protocol  Protocol  `json:"protocol,omitempty"`
	Rule      Rule      `json:"rule,omitempty"`
}

type Rule struct {
	Type string `json:"type"`
	*RuleGroup
	*RuleLeaf
}

type RuleGroup struct {
	Condition *Condition `json:"condition,omitempty"`
	Rules     []Rule     `json:"rules,omitempty"`
}

type RuleLeaf struct {
	ID       *string        `json:"id,omitempty"`
	Field    *string        `json:"field,omitempty"`
	Datatype *DataType      `json:"datatype,omitempty"`
	Input    *InputTypes    `json:"input,omitempty"`
	Operator *OperatorTypes `json:"operator,omitempty"`
	Key      *string        `json:"key,omitempty"`
	Value    *ValueTypes    `json:"value,omitempty"`
}

type DataType string
type InputTypes = string
type OperatorTypes = string
type ValueTypes interface {
}
type Protocol string

const (
	MYSQL Protocol = "MYSQL"
	HTTP  Protocol = "HTTP"
)

const (
	server TraceRole = "server"
	client TraceRole = "client"
)

type TraceRole string

const (
	AND Condition = "AND"
	OR  Condition = "OR"
)

type Condition string

type Filter struct {
	Type        string       `json:"type"`
	Condition   Condition    `json:"condition"`
	Filters     *[]Filter    `json:"filters,omitempty"`
	WorkloadIds *WorkloadIds `json:"workload_ids,omitempty"`
}

type WorkloadIds []string

type Rules []Rule

func (a Rules) Len() int      { return len(a) }
func (a Rules) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a Rules) Less(i, j int) bool {
	if a[i].RuleGroup == nil && a[j].RuleGroup != nil {
		sort.Sort(Rules(a[j].Rules))
		return true
	} else if a[i].RuleGroup != nil && a[j].RuleGroup == nil {
		sort.Sort(Rules(a[i].Rules))
		return false
	} else if a[i].RuleGroup != nil && a[j].RuleGroup != nil {
		var ret *bool
		if a[i].Condition != a[j].Condition {
			ret = utils.ToPtr(*a[i].Condition < *a[j].Condition)
		} else if a[i].Type != a[j].Type {
			ret = utils.ToPtr(a[i].Type < a[j].Type)
		}

		sort.Sort(Rules(a[i].Rules))
		sort.Sort(Rules(a[j].Rules))

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

		if a[i].Type != a[j].Type {
			return *a[i].ID < *a[j].ID
		}

		if *a[i].ID != *a[j].ID {
			return *a[i].ID < *a[j].ID
		}

		if *a[i].Field != *a[j].Field {
			return *a[i].Field < *a[j].Field
		}

		if *a[i].Datatype != *a[j].Datatype {
			return *a[i].Datatype < *a[j].Datatype
		}

		if *a[i].Input != *a[j].Input {
			return *a[i].Input < *a[j].Input
		}

		if *a[i].Operator != *a[j].Operator {
			return *a[i].Operator < *a[j].Operator
		}

		if *a[i].Key != *a[j].Key {
			return *a[i].Key < *a[j].Key
		}

		x := (*(a[i].Value)).(string)
		y := (*(a[j].Value)).(string)
		if x != y {
			return x < y
		}

		return true
	}
}
