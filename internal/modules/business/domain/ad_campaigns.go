package domain

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"github.com/google/uuid"
)

// AdCampaign đại diện cho bảng 'ad_campaigns' trong Postgres.
type AdCampaign struct {
	// 1. IDENTITY
	CampaignID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"campaign_id"`

	// 2. RELATIONS
	// Liên kết tới Ad_Accounts. Cần Index để lọc "Chiến dịch của tài khoản X"
	AccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`

	// 3. CORE INFO
	Name       string                 `gorm:"type:varchar(255);not null" json:"name"`
	Objective  enum.CampaignObjective `gorm:"type:varchar(50);not null" json:"objective"`   // reach, traffic...
	BuyingType enum.BuyingType        `gorm:"type:varchar(20);not null" json:"buying_type"` // auction...

	// 4. BUDGET (Financials)
	// Dùng Pointer (*float64) để cho phép giá trị NULL trong DB.
	// Logic: Thường chỉ set 1 trong 2.
	DailyBudget    *float64 `gorm:"type:decimal(15,2)" json:"daily_budget,omitempty"`
	LifetimeBudget *float64 `gorm:"type:decimal(15,2)" json:"lifetime_budget,omitempty"`

	// 5. SCHEDULE
	StartTime time.Time `gorm:"type:timestamp;not null" json:"start_time"`

	// EndTime có thể NULL (Chạy vô thời hạn) -> Dùng Pointer
	EndTime *time.Time `gorm:"type:timestamp" json:"end_time,omitempty"`

	// 6. STATUS
	// Index để lọc nhanh các chiến dịch đang chạy (Active)
	Status enum.CampaignStatus `gorm:"type:varchar(20);index;default:'paused'" json:"status"`
}

// TableName ghi đè tên bảng
func (AdCampaign) TableName() string {
	return "ad_campaigns"
}
