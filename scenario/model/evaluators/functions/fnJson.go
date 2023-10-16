package functions

import (
	"encoding/json"
	"fmt"
	"github.com/jmespath/go-jmespath"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
)

const (
	JsonExtract = "jsonExtract"
	LoggerTag   = "functions"
)

type ExtractJson struct {
	name         string
	args         []string
	attrStore    *stores.ExecutorAttrStore
	attrStoreKey *cache.AttribStoreKey
	ff           *FunctionFactory
}

func (fn ExtractJson) Execute(valueAtObject interface{}) (interface{}, bool) {

	path := fn.args[0]
	var newValueAtObject interface{}
	var ok bool
	path, newValueAtObject, ok = fn.transformAttribute(path, valueAtObject)
	if ok {
		path = fmt.Sprintf("%v", newValueAtObject)
	}

	valueAtObject, ok = fn.executeJson(path, valueAtObject)
	return valueAtObject, ok
}

func (fn ExtractJson) executeJson(path string, valueAtObject interface{}) (interface{}, bool) {

	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In ExtractJson.Execute: Recovered from panic: %v", r)
		}
	}()

	// check if valueAtObject is a string
	var err error
	var stringVal string

	// if valueAtObject is a string, convert it to json, else directly read the json
	var jsonObject interface{}
	var ok bool
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

	valueAtObject, err = jmespath.Search(path, jsonObject)

	if err != nil {
		zkLogger.ErrorF(LoggerTag, "Error evaluating jmespath at path:%s for store %v", path, jsonObject)
		return "", false
	}
	return valueAtObject, true
}

func (fn ExtractJson) GetName() string {
	return fn.name
}

func (fn ExtractJson) transformAttribute(path string, valueAtObject interface{}) (string, interface{}, bool) {

	// resolve the path from attribute store
	resolvedVal, ok := fn.attrStore.GetAttributeFromStore(*fn.attrStoreKey, path)
	if ok {
		path = resolvedVal
	}
	valueAtObject, ok = getValueFromStoreInternal(path, valueAtObject.(map[string]interface{}), fn.ff, fn.attrStoreKey, true)
	return path, valueAtObject, ok
}
