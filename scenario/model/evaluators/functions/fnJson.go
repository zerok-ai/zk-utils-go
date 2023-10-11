package functions

import (
	"encoding/json"
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
	path := fn.args[0]

	// resolve the path from attribute store
	resolvedVal := fn.attrStore.GetAttributeFromStore(*fn.attrStoreKey, path)
	if resolvedVal != nil {
		path = *resolvedVal
	}

	valueAtObject, err = jmespath.Search(path, jsonObject)

	if err != nil {
		zkLogger.ErrorF(LoggerTag, "Error evaluating jmespath at path:%s for store %v", path, jsonObject)
		return "", false
	}
	return valueAtObject, true
}
