package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/enum"
	"github.com/gocql/gocql"
)
const (
	TableNotifications = "notifications"
)
// Notification đại diện cho bảng 'notifications' trong Cassandra.
// Thiết kế Write-Heavy, tối ưu cho việc load danh sách thông báo theo user.
type Notification struct {
	// =========================================================================
	// 1. PRIMARY KEY
	// =========================================================================

	// PARTITION KEY: Gom tất cả noti của 1 user vào 1 node.
	UserID gocql.UUID `cql:"user_id" json:"user_id"`

	// CLUSTERING KEY 1: Sắp xếp thời gian (DESC -> Mới nhất lên đầu).
	CreatedAt time.Time `cql:"created_at" json:"created_at"`

	// CLUSTERING KEY 2: Đảm bảo tính duy nhất (TimeUUID).
	NotificationID gocql.UUID `cql:"notification_id" json:"notification_id"`

	// =========================================================================
	// 2. CORE INFO
	// =========================================================================

	// Enum: post_like, friend_request...
	Type enum.NotificationType `cql:"type" json:"type"`

	// =========================================================================
	// 3. ACTOR (Người gây ra hành động - Denormalization)
	// =========================================================================
	// Lưu cứng thông tin tại thời điểm tạo noti.
	// Nếu A đổi avatar sau này, noti cũ vẫn hiện avatar cũ (Chấp nhận được).

	ActorID     gocql.UUID `cql:"actor_id" json:"actor_id"`
	ActorName   string     `cql:"actor_name" json:"actor_name"`
	ActorAvatar string     `cql:"actor_avatar" json:"actor_avatar"`

	// =========================================================================
	// 4. TARGET (Đối tượng bị tác động)
	// =========================================================================

	// ID của đối tượng (PostID, CommentID...).
	// Lưu String để hỗ trợ cả UUID và MongoDB ObjectId.
	TargetID string `cql:"target_id" json:"target_id"`

	// Snapshot nội dung ngắn để hiển thị ngay mà không cần query lại Post/Comment.
	// VD: "Đã bình luận: 'Bài viết hay quá...'"
	TargetPreview string `cql:"target_preview" json:"target_preview"`

	// =========================================================================
	// 5. STATUS
	// =========================================================================

	IsRead    bool `cql:"is_read" json:"is_read"`       // Đã xem chưa? (Dấu chấm đỏ)
	IsClicked bool `cql:"is_clicked" json:"is_clicked"` // Đã click vào chưa?

	// =========================================================================
	// 6. GROUPING LOGIC
	// =========================================================================

	// Dùng để gộp các thông báo giống nhau.
	// VD Logic: Khi tạo noti Like, check xem Redis có key "POST_LIKE:post_123" chưa.
	// Nếu có -> Update count. Nếu chưa -> Insert noti mới.
	GroupKey string `cql:"group_key" json:"group_key"`
}

// TableName trả về tên bảng
func (Notification) TableName() string {
	return TableNotifications
}