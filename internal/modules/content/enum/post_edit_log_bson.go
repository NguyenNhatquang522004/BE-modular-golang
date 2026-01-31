package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e TargetCollection) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *TargetCollection) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for TargetCollection, got %v", t)
	}
	// Hàm TargetCollectionString được sinh ra bởi 'enumer'
	val, err := TargetCollectionString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}