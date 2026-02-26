package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

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
