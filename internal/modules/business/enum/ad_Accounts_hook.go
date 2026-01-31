package enum

import (
	"database/sql/driver"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)


// =============================================================================
// ACCOUNT STATUS HOOKS
// =============================================================================

// 1. SQL DRIVER: VALUE (Go -> Postgres)
// Giúp GORM biết cách lưu Enum int thành string 'active' xuống DB
func (e AccountStatus) Value() (driver.Value, error) {
	return e.String(), nil
}

// 2. SQL DRIVER: SCAN (Postgres -> Go)
// Giúp GORM biết cách đọc string 'active' từ DB lên thành Enum int
func (e *AccountStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	var strVal string
	switch v := value.(type) {
	case []byte:
		strVal = string(v)
	case string:
		strVal = v
	default:
		return fmt.Errorf("failed to scan AccountStatus: %v", value)
	}

	val, err := AccountStatusString(strVal)
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// 3. BSON: MARSHAL (Go -> Mongo - Nếu cần dùng)
func (e AccountStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
return bsontype.String, []byte(e.String()), nil 
}

// 4. BSON: UNMARSHAL (Mongo -> Go - Nếu cần dùng)
func (e *AccountStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for AccountStatus, got %v", t)
	}
	val, err := AccountStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}