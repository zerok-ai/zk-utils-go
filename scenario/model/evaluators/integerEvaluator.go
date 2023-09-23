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

func (re IntegerRuleEvaluator) init() RuleEvaluator {
	return re
}

func (re IntegerRuleEvaluator) EvalRule(r model.Rule, store DataStore) (bool, error) {

	// get the values assuming that the rule object is valid
	operator := string(*r.RuleLeaf.Operator)

	//	switch on operator
	switch operator {

	case operatorLessThan:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(r, store)
		if err != nil {
			return false, err
		}
		return valueFromStore < valueFromRule, nil
	case operatorLessThanEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(r, store)
		if err != nil {
			return false, err
		}
		return valueFromStore <= valueFromRule, nil
	case operatorGreaterThan:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(r, store)
		if err != nil {
			return false, err
		}
		return valueFromStore > valueFromRule, nil
	case operatorGreaterThanEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(r, store)
		if err != nil {
			return false, err
		}
		return valueFromStore >= valueFromRule, nil
	case operatorEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(r, store)
		if err != nil {
			return false, err
		}
		return valueFromStore == valueFromRule, nil
	case operatorNotEqual:
		valueFromRule, valueFromStore, err := re.valueFromRuleAndStore(r, store)
		if err != nil {
			return false, err
		}
		return valueFromStore != valueFromRule, nil

	case operatorBetween:
		return re.isValueInRange(r, store)
	case operatorNotBetween:
		valueInRange, err := re.isValueInRange(r, store)
		return !valueInRange, err

	case operatorIn:
		isPresent := re.isValuePresentInCSV(r, store)
		return isPresent, nil
	case operatorNotIn:
		isPresent := re.isValuePresentInCSV(r, store)
		return !isPresent, nil

	}

	return false, fmt.Errorf("integer: invalid operator: %s", operator)
}

func (re IntegerRuleEvaluator) getValuesFromCSString(csv string) []int64 {
	retArr := make([]int64, 0)
	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := getInteger64(part)
		if err != nil {
			logger.Error(loggerTag, "error converting %s to integer: %v", stringSet[0], err)
			continue
		}
		retArr = append(retArr, int64(number))
	}

	return retArr
}

func getInteger64(value interface{}) (int64, error) {
	// convert value to int64
	if value == nil {
		return 0, fmt.Errorf("invalid integer value: %s", value)
	}

	if intValue, ok := value.(int64); ok {
		return intValue, nil
	}

	strValue := fmt.Sprintf("%v", value)
	intValue, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value: %s", strValue)
	}
	return intValue, nil
}

func (re IntegerRuleEvaluator) isValueInRange(r model.Rule, store DataStore) (bool, error) {
	operator := string(*r.Operator)
	valueFromStore, err1 := getInteger64(store[*r.ID])
	if err1 != nil {
		return false, fmt.Errorf("range-int: error converting store value %s to integer: %v", string(*r.Value), err1)
	}
	numbers := re.getValuesFromCSString(string(*r.Value))
	if len(numbers) != 2 {
		return false, fmt.Errorf("invalid number of values for operator %s: %s", operator, string(*r.Value))
	}
	return valueFromStore >= numbers[0] && valueFromStore <= numbers[1], nil
}

func (re IntegerRuleEvaluator) isValuePresentInCSV(r model.Rule, store DataStore) bool {

	value, err1 := getInteger64(store[*r.ID])
	if err1 != nil {
		logger.Error(loggerTag, "csv: error converting store value %s to integer: %v", string(*r.Value), err1)
		return false
	}

	csv := string(*r.Value)
	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := getInteger64(part)
		if err != nil {
			logger.Error(loggerTag, "error converting %s to integer: %v", stringSet[0], err)
			continue
		}
		if number == value {
			return true
		}
	}
	return false
}

func (re IntegerRuleEvaluator) valueFromRuleAndStore(r model.Rule, store DataStore) (int64, int64, error) {
	valueFromRule, err := getInteger64(*r.RuleLeaf.Value)
	if err != nil {
		return 0, 0, fmt.Errorf("error converting rule value %s to integer: %v", string(*r.Value), err)
	}
	valueFromStore, err1 := getInteger64(getValueFromStore(r, store))
	if err1 != nil {
		return 0, 0, fmt.Errorf("value-int: error converting store value %s to integer: %v", string(*r.Value), err1)
	}
	return valueFromRule, valueFromStore, nil
}
