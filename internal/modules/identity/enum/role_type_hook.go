package enum

import (
	"database/sql/driver"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// =============================================================================
// ROLE TYPE HOOKS
// =============================================================================

// 1. SQL: Value (Go -> Postgres)
// GORM sẽ gọi hàm này để lấy string lưu xuống DB (VD: "admin", "user")
func (e RoleType) Value() (driver.Value, error) {
	return e.String(), nil
}

// 2. SQL: Scan (Postgres -> Go)
// GORM sẽ gọi hàm này để parse string từ DB lên thành Enum (VD: "admin" -> RoleTypeAdmin)
func (e *RoleType) Scan(value interface{}) error {
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
		return fmt.Errorf("failed to scan RoleType: %v", value)
	}

	// Hàm RoleTypeString được sinh ra bởi lệnh go generate (thư viện enumer)
	val, err := RoleTypeString(strVal)
	if err != nil {
		return err
	}
	*e = val
	return nil
}

// 3. BSON: Marshal (Go -> Mongo)
// Khi lưu vào Mongo, nó sẽ lưu dưới dạng String thay vì Int
func (e RoleType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, []byte(e.String()), nil
}

// 4. BSON: Unmarshal (Mongo -> Go)
// Khi đọc từ Mongo lên, nó sẽ parse String thành Enum
func (e *RoleType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for RoleType, got %v", t)
	}

	// data ở đây là raw bytes của string giá trị
	val, err := RoleTypeString(string(data))
	if err != nil {
		return err
	}
	*e = val
	return nil
}
