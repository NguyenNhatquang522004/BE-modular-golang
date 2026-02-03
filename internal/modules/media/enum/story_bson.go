package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// StoryMediaType
// -----------------------------------------------------------------------------
func (e StoryMediaType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *StoryMediaType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for StoryMediaType, got %v", t)
	}
	val, err := StoryMediaTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// StoryPrivacyType
// -----------------------------------------------------------------------------
func (e StoryPrivacyType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *StoryPrivacyType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for StoryPrivacyType, got %v", t)
	}
	val, err := StoryPrivacyTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// OverlayType
// -----------------------------------------------------------------------------
func (e OverlayType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

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