package res

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"

// PageFollowerRes chứa dữ liệu trả về cho client.
type PageFollowerRes struct {
	ID       string               `json:"id"`
	PageID   string               `json:"page_id"`
	UserID   string               `json:"user_id"`
	Settings *FollowerSettingsRes `json:"settings,omitempty"`
}

// FollowerSettingsRes chứa cấu hình thông báo trả về.
type FollowerSettingsRes struct {
	NotificationLevel sharedEnums.NotificationLevel `json:"notification_level"`
	IsFavorite        bool                          `json:"is_favorite"`
}
