package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
const (
    CollectionStories = "Stories"
)
type Story struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. OWNERSHIP
	// UserID từ Postgres (UUID) -> String
	// Index: { user_id: 1, created_at: -1 } -> Load story của 1 người
	UserID string `bson:"user_id" json:"user_id"`

	// 2. CONTENT
	Media    StoryMedia      `bson:"media" json:"media"`
	Overlays []*StoryOverlay `bson:"overlays,omitempty" json:"overlays,omitempty"`

	// 3. CONFIGURATION
	// [CẬP NHẬT] Dùng struct Privacy nâng cao thay vì string đơn giản
	Privacy  StoryPrivacy  `bson:"privacy" json:"privacy"`
	Settings StorySettings `bson:"settings" json:"settings"`

	// 4. INTERACTION CACHE
	// Lưu tối đa 3 người xem gần nhất để hiển thị thumbnail chồng lên nhau
	PreviewViewers []ViewerPreview `bson:"preview_viewers,omitempty" json:"preview_viewers,omitempty"`
	Stats          StoryStats      `bson:"stats" json:"stats"`

	// 5. LIFECYCLE (Logic 24h)
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	// Thời điểm hết hạn (thường là CreatedAt + 24h)
	// Index: TTL (Time To Live) có thể dùng ở đây NẾU muốn xóa thật.
	// Nhưng vì có tính năng Archive, ta dùng Index thường để lọc.
	// Index: { expires_at: 1, is_archived: 1 }
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`

	// True nếu đã quá 24h và được chuyển vào kho lưu trữ
	IsArchived bool `bson:"is_archived" json:"is_archived"`
}

// --- MEDIA ---
type StoryMedia struct {
	URL          string              `bson:"url" json:"url"` // Link SeaweedFS
	Type         enum.StoryMediaType `bson:"type" json:"type"`
	Duration     float64             `bson:"duration" json:"duration"` // Seconds
	ThumbnailURL string              `bson:"thumbnail_url" json:"thumbnail_url"`
	SizeBytes    int64               `bson:"size_bytes" json:"size_bytes"`
}

// --- PRIVACY ---
type StoryPrivacy struct {
	Type enum.StoryPrivacyType `bson:"type" json:"type"`

	// Danh sách UserID (Postgres UUID -> String)
	AllowList []string `bson:"allow_list,omitempty" json:"allow_list,omitempty"`
	BlockList []string `bson:"block_list,omitempty" json:"block_list,omitempty"`
}

// --- SETTINGS ---
type StorySettings struct {
	AllowReply bool `bson:"allow_reply" json:"allow_reply"`
	AllowShare bool `bson:"allow_share" json:"allow_share"`
}

// --- PREVIEW VIEWERS (Cache) ---
type ViewerPreview struct {
	UserID string `bson:"user_id" json:"user_id"` // UUID String
	Avatar string `bson:"avatar" json:"avatar"`
	Name   string `bson:"name" json:"name"`
}

// --- OVERLAYS (Sticker/Poll/Music) ---
type StoryOverlay struct {
	Type     enum.OverlayType `bson:"type" json:"type"`
	Position OverlayPosition  `bson:"position" json:"position"`

	// Dữ liệu động tùy theo Type
	// VD: Poll -> { question: "...", options: [] }
	// VD: Music -> { track_id: "...", start_time: 15 }
	Data map[string]interface{} `bson:"data" json:"data"`
}

type OverlayPosition struct {
	X        float64 `bson:"x" json:"x"`
	Y        float64 `bson:"y" json:"y"`
	Rotation float64 `bson:"rotation" json:"rotation"`
	Scale    float64 `bson:"scale" json:"scale"`
}

// --- STATS ---
type StoryStats struct {
	ViewsCount int `bson:"views_count" json:"views_count"`
	LikesCount int `bson:"likes_count" json:"likes_count"`
	ReplyCount int `bson:"reply_count" json:"reply_count"`
}
func (Story) CollectionName() string {
    return CollectionStories
}