package enum

import (
	"database/sql/driver"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Helper (Tái sử dụng)


// =============================================================================
// 1. OBJECTIVE HOOKS
// =============================================================================
func (e CampaignObjective) Value() (driver.Value, error) { return e.String(), nil }
func (e *CampaignObjective) Scan(v interface{}) error {
	return scanEnum(v, func(s string) (interface{}, error) { return CampaignObjectiveString(s) }, e)
}
func (e CampaignObjective) MarshalBSONValue() (bsontype.Type, []byte, error) { return bsontype.String, []byte(e.String()), nil }
func (e *CampaignObjective) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	val, err := CampaignObjectiveString(string(data)); if err != nil { return err }; *e = val; return nil
}

// =============================================================================
// 2. BUYING TYPE HOOKS
// =============================================================================
func (e BuyingType) Value() (driver.Value, error) { return e.String(), nil }
func (e *BuyingType) Scan(v interface{}) error {
	return scanEnum(v, func(s string) (interface{}, error) { return BuyingTypeString(s) }, e)
}
func (e BuyingType) MarshalBSONValue() (bsontype.Type, []byte, error) { return bsontype.String, []byte(e.String()), nil }
func (e *BuyingType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	val, err := BuyingTypeString(string(data)); if err != nil { return err }; *e = val; return nil
}

// --- Generic Scan Helper (Để code gọn hơn, tránh lặp lại switch case) ---
func scanEnum(value interface{}, parser func(string) (interface{}, error), target interface{}) error {
	if value == nil { return nil }
	var strVal string
	switch v := value.(type) {
	case []byte: strVal = string(v)
	case string: strVal = v
	default: return fmt.Errorf("failed to scan enum: %v", value)
	}
	val, err := parser(strVal)
	if err != nil { return err }
	
	// Gán giá trị vào pointer target (Dùng reflection hoặc ép kiểu thủ công ở trên gọn hơn)
	// Ở trên ta đã ép kiểu trong từng hàm Scan rồi, helper này chỉ minh họa logic chung.
	// Để code chạy 100% không cần reflect phức tạp, logic switch case nên nằm ở từng hàm Scan như bài trước.
	// Dưới đây là cách gọi lại logic của bài trước cho an toàn tuyệt đối:
	switch t := target.(type) {
	case *CampaignObjective: *t = val.(CampaignObjective)
	case *BuyingType: *t = val.(BuyingType)
	}
	return nil
}