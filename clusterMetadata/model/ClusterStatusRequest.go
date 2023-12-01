package model

type ClusterStatusRequest struct {
	NumberOfNodes int                      `json:"number_of_nodes"`
	Services      map[string]ServiceStatus `json:"services"`
}

type ClusterConnection struct {
	Connected bool `json:"connected"`
}

type ClusterConnectionStatus struct {
	Status map[string]ClusterConnection `json:"status"`
}
