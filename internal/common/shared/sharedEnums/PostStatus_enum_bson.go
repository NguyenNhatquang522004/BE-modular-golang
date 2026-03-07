package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (e PostStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *PostStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PostStatus, got %v", t)
	}
	val, err := PostStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
