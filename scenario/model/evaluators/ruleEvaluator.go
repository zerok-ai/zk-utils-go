package evaluators

import (
	"fmt"
	"github.com/jmespath/go-jmespath"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
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

type LeafRuleEvaluator interface {
	init() LeafRuleEvaluator
	evalRule(rule model.Rule, attributeNameOfID string, valueStore map[string]interface{}) (bool, error)
}

type GroupRuleEvaluator interface {
	init() GroupRuleEvaluator
	evalRule(rule model.Rule, attributeVersion string, protocol model.ProtocolName, valueStore map[string]interface{}) (bool, error)
}

type RuleEvaluator struct {
	executorName       model.ExecutorName
	dataSource         DataStore
	attributeNameStore *cache.AttributeCache
	leafRuleEvaluators map[string]LeafRuleEvaluator
	groupRuleEvaluator GroupRuleEvaluator
}

func NewRuleEvaluator(executorName model.ExecutorName, attributeNameStore *cache.AttributeCache) RuleEvaluator {
	return RuleEvaluator{
		executorName:       executorName,
		attributeNameStore: attributeNameStore,
		leafRuleEvaluators: make(map[string]LeafRuleEvaluator),
	}.init()
}

func (re RuleEvaluator) init() RuleEvaluator {
	re.groupRuleEvaluator = RuleGroupEvaluator{re}.init()

	re.leafRuleEvaluators[typeString] = StringRuleEvaluator{}.init()
	re.leafRuleEvaluators[typeInteger] = IntegerRuleEvaluator{}.init()
	re.leafRuleEvaluators[typeFloat] = FloatRuleEvaluator{}.init()
	re.leafRuleEvaluators[typeBool] = BooleanEvaluator{}.init()

	return re
}

func (re RuleEvaluator) EvalRule(rule model.Rule, attributeVersion string, protocol model.ProtocolName, valueStore map[string]interface{}) (bool, error) {

	result, err := re.evalRule(rule, attributeVersion, protocol, valueStore)
	return result, err
}

func (re RuleEvaluator) evalRule(rule model.Rule, attributeVersion string, protocol model.ProtocolName, valueStore map[string]interface{}) (bool, error) {

	handled, value := false, false
	var err error
	if rule.Type == model.RULE_GROUP {
		value, err = re.groupRuleEvaluator.evalRule(rule, attributeVersion, protocol, valueStore)
		zkLogger.DebugF(LoggerTag, "Evaluated value for group =%v, for condition=%s", value, *rule.RuleGroup.Condition)
	} else {
		err = re.validate(rule)
		if err != nil {
			return false, err
		}

		// replace id with actual attribute executorName
		attributeNameOfID := re.getAttributeName(rule, attributeVersion, protocol)

		zkLogger.DebugF(LoggerTag, "RuleId:- ruleID=%s, attributeName=%s", *rule.RuleLeaf.ID, *attributeNameOfID)

		handled, value, err = re.handleCommonOperators(rule, *attributeNameOfID, valueStore)
		leafEvaluatorType := string(*rule.RuleLeaf.Datatype)

		if !handled {
			ruleEvaluator := re.leafRuleEvaluators[leafEvaluatorType]
			if ruleEvaluator == nil {
				return false, fmt.Errorf("LeafRuleEvaluator not found for type: %s", leafEvaluatorType)
			}
			value, err = ruleEvaluator.evalRule(rule, *attributeNameOfID, valueStore)
		}
		zkLogger.DebugF(LoggerTag, "Evaluated value=%v, for attributeName=%v", value, *attributeNameOfID)
	}
	return value, err
}

func (re RuleEvaluator) getAttributeName(rule model.Rule, attributeVersion string, protocol model.ProtocolName) *string {

	attributeName := *rule.RuleLeaf.ID

	// get the actual id from the idStore. If not found, use the id as is
	attributeNameFromStore := re.attributeNameStore.Get(string(re.executorName), attributeVersion, protocol, *rule.RuleLeaf.ID)
	if attributeNameFromStore != nil {
		attributeName = *attributeNameFromStore
	}

	jsonPath := rule.RuleLeaf.JsonPath
	if jsonPath != nil {
		//add jsonPath to the attribute name using jsonExtract function
		jsonPathString := "#jsonExtract("
		for index, path := range *jsonPath {
			if index > 0 {
				jsonPathString += "."
			}
			jsonPathString += "\"" + path + "\""
		}
		jsonPathString += ")"

		attributeName += jsonPathString
	}

	return &attributeName
}

// handleCommonOperators is a helper function to handle common operators like exists. The function returns
// a bool indicating if the rule is handled, a bool indicating the value, if handled and an error if any.
func (re RuleEvaluator) handleCommonOperators(r model.Rule, attributeNameOfID string, store map[string]interface{}) (handled bool, returnValue bool, err error) {
	operator := string(*r.Operator)
	handled = false
	returnValue = false

	//	switch on operator
	switch operator {
	case operatorExists:
		handled = true
		value, err := jmespath.Search(attributeNameOfID, store)
		returnValue = err == nil && value != nil
	case operatorNotExists:
		handled = true
		value, err := jmespath.Search(attributeNameOfID, store)
		returnValue = err != nil || value == nil
	}

	return handled, returnValue, nil
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
