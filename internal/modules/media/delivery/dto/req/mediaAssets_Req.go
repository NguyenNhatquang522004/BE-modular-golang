package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB-STRUCTS ---

type MediaMetadataReq struct {
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Duration  float64 `json:"duration,omitempty"`
	SizeBytes int64   `json:"size_bytes"`
	MimeType  string  `json:"mime_type"`
}

type TagPositionReq struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MediaTagReq struct {
	UserID   string                       `json:"user_id" validate:"required"`
	Name     string                       `json:"name"`
	Position TagPositionReq               `json:"position"`
	Status   sharedEnums.ProcessingStatus `json:"status"`
}

type MediaReactionStatsReq struct {
	Total int `bson:"total" json:"total"`
	Like  int `bson:"like" json:"like"`
	Love  int `bson:"love" json:"love"`
	Haha  int `bson:"haha" json:"haha"`
	Wow   int `bson:"wow" json:"wow"`
	Sad   int `bson:"sad" json:"sad"`
	Angry int `bson:"angry" json:"angry"`
}

type MediaPrivacyReq struct {
	Level            sharedEnums.PrivacyScope `json:"level"`
	InheritFromAlbum bool                     `json:"inherit_from_album"`
}

// --- MAIN REQUEST DTO ---

// MediaAssetReq: Ánh xạ 100% Entity
type MediaAssetReq struct {
	ID             string                `json:"id,omitempty"`
	UserID         string                `json:"user_id" validate:"required"`
	AlbumID        string                `json:"album_id,omitempty"`
	PostID         string                `json:"post_id,omitempty"`
	GroupID        string                `json:"group_id,omitempty"`
	PageID         string                `json:"page_id,omitempty"`
	StoryID        string                `json:"story_id,omitempty"`
	StorageFileID  string                `json:"url" validate:"required"` // Map chuẩn JSON "url"
	OriginalURL    string                `json:"original_url"`
	ThumbnailURL   string                `json:"thumbnail_url"`
	AssetType      sharedEnums.MediaType `json:"asset_type"`
	Metadata       MediaMetadataReq      `json:"metadata"`
	Caption        string                `json:"caption,omitempty"`
	Hashtags       []string              `json:"hashtags,omitempty"`
	TaggedUsers    []MediaTagReq         `json:"tagged_users,omitempty"`
	ReactionsCount MediaReactionStatsReq `json:"reactions_count"`
	CommentCount   int                   `json:"comment_count"`
	Order          int                   `json:"order"`
	Privacy        *MediaPrivacyReq      `json:"privacy,omitempty"`
	CreatedAt      *time.Time            `json:"created_at,omitempty"`
	UpdatedAt      *time.Time            `json:"updated_at,omitempty"`
	DeletedAt      *time.Time            `json:"deleted_at,omitempty"`
}

// UpdateMediaAssetReq: Dùng pointer 100% cho Partial Update
type UpdateMediaAssetReq struct {
	AlbumID *string `json:"album_id,omitempty"`
	PostID  *string `json:"post_id,omitempty"`
	GroupID *string `json:"group_id,omitempty"`
	PageID  *string `json:"page_id,omitempty"`
	StoryID        *string                `json:"story_id,omitempty"`
	StorageFileID *string `json:"url,omitempty"`

	OriginalURL    *string                `json:"original_url,omitempty"`
	ThumbnailURL   *string                `json:"thumbnail_url,omitempty"`
	AssetType      *sharedEnums.MediaType `json:"asset_type,omitempty"`
	Metadata       *MediaMetadataReq      `json:"metadata,omitempty"`
	Caption        *string                `json:"caption,omitempty"`
	Hashtags       []string               `json:"hashtags,omitempty"`
	TaggedUsers    []MediaTagReq          `json:"tagged_users,omitempty"`
	ReactionsCount *MediaReactionStatsReq `json:"reactions_count,omitempty"`
	CommentCount   *int                   `json:"comment_count,omitempty"`
	Order          *int                   `json:"order,omitempty"`
	Privacy        *MediaPrivacyReq       `json:"privacy,omitempty"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
}
