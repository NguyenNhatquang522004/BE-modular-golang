package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionnamepostextension = "PostExtension"
)

type PostExtension struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Index: Unique. Một Post chỉ có 1 Extension document.
	PostID primitive.ObjectID `bson:"post_id" json:"post_id"`

	// Các trường mở rộng (Dùng Pointer để tiết kiệm bộ nhớ & omitempty)
	// Chỉ 1 trong các trường này sẽ có dữ liệu tại 1 thời điểm.
	ShareData      *ShareData      `bson:"share_data,omitempty" json:"share_data,omitempty"`
	BackgroundData *BackgroundData `bson:"background_data,omitempty" json:"background_data,omitempty"`
	QnAData        *QnAData        `bson:"qna_data,omitempty" json:"qna_data,omitempty"`
	ActivityData   *ActivityData   `bson:"activity_data,omitempty" json:"activity_data,omitempty"`
	LocationDetail *LocationDetail `bson:"location_detail,omitempty" json:"location_detail,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- A. SHARE DATA ---
type ShareData struct {
	OriginalPostID primitive.ObjectID `bson:"original_post_id" json:"original_post_id"` // Bài gốc
	ParentPostID   primitive.ObjectID `bson:"parent_post_id" json:"parent_post_id"`     // Bài cha trực tiếp

	// Snapshot: Cache thông tin bài gốc để hiển thị nhanh mà không cần join
	Snapshot ShareSnapshot `bson:"snapshot" json:"snapshot"`
}

type ShareSnapshot struct {
	// AuthorID là UUID từ Postgres -> Lưu String
	AuthorID string `bson:"author_id" json:"author_id"`

	AuthorName     string    `bson:"author_name" json:"author_name"`
	AuthorAvatar   string    `bson:"author_avatar" json:"author_avatar"`
	ContentExcerpt string    `bson:"content_excerpt" json:"content_excerpt"` // Trích đoạn nội dung ngắn
	MediaThumb     string    `bson:"media_thumb,omitempty" json:"media_thumb,omitempty"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}

// --- B. BACKGROUND DATA ---
type BackgroundData struct {
	ThemeID   string `bson:"theme_id" json:"theme_id"`
	TextColor string `bson:"text_color" json:"text_color"` // VD: "#FFFFFF"
}

// --- C. Q&A DATA ---
type QnAData struct {
	Question   string `bson:"question" json:"question"`
	ButtonText string `bson:"button_text" json:"button_text"` // VD: "Trả lời ngay"
}

// --- D. ACTIVITY DATA ---
type ActivityData struct {
	Type       enum.ActivityType `bson:"type" json:"type"`                               // Enum: watching, traveling...
	ObjectID   string            `bson:"object_id,omitempty" json:"object_id,omitempty"` // ID phim/sách (nếu có)
	ObjectName string            `bson:"object_name" json:"object_name"`                 // "Phim Mai", "Hà Nội"
}

// --- E. LOCATION DETAIL (GEOJSON) ---
type LocationDetail struct {
	// Type luôn là "Point" theo chuẩn GeoJSON MongoDB
	Type string `bson:"type" json:"type"`

	// Coordinates: [Longitude, Latitude] - Kinh độ trước, Vĩ độ sau
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`

	Address string `bson:"address" json:"address"`
	MapURL  string `bson:"map_url,omitempty" json:"map_url,omitempty"`
}

func (PostExtension) Collectionnamepostextension() string {
	return collectionnamepostextension
}
