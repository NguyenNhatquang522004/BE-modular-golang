package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e SavedTargetType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *SavedTargetType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for SavedTargetType, got %v", t)
	}
	// Hàm SavedTargetTypeString được sinh ra bởi 'enumer'
	val, err := SavedTargetTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}