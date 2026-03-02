package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
)

// PageReq ánh xạ 100% các trường từ Entity
type PageReq struct {
	ID            string               `json:"id,omitempty"` // Thường để rỗng khi Create, dùng khi Update
	CreatorUserID string               `json:"creator_user_id" validate:"required"`
	Name          string               `json:"name" validate:"required"`
	Slug          string               `json:"slug" validate:"required"`
	CategoryID    string               `json:"category_id"`
	IsVerified    bool                 `json:"is_verified"`
	Status        enum.PageStatus      `json:"status"`
	Avatar        PageAvatarReq        `json:"avatar"`
	Cover         PageCoverReq         `json:"cover"`
	Bio           string               `json:"bio"`
	Website       string               `json:"website"`
	Email         string               `json:"email" validate:"omitempty,email"`
	PhoneNumber   string               `json:"phone_number"`
	Address       PageAddressReq       `json:"address"`
	BusinessHours []BusinessHourReq    `json:"business_hours,omitempty"`
	CTAButton     *CTAButtonReq        `json:"cta_button,omitempty"`
	Settings      PageSettingsReq      `json:"settings"`
	Stats         PageStatsReq         `json:"stats"` // Best practice: stats nên read-only, nhưng đưa vào để đủ 100% mapping
	CreatedAt     time.Time            `json:"created_at,omitempty"`
	UpdatedAt     time.Time            `json:"updated_at,omitempty"`
}

type PageAvatarReq struct {
	URL string `json:"url"`
}

type PageCoverReq struct {
	URL       string  `json:"url"`
	PositionY float64 `json:"position_y"`
}

type PageAddressReq struct {
	Street      string    `json:"street"`
	City        string    `json:"city"`
	Zipcode     string    `json:"zipcode"`
	Coordinates []float64 `json:"coordinates,omitempty"` // [Longitude, Latitude]
}

type BusinessHourReq struct {
	Day   enum.DayOfWeek `json:"day"`
	Open  string         `json:"open"`
	Close string         `json:"close"`
}

type CTAButtonReq struct {
	Type  enum.CTAType `json:"type"`
	Value string       `json:"value"`
}

type PageSettingsReq struct {
	AllowVisitorPost bool                 `json:"allow_visitor_post"`
	ProfanityFilter  sharedEnums.UserBadge `json:"profanity_filter"`
	MessagingStatus  enum.MessagingStatus `json:"messaging_status"`
}

type PageStatsReq struct {
	FollowersCount int     `json:"followers_count"`
	LikesCount     int     `json:"likes_count"`
	RatingScore    float64 `json:"rating_score"`
	ReviewCount    int     `json:"review_count"`
}