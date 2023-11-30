package model

type ClusterStatusRequest struct {
	NumberOfNodes int                      `json:"number_of_nodes"`
	Services      map[string]ServiceStatus `json:"services"`
}
