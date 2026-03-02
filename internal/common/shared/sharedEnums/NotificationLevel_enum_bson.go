package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func (e NotificationLevel) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *NotificationLevel) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for NotificationLevel, got %v", t)
	}
	val, err := NotificationLevelString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
