package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

// --- SUB-STRUCTS ---

type StoryMediaReq struct {
	URL          string                `json:"url" validate:"required"`
	Type         sharedEnums.MediaType `json:"type"`
	Duration     float64               `json:"duration"`
	ThumbnailURL string                `json:"thumbnail_url"`
	SizeBytes    int64                 `json:"size_bytes"`
}

type StoryPrivacyReq struct {
	Type      sharedEnums.PrivacyScope `json:"type"`
	AllowList []string                 `json:"allow_list,omitempty"`
	BlockList []string                 `json:"block_list,omitempty"`
}

type StorySettingsReq struct {
	AllowReply bool `json:"allow_reply"`
	AllowShare bool `json:"allow_share"`
}

type ViewerPreviewReq struct {
	UserID string `json:"user_id"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

type OverlayPositionReq struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"rotation"`
	Scale    float64 `json:"scale"`
}

type StoryOverlayReq struct {
	Type     enum.OverlayType       `json:"type"`
	Position OverlayPositionReq     `json:"position"`
	Data     map[string]interface{} `json:"data"`
}

type StoryStatsReq struct {
	ViewsCount int `bson:"views_count" json:"views_count"`
	Likes      int `bson:"likes_count" json:"likes_count"`
	Love       int `bson:"love" json:"love"`
	Haha       int `bson:"haha" json:"haha"`
	Wow        int `bson:"wow" json:"wow"`
	Sad        int `bson:"sad" json:"sad"`
	Angry      int `bson:"angry" json:"angry"`
	ReplyCount int `bson:"reply_count" json:"reply_count"`
}

// --- MAIN REQUEST DTO ---

// StoryReq: Ánh xạ 100% Entity
type StoryReq struct {
	ID             string             `json:"id,omitempty"` // Trình bày dạng chuỗi cho Hex ObjectID
	UserID         string             `json:"user_id" validate:"required"`
	Media          StoryMediaReq      `json:"media"`
	Overlays       []*StoryOverlayReq `json:"overlays,omitempty"`
	Privacy        StoryPrivacyReq    `json:"privacy"`
	Settings       StorySettingsReq   `json:"settings"`
	PreviewViewers []ViewerPreviewReq `json:"preview_viewers,omitempty"`
	Stats          StoryStatsReq      `json:"stats"`
	CreatedAt      *time.Time         `json:"created_at,omitempty"`
	ExpiresAt      *time.Time         `json:"expires_at,omitempty"`
	IsArchived     bool               `json:"is_archived"` 
}

// UpdateStoryReq: Áp dụng 100% pointer để hỗ trợ Partial Update
type UpdateStoryReq struct {
	Media          *StoryMediaReq     `json:"media,omitempty"`
	Overlays       []*StoryOverlayReq `json:"overlays,omitempty"`
	Privacy        *StoryPrivacyReq   `json:"privacy,omitempty"`
	Settings       *StorySettingsReq  `json:"settings,omitempty"`
	PreviewViewers []ViewerPreviewReq `json:"preview_viewers,omitempty"`
	Stats          *StoryStatsReq     `json:"stats,omitempty"`
	ExpiresAt      *time.Time         `json:"expires_at,omitempty"`
	IsArchived     *bool              `json:"is_archived,omitempty"`
}
