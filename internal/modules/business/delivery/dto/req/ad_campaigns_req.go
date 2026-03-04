package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"github.com/google/uuid"
)

// AdCampaignReq chứa 100% các trường ánh xạ từ Entity AdCampaign.
type AdCampaignReq struct {
	CampaignID     uuid.UUID                     `json:"campaign_id"`
	AccountID      uuid.UUID                     `json:"account_id"`
	Name           string                        `json:"name"`
	Objective      *enum.CampaignObjective       `json:"objective"`
	BuyingType     *enum.BuyingType              `json:"buying_type"`
	DailyBudget    *float64                      `json:"daily_budget,omitempty"`
	LifetimeBudget *float64                      `json:"lifetime_budget,omitempty"`
	StartTime      time.Time                     `json:"start_time"`
	EndTime        *time.Time                    `json:"end_time,omitempty"`
	Status         *sharedEnums.ProcessingStatus `json:"status"`
	CreatedAt      time.Time                     `json:"created_at"`
	UpdatedAt      time.Time                     `json:"updated_at"`
	DeletedAt      *time.Time                    `json:"deleted_at,omitempty"`
}
