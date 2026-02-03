package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =============================================================================
// ROOT ENTITY: EDIT LOG
// =============================================================================
const (
	collectionnameposteditlog = "PostEntityEditLog"
)

type PostEntityEditLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. Target (Mục tiêu chỉnh sửa)
	// Enum xác định sửa cái gì: 'posts' hay 'comments'
	TargetCollection enum.TargetCollection `bson:"target_collection" json:"target_collection"`

	// ID của bài viết hoặc comment bị sửa (ObjectId vì nằm trong Mongo)
	// Index: Compound { target_id: 1, target_collection: 1, version: -1 }
	TargetID primitive.ObjectID `bson:"target_id" json:"target_id"`

	// 2. Meta Info
	Version  int       `bson:"version" json:"version"` // 1, 2, 3...
	EditedAt time.Time `bson:"edited_at" json:"edited_at"`

	// EditorID là UserID từ Postgres (UUID) -> Lưu String
	EditorID string `bson:"editor_id" json:"editor_id"`

	// 3. Nội dung thay đổi
	Diff LogDiff `bson:"diff" json:"diff"`

	// 4. Audit Info (Dùng cho bảo mật/tracking)
	IPAddress string `bson:"ip_address" json:"ip_address"`
	UserAgent string `bson:"user_agent" json:"user_agent"`
}

// =============================================================================
// SUB-STRUCT: DIFF
// =============================================================================

type LogDiff struct {
	// Nội dung cũ (Trước khi sửa)
	OldContent string `bson:"old_content" json:"old_content"`

	// Nội dung mới (Sau khi sửa)
	NewContent string `bson:"new_content" json:"new_content"`

	// Danh sách các field bị thay đổi. VD: ["content", "privacy"]
	ChangedFields []string `bson:"changed_fields" json:"changed_fields"`
}

func (PostEntityEditLog) Collectionnameposteditlog() string {
	return collectionnameposteditlog
}
