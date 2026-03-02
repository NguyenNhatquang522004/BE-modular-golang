package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/google/uuid"
)

// AdReq đại diện cho payload dữ liệu mà client gửi lên (Create/Update).
// Đã ánh xạ 100% các trường từ Entity theo yêu cầu.
type AdReq struct {
	AdID            uuid.UUID                     `json:"ad_id"`
	CampaignID      uuid.UUID                     `json:"campaign_id" validate:"required"`
	TargetPostID    string                        `json:"target_post_id" validate:"required"`
	BidAmount       float64                       `json:"bid_amount" validate:"required,gt=0"`
	Status          *sharedEnums.ProcessingStatus `json:"status"`
	RejectionReason *string                       `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`
	DeletedAt       *time.Time                    `json:"deleted_at,omitempty"`
}
