package evaluators

import (
	"fmt"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

type BooleanEvaluator struct {
	baseRuleEvaluator BaseRuleEvaluator
}

func (re BooleanEvaluator) init() RuleEvaluator {
	return re
}

func (re BooleanEvaluator) EvalRule(r model.Rule, store DataStore) (bool, error) {

	valueFromRule, err := getBooleanValue(string(*r.Value))
	if err != nil {
		return false, err
	}

	valueFromStore, err1 := getBooleanValue(store[*r.ID])
	if err1 != nil {
		return false, err1
	}

	//	switch on operator
	operator := string(*r.Operator)
	switch operator {

	case operatorEqual:
		return valueFromRule == valueFromStore, nil
	case operatorNotEqual:
		return valueFromRule != valueFromStore, nil

	}

	return false, fmt.Errorf("invalid operator: %s", operator)
}

func getBooleanValue(strValue string) (bool, error) {
	if strValue == "true" {
		return true, nil
	} else if strValue == "false" {
		return false, nil
	} else {
		return false, fmt.Errorf("invalid boolean value: %s", strValue)
	}
}
