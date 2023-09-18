package evaluators

import (
	"encoding/json"
	"fmt"
	"github.com/jmespath/go-jmespath"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

const (
	loggerTag = "scenario-evaluator"

	typeString             = "string"
	typeKeyValue           = "key-value"
	typeInteger            = "integer"
	typeFloat              = "float"
	typeBool               = "bool"
	typeWorkLoadIdentifier = "workload-identifier"

	operatorExists    = "exists"
	operatorNotExists = "not_exists"

	operatorMatches          = "matches"
	operatorDoesNotMatch     = "does_not_match"
	operatorEqual            = "equal"
	operatorNotEqual         = "not_equal"
	operatorContains         = "contains"
	operatorDoesNotContain   = "does_not_contain"
	operatorIn               = "in"
	operatorNotIn            = "not_in"
	operatorBeginsWith       = "begins_with"
	operatorDoesNotBeginWith = "does_not_begin_with"
	operatorEndsWith         = "ends_with"
	operatorDoesNotEndWith   = "does_not_end_with"

	operatorLessThan         = "less_than"
	operatorLessThanEqual    = "less_than_equal"
	operatorGreaterThan      = "greater_than"
	operatorGreaterThanEqual = "greater_than_equal"
	operatorBetween          = "between"
	operatorNotBetween       = "not_between"
)

type DataStore map[string]string
type RuleEvaluator interface {
	init() RuleEvaluator
	EvalRule(rule model.Rule, store DataStore) (bool, error)
}

type BaseRuleEvaluator struct {
	dataSource     DataStore
	ruleEvaluators map[string]RuleEvaluator
}

func NewRuleEvaluator() RuleEvaluator {
	return BaseRuleEvaluator{
		ruleEvaluators: make(map[string]RuleEvaluator),
	}.init()
}

func (re BaseRuleEvaluator) init() RuleEvaluator {
	re.ruleEvaluators[model.RULE_GROUP] = RuleGroupEvaluator{re}.init()

	re.ruleEvaluators[typeString] = StringRuleEvaluator{}.init()
	re.ruleEvaluators[typeInteger] = IntegerRuleEvaluator{}.init()
	re.ruleEvaluators[typeFloat] = FloatRuleEvaluator{}.init()
	re.ruleEvaluators[typeBool] = BooleanEvaluator{}.init()

	return re
}

func (re BaseRuleEvaluator) EvalRule(r model.Rule, store DataStore) (bool, error) {

	handled, value := false, false
	var err error
	var evaluator string
	if r.Type != model.RULE_GROUP {
		err = re.validate(r, store)
		if err != nil {
			return false, err
		}

		handled, value, err = re.handleCommonOperators(r, store)
		evaluator = string(*r.RuleLeaf.Datatype)
	} else {
		evaluator = model.RULE_GROUP
	}
	if !handled {

		r, store, err = re.handlePath(r, store)
		if err != nil {
			return false, err
		}

		ruleEvaluator := re.ruleEvaluators[evaluator]
		if ruleEvaluator == nil {
			return false, fmt.Errorf("ruleEvaluator not found for type: %s", r.Type)
		}
		return ruleEvaluator.EvalRule(r, store)
	}

	return value, err
}

// handleCommonOperators is a helper function to handle common operators like exists. The function returns
// a bool indicating if the rule is handled, a bool indicating the value, if handled and an error if any.
func (re BaseRuleEvaluator) handleCommonOperators(r model.Rule, store DataStore) (bool, bool, error) {
	operator := string(*r.Operator)
	//	switch on operator
	switch operator {
	case operatorExists:
		_, ok := store[*r.ID]
		return true, ok, nil
	case operatorNotExists:
		_, ok := store[*r.ID]
		return true, !ok, nil
	}

	return false, false, nil
}

func (re BaseRuleEvaluator) validate(r model.Rule, store DataStore) error {
	id := r.ID
	operator := r.Operator
	valueFromRule := r.Value
	dataType := r.Datatype
	if id == nil {
		return fmt.Errorf("id is nil")
	} else if operator == nil {
		return fmt.Errorf("operator is nil")
	} else if valueFromRule == nil {
		return fmt.Errorf("value is nil")
	} else if dataType == nil {
		return fmt.Errorf("datatype is nil")
	}

	return nil
}

func (re BaseRuleEvaluator) handlePath(r model.Rule, store DataStore) (model.Rule, DataStore, error) {

	var jsonPath, arrayIndex *string = nil, nil
	if r.RuleLeaf != nil {
		jsonPath = r.RuleLeaf.JsonPath
		arrayIndex = r.RuleLeaf.ArrayIndex
	}
	if jsonPath != nil {
		valueFromStore, ok := store[*r.ID]
		if !ok {
			return r, store, fmt.Errorf("value for id: %s not found in store", *r.ID)
		}

		// Define a map to store the parsed JSON data
		var data map[string]interface{}

		// Unmarshal the JSON data into the map
		err := json.Unmarshal([]byte(valueFromStore), &data)
		if err != nil {
			return r, store, fmt.Errorf("error unmarshalling json_path for id: %s  %v", *r.ID, err)
		}

		//load json from valueFromStore
		valueAtPath, err := jmespath.Search(*jsonPath, valueFromStore)
		if err != nil {
			return r, store, fmt.Errorf("value for id: %s not found in store at path:%s ", *r.ID, *jsonPath)
		}

		newId := *r.ID + *jsonPath
		r.ID = &newId
		store[*r.ID] = valueAtPath.(string)

		return r, store, nil
	}

	//TODO handle array index
	if arrayIndex != nil {
		return r, store, nil
	}

	return r, store, nil
}
