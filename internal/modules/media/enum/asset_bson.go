package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// AssetType Hook
// -----------------------------------------------------------------------------
func (e AssetType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *AssetType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for AssetType, got %v", t)
	}
	val, err := AssetTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// TagStatus Hook
// -----------------------------------------------------------------------------
func (e TagStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *TagStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for TagStatus, got %v", t)
	}
	val, err := TagStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}