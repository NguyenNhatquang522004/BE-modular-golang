package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionnamepostsetting = "PostSetting"
)

type PostSetting struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Index: Unique. Mỗi bài viết chỉ có 1 cấu hình setting.
	PostID primitive.ObjectID `bson:"post_id" json:"post_id"`

	// Dùng Pointer cho các nhóm cấu hình để có thể null (nil) nếu không set
	Schedule  *PostSchedule  `bson:"schedule,omitempty" json:"schedule,omitempty"`
	AdsInfo   *AdsInfo       `bson:"ads_info,omitempty" json:"ads_info,omitempty"`
	Targeting *PostTargeting `bson:"targeting,omitempty" json:"targeting,omitempty"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- 1. SCHEDULE (Quản lý đăng bài) ---
type PostSchedule struct {
	IsScheduled bool `bson:"is_scheduled" json:"is_scheduled"`

	// Thời điểm bài viết sẽ public
	PublishTime time.Time `bson:"publish_time" json:"publish_time"`

	// ID của người thực hiện hành động bấm nút "Schedule" (Postgres UUID)
	PublisherUserID string `bson:"publisher_user_id" json:"publisher_user_id"`

	// Vai trò của người đó tại thời điểm đăng (Snapshot lại để audit log)
	AuthorRoleSnapshot sharedEnums.RoleType `bson:"author_role_snapshot" json:"author_role_snapshot"`
}

// --- 2. ADS INFO (Quản lý quảng cáo) ---
type AdsInfo struct {
	// Link sang bảng Campaign (Collection Ads)
	CampaignID primitive.ObjectID `bson:"campaign_id" json:"campaign_id"`

	IsSponsored bool   `bson:"is_sponsored" json:"is_sponsored"`
	CTALink     string `bson:"cta_link" json:"cta_link"` // VD: "https://shop.com/buy-now"
}

// --- 3. TARGETING (Target Audience) ---
// Dùng cho thuật toán Feed hoặc Quảng cáo
type PostTargeting struct {
	// Địa điểm: Có thể là tên TP, hoặc ID địa điểm
	Locations []string `bson:"locations,omitempty" json:"locations,omitempty"`

	AgeMin int `bson:"age_min" json:"age_min"`
	AgeMax int `bson:"age_max" json:"age_max"`

	// Target giới tính: ["male", "female"]
	Genders []string `bson:"genders,omitempty" json:"genders,omitempty"`

	// Target ngôn ngữ: ["vi", "en"]
	Languages []string `bson:"languages,omitempty" json:"languages,omitempty"`

	// Sở thích: ["technology", "travel"]
	Interests []string `bson:"interests,omitempty" json:"interests,omitempty"`
}

func (PostSetting) CollectionNamePostsetting() string {
	return collectionnamepostsetting
}
