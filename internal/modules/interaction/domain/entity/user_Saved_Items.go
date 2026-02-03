package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =============================================================================
// ROOT ENTITY: USER SAVED ITEM
// =============================================================================
const (
	collectionnamUserSavedItem = "UserSavedItem"
)

type UserSavedItem struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. CHỦ SỞ HỮU
	// UserID từ Postgres (UUID) -> Lưu String
	// Index: Compound { user_id: 1, collection_name: 1 } (Để lọc theo bộ sưu tập)
	UserID string `bson:"user_id" json:"user_id"`

	// 2. ĐỐI TƯỢNG ĐƯỢC LƯU
	// ID của Post/Reel/Video (ObjectId vì nằm trong Mongo)
	// Index: Compound Unique { user_id: 1, target_id: 1, target_type: 1 }
	// -> Mục đích: Ngăn chặn user lưu trùng 1 bài 2 lần.
	TargetID   primitive.ObjectID   `bson:"target_id" json:"target_id"`
	TargetType enum.SavedTargetType `bson:"target_type" json:"target_type"`

	// 3. HIỂN THỊ NHANH (Denormalization)
	// Lưu snapshot để khi list ra không cần query ngược lại bảng Post
	Snapshot SavedItemSnapshot `bson:"snapshot" json:"snapshot"`

	// 4. PHÂN LOẠI
	// Tên bộ sưu tập. VD: "Món ăn ngon", "Tech news".
	// Nếu rỗng -> Coi như là "Mặc định"
	CollectionName string `bson:"collection_name" json:"collection_name"`

	// 5. META
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// =============================================================================
// SUB-STRUCT: SNAPSHOT
// =============================================================================

type SavedItemSnapshot struct {
	AuthorName     string `bson:"author_name" json:"author_name"`
	ContentPreview string `bson:"content_preview" json:"content_preview"` // Cắt ngắn 100 ký tự đầu
	ThumbnailURL   string `bson:"thumbnail_url" json:"thumbnail_url"`
}

func (UserSavedItem) CollectionnamUserSavedItem() string {
	return collectionnamUserSavedItem
}
