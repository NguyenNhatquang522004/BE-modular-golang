package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e EventLocationType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *EventLocationType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for EventLocationType, got %v", t)
	}
	val, err := EventLocationTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
