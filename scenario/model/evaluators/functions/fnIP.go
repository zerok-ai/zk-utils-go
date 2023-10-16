package functions

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/podDetails"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
)

const (
	getWorkloadFromIP = "getWorkloadFromIP"
)

type ExtractWorkLoadFromIP struct {
	name            string
	args            []string
	podDetailsStore *stores.LocalCacheHSetStore
	attrStore       *stores.ExecutorAttrStore
	attrStoreKey    *cache.AttribStoreKey
	ff              *FunctionFactory
}

func (fn ExtractWorkLoadFromIP) Execute(valueAtObject interface{}) (interface{}, bool) {

	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In Execute of ExtractWorkLoadFromIP: Recovered from panic: %v", r)
		}
	}()

	if len(fn.args) < 1 {
		return "", false
	}

	// get the path and ip
	path := fn.args[0]
	newValueAtObject, ok := fn.transformAttribute(path, valueAtObject)
	if ok {
		path = fmt.Sprintf("%v", newValueAtObject)
	}

	// get the workload for the ip
	serviceName := podDetails.GetServiceNameFromPodDetailsStore(path, fn.podDetailsStore)
	return serviceName, true
}

func (fn ExtractWorkLoadFromIP) GetName() string {
	return fn.name
}

func (fn ExtractWorkLoadFromIP) transformAttribute(path string, valueAtObject interface{}) (interface{}, bool) {

	// resolve the path from attribute store
	resolvedVal, ok := fn.attrStore.GetAttributeFromStore(*fn.attrStoreKey, path)
	if ok {
		path = resolvedVal
	}
	return getValueFromStoreInternal(path, valueAtObject.(map[string]interface{}), fn.ff, fn.attrStoreKey, true)
}
