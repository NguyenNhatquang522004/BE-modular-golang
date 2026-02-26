package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- SUB-STRUCTS ---

type RecordingSettingRes struct {
	IsRecorded bool   `json:"is_recorded"`
	ArchiveURL string `json:"archive_url,omitempty"`
}

type LiveStatsRes struct {
	PeakViewers   int `json:"peak_viewers"`
	TotalViews    int `json:"total_views"`
	TotalLikes    int `json:"total_likes"`
	TotalComments int `json:"total_comments"`
}

// --- MAIN RESPONSE DTO ---

// LiveSessionRes: Dữ liệu trả về cho client, 100% thuộc tính.
// Lưu ý: Thực tế StreamKey rất nhạy cảm, nhưng theo yêu cầu map 100% nên vẫn giữ lại.
type LiveSessionRes struct {
	ID               string                       `json:"id"`
	HostUserID       string                       `json:"host_user_id"`
	Title            string                       `json:"title"`
	Description      string                       `json:"description"`
	CategoryID       string                       `json:"category_id"`
	Status           sharedEnums.ProcessingStatus `json:"status"`
	StreamKey        string                       `json:"stream_key,omitempty"`
	PlaybackURL      string                       `json:"playback_url"`
	RecordingSetting RecordingSettingRes          `json:"recording_setting"`
	BannedUsers      []string                     `json:"banned_users,omitempty"`
	PinnedCommentID  string                       `json:"pinned_comment_id,omitempty"`
	StartedAt        *time.Time                   `json:"started_at,omitempty"`
	EndedAt          *time.Time                   `json:"ended_at,omitempty"`
	Stats            LiveStatsRes                 `json:"stats"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
}
