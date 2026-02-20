package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

// --- SUB-STRUCTS ---

type ReelVideoRes struct {
	URL           string  `json:"url"`
	ThumbnailURL  string  `json:"thumbnail_url"`
	PreviewGifURL string  `json:"preview_gif_url"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Duration      float64 `json:"duration"`
}

type AudioMetaRes struct {
	TrackID         string  `json:"track_id,omitempty"`
	IsOriginalAudio bool    `json:"is_original_audio"`
	VolumeAdjust    float64 `json:"volumn_adjust"`
	AudioStartTime  float64 `json:"audio_start_time"`
}

type RemixInfoRes struct {
	ParentReelID string         `json:"parent_reel_id"`
	Type         enum.RemixType `json:"type"`
	IsRemixable  bool           `json:"is_remixable"`
}

type ReelStatsRes struct {
	Views    int `json:"views"`
	Likes    int `json:"likes"`
	Shares   int `json:"shares"`
	Saves    int `json:"saves"`
	Comments int `json:"comments"`
}

// --- MAIN RESPONSE DTO ---

type ReelRes struct {
	ID               string                `json:"id"`
	UserID           string                `json:"user_id"`
	ProcessingStatus enum.ProcessingStatus `json:"processing_status"`
	Video            ReelVideoRes          `json:"video"`
	Caption          string                `json:"caption"`
	Hashtags         []string              `json:"hashtags,omitempty"`
	Mentions         []string              `json:"mentions,omitempty"`
	AudioMeta        AudioMetaRes          `json:"audio_meta"`
	RemixInfo        *RemixInfoRes         `json:"remix_info,omitempty"`
	Stats            ReelStatsRes          `json:"stats"`
	Privacy          enum.ReelPrivacy      `json:"privacy"`
	CreatedAt        time.Time             `json:"created_at"`
	DeletedAt        *time.Time            `json:"deleted_at,omitempty"`
}