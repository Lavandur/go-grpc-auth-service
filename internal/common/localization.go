package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// LocalizedString - type for localization by json ("en": "some")
type LocalizedString map[string]string

func (a *LocalizedString) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *LocalizedString) Scan(value interface{}) error {
	if value == nil {
		a = nil
		return nil
	}

	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	return json.Unmarshal([]byte(b), &a)
}
