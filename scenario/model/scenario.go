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
	"strings"
)

//+k8s:deepcopy-gen=true

var LogTag = "scenario_model"

// +k8s:deepcopy-gen=true
type Scenario struct {
	Version   string               `json:"version"`
	Id        string               `json:"scenario_id"`
	Title     string               `json:"scenario_title"`
	Type      string               `json:"scenario_type"`
	Enabled   bool                 `json:"enabled"`
	Workloads *map[string]Workload `json:"workloads"`
	Filter    Filter               `json:"filter"`
	GroupBy   []GroupBy            `json:"group_by"`
	RateLimit []RateLimit          `json:"rate_limit"`
}

func (s Scenario) Equals(otherInterface interfaces.ZKComparable) bool {

	other, ok := otherInterface.(Scenario)
	if !ok {
		return false
	}

	if s.Version != other.Version || s.Title != other.Title || s.Id != other.Id || s.Type != other.Type || s.Enabled != other.Enabled {
		return false
	}

	if (s.GroupBy == nil && other.GroupBy != nil) || (s.GroupBy != nil && other.GroupBy == nil) {
		return false
	}

	if s.GroupBy != nil && other.GroupBy != nil && (len(s.GroupBy) != len(other.GroupBy)) {
		return false
	}

	sort.Slice(s.GroupBy, func(i, j int) bool {
		return s.GroupBy[i].LessThan(s.GroupBy[j])
	})

	for i, groupBy := range s.GroupBy {
		otherGroupBy := other.GroupBy[i]
		if !groupBy.Equals(otherGroupBy) {
			return false
		}
	}

	if (s.RateLimit == nil && other.RateLimit != nil) || (s.RateLimit != nil && other.RateLimit == nil) {
		return false
	}

	if s.RateLimit != nil && other.RateLimit != nil && (len(s.RateLimit) != len(other.RateLimit)) {
		return false
	}

	sort.Slice(s.RateLimit, func(i, j int) bool {
		return s.RateLimit[i].LessThan(s.RateLimit[j])
	})

	for i, rateLimit := range s.RateLimit {
		otherRateLimit := other.RateLimit[i]
		if !rateLimit.Equals(otherRateLimit) {
			return false
		}
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

func (g GroupBy) Equals(other GroupBy) bool {
	if g.WorkloadId == other.WorkloadId && g.Title == other.Title && g.Hash == other.Hash {
		return true
	}

	return false
}

func (r RateLimit) Equals(o RateLimit) bool {
	if r.BucketMaxSize == o.BucketMaxSize && r.BucketRefillSize == o.BucketRefillSize && r.TickDuration == o.TickDuration {
		return true
	}

	return false
}

func (s Scenario) Less(other Scenario) bool {

	if strconv.FormatBool(s.Enabled) < strconv.FormatBool(other.Enabled) || s.Id < other.Id || s.Version < other.Version {
		return true
	}

	return false
}

// ref : https://github.com/kubernetes-sigs/controller-tools/issues/585
// +k8s:deepcopy-gen=true
type Workload struct {
	Executor  ExecutorName `json:"executor"`
	Service   string       `json:"service,omitempty"`
	TraceRole TraceRole    `json:"trace_role,omitempty"`
	Protocol  ProtocolName `json:"protocol,omitempty"`
	Rule      Rule         `json:"rule,omitempty"`
}

func (wr Workload) GetNamespaceAndWorkloadName() (string, string, error) {
	if wr.Service == "" {
		return "", "", fmt.Errorf("service name not present")
	}
	// split Service by '/' to obtain namespace and workload name
	serviceSplit := strings.Split(wr.Service, "/")
	if len(serviceSplit) != 2 {
		return "", "", fmt.Errorf("invalid service name")
	}
	return serviceSplit[0], serviceSplit[1], nil
}

func (wr Workload) Equals(other Workload) bool {
	if wr.Executor != other.Executor || wr.Service != other.Service || wr.TraceRole != other.TraceRole || wr.Protocol != other.Protocol {
		return false
	}

	if !wr.Rule.Equals(other.Rule) {
		return false
	}

	return true
}

// +k8s:deepcopy-gen=true
type Rule struct {
	Type       string `json:"type"`
	*RuleGroup `json:""`
	*RuleLeaf  `json:""`
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

// +k8s:deepcopy-gen=true
type RuleGroup struct {
	Condition *Condition `json:"condition,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Rules Rules `json:"rules,omitempty"`
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

// +k8s:deepcopy-gen=true
type RuleLeaf struct {
	ID       *string        `json:"id,omitempty"`
	Field    *string        `json:"field,omitempty"`
	Datatype *DataType      `json:"datatype,omitempty"`
	Input    *InputTypes    `json:"input,omitempty"`
	Operator *OperatorTypes `json:"operator,omitempty"`
	Value    *ValueTypes    `json:"value,omitempty"`
	JsonPath *[]string      `json:"json_path,omitempty"`
}

func (r RuleLeaf) String() string {
	return fmt.Sprintf("RuleLeaf{ID: %v, Field: %v, Datatype: %v, Input: %v, Operator: %v, Value: %v, JsonPath: %v}", *r.ID, *r.Field, *r.Datatype, *r.Input, *r.Operator, *r.Value, *r.JsonPath)
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

	if r.JsonPath == nil && other.JsonPath != nil {
		return true
	} else if r.JsonPath != nil && other.JsonPath == nil {
		return false
	} else if r.JsonPath != nil && other.JsonPath != nil && !reflect.DeepEqual(*r.JsonPath, *other.JsonPath) {
		rLength := len(*r.JsonPath)
		otherLength := len(*other.JsonPath)
		if rLength < otherLength {
			return true
		} else if rLength > otherLength {
			return false
		}

		for i := 0; i < rLength; i++ {
			if (*r.JsonPath)[i] < (*other.JsonPath)[i] {
				return true
			} else if (*r.JsonPath)[i] > (*other.JsonPath)[i] {
				return false
			}
		}
		return false
	}

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
type ProtocolName string
type ExecutorName string

const (
	ExecutorEbpf ExecutorName = "EBPF"
	ExecutorOTel ExecutorName = "OTEL"

	ProtocolHTTP       ProtocolName = "HTTP"
	ProtocolGRPC       ProtocolName = "GRPC"
	ProtocolGeneral    ProtocolName = "GENERAL"
	ProtocolIdentifier ProtocolName = "IDENTIFIER"
)

// +k8s:deepcopy-gen=true
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
	MYSQL      ProtocolName = "MYSQL"
	HTTP       ProtocolName = "HTTP"
	RULE       string       = "rule"
	RULE_GROUP string       = "rule_group"
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

// +k8s:deepcopy-gen=true
type Condition string

func WorkLoadUUID(w Workload) uuid.UUID {
	w.Rule.Rules.Sort()
	jStr, _ := json.Marshal(w)
	id := crypto.CalculateHashNewSHA2(string(jStr))
	return id
}

// +k8s:deepcopy-gen=true
type GroupBy struct {
	WorkloadId string `json:"workload_id"`
	Title      string `json:"title"`
	Hash       string `json:"hash"`
}

func (g GroupBy) LessThan(other GroupBy) bool {
	if g.WorkloadId < other.WorkloadId {
		return true
	}
	if g.WorkloadId == other.WorkloadId {
		if g.Title < other.Title {
			return true
		}
		if g.Title == other.Title {
			if g.Hash < other.Hash {
				return true
			}
		}
	}
	return false
}

// +k8s:deepcopy-gen=true
type RateLimit struct {
	BucketMaxSize    int    `json:"bucket_max_size"`
	BucketRefillSize int    `json:"bucket_refill_size"`
	TickDuration     string `json:"tick_duration"`
}

func (r RateLimit) LessThan(other RateLimit) bool {
	if r.BucketMaxSize < other.BucketMaxSize {
		return true
	}
	if r.BucketMaxSize == other.BucketMaxSize {
		if r.BucketRefillSize < other.BucketRefillSize {
			return true
		}
		if r.BucketRefillSize == other.BucketRefillSize {
			if r.TickDuration < other.TickDuration {
				return true
			}
		}
	}
	return false
}

//deep copy methods

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupBy) DeepCopyInto(out *GroupBy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupBy.
func (in *GroupBy) DeepCopy() *GroupBy {
	if in == nil {
		return nil
	}
	out := new(GroupBy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimit) DeepCopyInto(out *RateLimit) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimit.
func (in *RateLimit) DeepCopy() *RateLimit {
	if in == nil {
		return nil
	}
	out := new(RateLimit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	if in.RuleGroup != nil {
		in, out := &in.RuleGroup, &out.RuleGroup
		*out = new(RuleGroup)
		(*in).DeepCopyInto(*out)
	}
	if in.RuleLeaf != nil {
		in, out := &in.RuleLeaf, &out.RuleLeaf
		*out = new(RuleLeaf)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleGroup) DeepCopyInto(out *RuleGroup) {
	*out = *in
	if in.Condition != nil {
		in, out := &in.Condition, &out.Condition
		*out = new(Condition)
		**out = **in
	}
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make(Rules, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleGroup.
func (in *RuleGroup) DeepCopy() *RuleGroup {
	if in == nil {
		return nil
	}
	out := new(RuleGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuleLeaf) DeepCopyInto(out *RuleLeaf) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Field != nil {
		in, out := &in.Field, &out.Field
		*out = new(string)
		**out = **in
	}
	if in.Datatype != nil {
		in, out := &in.Datatype, &out.Datatype
		*out = new(DataType)
		**out = **in
	}
	if in.Input != nil {
		in, out := &in.Input, &out.Input
		*out = new(InputTypes)
		**out = **in
	}
	if in.Operator != nil {
		in, out := &in.Operator, &out.Operator
		*out = new(OperatorTypes)
		**out = **in
	}
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(ValueTypes)
		**out = **in
	}
	if in.JsonPath != nil {
		in, out := &in.JsonPath, &out.JsonPath
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuleLeaf.
func (in *RuleLeaf) DeepCopy() *RuleLeaf {
	if in == nil {
		return nil
	}
	out := new(RuleLeaf)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Rules) DeepCopyInto(out *Rules) {
	{
		in := &in
		*out = make(Rules, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rules.
func (in Rules) DeepCopy() Rules {
	if in == nil {
		return nil
	}
	out := new(Rules)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scenario) DeepCopyInto(out *Scenario) {
	*out = *in
	if in.Workloads != nil {
		in, out := &in.Workloads, &out.Workloads
		*out = new(map[string]Workload)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]Workload, len(*in))
			for key, val := range *in {
				(*out)[key] = *val.DeepCopy()
			}
		}
	}
	in.Filter.DeepCopyInto(&out.Filter)
	if in.GroupBy != nil {
		in, out := &in.GroupBy, &out.GroupBy
		*out = make([]GroupBy, len(*in))
		copy(*out, *in)
	}
	if in.RateLimit != nil {
		in, out := &in.RateLimit, &out.RateLimit
		*out = make([]RateLimit, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scenario.
func (in *Scenario) DeepCopy() *Scenario {
	if in == nil {
		return nil
	}
	out := new(Scenario)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Workload) DeepCopyInto(out *Workload) {
	*out = *in
	in.Rule.DeepCopyInto(&out.Rule)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Workload.
func (in *Workload) DeepCopy() *Workload {
	if in == nil {
		return nil
	}
	out := new(Workload)
	in.DeepCopyInto(out)
	return out
}
