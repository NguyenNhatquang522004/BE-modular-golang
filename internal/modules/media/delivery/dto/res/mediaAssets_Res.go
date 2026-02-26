package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB-STRUCTS ---

type MediaMetadataRes struct {
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Duration  float64 `json:"duration,omitempty"`
	SizeBytes int64   `json:"size_bytes"`
	MimeType  string  `json:"mime_type"`
}

type TagPositionRes struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MediaTagRes struct {
	UserID   string                       `json:"user_id"`
	Name     string                       `json:"name"`
	Position TagPositionRes               `json:"position"`
	Status   sharedEnums.ProcessingStatus `json:"status"`
}

type MediaReactionStatsRes struct {
	Total int `bson:"total" json:"total"`
	Like  int `bson:"like" json:"like"`
	Love  int `bson:"love" json:"love"`
	Haha  int `bson:"haha" json:"haha"`
	Wow   int `bson:"wow" json:"wow"`
	Sad   int `bson:"sad" json:"sad"`
	Angry int `bson:"angry" json:"angry"`
}

type MediaPrivacyRes struct {
	Level            sharedEnums.PrivacyScope `json:"level"`
	InheritFromAlbum bool                     `json:"inherit_from_album"`
}

// --- MAIN RESPONSE DTO ---

type MediaAssetRes struct {
	ID             string                `json:"id"`
	UserID         string                `json:"user_id"`
	AlbumID        string                `json:"album_id,omitempty"`
	PostID         string                `json:"post_id,omitempty"`
	GroupID        string                `json:"group_id,omitempty"`
	StorageFileID  string                `json:"url"` // Map với "url"
	OriginalURL    string                `json:"original_url"`
	ThumbnailURL   string                `json:"thumbnail_url"`
	AssetType      sharedEnums.MediaType `json:"asset_type"`
	Metadata       MediaMetadataRes      `json:"metadata"`
	Caption        string                `json:"caption,omitempty"`
	Hashtags       []string              `json:"hashtags,omitempty"`
	TaggedUsers    []MediaTagRes         `json:"tagged_users,omitempty"`
	ReactionsCount MediaReactionStatsRes `json:"reactions_count"`
	CommentCount   int                   `json:"comment_count"`
	Order          int                   `json:"order"`
	Privacy        *MediaPrivacyRes      `json:"privacy,omitempty"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	DeletedAt      *time.Time            `json:"deleted_at,omitempty"`
}
