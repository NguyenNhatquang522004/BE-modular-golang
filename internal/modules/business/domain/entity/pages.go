package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionPages = "Pages"
)

type Page struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// UserID người tạo/chủ sở hữu (Postgres UUID -> String)
	// Index: { creator_user_id: 1 } -> List pages của user
	CreatorUserID string `bson:"creator_user_id" json:"creator_user_id"`

	// 1. ĐỊNH DANH
	// Index: Text { name: "text" } -> Tìm kiếm Page
	Name string `bson:"name" json:"name"`

	// Index: Unique { slug: 1 } -> URL: facebook.com/cafebiz.vn
	Slug string `bson:"slug" json:"slug"`

	// CategoryID: String (Ref tới Collection Categories hoặc Hardcode Enum)
	// Index: { category_id: 1 } -> Lọc theo ngành hàng
	CategoryID string          `bson:"category_id" json:"category_id"`
	IsVerified bool            `bson:"is_verified" json:"is_verified"`
	Status     enum.PageStatus `bson:"status" json:"status"`

	// 2. BRANDING
	Avatar *PageAvatar `bson:"avatar" json:"avatar"`
	Cover  *PageCover  `bson:"cover" json:"cover"`

	// 3. BUSINESS INFO
	Bio         string       `bson:"bio" json:"bio"`
	Website     string       `bson:"website" json:"website"`
	Email       string       `bson:"email" json:"email"`
	PhoneNumber string       `bson:"phone_number" json:"phone_number"`
	Address     *PageAddress `bson:"address" json:"address"`

	// 4. GIỜ MỞ CỬA (Array)
	BusinessHours []BusinessHour `bson:"business_hours,omitempty" json:"business_hours,omitempty"`

	// 5. CTA
	CTAButton *CTAButton `bson:"cta_button,omitempty" json:"cta_button,omitempty"`

	// 6. CÀI ĐẶT
	Settings *PageSettings `bson:"settings" json:"settings"`

	// 7. METRICS
	// Index: { "stats.followers_count": -1 } -> Gợi ý page hot
	Stats *PageStats `bson:"stats" json:"stats"`

	// 8. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- BRANDING ---
type PageAvatar struct {
	URL string `bson:"url" json:"url"`
}

type PageCover struct {
	URL       string  `bson:"url" json:"url"`
	PositionY float64 `bson:"position_y" json:"position_y"` // 0.0 - 100.0
}

// --- ADDRESS (GeoJSON support) ---
type PageAddress struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	Zipcode string `bson:"zipcode" json:"zipcode"`

	// GeoJSON: [Longitude, Latitude]
	// Index: 2dsphere
	Coordinates []float64 `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
}

// --- BUSINESS HOURS ---
type BusinessHour struct {
	Day   enum.DayOfWeek `bson:"day" json:"day"`     // 'Monday'
	Open  string         `bson:"open" json:"open"`   // "08:00"
	Close string         `bson:"close" json:"close"` // "22:00"
}

// --- CALL TO ACTION ---
type CTAButton struct {
	Type  enum.CTAType `bson:"type" json:"type"`   // 'send_message'
	Value string       `bson:"value" json:"value"` // Link hoặc số điện thoại
}

// --- SETTINGS ---
type PageSettings struct {
	AllowVisitorPost bool                 `bson:"allow_visitor_post" json:"allow_visitor_post"`
	ProfanityFilter  enum.ProfanityFilter `bson:"profanity_filter" json:"profanity_filter"` // 'strong'
	MessagingStatus  enum.MessagingStatus `bson:"messaging_status" json:"messaging_status"` // 'online'
}

// --- METRICS ---
type PageStats struct {
	FollowersCount int     `bson:"followers_count" json:"followers_count"`
	LikesCount     int     `bson:"likes_count" json:"likes_count"`
	RatingScore    float64 `bson:"rating_score" json:"rating_score"` // VD: 4.8
	ReviewCount    int     `bson:"review_count" json:"review_count"`
}

func (Page) CollectionName() string {
	return CollectionPages
}
