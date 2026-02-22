package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"github.com/google/uuid"
)

// AdRes đại diện cho dữ liệu trả về cho client.
type AdRes struct {
	AdID            uuid.UUID     `json:"ad_id"`
	CampaignID      uuid.UUID     `json:"campaign_id"`
	TargetPostID    string        `json:"target_post_id"`
	BidAmount       float64       `json:"bid_amount"`
	Status          *enum.AdStatus `json:"status"`
	RejectionReason *string       `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       *time.Time    `json:"deleted_at,omitempty"`
}