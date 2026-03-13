package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (e ConversationScope) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *ConversationScope) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ConversationScope, got %v", t)
	}
	val, err := ConversationScopeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
