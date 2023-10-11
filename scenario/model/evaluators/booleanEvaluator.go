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

func (re *BooleanEvaluator) evalRule(rule model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In bool eval: Recovered from panic: %v", r)
		}
	}()
	valueFromRule, err := getBooleanValue(string(*rule.Value))
	if err != nil {
		return false, err
	}

	// get the value from the value store
	value, ok := GetValueFromStore(attributeNameOfID, valueStore, re.functionFactory, re.attrStoreKey)
	if !ok {
		return false, fmt.Errorf("value for attributeName: %s not found in valueStore", attributeNameOfID)
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
