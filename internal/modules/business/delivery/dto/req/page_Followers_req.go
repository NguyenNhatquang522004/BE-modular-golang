package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"

// PageFollowerReq chứa dữ liệu request từ client.
type PageFollowerReq struct {
	// ID có thể truyền hoặc không (thường tạo mới sẽ không có)
	ID       string               `json:"id,omitempty" validate:"omitempty"`
	PageID   string               `json:"page_id" validate:"required"`
	UserID   string               `json:"user_id" validate:"required,uuid"` // Đảm bảo là UUID theo chuẩn Postgres
	Settings *FollowerSettingsReq `json:"settings,omitempty"`
}

// FollowerSettingsReq chứa cấu hình thông báo từ request.
type FollowerSettingsReq struct {
	NotificationLevel sharedEnums.NotificationLevel `json:"notification_level" validate:"required"`
	IsFavorite        bool                          `json:"is_favorite"`
}
