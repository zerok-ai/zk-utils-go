package evaluators

import "github.com/zerok-ai/zk-utils-go/scenario/model"

type RuleGroupEvaluator struct {
	ruleEvaluator RuleEvaluator
}

func (re RuleGroupEvaluator) init() GroupRuleEvaluator {
	return re
}

func (re RuleGroupEvaluator) evalGroupRule(r model.Rule, idStore DataStore, valueStore map[string]interface{}) (bool, error) {

	// evaluate all the rules
	condition := *r.Condition
	result := true // default is true for both `AND` and `OR` to work
	if condition == model.AND {
		result = true
	} else {
		result = false
	}
	for _, rule := range r.Rules {
		ok, err := re.ruleEvaluator.EvalRule(rule, idStore, valueStore)
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
