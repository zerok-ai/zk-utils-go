package functions

import (
	"github.com/jmespath/go-jmespath"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/podDetails"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
)

const (
	getWorkloadFromIP = "getWorkloadFromIP"
)

type ExtractWorkLoadFromIP struct {
	Name            string
	Args            []string
	podDetailsStore *stores.LocalCacheHSetStore
}

func (fn ExtractWorkLoadFromIP) Execute(valueAtObject interface{}) (interface{}, bool) {

	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In Execute of ExtractWorkLoadFromIP: Recovered from panic: %v", r)
		}
	}()

	if len(fn.Args) < 1 {
		return "", false
	}

	// get the path and ip
	path := fn.Args[0]
	ip, err := jmespath.Search(path, valueAtObject)
	if err != nil || ip == nil || ip.(string) == "" {
		return "", false
	}

	// get the workload for the ip
	serviceName := podDetails.GetServiceNameFromPodDetailsStore(ip.(string), fn.podDetailsStore)
	return serviceName, true
}

func (fn ExtractWorkLoadFromIP) GetName() string {
	return fn.Name
}
