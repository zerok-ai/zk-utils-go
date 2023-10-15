package evaluators

import "github.com/zerok-ai/zk-utils-go/scenario/model"

type RuleGroupEvaluator struct {
	baseRuleEvaluator *RuleEvaluator
}

func (re *RuleGroupEvaluator) init() GroupRuleEvaluator {
	return re
}

func (re *RuleGroupEvaluator) evalRule(rule model.Rule, attributeVersion string, protocol model.ProtocolName, valueStore map[string]interface{}) (bool, error) {

	// evaluate all the rules
	condition := *rule.Condition
	result := true // default is true for both `AND` and `OR` to work
	if condition == model.AND {
		result = true
	} else {
		result = false
	}
	for _, childRule := range rule.Rules {
		ok, err := re.baseRuleEvaluator.evalRule(childRule, attributeVersion, protocol, valueStore)
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
