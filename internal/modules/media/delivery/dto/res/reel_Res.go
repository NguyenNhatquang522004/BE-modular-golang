package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
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
	Views    int `bson:"views" json:"views"`
	Total    int `bson:"total" json:"total"`
	Like     int `bson:"like" json:"like"`
	Love     int `bson:"love" json:"love"`
	Haha     int `bson:"haha" json:"haha"`
	Wow      int `bson:"wow" json:"wow"`
	Sad      int `bson:"sad" json:"sad"`
	Angry    int `bson:"angry" json:"angry"`
	Shares   int `bson:"shares" json:"shares"`
	Saves    int `bson:"saves" json:"saves"`
	Comments int `bson:"comments" json:"comments"`
}

// --- MAIN RESPONSE DTO ---

type ReelRes struct {
	ID               string                       `json:"id"`
	UserID           string                       `json:"user_id"`
	ProcessingStatus sharedEnums.ProcessingStatus `json:"processing_status"`
	Video            ReelVideoRes                 `json:"video"`
	Caption          string                       `json:"caption"`
	Hashtags         []string                     `json:"hashtags,omitempty"`
	Mentions         []string                     `json:"mentions,omitempty"`
	AudioMeta        AudioMetaRes                 `json:"audio_meta"`
	RemixInfo        *RemixInfoRes                `json:"remix_info,omitempty"`
	Stats            ReelStatsRes                 `json:"stats"`
	Privacy          sharedEnums.PrivacyScope     `json:"privacy"`
	CreatedAt        time.Time                    `json:"created_at"`
	DeletedAt        *time.Time                   `json:"deleted_at,omitempty"`
}
