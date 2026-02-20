package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

// --- SUB-STRUCTS ---

type ReelVideoReq struct {
	URL           string  `json:"url" validate:"required"`
	ThumbnailURL  string  `json:"thumbnail_url"`
	PreviewGifURL string  `json:"preview_gif_url"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Duration      float64 `json:"duration"`
}

type AudioMetaReq struct {
	TrackID         string  `json:"track_id,omitempty"` // String để client gửi Hex ObjectID
	IsOriginalAudio bool    `json:"is_original_audio"`
	VolumeAdjust    float64 `json:"volumn_adjust"` // Giữ nguyên đúng tag "volumn_adjust"
	AudioStartTime  float64 `json:"audio_start_time"`
}

type RemixInfoReq struct {
	ParentReelID string         `json:"parent_reel_id" validate:"required"`
	Type         enum.RemixType `json:"type"`
	IsRemixable  bool           `json:"is_remixable"`
}

type ReelStatsReq struct {
	Views    int `json:"views"`
	Likes    int `json:"likes"`
	Shares   int `json:"shares"`
	Saves    int `json:"saves"`
	Comments int `json:"comments"`
}

// --- MAIN REQUEST DTO ---

// ReelReq: Ánh xạ 100% các trường từ Entity
type ReelReq struct {
	ID               string                `json:"id,omitempty"`
	UserID           string                `json:"user_id" validate:"required"`
	ProcessingStatus enum.ProcessingStatus `json:"processing_status"`
	Video            ReelVideoReq          `json:"video"`
	Caption          string                `json:"caption"`
	Hashtags         []string              `json:"hashtags,omitempty"`
	Mentions         []string              `json:"mentions,omitempty"`
	AudioMeta        AudioMetaReq          `json:"audio_meta"`
	RemixInfo        *RemixInfoReq         `json:"remix_info,omitempty"`
	Stats            ReelStatsReq          `json:"stats"`
	Privacy          enum.ReelPrivacy      `json:"privacy"`
	CreatedAt        *time.Time            `json:"created_at,omitempty"`
	DeletedAt        *time.Time            `json:"deleted_at,omitempty"`
}

// UpdateReelReq: Dùng pointer 100% để hỗ trợ Partial Update hiệu quả
type UpdateReelReq struct {
	ProcessingStatus *enum.ProcessingStatus `json:"processing_status,omitempty"`
	Video            *ReelVideoReq          `json:"video,omitempty"`
	Caption          *string                `json:"caption,omitempty"`
	Hashtags         []string               `json:"hashtags,omitempty"`
	Mentions         []string               `json:"mentions,omitempty"`
	AudioMeta        *AudioMetaReq          `json:"audio_meta,omitempty"`
	RemixInfo        *RemixInfoReq          `json:"remix_info,omitempty"`
	Stats            *ReelStatsReq          `json:"stats,omitempty"`
	Privacy          *enum.ReelPrivacy      `json:"privacy,omitempty"`
	DeletedAt        *time.Time             `json:"deleted_at,omitempty"`
}