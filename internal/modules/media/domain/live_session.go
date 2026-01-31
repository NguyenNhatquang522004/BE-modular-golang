package domain

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LiveSession struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. OWNERSHIP
	// UserID từ Postgres (UUID) -> Lưu String
	// Index: { host_user_id: 1, status: 1 } -> Kiểm tra xem user này có đang live không
	HostUserID string `bson:"host_user_id" json:"host_user_id"`

	// 2. INFO
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`

	// ID danh mục (Game, Music...). Index: { category_id: 1, status: 1 }
	CategoryID string `bson:"category_id" json:"category_id"`

	Status enum.LiveStatus `bson:"status" json:"status"`

	// 3. TECHNICAL INFO
	// StreamKey: Tuyệt đối không expose ra public API list
	StreamKey        string           `bson:"stream_key" json:"stream_key"`
	PlaybackURL      string           `bson:"playback_url" json:"playback_url"` // HLS/FLV URL
	RecordingSetting RecordingSetting `bson:"recording_setting" json:"recording_setting"`

	// 4. ROOM MANAGEMENT
	// Danh sách UUID user bị chặn (Blacklist của phòng này)
	BannedUsers []string `bson:"banned_users,omitempty" json:"banned_users,omitempty"`

	// ID của Comment đang được ghim (Lưu string cho linh hoạt ID)
	PinnedCommentID string `bson:"pinned_comment_id,omitempty" json:"pinned_comment_id,omitempty"`

	// 5. TIMESTAMPS
	// Dùng pointer vì lúc tạo phòng (Pending) thì chưa có StartedAt/EndedAt
	StartedAt *time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`
	EndedAt   *time.Time `bson:"ended_at,omitempty" json:"ended_at,omitempty"`

	// 6. METRICS
	Stats LiveStats `bson:"stats" json:"stats"`

	// Timestamp tạo record (khác với lúc bắt đầu live)
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// --- RECORDING SETTINGS ---
type RecordingSetting struct {
	IsRecorded bool `bson:"is_recorded" json:"is_recorded"`
	// URL file video đã lưu sau khi live xong (SeaweedFS / S3)
	ArchiveURL string `bson:"archive_url,omitempty" json:"archive_url,omitempty"`
}

// --- STATS ---
// Các chỉ số này nên được update bất đồng bộ từ Redis -> Mongo
// (Ví dụ: Mỗi 30s worker sync 1 lần để giảm tải DB)
type LiveStats struct {
	PeakViewers   int `bson:"peak_viewers" json:"peak_viewers"` // Mắt xem đỉnh điểm
	TotalViews    int `bson:"total_views" json:"total_views"`   // Tổng lượt click vào xem
	TotalLikes    int `bson:"total_likes" json:"total_likes"`
	TotalComments int `bson:"total_comments" json:"total_comments"`
}
