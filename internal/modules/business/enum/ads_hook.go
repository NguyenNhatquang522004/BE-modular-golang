package enum

import (
	"database/sql/driver"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)


// =============================================================================
// AD STATUS HOOKS
// =============================================================================

// 1. SQL: Value (Go -> Postgres)
// GORM sẽ gọi hàm này để lấy string lưu xuống DB
func (e AdStatus) Value() (driver.Value, error) {
	return e.String(), nil
}

// 2. SQL: Scan (Postgres -> Go)
// GORM sẽ gọi hàm này để parse string từ DB lên thành Enum
func (e *AdStatus) Scan(value interface{}) error {
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
		return fmt.Errorf("failed to scan AdStatus: %v", value)
	}

	val, err := AdStatusString(strVal)
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// 3. BSON: Marshal (Go -> Mongo)
func (e AdStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
return bsontype.String, []byte(e.String()), nil 
}

// 4. BSON: Unmarshal (Mongo -> Go)
func (e *AdStatus) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for AdStatus, got %v", t)
	}
	val, err := AdStatusString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}