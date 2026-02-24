package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB DTOs ---
type MediaMetadataReq struct {
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
	Duration  float64 `json:"duration,omitempty"`
	SizeBytes int64   `json:"size_bytes" validate:"required"`
	MimeType  string  `json:"mime_type" validate:"required"`
}

type TaggedUserReq struct {
	UserID string  `json:"user_id" validate:"required"`
	Name   string  `json:"name" validate:"required"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

type MediaItemReq struct {
	ID           string                `json:"id,omitempty"` // ID của item (nếu có để update)
	MediaType    sharedEnums.MediaType `json:"media_type" validate:"required"`
	URL          string                `json:"url" validate:"required"`
	ThumbnailURL string                `json:"thumbnail_url"`
	Metadata     MediaMetadataReq      `json:"metadata"`
	Order        int                   `json:"order"`
	TaggedUsers  []TaggedUserReq       `json:"tagged_users,omitempty"`
}

// --- MAIN DTO ---
type PostMediaReq struct {
	ID        string          `json:"id,omitempty"`
	PostID    string          `json:"post_id" validate:"required"`
	Items     []*MediaItemReq `json:"items"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
}
