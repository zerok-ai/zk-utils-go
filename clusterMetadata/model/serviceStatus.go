package model

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Healthy   bool           `json:"healthy"`
	PodStatus map[string]int `json:"pod_status,omitempty"`
}
