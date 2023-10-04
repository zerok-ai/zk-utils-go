package evaluators

import (
	"context"
	"fmt"
	"github.com/zerok-ai/zk-utils-go/ds"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/storage/redis/config"
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
	evalRule(rule model.Rule, valueStore map[string]interface{}) (bool, error)
}

type GroupRuleEvaluator interface {
	init() GroupRuleEvaluator
	evalRule(rule model.Rule, attributeVersion, protocol string, valueStore map[string]interface{}) (bool, error)
}

type RuleEvaluator struct {
	executor           model.Executor
	dataSource         DataStore
	attributeNameStore *cache.AttributeCache
	leafRuleEvaluators map[string]LeafRuleEvaluator
	groupRuleEvaluator GroupRuleEvaluator
}

func NewRuleEvaluator(redisConfig config.RedisConfig, executorName model.Executor, ctx context.Context) RuleEvaluator {
	return RuleEvaluator{
		executor:           executorName,
		leafRuleEvaluators: make(map[string]LeafRuleEvaluator),
		attributeNameStore: getAttributeNamesStore(redisConfig, ctx),
	}.init()
}

func getAttributeNamesStore(redisConfig config.RedisConfig, ctx context.Context) *cache.AttributeCache {

	dbName := "attrNames"
	noExpiryCache := ds.GetCacheWithExpiry[map[string]string](ds.NoExpiry)
	redisClient := config.GetRedisConnection(dbName, redisConfig)

	localCache := cache.GetAttributeCache(redisClient, noExpiryCache, nil, ctx)
	return localCache
}

func (re RuleEvaluator) init() RuleEvaluator {
	re.groupRuleEvaluator = RuleGroupEvaluator{re}.init()

	re.leafRuleEvaluators[typeString] = StringRuleEvaluator{}.init()
	re.leafRuleEvaluators[typeInteger] = IntegerRuleEvaluator{}.init()
	re.leafRuleEvaluators[typeFloat] = FloatRuleEvaluator{}.init()
	re.leafRuleEvaluators[typeBool] = BooleanEvaluator{}.init()

	return re
}

func (re RuleEvaluator) EvalRule(rule model.Rule, attributeVersion string, protocol string, valueStore map[string]interface{}) (bool, error) {
	return re.evalRule(rule, attributeVersion, protocol, valueStore)
}

func (re RuleEvaluator) evalRule(rule model.Rule, attributeVersion string, protocol string, valueStore map[string]interface{}) (bool, error) {

	handled, value := false, false
	var err error
	if rule.Type == model.RULE_GROUP {
		return re.groupRuleEvaluator.evalRule(rule, attributeVersion, protocol, valueStore)
	} else {
		err = re.validate(rule)
		if err != nil {
			return false, err
		}

		// replace id with actual attribute executor
		rule.RuleLeaf.ID = re.getAttributeName(rule, attributeVersion, protocol)

		handled, value, err = re.handleCommonOperators(rule, valueStore)
		evaluator := string(*rule.RuleLeaf.Datatype)

		if !handled {
			ruleEvaluator := re.leafRuleEvaluators[evaluator]
			if ruleEvaluator == nil {
				return false, fmt.Errorf("LeafRuleEvaluator not found for type: %s", rule.Type)
			}
			return ruleEvaluator.evalRule(rule, valueStore)
		}
	}

	return value, err
}

func (re RuleEvaluator) getAttributeName(rule model.Rule, attributeVersion, protocol string) *string {

	// get the id
	id := *rule.RuleLeaf.ID
	attributeName := id

	// get the actual id from the idStore. If not found, use the id as is
	attributeNameFromStore := re.attributeNameStore.Get(string(re.executor), attributeVersion, protocol, attributeName)
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
