package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type GenericMap map[string]interface{}

func (genericMap GenericMap) Value() (driver.Value, error) {
	return json.Marshal(genericMap)
}

func (genericMap GenericMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &genericMap)
}