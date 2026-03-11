package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue implements the bson.ValueMarshaler interface for OverlayType.
func (e OverlayType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface for OverlayType.
func (e *OverlayType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for OverlayType, got %v", t)
	}
	val, err := OverlayTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
