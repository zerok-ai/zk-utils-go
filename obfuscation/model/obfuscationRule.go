package model

import "github.com/zerok-ai/zk-utils-go/interfaces"

type Rule struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Analyzer   Analyzer   `json:"analyzer"`
	Anonymizer Anonymizer `json:"anonymizer"`
	CreatedAt  int64      `json:"created_at"`
	UpdatedAt  int64      `json:"updated_at"`
	Enabled    bool       `json:"enabled"`
}

type RuleOperator struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Analyzer   Analyzer   `json:"analyzer"`
	Anonymizer Anonymizer `json:"anonymizer"`
	UpdatedAt  int64      `json:"updated_at"`
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

func (i RuleOperator) Equals(otherInterface interfaces.ZKComparable) bool {
	other, ok := otherInterface.(RuleOperator)
	if !ok {
		return false
	}

	// Compare Id and Name fields
	if i.Id != other.Id || i.Name != other.Name {
		return false
	}

	// Compare UpdatedAt field
	if i.UpdatedAt != other.UpdatedAt {
		return false
	}

	// Compare Analyzer fields
	if i.Analyzer.Type != other.Analyzer.Type || i.Analyzer.Pattern != other.Analyzer.Pattern {
		return false
	}

	// Compare Anonymizer fields
	if i.Anonymizer.Operator != other.Anonymizer.Operator {
		return false
	}

	// Compare AnonymizerParams fields
	if i.Anonymizer.Params.NewValue != other.Anonymizer.Params.NewValue {
		return false
	}

	return true
}
