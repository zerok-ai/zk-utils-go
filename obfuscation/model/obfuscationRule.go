package model

import "time"

type Rule struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Analyzer   Analyzer   `json:"analyzer"`
	Anonymizer Anonymizer `json:"anonymizer"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreatedBy  string     `json:"created_by"`
	Enabled    bool       `json:"enabled"`
}

type Analyzer struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

type Anonymizer struct {
	Operator string           `json:"operator"`
	Params   AnonymizerParams `json:"params"`
}

type AnonymizerParams struct {
	NewValue string `json:"new_value"`
}
