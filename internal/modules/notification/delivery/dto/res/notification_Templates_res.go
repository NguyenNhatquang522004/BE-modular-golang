package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- MAIN RESPONSE ---
type NotificationTemplateRes struct {
	ID         string                       `json:"id"`
	Type       sharedEnums.NotificationType `json:"type"`
	Template   map[string]string            `json:"template"`
	IconURL    string                       `json:"icon_url"`
	ActionLink string                       `json:"action_link"`
	CreatedAt  time.Time                    `json:"created_at"`
	UpdatedAt  time.Time                    `json:"updated_at"`
}
