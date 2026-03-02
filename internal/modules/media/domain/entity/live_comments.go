package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)
const (
    TableLiveComments = "live_comments"
)
// LiveComment đại diện cho bảng 'live_comments' trong Cassandra.
// Tối ưu cho High Write Throughput (Chat nhảy liên tục).
type LiveComment struct {
	// 1. PARTITION KEY
	// Gom tất cả comment của 1 buổi stream vào 1 node.
	StreamID gocql.UUID `cql:"stream_id" json:"stream_id"`

	// 2. CLUSTERING KEY 1
	// Sắp xếp giảm dần (DESC) -> Comment mới nhất hiện lên đầu.
	CreatedAt time.Time `cql:"created_at" json:"created_at"`

	// 3. CLUSTERING KEY 2
	// Đảm bảo tính duy nhất nếu có 2 comment cùng mili-giây.
	CommentID gocql.UUID `cql:"comment_id" json:"comment_id"`

	// 4. USER INFO (SNAPSHOT)
	// Lưu cứng thông tin user lúc comment. User đổi avatar sau này thì comment cũ vẫn giữ avatar cũ.
	UserID        gocql.UUID `cql:"user_id" json:"user_id"`
	UserNickname  string     `cql:"user_nickname" json:"user_nickname"`
	UserAvatarURL string     `cql:"user_avatar_url" json:"user_avatar_url"`

	// Badges: Cassandra SET<TEXT> map về []string trong Go.
	// Logic: Service sẽ convert []enum.UserBadge -> []string trước khi gán vào đây.
	UserBadges []*sharedEnums.UserBadge `cql:"user_badges" json:"user_badges"`

	// 5. CONTENT & STATUS
	Content  string `cql:"content" json:"content"`
	IsPinned bool   `cql:"is_pinned" json:"is_pinned"`
}

// TableName trả về tên bảng
func (LiveComment) TableName() string {
    return TableLiveComments
}