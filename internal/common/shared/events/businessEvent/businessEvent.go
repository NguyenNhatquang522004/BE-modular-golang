package businessEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/enum"
)

type StatsPagePayload struct {
	UserID         string    `json:"user_id"`
	PageID         string    `json:"page_id"`
	FollowersCount int       `json:"followers_count"`
	LikesCount     int       `json:"likes_count"`
	RatingScore    float64   `json:"rating_score"`
	ReviewCount    int       `json:"review_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type FollowerPagePayload struct {
	ID       string                   `json:"id,omitempty" validate:"omitempty"`
	PageID   string                   `json:"page_id" validate:"required"`
	UserID   string                   `json:"user_id" validate:"required,uuid"` // Đảm bảo là UUID theo chuẩn Postgres
	Settings *FollowerSettingsPayload `json:"settings,omitempty"`
}
type FollowerSettingsPayload struct {
	NotificationLevel sharedEnums.NotificationLevel `json:"notification_level" validate:"required"`
	IsFavorite        bool                          `json:"is_favorite"`
}

type TopicDailyMetricsPagePayload struct {
	ID         string    `json:"page_id" binding:"required,uuid"`
	MetricDate time.Time `json:"metric_date" binding:"required"`

	ReachTotal   int64 `json:"reach_total" binding:"gte=0"`
	ReachPaid    int64 `json:"reach_paid" binding:"gte=0"`
	ReachOrganic int64 `json:"reach_organic" binding:"gte=0"`

	ImpressionsTotal int64 `json:"impressions_total" binding:"gte=0"`
	NewFollowers     int   `json:"new_followers" binding:"gte=0"`
	Unfollows        int   `json:"unfollows" binding:"gte=0"`

	ProfileViews  int `json:"profile_views" binding:"gte=0"`
	WebsiteClicks int `json:"website_clicks" binding:"gte=0"`
	CTAClicks     int `json:"cta_clicks" binding:"gte=0"`
}

type AvatarPayload struct {
	ID  string `json:"id,omitempty"` // ID của media đã tồn tại, nếu có
	URL string `json:"url" validate:"required,url"`
}

type CoverPayload struct {
	ID        string  `json:"id,omitempty"` // ID của media đã tồn tại, nếu có
	URL       string  `json:"url" validate:"required,url"`
	PositionY float64 `json:"position_y" validate:"min=0,max=100"`
}

type AddressPayload struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	Zipcode string `json:"zipcode" validate:"omitempty"`
	// Yêu cầu mảng chính xác 2 phần tử: [Longitude, Latitude]
	Coordinates []float64 `json:"coordinates,omitempty" validate:"omitempty,len=2"`
}

type BusinessHourPayload struct {
	Day   enum.DayOfWeek `json:"day" validate:"required"`
	Open  string         `json:"open" validate:"required,datetime=15:04"`  // Format HH:mm
	Close string         `json:"close" validate:"required,datetime=15:04"` // Format HH:mm
}

type CTAPayload struct {
	Type  enum.CTAType `json:"type" validate:"required"`
	Value string       `json:"value" validate:"required"`
}

type SettingsPayload struct {
	// Dùng con trỏ *bool để phân biệt false và không gửi (nil)
	AllowVisitorPost *bool                 `json:"allow_visitor_post" validate:"required"`
	ProfanityFilter  sharedEnums.UserBadge `json:"profanity_filter" validate:"required"`
	MessagingStatus  enum.MessagingStatus  `json:"messaging_status" validate:"required"`
}

// ==========================================
// 2. CREATE PAGE PAYLOAD
// Dùng cho POST /pages
// ==========================================

type CreatePagePayload struct {
	PageID string `json:"page_id,omitempty"` // Tự sinh UUID nếu không gửi
	UserID string `json:"user_id" validate:"required,uuid"`
	// CÁC TRƯỜNG BẮT BUỘC (Required)
	Name       string `json:"name" validate:"required,min=2,max=255"`
	CategoryID string `json:"category_id" validate:"required"` // Có thể thêm validate uuid/objectid tuỳ logic

	// CÁC TRƯỜNG TUỲ CHỌN (Optional)
	Slug        string `json:"slug,omitempty" validate:"omitempty,min=3,max=100"`
	Bio         string `json:"bio,omitempty" validate:"omitempty,max=1000"`
	Website     string `json:"website,omitempty" validate:"omitempty,url"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"omitempty,min=8,max=20"`

	Avatar        *AvatarPayload        `json:"avatar,omitempty" validate:"omitempty"`
	Cover         *CoverPayload         `json:"cover,omitempty" validate:"omitempty"`
	Address       *AddressPayload       `json:"address,omitempty" validate:"omitempty"`
	BusinessHours []BusinessHourPayload `json:"business_hours,omitempty" validate:"omitempty,dive"` // dive: validate từng phần tử trong mảng
	CTAButton     *CTAPayload           `json:"cta_button,omitempty" validate:"omitempty"`
	Settings      *SettingsPayload      `json:"settings,omitempty" validate:"omitempty"`
}
type UpdatePagePayload struct {
	UserActionID string `json:"user_action" validate:"required"` // "update_by_user" hoặc "auto_update"
	// TẤT CẢ CÁC TRƯỜNG LÀ CON TRỎ (Pointer)
	// Lý do: Để biết chính xác Client muốn update trường nào.
	// Nếu Client không gửi Name, Name sẽ là nil -> Không update.
	// Nếu Client gửi Name = "", Name != nil và value là "" -> Update thành chuỗi rỗng.
	PageID     string  `json:"page_id" validate:"required,uuid"`
	Name       *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Slug       *string `json:"slug,omitempty" validate:"omitempty,min=3,max=100"`
	CategoryID *string `json:"category_id,omitempty" validate:"omitempty"`

	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=1000"`
	Website     *string `json:"website,omitempty" validate:"omitempty,url"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,min=8,max=20"`

	Avatar        *AvatarPayload        `json:"avatar,omitempty" validate:"omitempty"`
	Cover         *CoverPayload         `json:"cover,omitempty" validate:"omitempty"`
	Address       *AddressPayload       `json:"address,omitempty" validate:"omitempty"`
	BusinessHours []BusinessHourPayload `json:"business_hours,omitempty" validate:"omitempty,dive"`
	CTAButton     *CTAPayload           `json:"cta_button,omitempty" validate:"omitempty"`
	Settings      *SettingsPayload      `json:"settings,omitempty" validate:"omitempty"`
}

type DeletedPagePayload struct {
	UserActionID string `json:"user_action" validate:"required"` // "delete_by_user" hoặc "auto_remove"
	PageID       string `json:"page_id" validate:"required,uuid"`
	UserID       string `json:"user_id" validate:"required,uuid"`
}

type CreatePageRolePayload struct {
	// Dùng string ở payload, Service layer sẽ dùng primitive.ObjectIDFromHex() để convert
	// tag 'mongodb' đảm bảo string truyền vào đúng format 24 hex characters
	PageID string `json:"page_id" binding:"required,mongodb"`

	// tag 'uuid' đảm bảo string truyền vào đúng chuẩn UUID của Postgres
	UserID     string `json:"user_id" binding:"required,uuid"`
	AssignedBy string `json:"-"` // Không nhận từ payload, sẽ set trong Service layer dựa trên user đang login
	// Role bắt buộc phải có
	Role sharedEnums.RoleType `json:"role" binding:"required"`

	// Optional. Tag 'dive,required' đảm bảo nếu có mảng, thì các phần tử trong mảng không được rỗng ("")
	CustomPermissions []string `json:"custom_permissions,omitempty" binding:"omitempty,dive,required"`
}

type UpdatePageRolePayload struct {
	UserActionID string `json:"user_action" binding:"required"` // "update_by_user" hoặc "auto_update"
	PageRoleID   string `json:"page_role_id" binding:"required,uuid"`
	// Dùng pointer (*sharedEnums.RoleType) để hỗ trợ partial update (PATCH).
	// Nếu dùng value thường, ta sẽ không biết là client gửi zero-value hay không gửi.
	Role *sharedEnums.RoleType `json:"role,omitempty" binding:"omitempty"`

	// Slice bản chất đã là reference type (có thể nil), nên không cần pointer.
	// Nếu client gửi mảng rỗng [], ta sẽ ghi đè xóa hết permission.
	CustomPermissions []string `json:"custom_permissions,omitempty" binding:"omitempty,dive,required"`
}
type DeletedPageRolePayload struct {
	UserAction string `json:"user_action" binding:"required"` // "delete_by_user" hoặc "auto_remove"
	PageRoleID string `json:"page_role_id" binding:"required,uuid"`
	UserID     string `json:"user_id" binding:"required,uuid"`
	DeLeteAll  bool   `json:"delete_all"` // Nếu true, sẽ xóa tất cả role của user này trên page (dùng khi xóa user hoặc xóa page)
}
type CreatePageFollowerPayload struct {
	// Dùng string kết hợp tag validate "mongodb" để an toàn khi parse JSON.
	// Bắt buộc phải có để biết follow page nào.
	PageID string `json:"page_id" validate:"required,mongodb"`

	// Dùng string kết hợp tag validate "uuid".
	// Bắt buộc phải có để biết user nào đang follow.
	UserID string `json:"user_id" validate:"required,uuid"`

	// Cài đặt có thể có hoặc không lúc mới follow. Nếu không có, service sẽ set default.
	Settings *CreateFollowerSettingsPayload `json:"settings,omitempty" validate:"omitempty"`
}

type CreateFollowerSettingsPayload struct {
	// Require nếu client có gửi object Settings lên.
	NotificationLevel sharedEnums.NotificationLevel `json:"notification_level" validate:"required"`

	// Không dùng pointer ở Create vì nếu không truyền, mặc định (false) là hợp lý.
	IsFavorite bool `json:"is_favorite"`
}
type UpdatePageFollowerPayload struct {
	PageID string `json:"page_id" validate:"required,mongodb"`
	UserID string `json:"user_id" validate:"required,uuid"`
	// Bỏ qua ID, PageID, UserID vì chúng nằm ở URL params hoặc token, và không được phép thay đổi.

	// Settings bắt buộc phải có nếu gọi API update, nếu không thì gọi API làm gì?
	Settings *UpdateFollowerSettingsPayload `json:"settings" validate:"required"`
}

type UpdateFollowerSettingsPayload struct {
	// SỬ DỤNG POINTER (*):
	// - Nếu client không truyền `notification_level` -> giá trị là nil -> Backend bỏ qua, giữ nguyên DB.
	// - Nếu client truyền -> backend sẽ lấy giá trị được trỏ tới để update.
	NotificationLevel *sharedEnums.NotificationLevel `json:"notification_level,omitempty" validate:"omitempty"`

	// SỬ DỤNG POINTER (*):
	// - Cực kỳ quan trọng với kiểu bool. Nếu dùng `bool` thường, khi client không gửi, nó mặc định là `false`.
	// Backend sẽ không biết client muốn update thành `false` hay client quên không gửi. Dùng `*bool` sẽ giải quyết được.
	IsFavorite *bool `json:"is_favorite,omitempty" validate:"omitempty"`
}

type DeletedPageFollowerPayload struct {
	PageID    string `json:"page_id" validate:"required,mongodb"`
	UserID    string `json:"user_id" validate:"required,uuid"`
	DeleteAll bool   `json:"delete_all"` // Nếu true, sẽ xóa tất cả follow của user này trên page (dùng khi xóa user hoặc xóa page)
}

type PageDailyMetricPayload struct {
	// Dùng string thay vì gocql.UUID ở tầng giao tiếp API.
	// Yêu cầu client bắt buộc gửi và phải đúng định dạng UUID.
	PageID string `json:"page_id" validate:"required,uuid"`

	// Ngày thống kê là bắt buộc.
	MetricDate time.Time `json:"metric_date" validate:"required"`

	// --- REACH METRICS ---
	// Dùng con trỏ kèm gte=0 để đảm bảo data không bị âm (Greater Than or Equal to 0).
	ReachTotal   *int64 `json:"reach_total,omitempty" validate:"omitempty,gte=0"`
	ReachPaid    *int64 `json:"reach_paid,omitempty" validate:"omitempty,gte=0"`
	ReachOrganic *int64 `json:"reach_organic,omitempty" validate:"omitempty,gte=0"`

	// --- ENGAGEMENT METRICS ---
	ImpressionsTotal *int64 `json:"impressions_total,omitempty" validate:"omitempty,gte=0"`
	NewFollowers     *int   `json:"new_followers,omitempty" validate:"omitempty,gte=0"`
	Unfollows        *int   `json:"unfollows,omitempty" validate:"omitempty,gte=0"`

	// --- CONVERSION METRICS ---
	ProfileViews  *int `json:"profile_views,omitempty" validate:"omitempty,gte=0"`
	WebsiteClicks *int `json:"website_clicks,omitempty" validate:"omitempty,gte=0"`
	CTAClicks     *int `json:"cta_clicks,omitempty" validate:"omitempty,gte=0"`
	DeleteALL     bool `json:"delete_all"` // Nếu true, sẽ xóa tất cả metric của page này (dùng khi xóa page)
}
