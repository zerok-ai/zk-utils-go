package evaluators

import (
	"fmt"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"strconv"
	"strings"
)

type IntegerRuleEvaluator struct {
}

func (re IntegerRuleEvaluator) init() LeafRuleEvaluator {
	return re
}

func (re IntegerRuleEvaluator) evalRule(rule model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (bool, error) {

	// get the values assuming that the rule object is valid
	operator := string(*rule.Operator)

	//	switch on operator
	switch operator {

	case operatorLessThan:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeNameOfID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore < valueFromRule, nil
	case operatorLessThanEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeNameOfID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore <= valueFromRule, nil
	case operatorGreaterThan:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeNameOfID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore > valueFromRule, nil
	case operatorGreaterThanEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeNameOfID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore >= valueFromRule, nil
	case operatorEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeNameOfID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore == valueFromRule, nil
	case operatorNotEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(rule, attributeNameOfID, valueStore)
		if err != nil {
			return false, err
		}
		return valueFromStore != valueFromRule, nil

	case operatorBetween:
		return re.isValueInRange(rule, attributeNameOfID, valueStore)
	case operatorNotBetween:
		valueInRange, err := re.isValueInRange(rule, attributeNameOfID, valueStore)
		return !valueInRange, err

	case operatorIn:
		isPresent := re.isValuePresentInCSV(rule, attributeNameOfID, valueStore)
		return isPresent, nil
	case operatorNotIn:
		isPresent := re.isValuePresentInCSV(rule, attributeNameOfID, valueStore)
		return !isPresent, nil

	}

	return false, fmt.Errorf("integer: invalid operator: %s", operator)
}

func (re IntegerRuleEvaluator) getValuesFromCSString(csv string) []int {
	retArr := make([]int, 0)
	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.Atoi(part)
		if err != nil {
			logger.Error(LoggerTag, "error converting %s to integer: %v", stringSet[0], err)
			continue
		}
		retArr = append(retArr, number)
	}

	return retArr
}

func (re IntegerRuleEvaluator) isValueInRange(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (bool, error) {
	operator := string(*r.Operator)
	valueFromStore, err := re.valueFromStore(r, attributeNameOfID, valueStore)
	if err != nil {
		return false, err
	}

	numbers := re.getValuesFromCSString(string(*r.Value))
	if len(numbers) != 2 {
		return false, fmt.Errorf("invalid number of values for operator %s: %s", operator, string(*r.Value))
	}
	return valueFromStore >= numbers[0] && valueFromStore <= numbers[1], nil
}

func (re IntegerRuleEvaluator) isValuePresentInCSV(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) bool {

	csv := string(*r.Value)
	value, err := re.valueFromStore(r, attributeNameOfID, valueStore)
	if err != nil {
		logger.Error(LoggerTag, "%v", err)
		return false
	}

	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.Atoi(part)
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

func (re IntegerRuleEvaluator) valueFromRuleAndStore(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (int, int, error) {
	valueFromRule, err := strconv.Atoi(string(*r.Value))
	if err != nil {
		return 0, 0, fmt.Errorf("error converting rule value %s to integer: %v", string(*r.Value), err)
	}

	valueFromStore, err := re.valueFromStore(r, attributeNameOfID, valueStore)
	if err != nil {
		return 0, 0, err
	}
	return valueFromRule, valueFromStore, nil
}

func (re IntegerRuleEvaluator) valueFromStore(r model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (int, error) {

	valueInterface, ok := GetValueFromStore(attributeNameOfID, valueStore)
	if !ok || valueInterface == nil {
		return 0, fmt.Errorf("value not found for id %s", attributeNameOfID)
	}

	valueFromStore, err1 := strconv.Atoi(fmt.Sprintf("%v", valueInterface))
	if err1 != nil {
		return 0, fmt.Errorf("error converting valueStore value %s to integer: %v", string(*r.Value), err1)
	}
	return valueFromStore, nil
}
