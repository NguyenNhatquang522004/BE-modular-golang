package entity

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
const (
	CollectionPageFollowers = "page_followers"
)
// --- SETTINGS ---
type FollowerSettings struct {
	// Mức độ thông báo: all, highlight, off
	NotificationLevel enum.NotificationLevel `bson:"notification_level" json:"notification_level"`

	// "See First" (Xem trước): Ưu tiên hiển thị bài của Page này trên Newsfeed user
	IsFavorite bool `bson:"is_favorite" json:"is_favorite"`
}
func (FollowerSettings) CollectionName() string {
	return CollectionPageFollowers
}
