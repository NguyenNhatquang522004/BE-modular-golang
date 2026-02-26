package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e PublisherRole) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *PublisherRole) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PublisherRole, got %v", t)
	}
	// Hàm PublisherRoleString được sinh ra bởi 'enumer'
	val, err := PublisherRoleString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
