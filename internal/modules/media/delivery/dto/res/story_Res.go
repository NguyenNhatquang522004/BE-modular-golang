package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

// --- SUB-STRUCTS ---

type StoryMediaRes struct {
	URL          string             `json:"url"`
	Type         enum.StoryMediaType `json:"type"`
	Duration     float64            `json:"duration"`
	ThumbnailURL string             `json:"thumbnail_url"`
	SizeBytes    int64              `json:"size_bytes"`
}

type StoryPrivacyRes struct {
	Type      enum.StoryPrivacyType `json:"type"`
	AllowList []string              `json:"allow_list,omitempty"`
	BlockList []string              `json:"block_list,omitempty"`
}

type StorySettingsRes struct {
	AllowReply bool `json:"allow_reply"`
	AllowShare bool `json:"allow_share"`
}

type ViewerPreviewRes struct {
	UserID string `json:"user_id"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

type OverlayPositionRes struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"rotation"`
	Scale    float64 `json:"scale"`
}

type StoryOverlayRes struct {
	Type     enum.OverlayType       `json:"type"`
	Position OverlayPositionRes     `json:"position"`
	Data     map[string]interface{} `json:"data"`
}

type StoryStatsRes struct {
	ViewsCount int `json:"views_count"`
	LikesCount int `json:"likes_count"`
	ReplyCount int `json:"reply_count"`
}

// --- MAIN RESPONSE DTO ---

// StoryRes: Chuẩn trả về API
type StoryRes struct {
	ID             string             `json:"id"`
	UserID         string             `json:"user_id"`
	Media          StoryMediaRes      `json:"media"`
	Overlays       []*StoryOverlayRes `json:"overlays,omitempty"`
	Privacy        StoryPrivacyRes    `json:"privacy"`
	Settings       StorySettingsRes   `json:"settings"`
	PreviewViewers []ViewerPreviewRes `json:"preview_viewers,omitempty"`
	Stats          StoryStatsRes      `json:"stats"`
	CreatedAt      time.Time          `json:"created_at"`
	ExpiresAt      time.Time          `json:"expires_at"`
	IsArchived     bool               `json:"is_archived"`
}