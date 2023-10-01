package evaluators

import (
	"fmt"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

const (
	LoggerTag = "scenario-evaluator"

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

func (ds DataStore) String() string {
	str := "{"
	for k, v := range ds {
		str += fmt.Sprintf("%s:%s,", k, v)
	}
	str += "}"
	return str
}

type RuleEvaluatorInternal interface {
	init() RuleEvaluatorInternal
	EvalRule(rule model.Rule, valueStore map[string]interface{}) (bool, error)
}

type RuleEvaluator struct {
	dataSource     DataStore
	ruleEvaluators map[string]RuleEvaluatorInternal
}

func getIDFromIDStore(rule model.Rule, idStore DataStore) string {

	// get the id
	id := *rule.RuleLeaf.ID

	// get the actual id from the idStore. If not found, use the id as is
	value, ok := idStore[id]
	if !ok {
		value = id
	}

	jsonPath := rule.RuleLeaf.JsonPath
	if jsonPath != nil {
		value = value + *jsonPath
	}

	return value
}

func NewRuleEvaluator() RuleEvaluator {
	return RuleEvaluator{
		ruleEvaluators: make(map[string]RuleEvaluatorInternal),
	}.init()
}

func (re RuleEvaluator) init() RuleEvaluator {
	re.ruleEvaluators[model.RULE_GROUP] = RuleGroupEvaluator{re}.init()

	re.ruleEvaluators[typeString] = StringRuleEvaluator{}.init()
	re.ruleEvaluators[typeInteger] = IntegerRuleEvaluator{}.init()
	re.ruleEvaluators[typeFloat] = FloatRuleEvaluator{}.init()
	re.ruleEvaluators[typeBool] = BooleanEvaluator{}.init()

	return re
}

func (re RuleEvaluator) EvalRule(rule model.Rule, idStore DataStore, valueStore map[string]interface{}) (bool, error) {

	handled, value := false, false
	var err error
	var evaluator string
	if rule.Type != model.RULE_GROUP {
		err = re.validate(rule)
		if err != nil {
			return false, err
		}

		handled, value, err = re.handleCommonOperators(rule, valueStore)
		evaluator = string(*rule.RuleLeaf.Datatype)
	} else {
		newID := getIDFromIDStore(rule, idStore)
		rule.RuleLeaf.ID = &newID
		evaluator = model.RULE_GROUP
	}
	if !handled {
		ruleEvaluator := re.ruleEvaluators[evaluator]
		if ruleEvaluator == nil {
			return false, fmt.Errorf("ruleEvaluator not found for type: %s", rule.Type)
		}
		return ruleEvaluator.EvalRule(rule, valueStore)
	}

	return value, err
}

// handleCommonOperators is a helper function to handle common operators like exists. The function returns
// a bool indicating if the rule is handled, a bool indicating the value, if handled and an error if any.
func (re RuleEvaluator) handleCommonOperators(r model.Rule, store map[string]interface{}) (bool, bool, error) {
	operator := string(*r.Operator)
	//	switch on operator
	switch operator {
	case operatorExists:
		_, ok := store[*r.RuleLeaf.ID]
		return true, ok, nil
	case operatorNotExists:
		_, ok := store[*r.RuleLeaf.ID]
		return true, !ok, nil
	}

	return false, false, nil
}

func (re RuleEvaluator) validate(r model.Rule) error {
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
