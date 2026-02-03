package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// GroupPrivacy
// -----------------------------------------------------------------------------
func (e GroupPrivacy) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *GroupPrivacy) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for GroupPrivacy, got %v", t)
	}
	val, err := GroupPrivacyString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// GroupApprover
// -----------------------------------------------------------------------------
func (e GroupApprover) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *GroupApprover) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for GroupApprover, got %v", t)
	}
	val, err := GroupApproverString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
