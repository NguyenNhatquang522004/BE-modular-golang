package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue: Int -> String
func (e AlbumType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *AlbumType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for AlbumType, got %v", t)
	}
	// Hàm AlbumTypeString được sinh ra bởi 'enumer'
	val, err := AlbumTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}