package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
)

// --- SUB-STRUCTS ---

type RecordingSettingReq struct {
	IsRecorded bool   `json:"is_recorded"`
	ArchiveURL string `json:"archive_url,omitempty"`
}

type LiveStatsReq struct {
	PeakViewers   int `json:"peak_viewers"`
	TotalViews    int `json:"total_views"`
	TotalLikes    int `json:"total_likes"`
	TotalComments int `json:"total_comments"`
}

// --- MAIN REQUEST DTO ---

// LiveSessionReq: Ánh xạ đủ 100% các trường của Entity.
type LiveSessionReq struct {
	ID               string              `json:"id,omitempty"` // String để client gửi Hex ID
	HostUserID       string              `json:"host_user_id" validate:"required"`
	Title            string              `json:"title" validate:"required"`
	Description      string              `json:"description"`
	CategoryID       string              `json:"category_id"`
	Status           enum.LiveStatus     `json:"status"`
	StreamKey        string              `json:"stream_key"`
	PlaybackURL      string              `json:"playback_url"`
	RecordingSetting RecordingSettingReq `json:"recording_setting"`
	BannedUsers      []string            `json:"banned_users,omitempty"`
	PinnedCommentID  string              `json:"pinned_comment_id,omitempty"`
	StartedAt        *time.Time          `json:"started_at,omitempty"`
	EndedAt          *time.Time          `json:"ended_at,omitempty"`
	Stats            LiveStatsReq        `json:"stats"`
	CreatedAt        *time.Time          `json:"created_at,omitempty"`
	UpdatedAt        *time.Time          `json:"updated_at,omitempty"`
}

// UpdateLiveSessionReq: Dùng con trỏ (pointer) 100% cho Partial Update
type UpdateLiveSessionReq struct {
	HostUserID       *string              `json:"host_user_id,omitempty"`
	Title            *string              `json:"title,omitempty"`
	Description      *string              `json:"description,omitempty"`
	CategoryID       *string              `json:"category_id,omitempty"`
	Status           *enum.LiveStatus     `json:"status,omitempty"`
	StreamKey        *string              `json:"stream_key,omitempty"`
	PlaybackURL      *string              `json:"playback_url,omitempty"`
	RecordingSetting *RecordingSettingReq `json:"recording_setting,omitempty"`
	BannedUsers      []string             `json:"banned_users,omitempty"`
	PinnedCommentID  *string              `json:"pinned_comment_id,omitempty"`
	StartedAt        *time.Time           `json:"started_at,omitempty"`
	EndedAt          *time.Time           `json:"ended_at,omitempty"`
	Stats            *LiveStatsReq        `json:"stats,omitempty"`
}
