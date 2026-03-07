package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e ContextType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *ContextType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ContextType, got %v", t)
	}
	val, err := ContextTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
