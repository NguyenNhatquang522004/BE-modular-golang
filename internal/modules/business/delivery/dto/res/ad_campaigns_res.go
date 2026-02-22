package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"github.com/google/uuid"
)

// AdCampaignRes dùng để trả dữ liệu về cho client.
type AdCampaignRes struct {
	CampaignID     uuid.UUID              `json:"campaign_id"`
	AccountID      uuid.UUID              `json:"account_id"`
	Name           string                 `json:"name"`
	Objective      *enum.CampaignObjective `json:"objective"`
	BuyingType     *enum.BuyingType        `json:"buying_type"`
	DailyBudget    *float64               `json:"daily_budget,omitempty"`
	LifetimeBudget *float64               `json:"lifetime_budget,omitempty"`
	StartTime      time.Time              `json:"start_time"`
	EndTime        *time.Time             `json:"end_time,omitempty"`
	Status         *enum.CampaignStatus    `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
}