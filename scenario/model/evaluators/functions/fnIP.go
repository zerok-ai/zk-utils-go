package functions

import (
	"github.com/jmespath/go-jmespath"
	zkRedis "github.com/zerok-ai/zk-utils-go/storage/redis"
)

const (
	getWorkloadFromIP = "getWorkloadFromIP"
)

type ExtractWorkLoadFromIP struct {
	Name           string
	Args           []string
	serviceIPStore *zkRedis.LocalCacheHSetStore
}

func (fn ExtractWorkLoadFromIP) Execute(valueAtObject interface{}) (interface{}, bool) {

	if len(fn.Args) < 1 || fn.serviceIPStore == nil {
		return "", false
	}

	// get the path and ip
	path := fn.Args[0]
	ip, err := jmespath.Search(path, valueAtObject)
	if err != nil || ip == nil || ip.(string) == "" {
		return "", false
	}

	// get the workload for the ip
	workloadDetailsPtr, _ := (*fn.serviceIPStore).Get(ip.(string))
	//workloadDetailsPtr, _ := scenarioManager.serviceIPStore.Get("10.60.1.53")
	podDetails := LoadIPDetailsIntoHashmap(ip.(string), workloadDetailsPtr)

	serviceName, err := jmespath.Search("Metadata.ServiceName", podDetails)
	if err != nil || serviceName == nil {
		return ip.(string), false
	}
	return serviceName.(string), false
}
