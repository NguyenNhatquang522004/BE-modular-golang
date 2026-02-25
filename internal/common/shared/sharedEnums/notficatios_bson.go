package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MarshalBSONValue implement interface bson.ValueMarshaler cho NotificationType
func (n NotificationType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(n.String()), nil
}

// UnmarshalBSONValue implement interface bson.ValueUnmarshaler cho NotificationType
func (n *NotificationType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for NotificationType, got %v", t)
	}

	val, err := NotificationTypeString(string(data))
	if err != nil {
		return err
	}

	*n = val
	return nil
}
