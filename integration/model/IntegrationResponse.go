package model

import (
	"encoding/json"
	"github.com/zerok-ai/zk-utils-go/interfaces"
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
	Authentication json.RawMessage `json:"authentication,omitempty"`
	Level          Level           `json:"level"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Deleted        bool            `json:"deleted"`
	Disabled       bool            `json:"disabled"`
	MetricServer   *bool           `json:"metric_server"`
}

func (i IntegrationResponseObj) Equals(otherInterface interfaces.ZKComparable) bool {
	other, ok := otherInterface.(IntegrationResponseObj)
	if !ok {
		return false
	}
	return i.ID == other.ID
}
