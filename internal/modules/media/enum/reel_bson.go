package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// ProcessingStatus
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// RemixType
// -----------------------------------------------------------------------------
func (e RemixType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *RemixType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for RemixType, got %v", t)
	}
	val, err := RemixTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// ReelPrivacy
// -----------------------------------------------------------------------------
func (e ReelPrivacy) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *ReelPrivacy) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ReelPrivacy, got %v", t)
	}
	val, err := ReelPrivacyString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}