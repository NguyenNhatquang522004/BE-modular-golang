package enum

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// -----------------------------------------------------------------------------
// FileType
// -----------------------------------------------------------------------------
func (e FileType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// Gọi hàm helper marshalEnum từ file common_bson.go (nếu đã tạo)
	// Hoặc dùng: return bsontype.String, []byte(e.String()), nil
	return bsontype.String, []byte(e.String()), nil
}

func (e *FileType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("expected string for FileType, got %v", t)
	}
	val, err := FileTypeString(string(data))
	if err != nil {
		// Fallback: Nếu gặp file lạ chưa có trong enum, trả về Other thay vì lỗi
		*e = FileTypeOther
		return nil
	}
	*e = val
	return nil
}