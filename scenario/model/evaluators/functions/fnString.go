package functions

import "strings"

const (
	toLowerCase = "toLowerCase"
	toUpperCase = "toUpperCase"
)

type UpperCase struct {
	Name string
	Args []string
}

func (fn UpperCase) Execute(valueAtObject interface{}) (interface{}, bool) {
	stringVal, ok := valueAtObject.(string)
	if ok {
		return strings.ToUpper(stringVal), true
	}
	return "", false
}

type LowerCase struct {
	Name string
	Args []string
}

func (fn LowerCase) Execute(valueAtObject interface{}) (interface{}, bool) {
	stringVal, ok := valueAtObject.(string)
	if ok {
		return strings.ToLower(stringVal), true
	}
	return "", false
}
