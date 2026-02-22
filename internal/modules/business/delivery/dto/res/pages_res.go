package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
)

// PageRes trả về dữ liệu cho Client (Object ID được chuyển thành string)
type PageRes struct {
	ID            string            `json:"id"`
	CreatorUserID string            `json:"creator_user_id"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	CategoryID    string            `json:"category_id"`
	IsVerified    bool              `json:"is_verified"`
	Status        enum.PageStatus   `json:"status"`
	Avatar        *PageAvatarRes    `json:"avatar"`
	Cover         *PageCoverRes     `json:"cover"`
	Bio           string            `json:"bio"`
	Website       string            `json:"website"`
	Email         string            `json:"email"`
	PhoneNumber   string            `json:"phone_number"`
	Address       *PageAddressRes   `json:"address"`
	BusinessHours []BusinessHourRes `json:"business_hours,omitempty"`
	CTAButton     *CTAButtonRes     `json:"cta_button,omitempty"`
	Settings      *PageSettingsRes  `json:"settings"`
	Stats         *PageStatsRes     `json:"stats"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type PageAvatarRes struct {
	URL string `json:"url"`
}

type PageCoverRes struct {
	URL       string  `json:"url"`
	PositionY float64 `json:"position_y"`
}

type PageAddressRes struct {
	Street      string    `json:"street"`
	City        string    `json:"city"`
	Zipcode     string    `json:"zipcode"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

type BusinessHourRes struct {
	Day   enum.DayOfWeek `json:"day"`
	Open  string         `json:"open"`
	Close string         `json:"close"`
}

type CTAButtonRes struct {
	Type  enum.CTAType `json:"type"`
	Value string       `json:"value"`
}

type PageSettingsRes struct {
	AllowVisitorPost bool                 `json:"allow_visitor_post"`
	ProfanityFilter  enum.ProfanityFilter `json:"profanity_filter"`
	MessagingStatus  enum.MessagingStatus `json:"messaging_status"`
}

type PageStatsRes struct {
	FollowersCount int     `json:"followers_count"`
	LikesCount     int     `json:"likes_count"`
	RatingScore    float64 `json:"rating_score"`
	ReviewCount    int     `json:"review_count"`
}
