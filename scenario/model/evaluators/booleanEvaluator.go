package evaluators

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
)

type BooleanEvaluator struct {
	functionFactory *functions.FunctionFactory
	attrStoreKey    *cache.AttribStoreKey
}

func (re *BooleanEvaluator) init() LeafRuleEvaluator {
	return re
}

func NewBooleanEvaluator(functionFactory *functions.FunctionFactory) LeafRuleEvaluator {
	return (&BooleanEvaluator{functionFactory: functionFactory}).init()
}

func (re *BooleanEvaluator) evalRule(rule model.Rule, valueStore map[string]interface{}) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In bool eval: Recovered from panic: %v", r)
		}
	}()

	attributeID := *rule.RuleLeaf.ID

	operator := string(*rule.Operator)

	// get the value from the value store
	value, ok := re.functionFactory.EvaluateString(attributeID, valueStore, re.attrStoreKey)

	switch operator {
	case operatorExists:
		if !ok || value == nil {
			return false, nil
		}
		return true, nil
	case operatorNotExists:
		if ok && value != nil {
			return false, nil
		}
		return true, nil
	}

	if !ok {
		return false, fmt.Errorf("value for attributeName: %s not found in valueStore", attributeID)
	}

	valueFromStore, err1 := getBooleanValue(value)
	if err1 != nil {
		return false, err1
	}

	valueFromRule, err := getBooleanValue(string(*rule.Value))
	if err != nil {
		return false, err
	}

	//	switch on operator
	switch operator {

	case operatorEqual:
		return valueFromRule == valueFromStore, nil
	case operatorNotEqual:
		return valueFromRule != valueFromStore, nil

	}

	return false, fmt.Errorf("bool: invalid operator: %s", operator)
}

func (re *BooleanEvaluator) setAttrStoreKey(attrStoreKey *cache.AttribStoreKey) {
	re.attrStoreKey = attrStoreKey
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
