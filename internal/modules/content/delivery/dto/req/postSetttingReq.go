package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// --- NESTED DTOs ---
type PostScheduleReq struct {
	IsScheduled        bool                      `json:"is_scheduled"`
	PublishTime        time.Time                 `json:"publish_time" validate:"required"`
	PublisherUserID    string                    `json:"publisher_user_id" validate:"required"`
	AuthorRoleSnapshot sharedEnums.PublisherRole `json:"author_role_snapshot" validate:"required"`
}

type AdsInfoReq struct {
	CampaignID  string `json:"campaign_id" validate:"required"` // String thay vì ObjectID
	IsSponsored bool   `json:"is_sponsored"`
	CTALink     string `json:"cta_link,omitempty"`
}

type PostTargetingReq struct {
	Locations []string `json:"locations,omitempty"`
	AgeMin    int      `json:"age_min"`
	AgeMax    int      `json:"age_max"`
	Genders   []string `json:"genders,omitempty"`
	Languages []string `json:"languages,omitempty"`
	Interests []string `json:"interests,omitempty"`
}

// --- MAIN DTO ---
type PostSettingReq struct {
	ID        string            `json:"id,omitempty"` // String thay vì ObjectID cho HTTP Request
	PostID    string            `json:"post_id" validate:"required"`
	Schedule  *PostScheduleReq  `json:"schedule,omitempty"`
	AdsInfo   *AdsInfoReq       `json:"ads_info,omitempty"`
	Targeting *PostTargetingReq `json:"targeting,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt *time.Time        `json:"deleted_at,omitempty"`
}
