package functions

import (
	"encoding/json"
	"github.com/jmespath/go-jmespath"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
)

const (
	JsonExtract = "jsonExtract"
	LoggerTag   = "functions"
)

type ExtractJson struct {
	Name string
	Args []string
}

func (fn ExtractJson) Execute(valueAtObject interface{}) (value interface{}, ok bool) {

	// check if valueAtObject is a string
	var err error
	var stringVal string

	// if valueAtObject is a string, convert it to json, else directly read the json
	var jsonObject interface{}
	stringVal, ok = valueAtObject.(string)
	if ok {
		// convert string to json
		err = json.Unmarshal([]byte(stringVal), &jsonObject)
		if err != nil {
			zkLogger.Error(LoggerTag, "Error marshalling string:", stringVal)
			return "", false
		}
	} else {
		jsonObject = valueAtObject
	}
	path := fn.Args[0]
	valueAtObject, err = jmespath.Search(path, jsonObject)

	if err != nil {
		zkLogger.ErrorF(LoggerTag, "Error evaluating jmespath at path:%s for store %v", path, jsonObject)
		return "", false
	}
	return valueAtObject.(string), true
}
