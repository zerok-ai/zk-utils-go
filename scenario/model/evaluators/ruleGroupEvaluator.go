package evaluators

import (
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
)

type RuleGroupEvaluator struct {
	baseRuleEvaluator *RuleEvaluator
}

func (re *RuleGroupEvaluator) init() GroupRuleEvaluator {
	return re
}

func (re *RuleGroupEvaluator) evalRule(rule model.Rule, attrStoreKey cache.AttribStoreKey, valueStore map[string]interface{}) (bool, error) {

	// evaluate all the rules
	condition := *rule.Condition
	result := true // default is true for both `AND` and `OR` to work
	if condition == model.AND {
		result = true
	} else {
		result = false
	}
	for _, childRule := range rule.Rules {
		ok, err := re.baseRuleEvaluator.evalRule(childRule, attrStoreKey, valueStore)
		if err != nil {
			return false, err
		}

		if condition == model.AND {
			result = result && ok
			if !result {
				break
			}
		}

		if condition == model.OR {
			result = result || ok
			if result {
				break
			}
		}
	}
	return result, nil
}
