package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cursor struct {
	CreatedAt time.Time `json:"t"` // Cột sort chính (time)
	ID        uuid.UUID `json:"i"` // Cột tie-breaker (unique ID)
}

func EncodeCursor(t time.Time, id uuid.UUID) string {
	// Tạo struct
	c := Cursor{
		CreatedAt: t,
		ID:        id,
	}

	// 1. Marshal sang JSON
	data, err := json.Marshal(c)
	if err != nil {
		// Trong thực tế, lỗi marshal với struct đơn giản gần như không xảy ra.
		// Có thể log error hoặc return chuỗi rỗng tùy strategy.
		return ""
	}

	// 2. Encode sang Base64 URL-safe (Raw = không padding '=')
	return base64.RawURLEncoding.EncodeToString(data)
}
func DecodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	// Check rỗng
	if encodedCursor == "" {
		return time.Time{}, uuid.Nil, nil
	}

	// 1. Decode Base64
	data, err := base64.RawURLEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, uuid.Nil, errors.New("invalid cursor format (base64)")
	}

	// 2. Unmarshal JSON
	var c Cursor
	if err := json.Unmarshal(data, &c); err != nil {
		return time.Time{}, uuid.Nil, errors.New("invalid cursor format (json)")
	}

	// 3. Validate dữ liệu (Tuỳ chọn nhưng nên làm)
	if c.ID == uuid.Nil || c.CreatedAt.IsZero() {
		return time.Time{}, uuid.Nil, errors.New("invalid cursor data")
	}

	return c.CreatedAt, c.ID, nil
}

// / mongodb
// CursorMongodb chứa dữ liệu đã giải mã
type CursorMongodb struct {
	CreatedAt time.Time
	PostID    primitive.ObjectID
}

// EncodeCursor: Tạo chuỗi cursor từ bài viết cuối cùng
// Format: "RFCNanoTimestamp,HexID" (base64 encoded)
func EncodeCursorMongodb(t time.Time, id primitive.ObjectID) string {
	// Dùng RFC3339Nano để giữ độ chính xác tối đa của MongoDB Date
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), id.Hex())
	return base64.StdEncoding.EncodeToString([]byte(key))
}

// DecodeCursor: Giải mã chuỗi cursor từ Client gửi lên
func DecodeCursorMongodb(cursor string) (*CursorMongodb, error) {
	if cursor == "" {
		return nil, nil // Cursor rỗng -> Trang đầu tiên
	}

	bytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, errors.New("invalid cursor format")
	}

	parts := strings.Split(string(bytes), ",")
	if len(parts) != 2 {
		return nil, errors.New("invalid cursor data")
	}

	// 1. Parse Time
	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, errors.New("invalid cursor time")
	}

	// 2. Parse ID
	oid, err := primitive.ObjectIDFromHex(parts[1])
	if err != nil {
		return nil, errors.New("invalid cursor id")
	}

	return &CursorMongodb{
		CreatedAt: createdAt,
		PostID:    oid,
	}, nil
}

// EncodeCursorCassandra mã hóa Paging State của gocql thành chuỗi Base64
func EncodeCursorCassandra(pageState []byte) string {
	if len(pageState) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(pageState)
}

// DecodeCursorCassandra giải mã chuỗi Base64 từ Client thành Paging State
func DecodeCursorCassandra(cursor string) ([]byte, error) {
	if cursor == "" {
		return nil, nil // Trang đầu tiên
	}
	bytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, errors.New("invalid cursor format")
	}
	return bytes, nil
}