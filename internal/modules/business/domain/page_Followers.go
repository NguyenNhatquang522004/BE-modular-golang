package domain

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"

// --- SETTINGS ---
type FollowerSettings struct {
	// Mức độ thông báo: all, highlight, off
	NotificationLevel enum.NotificationLevel `bson:"notification_level" json:"notification_level"`

	// "See First" (Xem trước): Ưu tiên hiển thị bài của Page này trên Newsfeed user
	IsFavorite bool `bson:"is_favorite" json:"is_favorite"`
}
