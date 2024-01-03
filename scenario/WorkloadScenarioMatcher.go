package scenario

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
)

var LogTag = "scenario_match_handler"

func FindMatchingScenarios(workloadIds []string, scenarios map[string]*model.Scenario) ([]string, error) {
	var matchingScenarios []string
	var workloadIdMap = make(map[string]bool)
	for _, workloadId := range workloadIds {
		workloadIdMap[workloadId] = true
	}
	for _, scenario := range scenarios {
		val, err := evaluateFilter(scenario.Filter, workloadIdMap)
		if err != nil {
			zkLogger.Error(LogTag, "Error while evaluating the filter for: ", scenario.Title, err)
		} else if val {
			matchingScenarios = append(matchingScenarios, scenario.Title)
		}
	}
	return matchingScenarios, nil
}

func evaluateFilter(filter model.Filter, workloadIdMap map[string]bool) (bool, error) {
	defaultValue := getDefaultValue(filter.Condition)
	if filter.Type == model.WORKLOAD {
		for _, workloadId := range *filter.WorkloadIds {
			val := getValueForWorkload(workloadId, workloadIdMap)
			if filter.Condition == model.CONDITION_AND && !val {
				return false, nil
			} else if filter.Condition == model.CONDITION_OR && val {
				return true, nil
			}
		}
	} else if filter.Type == model.FILTER {
		for _, f := range *filter.Filters {
			val, err := evaluateFilter(f, workloadIdMap)
			if err != nil {
				return false, err
			}
			if filter.Condition == model.CONDITION_AND && !val {
				return false, nil
			} else if filter.Condition == model.CONDITION_OR && val {
				return true, nil
			}
		}
	}
	return defaultValue, nil
}

func getValueForWorkload(workloadId string, workloadIdMap map[string]bool) bool {
	val, ok := workloadIdMap[workloadId]
	if !ok {
		return false
	}
	return val
}

func getDefaultValue(condition model.Condition) bool {
	if condition == model.CONDITION_OR {
		return false
	}
	return true
}
