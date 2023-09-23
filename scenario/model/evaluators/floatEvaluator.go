package evaluators

import (
	"fmt"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"strconv"
	"strings"
)

type FloatRuleEvaluator struct {
}

func (re FloatRuleEvaluator) init() RuleEvaluator {
	return re
}

func (re FloatRuleEvaluator) EvalRule(r model.Rule, store DataStore) (bool, error) {

	// get the values assuming that the rule object is valid
	operator := string(*r.Operator)

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

	return false, fmt.Errorf("float: invalid operator: %s", operator)
}

func (re FloatRuleEvaluator) getValuesFromCSString(csv string) []float64 {
	retArr := make([]float64, 0)
	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.ParseFloat(part, 64)
		if err != nil {
			logger.Error(loggerTag, "error converting %s to integer: %v", stringSet[0], err)
			continue
		}
		retArr = append(retArr, number)
	}

	return retArr
}

func getFloat64(value interface{}) (float64, error) {

	// convert value to float
	if value == nil {
		return 0, fmt.Errorf("invalid float value: %s", value)
	}

	if floatValue, ok := value.(float64); ok {
		return floatValue, nil
	}

	strValue := fmt.Sprintf("%v", value)
	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float value: %s", strValue)
	}

	return floatValue, nil
}

func (re FloatRuleEvaluator) isValueInRange(r model.Rule, store DataStore) (bool, error) {
	operator := string(*r.Operator)
	valueFromStore, err1 := getFloat64(getValueFromStore(r, store))

	if err1 != nil {
		return false, fmt.Errorf("range: error converting store value %s to integer: %v", string(*r.RuleLeaf.Value), err1)
	}
	numbers := re.getValuesFromCSString(string(*r.RuleLeaf.Value))
	if len(numbers) != 2 {
		return false, fmt.Errorf("invalid number of values for operator %s: %s", operator, string(*r.RuleLeaf.Value))
	}
	return valueFromStore >= numbers[0] && valueFromStore <= numbers[1], nil
}

func (re FloatRuleEvaluator) isValuePresentInCSV(r model.Rule, store DataStore) bool {

	csv := string(*r.RuleLeaf.Value)
	value, err1 := getFloat64(getValueFromStore(r, store))
	if err1 != nil {
		logger.Error(loggerTag, "csv-float: error converting store value %s to integer: %v", string(*r.RuleLeaf.Value), err1)
		return false
	}

	stringSet := strings.Split(csv, ",")

	for _, part := range stringSet {
		number, err := strconv.ParseFloat(part, 64)
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

func (re FloatRuleEvaluator) valueFromRuleAndStore(r model.Rule, store DataStore) (float64, float64, error) {
	valueFromRule, err := getFloat64(string(*r.RuleLeaf.Value))
	if err != nil {
		return 0, 0, fmt.Errorf("error converting rule value %s to float: %v", string(*r.RuleLeaf.Value), err)
	}
	valueFromStore, err1 := getFloat64(getValueFromStore(r, store))
	if err1 != nil {
		return 0, 0, fmt.Errorf("value-float: error converting store value %s to float: %v", string(*r.RuleLeaf.Value), err1)
	}
	return valueFromRule, valueFromStore, nil
}
