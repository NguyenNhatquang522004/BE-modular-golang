package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- CREATE REQ ---
type CreateNotificationTemplateReq struct {
	Type       sharedEnums.NotificationType `json:"type" binding:"required"`
	Template   map[string]string            `json:"template" binding:"required"`
	IconURL    string                       `json:"icon_url"`
	ActionLink string                       `json:"action_link"`
	CreatedAt  time.Time                    `json:"created_at"` // Ánh xạ 100% theo entity
	UpdatedAt  time.Time                    `json:"updated_at"` // Ánh xạ 100% theo entity
}

// --- UPDATE REQ ---
type UpdateNotificationTemplateReq struct {
	Type       *sharedEnums.NotificationType `json:"type,omitempty"`
	Template   map[string]string             `json:"template,omitempty"`
	IconURL    *string                       `json:"icon_url,omitempty"`
	ActionLink *string                       `json:"action_link,omitempty"`
	UpdatedAt  *time.Time                    `json:"updated_at,omitempty"`
}
