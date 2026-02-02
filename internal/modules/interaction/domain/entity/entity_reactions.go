package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
	"github.com/gocql/gocql"
)

const (
	CassandratableEntityReaction = "EntityReaction"
)

// EntityReaction đại diện cho bảng 'entity_reactions' trong Cassandra
// Dùng để lưu trữ tương tác (High Volume Write)
type EntityReaction struct {
	// 1. PRIMARY KEY COMPONENT 1: PARTITION KEY
	// ID của Post hoặc Comment (Lấy từ MongoDB -> String)
	TargetID string `cql:"target_id" json:"target_id"`

	// 2. PRIMARY KEY COMPONENT 2: CLUSTERING KEY
	// ID của User (Lấy từ Postgres -> UUID)
	UserID gocql.UUID `cql:"user_id" json:"user_id"`

	// 3. DATA FIELDS
	// Enum xác định loại target ('post' hay 'comment')
	TargetType enum.ReactionTarget `cql:"target_type" json:"target_type"`

	// Enum xác định loại cảm xúc ('like', 'love',...)
	ReactionCode enum.ReactionCode `cql:"reaction_code" json:"reaction_code"`

	CreatedAt time.Time `cql:"created_at" json:"created_at"`
}

// TableName trả về tên bảng trong Cassandra
func (EntityReaction) CassandratableEntityReaction() string {
	return CassandratableEntityReaction
}
