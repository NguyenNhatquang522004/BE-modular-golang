package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =============================================================================
// ROOT ENTITY: EDIT LOG
// =============================================================================
const (
	collectionnamCommentEditLog = "CommentEntityEditLog"
)

type CommentEntityEditLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. ĐỊNH DANH & MAPPING
	// Enum: 'comments' hoặc 'posts'
	TargetCollection enum.TargetCollection `bson:"target_collection" json:"target_collection"`

	// ID của Comment/Post gốc (ObjectId vì nằm trong Mongo)
	// Index: Compound { target_id: 1, target_collection: 1, version: -1 }
	TargetID primitive.ObjectID `bson:"target_id" json:"target_id"`

	// 2. META DATA
	Version  int       `bson:"version" json:"version"` // 1, 2, 3...
	EditedAt time.Time `bson:"edited_at" json:"edited_at"`

	// EditorID là UserID từ Postgres (UUID) -> Bắt buộc lưu String
	EditorID string `bson:"editor_id" json:"editor_id"`

	// 3. NỘI DUNG THAY ĐỔI (DIFF)
	Diff LogDiff `bson:"diff" json:"diff"`

	// 4. AUDIT INFO
	IPAddress string `bson:"ip_address" json:"ip_address"`
	UserAgent string `bson:"user_agent" json:"user_agent"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// =============================================================================
// SUB-STRUCTS
// =============================================================================

// LogDiff chứa thông tin so sánh cũ/mới
type LogDiff struct {
	// Nội dung text
	OldContent string `bson:"old_content" json:"old_content"`
	NewContent string `bson:"new_content" json:"new_content"`

	// Media (Ảnh/Video)
	// Dùng Pointer để nếu không sửa ảnh thì field này null (omitempty) -> Tiết kiệm ổ cứng
	OldMedia *MediaSnapshot `bson:"old_media,omitempty" json:"old_media,omitempty"`
	NewMedia *MediaSnapshot `bson:"new_media,omitempty" json:"new_media,omitempty"`
}

// MediaSnapshot lưu trạng thái file media tại thời điểm sửa
type MediaSnapshot struct {
	Type string `bson:"type" json:"type"` // "image", "video", "gif"
	URL  string `bson:"url" json:"url"`
}

func (CommentEntityEditLog) CollectionnamCommentEditLog() string {
	return collectionnamCommentEditLog
}
