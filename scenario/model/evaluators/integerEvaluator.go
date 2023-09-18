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
	operator := string(*r.Operator)
	_, ok := store[*r.ID]
	if !ok {
		return false, fmt.Errorf("value for id: %s not found in store", *r.ID)
	}

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

func (re IntegerRuleEvaluator) getValuesFromCSString(csv string) []int {
	retArr := make([]int, 0)
	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.Atoi(part)
		if err != nil {
			logger.Error(loggerTag, "error converting %s to integer: %v", stringSet[0], err)
			continue
		}
		retArr = append(retArr, number)
	}

	return retArr
}

func (re IntegerRuleEvaluator) isValueInRange(r model.Rule, store DataStore) (bool, error) {
	operator := string(*r.Operator)
	valueFromStore, err1 := strconv.Atoi(store[*r.ID])
	if err1 != nil {
		return false, fmt.Errorf("error converting store value %s to integer: %v", string(*r.Value), err1)
	}
	numbers := re.getValuesFromCSString(string(*r.Value))
	if len(numbers) != 2 {
		return false, fmt.Errorf("invalid number of values for operator %s: %s", operator, string(*r.Value))
	}
	return valueFromStore >= numbers[0] && valueFromStore <= numbers[1], nil
}

func (re IntegerRuleEvaluator) isValuePresentInCSV(r model.Rule, store DataStore) bool {

	csv := string(*r.Value)
	value, err1 := strconv.Atoi(store[*r.ID])
	if err1 != nil {
		logger.Error(loggerTag, "error converting store value %s to integer: %v", string(*r.Value), err1)
		return false
	}

	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.Atoi(part)
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

func (re IntegerRuleEvaluator) valueFromRuleAndStore(r model.Rule, store DataStore) (int, int, error) {
	valueFromRule, err := strconv.Atoi(string(*r.Value))
	if err != nil {
		return 0, 0, fmt.Errorf("error converting rule value %s to integer: %v", string(*r.Value), err)
	}
	valueFromStore, err1 := strconv.Atoi(store[*r.ID])
	if err1 != nil {
		return 0, 0, fmt.Errorf("error converting store value %s to integer: %v", string(*r.Value), err1)
	}
	return valueFromRule, valueFromStore, nil
}
