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

