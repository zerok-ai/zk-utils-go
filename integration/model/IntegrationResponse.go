package model

import (
	"encoding/json"
	"time"
)

type Level string
type Type string

const (
	Org        Level = "ORG"
	Cluster    Level = "CLUSTER"
	Prometheus Type  = "PROMETHEUS"
)

type IntegrationResponseObj struct {
	ID             string          `json:"id"`
	ClusterId      string          `json:"cluster_id,omitempty"`
	Alias          string          `json:"alias"`
	Type           Type            `json:"type"`
	URL            string          `json:"url"`
	Authentication json.RawMessage `json:"authentication"`
	Level          Level           `json:"level"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Deleted        bool            `json:"deleted"`
	Disabled       bool            `json:"disabled"`
	MetricServer   bool            `json:"metric_server"`
}
