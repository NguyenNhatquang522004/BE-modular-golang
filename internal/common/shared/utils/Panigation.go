package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
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
