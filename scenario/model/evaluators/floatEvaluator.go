package evaluators

import (
	"fmt"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/functions"
	"strconv"
	"strings"
)

type FloatRuleEvaluator struct {
	functionFactory *functions.FunctionFactory
	attrStoreKey    *cache.AttribStoreKey
}

func (re *FloatRuleEvaluator) init() LeafRuleEvaluator {
	return re
}

func NewFloatRuleEvaluator(functionFactory *functions.FunctionFactory) LeafRuleEvaluator {
	return (&FloatRuleEvaluator{functionFactory: functionFactory}).init()
}

func (re *FloatRuleEvaluator) evalRule(rule model.Rule, valueStore map[string]interface{}) (bool, error) {

	defer func() {
		if r := recover(); r != nil {
			logger.ErrorF(LoggerTag, "In float eval: Recovered from panic: %v", r)
		}
	}()

	attributeID := *rule.RuleLeaf.ID

	// get the values assuming that the rule object is valid
	operator := string(*rule.Operator)

	//	switch on operator
	switch operator {

	case operatorExists:
		valueInterface, ok := re.functionFactory.EvaluateString(attributeID, valueStore, re.attrStoreKey)
		if !ok || valueInterface == nil {
			return false, nil
		}
		return true, nil
	case operatorNotExists:
		valueInterface, ok := re.functionFactory.EvaluateString(attributeID, valueStore, re.attrStoreKey)
		if ok && valueInterface != nil {
			return false, nil
		}
		return true, nil
	case operatorLessThan:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore < valueFromRule, nil
	case operatorLessThanEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore <= valueFromRule, nil
	case operatorGreaterThan:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore > valueFromRule, nil
	case operatorGreaterThanEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore >= valueFromRule, nil
	case operatorEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore == valueFromRule, nil
	case operatorNotEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore != valueFromRule, nil

	case operatorBetween:
		return re.isValueInRange(rule, attributeID, valueStore)
	case operatorNotBetween:
		valueInRange, err := re.isValueInRange(rule, attributeID, valueStore)
		return !valueInRange, err

	case operatorIn:
		isPresent := re.isValuePresentInCSV(rule, attributeID, valueStore)
		return isPresent, nil
	case operatorNotIn:
		isPresent := re.isValuePresentInCSV(rule, attributeID, valueStore)
		return !isPresent, nil

	}

	return false, fmt.Errorf("float: invalid operator: %s", operator)
}

func (re *FloatRuleEvaluator) getValuesFromCSString(csv string) []float64 {
	retArr := make([]float64, 0)
	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.ParseFloat(part, 64)
		if err != nil {
			logger.Error(LoggerTag, "error converting %s to integer: %v", stringSet[0], err)
			continue
		}
		retArr = append(retArr, number)
	}

	return retArr
}

func (re *FloatRuleEvaluator) isValueInRange(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (bool, error) {
	operator := string(*r.Operator)
	valueFromStore, err := re.valueFromStore(r, attributeNameOfID, valueStore)
	if err != nil {
		return false, err
	}

	numbers := re.getValuesFromCSString(string(*r.RuleLeaf.Value))
	if len(numbers) != 2 {
		return false, fmt.Errorf("invalid number of values for operator %s: %s", operator, string(*r.Value))
	}
	return valueFromStore >= numbers[0] && valueFromStore <= numbers[1], nil
}

func (re *FloatRuleEvaluator) isValuePresentInCSV(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) bool {

	csv := string(*r.Value)

	value, err := re.valueFromStore(r, attributeNameOfID, valueStore)
	if err != nil {
		logger.Error(LoggerTag, "%v", err)
		return false
	}

	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.ParseFloat(part, 64)
		if err != nil {
			logger.Error(LoggerTag, "error converting %s to integer: %v", stringSet[0], err)
			continue
		}
		if number == value {
			return true
		}
	}
	return false
}

func (re *FloatRuleEvaluator) valueFromRuleAndStore(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (float64, float64, error) {
	valueFromRule, err := strconv.ParseFloat(string(*r.Value), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error converting rule value %s to float: %v", string(*r.Value), err)
	}

	valueFromStore, err := re.valueFromStore(r, attributeNameOfID, valueStore)
	if err != nil {
		return 0, 0, err
	}

	return valueFromRule, valueFromStore, nil
}

func (re *FloatRuleEvaluator) valueFromStore(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (float64, error) {

	valueInterface, ok := re.functionFactory.EvaluateString(attributeNameOfID, valueStore, re.attrStoreKey)
	if !ok || valueInterface == nil {
		return 0, fmt.Errorf("value not found for id %s", attributeNameOfID)
	}

	valueFromStore, err1 := strconv.ParseFloat(fmt.Sprintf("%v", valueInterface), 64)
	if err1 != nil {
		return 0, fmt.Errorf("error converting valueStore value %s to float: %v", string(*r.Value), err1)
	}

	return valueFromStore, nil
}

func (re *FloatRuleEvaluator) setAttrStoreKey(attrStoreKey *cache.AttribStoreKey) {
	re.attrStoreKey = attrStoreKey
}
