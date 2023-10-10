package functions

import (
	"github.com/jmespath/go-jmespath"
	"github.com/zerok-ai/zk-utils-go/podDetails"
	zkRedis "github.com/zerok-ai/zk-utils-go/storage/redis"
)

const (
	getWorkloadFromIP = "getWorkloadFromIP"
)

type ExtractWorkLoadFromIP struct {
	Name            string
	Args            []string
	podDetailsStore *zkRedis.LocalCacheHSetStore
}

func (fn ExtractWorkLoadFromIP) Execute(valueAtObject interface{}) (interface{}, bool) {

	if len(fn.Args) < 1 || fn.podDetailsStore == nil {
		return "", false
	}

	// get the path and ip
	path := fn.Args[0]
	ip, err := jmespath.Search(path, valueAtObject)
	if err != nil || ip == nil || ip.(string) == "" {
		return "", false
	}

	// get the workload for the ip
	serviceName := podDetails.GetServiceNameFromPodDetails(ip.(string), fn.podDetailsStore)
	return serviceName, true
}
