package server

import (
	"encoding/json"
	"fmt"
)

func toJson(val any) ([]byte, error) {
	content, err := json.Marshal(val)
	if err != nil {
		return []byte{}, fmt.Errorf("to-json error: %w", err)
	}

	return content, err
}

func fromJsonTo(content []byte, val any) error {
	err := json.Unmarshal(content, &val)
	if err != nil {
		return fmt.Errorf(" from-json error: %w", err)
	}

	return nil
}
