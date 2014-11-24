package services

import (
	"encoding/json"
)

func DecodeFromJSON(content string, value interface{}) error {
	return json.Unmarshal([]byte(content), value)
}

func EncodeJSON(content interface{}) (string, error) {
	result, err := json.Marshal(content)
	return string(result), err
}
