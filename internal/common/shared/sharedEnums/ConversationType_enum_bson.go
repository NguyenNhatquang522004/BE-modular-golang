package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e ConversationType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *ConversationType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ConversationType, got %v", t)
	}
	val, err := ConversationTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
