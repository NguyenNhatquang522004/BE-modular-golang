package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionMediaAssets = "MediaAssets"
)

type MediaAsset struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. OWNERSHIP & REFERENCES
	// Người upload (UUID Postgres)
	// Index: { user_id: 1, created_at: -1 } -> Load thư viện ảnh cá nhân
	UserID string `bson:"user_id" json:"user_id"`

	// Album chứa ảnh (Mongo ID)
	// Index: { album_id: 1, order: 1 } -> Load ảnh trong album
	AlbumID primitive.ObjectID `bson:"album_id,omitempty" json:"album_id,omitempty"`

	// Post chứa ảnh (Mongo ID)
	PostID primitive.ObjectID `bson:"post_id,omitempty" json:"post_id,omitempty"`

	// Group chứa ảnh (Mongo ID - Theo yêu cầu đề bài)
	// Index: { group_id: 1, created_at: -1 } -> Tab Media của nhóm
	GroupID primitive.ObjectID `bson:"group_id,omitempty" json:"group_id,omitempty"`

	PageID primitive.ObjectID `bson:"page_id,omitempty" json:"page_id,omitempty"`

	CommentID primitive.ObjectID `bson:"comment_id,omitempty" json:"comment_id,omitempty"`

	StoryID primitive.ObjectID `bson:"story_id,omitempty" json:"story_id,omitempty"`

	ReelID primitive.ObjectID `bson:"reel_id,omitempty" json:"reel_id,omitempty"`

	// 2. STORAGE LINKS
	// ID file trong hệ thống lưu trữ (GridFS / SeaweedFS / S3)
	StorageFileID string `bson:"url" json:"url"` // Mapping với trường "url" trong JSON đề bài

	OriginalURL  string `bson:"original_url" json:"original_url"`   // Link full HD
	ThumbnailURL string `bson:"thumbnail_url" json:"thumbnail_url"` // Link ảnh nhỏ

	AssetType sharedEnums.MediaType `bson:"asset_type" json:"asset_type"`

	// 3. METADATA & CONTENT
	Metadata MediaMetadata `bson:"metadata" json:"metadata"`
	Caption  string        `bson:"caption,omitempty" json:"caption,omitempty"`

	// Index: Multikey { hashtags: 1 }
	Hashtags []string `bson:"hashtags,omitempty" json:"hashtags,omitempty"`

	// 4. SOCIAL FEATURES
	TaggedUsers []MediaTag `bson:"tagged_users,omitempty" json:"tagged_users,omitempty"`

	// Stats (Nên tách struct để dễ mở rộng)
	ReactionsCount MediaReactionStats `bson:"reactions_count" json:"reactions_count"`
	CommentCount   int                `bson:"comment_count" json:"comment_count"`

	// 5. SETTINGS
	Order   int           `bson:"order" json:"order"`
	Privacy *MediaPrivacy `bson:"privacy,omitempty" json:"privacy,omitempty"`

	// 6. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
type MediaAssetsContext struct {
	Type sharedEnums.ContextType `bson:"type" json:"type"`

	// TargetID có thể là GroupID (ObjectId) hoặc UserID (UUID) tùy context.
	// Lưu String là an toàn nhất để chứa cả 2 loại.
	TargetID string `bson:"target_id,omitempty" json:"target_id,omitempty"`
}
// --- METADATA ---
type MediaMetadata struct {
	Width     int     `bson:"width" json:"width"`
	Height    int     `bson:"height" json:"height"`
	Duration  float64 `bson:"duration,omitempty" json:"duration,omitempty"` // Seconds
	SizeBytes int64   `bson:"size_bytes" json:"size_bytes"`
	MimeType  string  `bson:"mime_type" json:"mime_type"` // "image/jpeg"
}

// --- TAGGED USER ---
type MediaTag struct {
	// UserID từ Postgres (UUID) -> Lưu String
	UserID string `bson:"user_id" json:"user_id"`

	Name string `bson:"name" json:"name"` // Cache tên

	// Tọa độ gắn thẻ (0.0 -> 1.0)
	Position TagPosition `bson:"position" json:"position"`

	Status sharedEnums.ProcessingStatus `bson:"status" json:"status"`
}

type TagPosition struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}

// --- REACTIONS COUNT ---
type MediaReactionStats struct {
	Total int `bson:"total" json:"total"`
	Like  int `bson:"like" json:"like"`
	Love  int `bson:"love" json:"love"`
	Haha  int `bson:"haha" json:"haha"`
	Wow   int `bson:"wow" json:"wow"`
	Sad   int `bson:"sad" json:"sad"`
	Angry int `bson:"angry" json:"angry"`
}

// --- PRIVACY SETTINGS ---
type MediaPrivacy struct {
	// Level string hoặc dùng Enum PrivacyScope tái sử dụng
	Level            sharedEnums.PrivacyScope `bson:"level" json:"level"`
	InheritFromAlbum bool                     `bson:"inherit_from_album" json:"inherit_from_album"`
}

func (MediaAsset) CollectionName() string {
	return CollectionMediaAssets
}
