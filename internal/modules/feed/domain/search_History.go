package domain

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchHistory struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Người thực hiện tìm kiếm (Postgres UUID -> String)
	// Index: { user_id: 1, created_at: -1 } -> Lấy lịch sử tìm kiếm gần đây
	UserID string `bson:"user_id" json:"user_id"`

	// Từ khóa tìm kiếm
	Keyword string `bson:"keyword" json:"keyword"`

	// --- TARGET INFO (Optional) ---
	// Chỉ lưu nếu user click vào một kết quả cụ thể.
	// Dùng Pointer + omitempty để không lưu field này nếu null -> Tiết kiệm ổ cứng
	TargetID   *string                `bson:"target_id,omitempty" json:"target_id,omitempty"`
	TargetType *enum.SearchTargetType `bson:"target_type,omitempty" json:"target_type,omitempty"`

	// --- TIMESTAMPS & TTL ---
	// Index TTL: { created_at: 1 }, expireAfterSeconds: 2592000 (30 ngày)
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
