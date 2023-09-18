package evaluators

import "github.com/zerok-ai/zk-utils-go/scenario/model"

type RuleGroupEvaluator struct {
	baseRuleEvaluator BaseRuleEvaluator
}

func (re RuleGroupEvaluator) init() RuleEvaluator {
	return re
}

func (re RuleGroupEvaluator) EvalRule(r model.Rule, store DataStore) (bool, error) {

	// evaluate all the rules
	condition := *r.Condition
	result := true // default is true for both `AND` and `OR` to work
	if condition == model.AND {
		result = true
	} else {
		result = false
	}
	for _, rule := range r.Rules {
		ok, err := re.baseRuleEvaluator.EvalRule(rule, store)
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
