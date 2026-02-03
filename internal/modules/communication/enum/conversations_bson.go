package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Helper function để giảm code lặp
func marshalEnum(s string) (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(s), nil
}

// -----------------------------------------------------------------------------
// ConversationType
// -----------------------------------------------------------------------------
func (e ConversationType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return marshalEnum(e.String())
}

func (e *ConversationType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ConversationType, got %v", t)
	}
	val, err := ConversationTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// ConversationScope
// -----------------------------------------------------------------------------
func (e ConversationScope) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return marshalEnum(e.String())
}

func (e *ConversationScope) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ConversationScope, got %v", t)
	}
	val, err := ConversationScopeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// ConversationStatus
// -----------------------------------------------------------------------------
func (e ConversationStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return marshalEnum(e.String())
}

func (e *ConversationStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ConversationStatus, got %v", t)
	}
	val, err := ConversationStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// PermissionLevel
// -----------------------------------------------------------------------------
func (e PermissionLevel) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return marshalEnum(e.String())
}

func (e *PermissionLevel) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PermissionLevel, got %v", t)
	}
	val, err := PermissionLevelString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

