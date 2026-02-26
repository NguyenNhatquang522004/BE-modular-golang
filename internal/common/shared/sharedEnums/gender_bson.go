package sharedEnums

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// --- 1. MARSHAL: Ghi vào MongoDB ---
// Hook này được gọi khi bạn lưu struct vào DB.
// Nhiệm vụ: Biến Int (0, 1) -> String ("male", "female")
func (e Gender) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// e.String() là hàm do 'enumer' tự sinh ra.
	// Nó trả về chuỗi đã transform snake_case (vd: "male")
	return bsontype.String, []byte(e.String()), nil
}

// --- 2. UNMARSHAL: Đọc từ MongoDB ---
// Hook này được gọi khi bạn Find dữ liệu từ DB lên struct.
// Nhiệm vụ: Biến String ("male") -> Int (0)
func (e *Gender) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	// Kiểm tra xem dữ liệu trong DB có phải là String không
	if t != bsontype.String {
		return fmt.Errorf("expected string for Gender, got %v", t)
	}

	// GenderString() là hàm do 'enumer' tự sinh ra.
	// Nó tìm chuỗi "male" và trả về số 0.
	val, err := GenderString(string(data))
	if err != nil {
		return err
	}

	*e = val
	return nil
}
