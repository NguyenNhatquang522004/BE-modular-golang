package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionPageFollowers = "page_followers"
)

type PageFollower struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. LIÊN KẾT (RELATIONSHIPS)
	// Page nào? (Mongo ID)
	PageID primitive.ObjectID `bson:"page_id" json:"page_id"`

	// User nào? (Postgres UUID -> String)
	UserID string `bson:"user_id" json:"user_id"`

	// 2. CÀI ĐẶT RIÊNG CHO FOLLOWER NÀY
	Settings  *FollowerSettings `bson:"settings,omitempty" json:"settings,omitempty"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time        `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- SETTINGS ---
type FollowerSettings struct {
	// Mức độ thông báo: all, highlight, off
	NotificationLevel enum.NotificationLevel `bson:"notification_level" json:"notification_level"`

	// "See First" (Xem trước): Ưu tiên hiển thị bài của Page này trên Newsfeed user
	IsFavorite bool `bson:"is_favorite" json:"is_favorite"`
}

func (PageFollower) CollectionName() string {
	return CollectionPageFollowers
}
