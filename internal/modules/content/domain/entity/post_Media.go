package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionnamepostmedia = "PostMedia"
)

type PostMedia struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Index: Unique. Một Post chỉ có 1 document PostMedia tương ứng.
	PostID primitive.ObjectID `bson:"post_id" json:"post_id"`

	// Danh sách các file media (ảnh, video)
	Items     []*MediaItem `bson:"items" json:"items"`
	CreatedAt time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time    `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time   `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- 2. Sub-Struct: MEDIA ITEM ---
type MediaItem struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	MediaType    sharedEnums.MediaType `bson:"media_type" json:"media_type"`
	URL          string                `bson:"url" json:"url"`                     // Ảnh gốc / Video gốc
	ThumbnailURL string                `bson:"thumbnail_url" json:"thumbnail_url"` // Ảnh nhỏ load cho nhanh

	Metadata MediaMetadata `bson:"metadata" json:"metadata"`
	Order    int           `bson:"order" json:"order"` // 1, 2, 3...

	// Danh sách người được tag trong ảnh này
	TaggedUsers []TaggedUser `bson:"tagged_users,omitempty" json:"tagged_users,omitempty"`
}

// --- 3. Sub-Struct: METADATA ---
type MediaMetadata struct {
	Width     int     `bson:"width,omitempty" json:"width,omitempty"`       // Pixel
	Height    int     `bson:"height,omitempty" json:"height,omitempty"`     // Pixel
	Duration  float64 `bson:"duration,omitempty" json:"duration,omitempty"` // Seconds (dùng float cho video chính xác)
	SizeBytes int64   `bson:"size_bytes" json:"size_bytes"`                 // Byte
	MimeType  string  `bson:"mime_type" json:"mime_type"`                   // image/jpeg, video/mp4
}

// --- 4. Sub-Struct: TAGGED USER ---
type TaggedUser struct {
	// UserID từ Postgres (UUID) -> Bắt buộc lưu String
	UserID string `bson:"user_id" json:"user_id"`

	// Cache tên user tại thời điểm tag để hiển thị nhanh (khỏi query ngược về Identity)
	Name string `bson:"name" json:"name"`

	// Tọa độ tag trên ảnh (tỉ lệ 0.0 -> 1.0)
	// Ví dụ: x=0.5, y=0.5 là chính giữa ảnh
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}

func (PostMedia) CollectionNamePostMedia() string {
	return collectionnamepostmedia
}
