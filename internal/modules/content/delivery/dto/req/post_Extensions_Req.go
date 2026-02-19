package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

// CreatePostExtensionReq: Dùng khi tạo mới extension (thường đi kèm tạo Post)
type CreatePostExtensionReq struct {
	PostID string `json:"post_id" binding:"required,mongoId"`

	// Các field con là Pointer để client có thể gửi null
	ShareData      *ShareDataReq      `json:"share_data,omitempty" binding:"omitempty"`
	BackgroundData *BackgroundDataReq `json:"background_data,omitempty" binding:"omitempty"`
	QnAData        *QnADataReq        `json:"qna_data,omitempty" binding:"omitempty"`
	ActivityData   *ActivityDataReq   `json:"activity_data,omitempty" binding:"omitempty"`
	LocationDetail *LocationDetailReq `json:"location_detail,omitempty" binding:"omitempty"`
}

// UpdatePostExtensionReq: Dùng khi edit bài viết, có thể thay đổi extension
type UpdatePostExtensionReq struct {
	// Gửi field nào thì update field đó. Gửi nil thì giữ nguyên.
	ShareData      *ShareDataReq      `json:"share_data,omitempty"`
	BackgroundData *BackgroundDataReq `json:"background_data,omitempty"`
	QnAData        *QnADataReq        `json:"qna_data,omitempty"`
	ActivityData   *ActivityDataReq   `json:"activity_data,omitempty"`
	LocationDetail *LocationDetailReq `json:"location_detail,omitempty"`
}

// --- Nested Structs ---

type ShareDataReq struct {
	OriginalPostID string           `json:"original_post_id" binding:"required,mongoId"`
	ParentPostID   string           `json:"parent_post_id" binding:"required,mongoId"`
	Snapshot       ShareSnapshotReq `json:"snapshot" binding:"required"`
}

type ShareSnapshotReq struct {
	AuthorID       string    `json:"author_id" binding:"required,uuid"` // Validate UUID Postgres
	AuthorName     string    `json:"author_name" binding:"required"`
	AuthorAvatar   string    `json:"author_avatar" binding:"omitempty,url"`
	ContentExcerpt string    `json:"content_excerpt" binding:"omitempty,max=200"` // Giới hạn độ dài preview
	MediaThumb     string    `json:"media_thumb,omitempty" binding:"omitempty,url"`
	CreatedAt      time.Time `json:"created_at" binding:"required"`
}

type BackgroundDataReq struct {
	ThemeID   string `json:"theme_id" binding:"required"`
	TextColor string `json:"text_color" binding:"required,hexcolor"` // Validate mã màu (VD: #FFFFFF)
}

type QnADataReq struct {
	Question   string `json:"question" binding:"required,max=500"`
	ButtonText string `json:"button_text" binding:"required,max=50"`
}

type ActivityDataReq struct {
	Type       enum.ActivityType `json:"type" binding:"required"` // Validate enum value
	ObjectID   string            `json:"object_id,omitempty"`
	ObjectName string            `json:"object_name" binding:"required"`
}

type LocationDetailReq struct {
	// Coordinates: [Longitude, Latitude]
	// Validate độ dài mảng phải là 2 và giá trị hợp lệ
	Coordinates []float64 `json:"coordinates" binding:"required,len=2,dive,number"` 
	Address     string    `json:"address" binding:"required"`
	MapURL      string    `json:"map_url,omitempty" binding:"omitempty,url"`
}