package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

// CreatePostSettingReq dùng khi tạo mới settings cho một bài post
type CreatePostSettingReq struct {
	PostID string `json:"post_id" binding:"required,mongoId"` // Validate format ObjectID

	// Các nhóm field là optional (pointer) để client có thể gửi hoặc không
	Schedule  *PostScheduleReq  `json:"schedule,omitempty"`
	AdsInfo   *AdsInfoReq       `json:"ads_info,omitempty"`
	Targeting *PostTargetingReq `json:"targeting,omitempty"`
}

// UpdatePostSettingReq dùng khi cập nhật (PUT/PATCH)
type UpdatePostSettingReq struct {
	// Schedule: Gửi object để update, gửi nil thì giữ nguyên (hoặc logic xóa tùy nghiệp vụ)
	Schedule  *PostScheduleReq  `json:"schedule,omitempty"`
	AdsInfo   *AdsInfoReq       `json:"ads_info,omitempty"`
	Targeting *PostTargetingReq `json:"targeting,omitempty"`
}

// --- Nested Structs cho Request ---

type PostScheduleReq struct {
	IsScheduled        bool               `json:"is_scheduled"`
	PublishTime        time.Time          `json:"publish_time" binding:"required_if=IsScheduled true"`               // Bắt buộc nếu IsScheduled = true
	PublisherUserID    string             `json:"publisher_user_id" binding:"required,uuid"`                         // Validate UUID Postgres
	AuthorRoleSnapshot enum.PublisherRole `json:"author_role_snapshot" binding:"required,oneof=admin editor author"` // Validate enum value
}

type AdsInfoReq struct {
	CampaignID  string `json:"campaign_id" binding:"required,mongoId"`
	IsSponsored bool   `json:"is_sponsored"`
	CTALink     string `json:"cta_link" binding:"omitempty,url"` // Validate URL format
}

type PostTargetingReq struct {
	Locations []string `json:"locations,omitempty"`
	AgeMin    int      `json:"age_min" binding:"omitempty,min=0,ltefield=AgeMax"`                  // Min <= Max
	AgeMax    int      `json:"age_max" binding:"omitempty,min=0,gtefield=AgeMin"`                  // Max >= Min
	Genders   []string `json:"genders,omitempty" binding:"omitempty,dive,oneof=male female other"` // Validate từng phần tử mảng
	Languages []string `json:"languages,omitempty" binding:"omitempty,dive,len=2"`                 // VD: "vi", "en"
	Interests []string `json:"interests,omitempty"`
}
