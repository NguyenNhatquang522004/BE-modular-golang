package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// LiveStatus
// -----------------------------------------------------------------------------
func (e LiveStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *LiveStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for LiveStatus, got %v", t)
	}
	val, err := LiveStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}