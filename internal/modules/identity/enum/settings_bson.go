package enum

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// ThemeMode
// -----------------------------------------------------------------------------
func (e ThemeMode) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *ThemeMode) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ThemeMode, got %v", t)
	}
	val, err := ThemeModeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// FontSize
// -----------------------------------------------------------------------------
func (e FontSize) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *FontSize) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for FontSize, got %v", t)
	}
	val, err := FontSizeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// PrivacyLevel
// -----------------------------------------------------------------------------
func (e PrivacyLevel) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *PrivacyLevel) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PrivacyLevel, got %v", t)
	}
	val, err := PrivacyLevelString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// LangCode
// -----------------------------------------------------------------------------
func (e LangCode) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *LangCode) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for LangCode, got %v", t)
	}
	val, err := LangCodeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// -----------------------------------------------------------------------------
// EmailFrequency
// -----------------------------------------------------------------------------
func (e EmailFrequency) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

func (e *EmailFrequency) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for EmailFrequency, got %v", t)
	}
	val, err := EmailFrequencyString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}