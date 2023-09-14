package evaluators

import (
	"fmt"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

const (
	loggerTag = "scenario-evaluator"

	typeString             = "string"
	typeKeyValue           = "key-value"
	typeInteger            = "integer"
	typeFloat              = "float"
	typeWorkLoadIdentifier = "workload-identifier"

	operatorExists = "exists"

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
type RuleEvaluator interface {
	init() RuleEvaluator
	EvalRule(rule model.Rule, store DataStore) (bool, error)
}

type BaseRuleEvaluator struct {
	dataSource     DataStore
	ruleEvaluators map[string]RuleEvaluator
}

func NewRuleEvaluator(dataSource DataStore) RuleEvaluator {
	return BaseRuleEvaluator{
		dataSource:     dataSource,
		ruleEvaluators: make(map[string]RuleEvaluator),
	}.init()
}

func (re BaseRuleEvaluator) init() RuleEvaluator {
	re.ruleEvaluators[model.RULE_GROUP] = RuleGroupEvaluator{re}.init()

	re.ruleEvaluators[typeString] = StringRuleEvaluator{}.init()
	re.ruleEvaluators[typeInteger] = IntegerRuleEvaluator{}.init()
	re.ruleEvaluators[typeFloat] = FloatRuleEvaluator{}.init()

	return re
}

func (re BaseRuleEvaluator) EvalRule(r model.Rule, store DataStore) (bool, error) {

	ruleEvaluator := re.ruleEvaluators[r.Type]
	if ruleEvaluator == nil {
		return false, fmt.Errorf("ruleEvaluator not found for type: %s", r.Type)
	}

	err := re.validate(r, store)
	if err != nil {
		return false, err
	}

	handled, value, err := re.handleCommonOperators(r, store)
	if !handled {
		return ruleEvaluator.EvalRule(r, store)
	}

	return value, err
}

// handleCommonOperators is a helper function to handle common operators like exists. The function returns
// a bool indicating if the rule is handled, a bool indicating the value, if handled and an error if any.
func (re BaseRuleEvaluator) handleCommonOperators(r model.Rule, store DataStore) (bool, bool, error) {
	operator := string(*r.Operator)
	//	switch on operator
	switch operator {
	case operatorExists:
		_, ok := store[*r.ID]
		return true, ok, nil
	}

	return false, false, nil
}

func (re BaseRuleEvaluator) validate(r model.Rule, store DataStore) error {
	id := r.ID
	operator := r.Operator
	valueFromRule := r.Value
	if id == nil || valueFromRule == nil || operator == nil {
		return fmt.Errorf("id or value or operator is nil")
	}

	_, ok := store[*id]
	if !ok {
		return fmt.Errorf("value for id: %s not found in store", *id)
	}
	return nil
}
