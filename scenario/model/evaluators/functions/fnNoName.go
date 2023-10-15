package functions

import (
	"github.com/jmespath/go-jmespath"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
)

const (
	NoName = "noName"
)

type NoNameFunction struct {
	name         string
	args         []string
	attrStore    *stores.ExecutorAttrStore
	attrStoreKey *cache.AttribStoreKey
	ff           *FunctionFactory
}

func (fn NoNameFunction) Execute(valueAtObject interface{}) (interface{}, bool) {

	// try to create functions for the args
	path := fn.args[0]

	newValueAtObject, ok := fn.transformAttribute(path, valueAtObject)
	if ok {
		return newValueAtObject, true
	} else {
		returnVal, err := jmespath.Search(path, valueAtObject)
		if err != nil {
			zkLogger.ErrorF(LoggerTag, "Error evaluating jmespath at path:%s for store %v", path, valueAtObject)
			return "", false
		}
		return returnVal, true
	}
}

func (fn NoNameFunction) GetName() string {
	return fn.name
}

func (fn NoNameFunction) transformAttribute(path string, valueAtObject interface{}) (interface{}, bool) {

	// resolve the path from attribute store
	resolvedVal, ok := fn.attrStore.GetAttributeFromStore(*fn.attrStoreKey, path)
	if ok {
		path = resolvedVal
	}
	return getValueFromStoreInternal(resolvedVal, valueAtObject.(map[string]interface{}), fn.ff, fn.attrStoreKey, false)

}
