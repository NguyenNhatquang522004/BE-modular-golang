package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// FollowType
// -----------------------------------------------------------------------------
func (e FollowType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *FollowType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for FollowType, got %v", t)
	}
	val, err := FollowTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// NotificationLevel
// -----------------------------------------------------------------------------
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
