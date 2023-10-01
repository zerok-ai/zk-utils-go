package evaluators

import (
	"fmt"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"regexp"
	"strings"
)

type StringRuleEvaluator struct {
}

func (re StringRuleEvaluator) init() LeafRuleEvaluator {
	return re
}

func (re StringRuleEvaluator) evalRule(rule model.Rule, valueStore map[string]interface{}) (bool, error) {

	// get the values assuming that the rule object is valid
	operator := string(*rule.Operator)
	valueFromRule := string(*rule.Value)

	valueFromStoreI, ok := valueStore[*rule.RuleLeaf.ID]
	if !ok {
		return false, fmt.Errorf("value for id: %s not found in valueStore", *rule.RuleLeaf.ID)
	}
	valueFromStore := fmt.Sprintf("%v", valueFromStoreI)

	//	switch on operator
	switch operator {

	case operatorMatches:
		matched, _ := regexp.MatchString(valueFromRule, valueFromStore)
		return matched, nil
	case operatorDoesNotMatch:
		matched, _ := regexp.MatchString(valueFromRule, valueFromStore)
		return !matched, nil
	case operatorEqual:
		return valueFromStore == valueFromRule, nil
	case operatorNotEqual:
		return valueFromStore != valueFromRule, nil
	case operatorContains:
		return strings.Contains(valueFromStore, valueFromRule), nil
	case operatorDoesNotContain:
		return !strings.Contains(valueFromStore, valueFromRule), nil
	case operatorIn:
		stringSet := strings.Split(valueFromRule, ",")
		for _, value := range stringSet {
			if valueFromStore == value {
				return true, nil
			}
		}
		return false, nil
	case operatorNotIn:
		stringSet := strings.Split(valueFromRule, ",")
		for _, value := range stringSet {
			if valueFromStore == value {
				return false, nil
			}
		}
		return true, nil
	case operatorBeginsWith:
		return strings.HasPrefix(valueFromStore, valueFromRule), nil
	case operatorDoesNotBeginWith:
		return !strings.HasPrefix(valueFromStore, valueFromRule), nil
	case operatorEndsWith:
		return strings.HasSuffix(valueFromStore, valueFromRule), nil
	case operatorDoesNotEndWith:
		return !strings.HasSuffix(valueFromStore, valueFromRule), nil

	}

	return false, fmt.Errorf("string: invalid operator: %s", operator)
}
