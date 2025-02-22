package evaluators

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
)

const (
	LoggerTag = "scenario-evaluator"

	typeString             = "string"
	typeKeyValue           = "key-value"
	typeInteger            = "integer"
	typeFloat              = "float"
	typeBool               = "bool"
	typeWorkLoadIdentifier = "workload-identifier"

	operatorExists    = "exists"
	operatorNotExists = "not_exists"

	operatorMatches          = "matches"
	operatorDoesNotMatch     = "does_not_match"
	operatorEqual            = "equal"
	operatorNotEqual         = "not_equal"
	operatorContains         = "contains"
	operatorDoesNotContain   = "does_not_contain"
	operatorIn               = "in"
	operatorNotIn            = "not_in"
	operatorBeginsWith       = "begins_with"
	operatorDoesNotBeginWith = "does_not_begin_with"
	operatorEndsWith         = "ends_with"
	operatorDoesNotEndWith   = "does_not_end_with"

	operatorLessThan         = "less_than"
	operatorLessThanEqual    = "less_than_equal"
	operatorGreaterThan      = "greater_than"
	operatorGreaterThanEqual = "greater_than_equal"
	operatorBetween          = "between"
	operatorNotBetween       = "not_between"
)

type DataStore map[string]string

func (ds DataStore) String() string {
	str := "{"
	for k, v := range ds {
		str += fmt.Sprintf("%s:%s,", k, v)
	}
	str += "}"
	return str
}

type LeafRuleEvaluator interface {
	init() LeafRuleEvaluator
	setAttrStoreKey(attrStoreKey *cache.AttribStoreKey)
	evalRule(rule model.Rule, valueStore map[string]interface{}) (bool, error)
}

type GroupRuleEvaluator interface {
	init() GroupRuleEvaluator
	evalRule(rule model.Rule, attrStoreKey cache.AttribStoreKey, valueStore map[string]interface{}) (bool, error)
}

type RuleEvaluator struct {
	executorAttrStore *stores.ExecutorAttrStore
	podDetailsStore   *stores.LocalCacheHSetStore
	functionFactory   *functions.FunctionFactory

	leafRuleEvaluators map[string]LeafRuleEvaluator
	groupRuleEvaluator GroupRuleEvaluator
}

func NewRuleEvaluator(executorAttrStore *stores.ExecutorAttrStore, podDetailsStore *stores.LocalCacheHSetStore) *RuleEvaluator {
	return (&RuleEvaluator{
		executorAttrStore:  executorAttrStore,
		podDetailsStore:    podDetailsStore,
		leafRuleEvaluators: make(map[string]LeafRuleEvaluator),
	}).init()
}

func (re *RuleEvaluator) init() *RuleEvaluator {
	re.groupRuleEvaluator = (&RuleGroupEvaluator{re}).init()
	re.functionFactory = functions.NewFunctionFactory(re.podDetailsStore, re.executorAttrStore)

	re.leafRuleEvaluators[typeString] = NewStringRuleEvaluator(re.functionFactory)
	re.leafRuleEvaluators[typeInteger] = NewFloatRuleEvaluator(re.functionFactory)
	re.leafRuleEvaluators[typeFloat] = NewFloatRuleEvaluator(re.functionFactory)
	re.leafRuleEvaluators[typeBool] = NewBooleanEvaluator(re.functionFactory)

	return re
}

func (re *RuleEvaluator) EvalRule(rule model.Rule, attrStoreKey cache.AttribStoreKey, valueStore map[string]interface{}) (bool, error) {

	// reset the new attrStoreKey in all the leafRuleEvaluators. This pushes the new protocol version to all the leafRuleEvaluators
	for _, leafEvaluator := range re.leafRuleEvaluators {
		leafEvaluator.setAttrStoreKey(&attrStoreKey)
	}

	result, err := re.evalRule(rule, attrStoreKey, valueStore)
	return result, err
}

func (re *RuleEvaluator) evalRule(rule model.Rule, attrStoreKey cache.AttribStoreKey, valueStore map[string]interface{}) (bool, error) {

	value := false
	var err error
	if rule.Type == model.RULE_GROUP {
		value, err = re.groupRuleEvaluator.evalRule(rule, attrStoreKey, valueStore)
		zkLogger.DebugF(LoggerTag, "Evaluated value for group =%v, for condition=%s", value, *rule.RuleGroup.Condition)
	} else {
		err = re.validate(rule)
		if err != nil {
			return false, err
		}

		attributeID := *rule.RuleLeaf.ID
		zkLogger.DebugF(LoggerTag, "RuleId:- ruleID=%s, attributeID=%s", *rule.RuleLeaf.ID, attributeID)

		leafEvaluatorType := string(*rule.RuleLeaf.Datatype)

		ruleEvaluator := re.leafRuleEvaluators[leafEvaluatorType]
		if ruleEvaluator == nil {
			return false, fmt.Errorf("LeafRuleEvaluator not found for type: %s", leafEvaluatorType)
		}
		value, err = ruleEvaluator.evalRule(rule, valueStore)
		zkLogger.DebugF(LoggerTag, "Evaluated value=%v, for attributeID=%v", value, attributeID)
	}
	return value, err
}

func (re *RuleEvaluator) validate(r model.Rule) error {
	id := r.ID
	operator := r.Operator
	valueFromRule := r.Value
	dataType := r.Datatype
	if id == nil {
		return fmt.Errorf("id is nil")
	} else if operator == nil {
		return fmt.Errorf("operator is nil")
	} else if valueFromRule == nil {
		return fmt.Errorf("value is nil")
	} else if dataType == nil {
		return fmt.Errorf("datatype is nil")
	}

	return nil
}
