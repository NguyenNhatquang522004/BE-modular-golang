package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB DTOs ---
type ShareSnapshotReq struct {
	AuthorID       string    `json:"author_id" validate:"required"`
	AuthorName     string    `json:"author_name"`
	AuthorAvatar   string    `json:"author_avatar"`
	ContentExcerpt string    `json:"content_excerpt"`
	MediaThumb     string    `json:"media_thumb,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ShareDataReq struct {
	OriginalPostID string           `json:"original_post_id" validate:"required"`
	ParentPostID   string           `json:"parent_post_id" validate:"required"`
	Snapshot       ShareSnapshotReq `json:"snapshot"`
}

type BackgroundDataReq struct {
	ThemeID   string `json:"theme_id" validate:"required"`
	TextColor string `json:"text_color" validate:"required"`
}

type QnADataReq struct {
	Question   string `json:"question" validate:"required"`
	ButtonText string `json:"button_text" validate:"required"`
}

type ActivityDataReq struct {
	Type       sharedEnums.ActivityType `json:"type" validate:"required"`
	ObjectID   string            `json:"object_id,omitempty"`
	ObjectName string            `json:"object_name" validate:"required"`
}

type LocationDetailReq struct {
	Type        string    `json:"type" validate:"required"` // Thường là "Point"
	Coordinates []float64 `json:"coordinates" validate:"required"`
	Address     string    `json:"address" validate:"required"`
	MapURL      string    `json:"map_url,omitempty"`
}

// --- MAIN DTO ---
type PostExtensionReq struct {
	ID             string             `json:"id,omitempty"`
	PostID         string             `json:"post_id" validate:"required"`
	ShareData      *ShareDataReq      `json:"share_data,omitempty"`
	BackgroundData *BackgroundDataReq `json:"background_data,omitempty"`
	QnAData        *QnADataReq        `json:"qna_data,omitempty"`
	ActivityData   *ActivityDataReq   `json:"activity_data,omitempty"`
	LocationDetail *LocationDetailReq `json:"location_detail,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
