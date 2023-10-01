package evaluators

import (
	"fmt"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

type BooleanEvaluator struct {
	baseRuleEvaluator RuleEvaluator
}

func (re BooleanEvaluator) init() LeafRuleEvaluator {
	return re
}

func (re BooleanEvaluator) evalRule(rule model.Rule, valueStore map[string]interface{}) (bool, error) {

	valueFromRule, err := getBooleanValue(string(*rule.Value))
	if err != nil {
		return false, err
	}

	// get the value from the value store
	value, ok := GetValueFromStore(*rule.RuleLeaf.ID, valueStore)
	if !ok {
		return false, fmt.Errorf("value for id: %s not found in valueStore", *rule.RuleLeaf.ID)
	}
	valueFromStore, err1 := getBooleanValue(value)
	if err1 != nil {
		return false, err1
	}

	//	switch on operator
	operator := string(*rule.Operator)
	switch operator {

	case operatorEqual:
		return valueFromRule == valueFromStore, nil
	case operatorNotEqual:
		return valueFromRule != valueFromStore, nil

	}

	return false, fmt.Errorf("bool: invalid operator: %s", operator)
}

func getBooleanValue(value interface{}) (bool, error) {
	strValue := fmt.Sprintf("%v", value)
	if strValue == "true" {
		return true, nil
	} else if strValue == "false" {
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean value: %s", strValue)
}
