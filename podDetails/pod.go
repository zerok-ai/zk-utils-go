package podDetails

import (
	"encoding/json"
	"fmt"
	"github.com/jmespath/go-jmespath"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
)

const LoggerTag = "ip"

type Set map[string]bool

func (s Set) Add(item string) {
	s[item] = true
}

func (s Set) Contains(item string) bool {
	return s[item]
}

type ProcessDetails struct {
	ProcessID   int                 `json:"pid"`
	ExeName     string              `json:"exe"`
	CmdLine     string              `json:"cmd"`
	Runtime     ProgrammingLanguage `json:"runtime"`
	ProcessName string              `json:"pname"`
	EnvMap      map[string]string   `json:"env"`
}

type ContainerDetails struct {
	Name                  string          `json:"container_name"`
	Image                 string          `json:"container_image"`
	ProcessExecutablePath []string        `json:"process.executable_path"`
	ProcessCommandArgs    []string        `json:"process.command_args"`
	Ports                 []ContainerPort `json:"ports"`
}

type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
}

type PodDetails struct {
	Metadata  PodMetadata      `json:"metadata"`
	Spec      PodSpec          `json:"spec"`
	Status    PodStatus        `json:"status"`
	Telemetry TelemetryDetails `json:"telemetry"`
}

type TelemetryDetails struct {
	TelemetryAutoVersion string `json:"telemetry_auto_version"`
	TelemetrySdkLanguage string `json:"telemetry_sdk_language"`
	TelemetrySdkName     string `json:"telemetry_sdk_name"`
	TelemetrySdkVersion  string `json:"telemetry_sdk_version"`
	ServiceName          string `json:"service_name"`
	ServiceVersion       string `json:"service_version"`
}

type PodMetadata struct {
	Namespace    string `json:"namespace"`
	PodName      string `json:"pod_name"`
	PodId        string `json:"pod_id"`
	WorkloadName string `json:"workload_name"`
	WorkloadKind string `json:"workload_kind"`
	ServiceName  string `json:"service_name"`
	CreateTS     string `json:"create_ts"`
}

type PodSpec struct {
	ServiceAccountName string             `json:"service_account_name"`
	NodeName           string             `json:"node_name"`
	Containers         []ContainerDetails `json:"containers"`
}

type PodStatus struct {
	Phase string `json:"phase"`
	PodIP string `json:"pod_ip"`
}
type ProgrammingLanguage string

type ContainerRuntime struct {
	Image    string            `json:"image"`
	ImageID  string            `json:"imageId"`
	Language []string          `json:"language"`
	Process  string            `json:"process,omitempty"`
	Cmd      []string          `json:"cmd,omitempty"`
	EnvMap   map[string]string `json:"env"`
}

func (cr ContainerRuntime) Equals(newContainerRuntime ContainerRuntime) bool {

	if cr.Image != newContainerRuntime.Image {
		return false
	}

	if cr.ImageID != newContainerRuntime.ImageID {
		return false
	}

	if len(cr.Language) != len(newContainerRuntime.Language) {
		return false
	}

	// collect all the elements for `cr` in a set and the languages may not be in order
	langSet := make(Set)
	for _, lang := range cr.Language {
		langSet.Add(lang)
	}

	// check if all the elements of the new array are present in the old array
	for index, _ := range cr.Language {
		if !langSet.Contains(newContainerRuntime.Language[index]) {
			return false
		}
	}

	return true
}

func (cr ContainerRuntime) String() string {

	stCr := fmt.Sprintf("%s:[", cr.Image)
	for _, lang := range cr.Language {
		stCr += lang + ", "
	}
	stCr += "]"

	return stCr
}

type RuntimeSyncRequest struct {
	RuntimeDetails []ContainerRuntime `json:"details"`
}

var serviceNamePaths = []string{"Metadata.ServiceName", "Telemetry.ServiceName"}

func GetServiceNameFromPodDetailsStore(ip string, podDetailsStore *stores.LocalCacheHSetStore) string {
	workloadDetailsPtr, _ := (*podDetailsStore).Get(ip)
	podDetails := loadPodDetailsIntoHashmap(ip, workloadDetailsPtr)

	var serviceName string
	for _, serviceNamePath := range serviceNamePaths {
		valAtPath, err := jmespath.Search(serviceNamePath, podDetails)
		if err == nil || valAtPath != nil {
			serviceName = valAtPath.(string)
			break
		}
	}
	return serviceName
}

const (
	status    = "status"
	metadata  = "metadata"
	spec      = "spec"
	telemetry = "telemetry"
)

func loadPodDetailsIntoHashmap(ip string, input *map[string]string) *PodDetails {
	podDetails := PodDetails{}
	if input == nil || len(*input) == 0 {
		zkLogger.ErrorF(LoggerTag, "Error getting service for ip = %s \n", ip)
		return &podDetails
	}

	//load status
	var podStatus PodStatus
	stringValue, ok := (*input)[status]
	if ok {
		err := json.Unmarshal([]byte(stringValue), &podStatus)
		if err != nil {
			zkLogger.ErrorF(LoggerTag, "Error marshalling status for ip = %s\n", ip)
		}
	}

	//load metadata
	var podMetadata PodMetadata
	stringValue, ok = (*input)[metadata]
	if ok {
		err := json.Unmarshal([]byte(stringValue), &podMetadata)
		if err != nil {
			zkLogger.ErrorF(LoggerTag, "Error marshalling metadata for ip = %s\n", ip)
		}
	}

	//load spec
	var podSpec PodSpec
	stringValue, ok = (*input)[spec]
	if ok {
		err := json.Unmarshal([]byte(stringValue), &podSpec)
		if err != nil {
			zkLogger.ErrorF(LoggerTag, "Error marshalling spec for ip = %s err:%v\n", ip, err)
		}
	}

	//load telemetry details
	var telemetryDetails TelemetryDetails
	stringValue, ok = (*input)[telemetry]
	if ok {
		err := json.Unmarshal([]byte(stringValue), &telemetryDetails)
		if err != nil {
			zkLogger.ErrorF(LoggerTag, "Error marshalling telemetry for ip = %s err:%v\n", ip, err)
		}
	}

	podDetails.Spec = podSpec
	podDetails.Status = podStatus
	podDetails.Metadata = podMetadata
	podDetails.Telemetry = telemetryDetails

	return &podDetails
}
