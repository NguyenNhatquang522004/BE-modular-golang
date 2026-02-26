package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (e ProcessingStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *ProcessingStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ProcessingStatus, got %v", t)
	}
	val, err := ProcessingStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
