package evaluators

import (
	"fmt"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
	"regexp"
	"strings"
)

type StringRuleEvaluator struct {
	functionFactory *functions.FunctionFactory
	attrStoreKey    *cache.AttribStoreKey
}

func (re *StringRuleEvaluator) init() LeafRuleEvaluator {
	return re
}

func NewStringRuleEvaluator(functionFactory *functions.FunctionFactory) LeafRuleEvaluator {
	return (&StringRuleEvaluator{functionFactory: functionFactory}).init()
}

func (re *StringRuleEvaluator) evalRule(rule model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (bool, error) {

	// get the values assuming that the rule object is valid
	operator := string(*rule.Operator)
	valueFromRule := string(*rule.Value)

	valueFromStoreI, ok := GetValueFromStore(attributeNameOfID, valueStore, re.functionFactory, re.attrStoreKey)
	if !ok {
		return false, fmt.Errorf("value for attributeName: %s not found in valueStore", attributeNameOfID)
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

func (re *StringRuleEvaluator) setAttrStoreKey(attrStoreKey *cache.AttribStoreKey) {
	re.attrStoreKey = attrStoreKey
}
