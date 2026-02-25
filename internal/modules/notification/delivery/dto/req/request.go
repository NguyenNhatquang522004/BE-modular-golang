package req

type CreateUserNotificationSettingsRequest struct {
	*CreateUserNotificationSettingReq
}
type UpdateUserNotificationSettingsRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	*UpdateUserNotificationSettingReq
}
type DeleteUserNotificationSettingsRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
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
