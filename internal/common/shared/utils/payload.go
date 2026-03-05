package utils

import (
	"encoding/json"
)

func ParsePayload[T any](payload any) (*T, error) {
	// Cách dễ nhất và an toàn nhất: Marshal về byte, rồi Unmarshal thẳng vào Struct
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var data T
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
