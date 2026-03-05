package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/notificationEvent"
)

type UseSettingNotificationRequest struct {
	*notificationEvent.NotificationChangePayload
	EventType constants.EventType `json:"event_type" binding:"required"`
}
type CreateDeleteNotificationTemplateRequest struct {
	Type       int    `json:"type" binding:"required"`
	TemplateID string `json:"template_id" binding:"required,uuid"`
	*CreateNotificationTemplateReq
}
type UpdateNotificationTemplateRequest struct {
	TemplateID string `json:"template_id" binding:"required,uuid"`
	*UpdateNotificationTemplateReq
}
type SendOTPRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	OTP    string `json:"otp" binding:"required"`
}
type SendPublishPostNotificationRequest struct {
	PostID string `json:"post_id" binding:"required,uuid"`
}
type SendSystemAlertRequest struct {
	AlertID string `json:"alert_id" binding:"required,uuid"`
	Preview string `json:"preview" binding:"required"`
}
