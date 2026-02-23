package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

// --- NESTED DTOs ---
type PostScheduleRes struct {
	IsScheduled        bool               `json:"is_scheduled"`
	PublishTime        time.Time          `json:"publish_time"`
	PublisherUserID    string             `json:"publisher_user_id"`
	AuthorRoleSnapshot enum.PublisherRole `json:"author_role_snapshot"`
}

type AdsInfoRes struct {
	CampaignID  string `json:"campaign_id"` // Trả về dạng string cho client
	IsSponsored bool   `json:"is_sponsored"`
	CTALink     string `json:"cta_link,omitempty"`
}

type PostTargetingRes struct {
	Locations []string `json:"locations,omitempty"`
	AgeMin    int      `json:"age_min"`
	AgeMax    int      `json:"age_max"`
	Genders   []string `json:"genders,omitempty"`
	Languages []string `json:"languages,omitempty"`
	Interests []string `json:"interests,omitempty"`
}

// --- MAIN DTO ---
type PostSettingRes struct {
	ID        string            `json:"id"`
	PostID    string            `json:"post_id"`
	Schedule  *PostScheduleRes  `json:"schedule,omitempty"`
	AdsInfo   *AdsInfoRes       `json:"ads_info,omitempty"`
	Targeting *PostTargetingRes `json:"targeting,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt *time.Time        `json:"deleted_at,omitempty"`
}