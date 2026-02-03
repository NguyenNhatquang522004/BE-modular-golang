package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// =============================================================================
// 1. ENUM: PostType
// =============================================================================

// MarshalBSONValue: Int -> String
func (e PostType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *PostType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PostType, got %v", t)
	}
	// Hàm PostTypeString được sinh ra bởi 'enumer'
	val, err := PostTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// =============================================================================
// 2. ENUM: ContextType
// =============================================================================

// MarshalBSONValue: Int -> String
func (e ContextType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *ContextType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for ContextType, got %v", t)
	}
	val, err := ContextTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// =============================================================================
// 3. ENUM: PrivacyScope
// =============================================================================

// MarshalBSONValue: Int -> String
func (e PrivacyScope) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *PrivacyScope) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PrivacyScope, got %v", t)
	}
	// Lưu ý: Tên hàm này dựa trên tên Type bạn đặt (PrivacyScope -> PrivacyScopeString)
	val, err := PrivacyScopeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// =============================================================================
// 4. ENUM: PostStatus
// =============================================================================

// MarshalBSONValue: Int -> String
func (e PostStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// UnmarshalBSONValue: String -> Int
func (e *PostStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for PostStatus, got %v", t)
	}
	val, err := PostStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}