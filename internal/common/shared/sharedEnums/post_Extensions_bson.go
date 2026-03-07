package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e ActivityType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *ActivityType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ActivityType, got %v", t)
	}
	// Hàm ActivityTypeString được sinh ra bởi lệnh 'go generate'
	val, err := ActivityTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}