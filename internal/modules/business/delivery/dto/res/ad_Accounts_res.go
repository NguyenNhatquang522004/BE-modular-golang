package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"github.com/google/uuid"
)

// AdAccountRes dùng để trả dữ liệu về cho client.
type AdAccountRes struct {
	AccountID   uuid.UUID          `json:"account_id"`
	OwnerUserID uuid.UUID          `json:"owner_user_id"`
	Currency    string             `json:"currency"`
	Timezone    string             `json:"timezone"`
	Balance     float64            `json:"balance"`
	CreditLimit float64            `json:"credit_limit"`
	Status      *enum.AccountStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at,omitempty"`
}