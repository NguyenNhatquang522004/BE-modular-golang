package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/enum"
)

// --- 1. CREATE REQUEST ---

type CreatePostRequest struct {
	// Không nhận ID hay Status từ client khi tạo (xử lý logic ở service)
	Type    enum.PostType       `json:"type" validate:"required"`
	Context *PostContextRequest `json:"context,omitempty"` // Có thể null nếu là post thường

	Content string              `json:"content" validate:"required,max=5000"`
	Summary *PostSummaryRequest `json:"summary,omitempty"`

	Privacy PostPrivacyRequest `json:"privacy" validate:"required"`

	// Hashtag và Mention là mảng string/uuid từ client
	Hashtags []string `json:"hashtags,omitempty"`
	Mentions []string `json:"mentions,omitempty"`

	// Nếu client muốn hẹn giờ đăng (Optional)
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// --- 2. UPDATE REQUEST ---

type UpdatePostRequest struct {
	// Không cho phép sửa Type, Context sau khi tạo để đảm bảo tính toàn vẹn
	Content string              `json:"content,omitempty"`
	Summary *PostSummaryRequest `json:"summary,omitempty"`
	Privacy *PostPrivacyRequest `json:"privacy,omitempty"` // Pointer để check nếu user muốn update privacy

	Status   *enum.PostStatus `json:"status,omitempty"`    // Dùng để Archive/Hide bài viết
	IsPinned *bool           `json:"is_pinned,omitempty"` // Pointer để phân biệt false và nil

	Hashtags []string `json:"hashtags,omitempty"`
	Mentions []string `json:"mentions,omitempty"`
}

// --- SUB-STRUCTS FOR REQUEST ---

type PostContextRequest struct {
	Type     enum.ContextType `json:"type" validate:"required"`
	TargetID string           `json:"target_id" validate:"required"`
}

type PostSummaryRequest struct {
	FeelingIcon  string `json:"feeling_icon,omitempty"`
	FeelingName  string `json:"feeling_name,omitempty"`
	LocationName string `json:"location_name,omitempty"`

	// Media thường được upload qua API riêng và trả về ID/URL,
	// client gửi metadata media vào đây để lưu snapshot.
	HasMedia          bool   `json:"has_media"`
	MediaCount        int    `json:"media_count"`
	ThumbnailURL      string `json:"thumbnail_url,omitempty"`
	BackgroundThemeID string `json:"background_theme_id,omitempty"`
}

type PostPrivacyRequest struct {
	Scope        enum.PrivacyScope `json:"scope" validate:"required"`
	AllowComment bool              `json:"allow_comment"`
	AllowShare   bool              `json:"allow_share"`
}
