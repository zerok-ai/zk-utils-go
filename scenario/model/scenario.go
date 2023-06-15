package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zerok-ai/zk-utils-go/crypto"
	"github.com/zerok-ai/zk-utils-go/interfaces"
	"reflect"
	"sort"
)

type Scenario struct {
	Version    string               `json:"version"`
	ScenarioId string               `json:"scenario_id"`
	Enabled    bool                 `json:"enabled"`
	Workloads  *map[string]Workload `json:"workloads"`
	Filter     Filter               `json:"filter"`
}

func (s Scenario) Equals(otherInterface interfaces.ZKComparable) bool {

	other, ok := otherInterface.(Scenario)
	if !ok {
		return false
	}

	if s.Version != other.Version || s.ScenarioId != other.ScenarioId || s.Enabled != other.Enabled {
		return false
	}

	if len(*s.Workloads) != len(*other.Workloads) {
		return false
	}
	for key, value := range *other.Workloads {
		val, ok := (*s.Workloads)[key]
		if !ok {
			return false
		}
		if !val.Equals(value) {
			return false
		}
	}

	if !s.Filter.Equals(other.Filter) {
		return false
	}

	return true
}

func (s Scenario) Less(other Scenario) bool {
	return s != other
}

type Workload struct {
	Service   string    `json:"service,omitempty"`
	TraceRole TraceRole `json:"trace_role,omitempty"`
	Protocol  Protocol  `json:"protocol,omitempty"`
	Rule      Rule      `json:"rule,omitempty"`
}

func (wr Workload) Equals(other Workload) bool {
	if wr.Service != other.Service || wr.TraceRole != other.TraceRole || wr.Protocol != other.Protocol {
		return false
	}

	if !wr.Rule.Equals(other.Rule) {
		return false
	}

	return true
}

type Rule struct {
	Type string `json:"type"`
	*RuleGroup
	*RuleLeaf
}

func (r Rule) Equals(other Rule) bool {
	if r.Type != other.Type {
		return false
	}

	if (r.RuleGroup == nil && other.RuleGroup != nil) || (r.RuleGroup != nil && other.RuleGroup == nil) {
		return false
	}

	if (r.RuleLeaf == nil && other.RuleLeaf != nil) || (r.RuleLeaf != nil && other.RuleLeaf == nil) {
		return false
	}

	if r.RuleGroup != nil && !(*r.RuleGroup).Equals(*other.RuleGroup) {
		return false
	}

	if r.RuleLeaf != nil && !(*r.RuleLeaf).Equals(*other.RuleLeaf) {
		return false
	}

	return true
}

type RuleGroup struct {
	Condition *Condition `json:"condition,omitempty"`
	Rules     Rules      `json:"rules,omitempty"`
}

func (r RuleGroup) Equals(other RuleGroup) bool {

	if *r.Condition != *other.Condition {
		return false
	}

	if (r.Rules == nil && other.Rules != nil) || (r.Rules != nil && other.Rules == nil) {
		return false
	} else if r.Rules != nil && other.Rules != nil {
		return r.Rules.Equals(other.Rules)
	}

	return true
}

func (r RuleGroup) LessThan(other RuleGroup) bool {
	if *r.Condition < *other.Condition {
		return true
	}

	if len(r.Rules) != len(other.Rules) {
		return len(r.Rules) < len(other.Rules)
	}

	for i := 0; i < len(r.Rules) && i < len(other.Rules); i++ {
		if r.Rules[i].Type < other.Rules[i].Type {
			return true
		} else if r.Rules[i].Type == other.Rules[i].Type {
			if r.Rules[i].RuleLeaf != nil && (*r.Rules[i].RuleLeaf).LessThan(*other.Rules[i].RuleLeaf) {
				return true
			}

			if r.Rules[i].RuleGroup != nil && (*r.Rules[i].RuleGroup).LessThan(*other.Rules[i].RuleGroup) {
				return true
			}
		}
	}

	return false
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

func (r RuleLeaf) Equals(other RuleLeaf) bool {
	return reflect.DeepEqual(r, other)
}

func (r RuleLeaf) LessThan(other RuleLeaf) bool {

	if *r.ID != *other.ID {
		return *r.ID < *other.ID
	}

	if *r.Field != *other.Field {
		return *r.Field < *other.Field
	}

	if *r.Datatype != *other.Datatype {
		return *r.Datatype < *other.Datatype
	}

	if *r.Input != *other.Input {
		return *r.Input < *other.Input
	}

	if *r.Operator != *other.Operator {
		return *r.Operator < *other.Operator
	}

	if *r.Key != *other.Key {
		return *r.Key < *other.Key
	}

	if *r.Value != *other.Value {
		return *r.Value < *other.Value
	}

	return false
}

type DataType string
type InputTypes string
type OperatorTypes string
type ValueTypes string
type Protocol string

type Rules []Rule

func (r Rules) Len() int      { return len(r) }
func (r Rules) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Rules) Less(i, j int) bool {
	rule := r[i]
	other := r[j]

	if rule.Type < other.Type {
		return true
	}

	if rule.Type == other.Type {
		if rule.RuleLeaf != nil {
			return (*rule.RuleLeaf).LessThan(*other.RuleLeaf)
		} else if rule.RuleGroup != nil {
			return (*rule.RuleGroup).LessThan(*other.RuleGroup)
		}
	}

	return false
}

func (r Rules) Sort() {
	for i := 0; i < len(r); i++ {
		if r[i].Type == RULE_GROUP {
			r[i].Rules.Sort()
		}
	}
	sort.Sort(r)
}

func (r Rules) Equals(other Rules) bool {

	if len(r) != len(other) {
		return false
	}

	// sort the arrays first
	r.Sort()
	other.Sort()

	// compare the arrays
	for index, value := range other {
		val := r[index]
		if !value.Equals(val) {
			fmt.Printf("in Rules Equals: Rule at index %d is not same\n", index)
			return false
		}
	}
	return true
}

const (
	MYSQL      Protocol = "MYSQL"
	HTTP       Protocol = "HTTP"
	RULE       string   = "rule"
	RULE_GROUP string   = "rule_group"
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

func WorkLoadUUID(w Workload) uuid.UUID {
	w.Rule.Rules.Sort()
	jStr, _ := json.Marshal(w)
	id := crypto.CalculateHash(string(jStr))
	return id
}
