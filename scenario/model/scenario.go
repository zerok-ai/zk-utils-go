package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zerok-ai/zk-utils-go/crypto"
	"github.com/zerok-ai/zk-utils-go/interfaces"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"reflect"
	"sort"
	"strconv"
)

var LogTag = "scenario_model"

type Scenario struct {
	Version   string               `json:"version"`
	Id        string               `json:"scenario_id"`
	Title     string               `json:"scenario_title"`
	Type      string               `json:"scenario_type"`
	Enabled   bool                 `json:"enabled"`
	Workloads *map[string]Workload `json:"workloads"`
	Filter    Filter               `json:"filter"`
}

func (s Scenario) Equals(otherInterface interfaces.ZKComparable) bool {

	other, ok := otherInterface.(Scenario)
	if !ok {
		return false
	}

	if s.Version != other.Version || s.Title != other.Title || s.Id != other.Id || s.Type != other.Type || s.Enabled != other.Enabled {
		return false
	}

	if (s.Workloads == nil && other.Workloads != nil) || (s.Workloads != nil && other.Workloads == nil) {
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

	if strconv.FormatBool(s.Enabled) < strconv.FormatBool(other.Enabled) || s.Id < other.Id || s.Version < other.Version {
		return true
	}

	return false
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

func (r Rule) String() string {
	return fmt.Sprintf("Rule{Type: %s, RuleGroup: %v, RuleLeaf: %v}", r.Type, r.RuleGroup, r.RuleLeaf)
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

	if r.RuleGroup != nil && other.RuleGroup != nil && !(*r.RuleGroup).Equals(*other.RuleGroup) {
		return false
	}

	if r.RuleLeaf != nil && other.RuleLeaf != nil && !(*r.RuleLeaf).Equals(*other.RuleLeaf) {
		return false
	}

	return true
}

type RuleGroup struct {
	Condition *Condition `json:"condition,omitempty"`
	Rules     Rules      `json:"rules,omitempty"`
}

func (r RuleGroup) Equals(other RuleGroup) bool {

	if r.Condition == nil && other.Condition != nil {
		return false
	} else if r.Condition != nil && other.Condition == nil {
		return false
	}

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

	if r.Condition == nil && other.Condition != nil {
		return true
	} else if r.Condition != nil && other.Condition == nil {
		return false
	}

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

			if r.Rules[i].RuleGroup != nil && other.Rules[i].RuleGroup == nil {
				return false
			} else if r.Rules[i].RuleGroup == nil && other.Rules[i].RuleGroup != nil {
				return true
			} else if r.Rules[i].RuleGroup != nil && other.Rules[i].RuleGroup != nil && (*r.Rules[i].RuleGroup).LessThan(*other.Rules[i].RuleGroup) {
				return true
			}

			if r.Rules[i].RuleLeaf != nil && other.Rules[i].RuleLeaf == nil {
				return true
			} else if r.Rules[i].RuleLeaf == nil && other.Rules[i].RuleLeaf != nil {
				return false
			} else if r.Rules[i].RuleLeaf != nil && other.Rules[i].RuleLeaf != nil && (*r.Rules[i].RuleLeaf).LessThan(*other.Rules[i].RuleLeaf) {
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
	Value    *ValueTypes    `json:"value,omitempty"`
}

func (r RuleLeaf) String() string {
	return fmt.Sprintf("RuleLeaf{ID: %v, Field: %v, Datatype: %v, Input: %v, Operator: %v, Value: %v}", *r.ID, *r.Field, *r.Datatype, *r.Input, *r.Operator, *r.Value)
}

func (r RuleLeaf) Equals(other RuleLeaf) bool {
	return reflect.DeepEqual(r, other)
}

func (r RuleLeaf) LessThan(other RuleLeaf) bool {

	comparison := stringCompare(r.ID, other.ID)
	if comparison != 0 {
		return comparison < 0
	}

	comparison = stringCompare(r.Field, other.Field)
	if comparison != 0 {
		return comparison < 0
	}

	if r.Datatype == nil && other.Datatype != nil {
		return true
	} else if r.Datatype != nil && other.Datatype == nil {
		return false
	} else if r.Datatype != nil && other.Datatype != nil && *r.Datatype != *other.Datatype {
		if *r.Datatype < *other.Datatype {
			return true
		}
		return false
	}

	if r.Input == nil && other.Input != nil {
		return true
	} else if r.Input != nil && other.Input == nil {
		return false
	} else if r.Input != nil && other.Input != nil && *r.Input != *other.Input {
		if *r.Input < *other.Input {
			return true
		}
		return false
	}

	if r.Operator == nil && other.Operator != nil {
		return true
	} else if r.Operator != nil && other.Operator == nil {
		return false
	} else if r.Operator != nil && other.Operator != nil && *r.Operator != *other.Operator {
		if *r.Operator < *other.Operator {
			return true
		}
		return false
	}

	if r.Value == nil && other.Value != nil {
		return true
	} else if r.Value != nil && other.Value == nil {
		return false
	} else if r.Value != nil && other.Value != nil && *r.Value != *other.Value {
		if *r.Value < *other.Value {
			return true
		}
		return false
	}
	fmt.Println("before returning false")
	return false
}

func stringCompare(a *string, b *string) int {
	if a == nil && b != nil {
		return -1
	} else if a != nil && b == nil {
		return 1
	} else if (a == nil && b == nil) || *a == *b {
		return 0
	}

	if *a < *b {
		return -1
	}
	return 1
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

		if rule.RuleGroup == nil && other.RuleGroup != nil {
			return true
		} else if rule.RuleGroup != nil && other.RuleGroup == nil {
			return false
		} else if rule.RuleGroup != nil && other.RuleGroup != nil {
			return (*rule.RuleGroup).LessThan(*other.RuleGroup)
		}

		if rule.RuleLeaf == nil && other.RuleLeaf != nil {
			return true
		} else if rule.RuleLeaf != nil && other.RuleLeaf == nil {
			return false
		} else if rule.RuleLeaf != nil && other.RuleLeaf != nil {
			return (*rule.RuleLeaf).LessThan(*other.RuleLeaf)
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
			zkLogger.Error(LogTag, "in Rules Equals: Rule at index %d is not same\n", index)
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
	id := crypto.CalculateHashNewSHA2(string(jStr))
	return id
}
