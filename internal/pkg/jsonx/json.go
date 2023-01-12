package jsonx

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

func Marshal(v any) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UserRegisteredEvent: %w", err)
	}
	return jsonBytes, nil
}
