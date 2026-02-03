package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
	"github.com/gocql/gocql"
)

const (
	CassandratableUserReactionHistory = "UserReactionHistory"
)

// UserReactionHistory đại diện cho bảng 'user_reaction_history' trong Cassandra.
// Bảng này phục vụ tính năng "Nhật ký hoạt động" (Activity Log).
// Query pattern: Select * from user_reaction_history WHERE user_id = ?
type UserReactionHistory struct {
	// 1. PARTITION KEY
	// Gom tất cả lịch sử của 1 user vào cùng 1 partition node -> Query cực nhanh.
	UserID gocql.UUID `cql:"user_id" json:"user_id"`

	// 2. CLUSTERING KEY
	// Sắp xếp giảm dần theo thời gian (Mới nhất lên đầu).
	CreatedAt time.Time `cql:"created_at" json:"created_at"`

	// 3. DATA FIELDS
	// ID của bài viết/comment mà user đã like (Lưu String vì gốc là Mongo ObjectId)
	TargetID string `cql:"target_id" json:"target_id"`

	// Enum tái sử dụng từ package enum (Post/Comment)
	TargetType enum.ReactionTarget `cql:"target_type" json:"target_type"`

	// Enum tái sử dụng từ package enum (Like/Love/Haha...)
	ReactionCode enum.ReactionCode `cql:"reaction_code" json:"reaction_code"`
}

// TableName trả về tên bảng trong Cassandra
func (UserReactionHistory) CassandratableUserReactionHistory() string {
	return CassandratableUserReactionHistory
}
