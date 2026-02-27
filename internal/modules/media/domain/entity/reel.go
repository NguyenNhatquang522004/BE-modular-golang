package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionReels = "Reels"
)

type Reel struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. OWNERSHIP
	// UserID từ Postgres (UUID) -> Lưu String
	// Index: { user_id: 1, created_at: -1 } -> Trang cá nhân
	UserID string `bson:"user_id" json:"user_id"`

	// 2. PROCESSING FLOW
	// Worker sẽ update field này. Client chỉ load các reel có status = 'active'
	ProcessingStatus sharedEnums.ProcessingStatus `bson:"processing_status" json:"processing_status"`

	// 3. CONTENT
	Video   ReelVideo `bson:"video" json:"video"`
	Caption string    `bson:"caption" json:"caption"`

	// Index: { hashtags: 1 } -> Tìm kiếm trend
	Hashtags []string `bson:"hashtags,omitempty" json:"hashtags,omitempty"`

	// Danh sách UserID (UUID String) được tag
	Mentions []string `bson:"mentions,omitempty" json:"mentions,omitempty"`

	// 4. FEATURES
	AudioMeta AudioMeta `bson:"audio_meta" json:"audio_meta"`

	// Pointer để tiết kiệm nếu không phải remix
	RemixInfo *RemixInfo `bson:"remix_info,omitempty" json:"remix_info,omitempty"`

	// 5. METADATA
	Stats   ReelStats                `bson:"stats" json:"stats"`
	Privacy sharedEnums.PrivacyScope `bson:"privacy" json:"privacy"`

	// 6. TIMESTAMPS
	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	// Soft Delete: Nếu field này != null thì coi như đã xóa
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- VIDEO METADATA ---
type ReelVideo struct {
	URL           string  `bson:"url" json:"url"` // HLS/DASH Stream URL
	ThumbnailURL  string  `bson:"thumbnail_url" json:"thumbnail_url"`
	PreviewGifURL string  `bson:"preview_gif_url" json:"preview_gif_url"` // Hover preview
	Width         int     `bson:"width" json:"width"`
	Height        int     `bson:"height" json:"height"`
	Duration      float64 `bson:"duration" json:"duration"` // Seconds
}

// --- AUDIO INFO ---
type AudioMeta struct {
	// ID của bài hát trong Music Library (Mongo).
	// Nullable nếu user không chọn nhạc nền.
	// Index: { "audio_meta.track_id": 1 } -> Để làm trang "Các video dùng âm thanh này"
	TrackID *primitive.ObjectID `bson:"track_id,omitempty" json:"track_id,omitempty"`

	IsOriginalAudio bool    `bson:"is_original_audio" json:"is_original_audio"`
	VolumeAdjust    float64 `bson:"volumn_adjust" json:"volumn_adjust"`       // 0.0 -> 1.0 (Giữ nguyên tên field volumn theo đề bài dù tiếng Anh chuẩn là volume)
	AudioStartTime  float64 `bson:"audio_start_time" json:"audio_start_time"` // Bắt đầu từ giây thứ mấy của bài nhạc
}

// --- REMIX INFO ---
type RemixInfo struct {
	// Video gốc (Mongo ID)
	ParentReelID primitive.ObjectID `bson:"parent_reel_id" json:"parent_reel_id"`
	Type         enum.RemixType     `bson:"type" json:"type"`
	IsRemixable  bool               `bson:"is_remixable" json:"is_remixable"`
}

// --- STATS ---
type ReelStats struct {
	Views    int `bson:"views" json:"views"`
	Total    int `bson:"total" json:"total"`
	Like     int `bson:"like" json:"like"`
	Love     int `bson:"love" json:"love"`
	Haha     int `bson:"haha" json:"haha"`
	Wow      int `bson:"wow" json:"wow"`
	Sad      int `bson:"sad" json:"sad"`
	Angry    int `bson:"angry" json:"angry"` 
	Shares   int `bson:"shares" json:"shares"`     //chưa làm
	Saves    int `bson:"saves" json:"saves"`       //chưa làm
	Comments int `bson:"comments" json:"comments"` //chưa làm
}

func (Reel) CollectionName() string {
	return CollectionReels
}
