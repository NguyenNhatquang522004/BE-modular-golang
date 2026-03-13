package contentEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type DeletePostPayload struct {
	PostID string `json:"post_id"`
	Reason string `json:"reason,omitempty"`
}
type UpdatePostPayload struct {
	PostID      string                `json:"post_id"`
	Content     *string               `json:"content,omitempty"`
	Privacy     *UpdatePrivacyPayload `json:"privacy,omitempty"`
	Hashtags    *[]string             `json:"hashtags,omitempty" validate:"omitempty,dive,max=50"`
	Mentions    *[]string             `json:"mentions,omitempty" validate:"omitempty,dive,uuid"`
	Deletemedia *[]string             `json:"delete_media,omitempty"` // Danh sách media_id cần xoá khỏi post, nếu client muốn xoá media nào đó thì gửi lên ID của nó, backend sẽ xoá khỏi post. Nếu client muốn xoá tất cả media thì gửi lên tất cả media_id hiện có của post.
	// Media thay vì update lẻ tẻ, thường Best Practice là client sẽ gửi lại MẢNG MỚI HOÀN TOÀN
	// Backend sẽ xoá media cũ và insert media mới để tránh rác logic
	Media     *[]MediaItemPayload           `json:"media,omitempty" validate:"omitempty,max=10"`
	Status    *sharedEnums.ProcessingStatus `json:"status,omitempty"`
	Extension *ExtensionPayload             `json:"extension,omitempty"`
	Setting   *SettingPayload               `json:"setting,omitempty"`
}

type UpdatePrivacyPayload struct {
	Scope        *sharedEnums.PrivacyScope `json:"scope,omitempty"`
	AllowComment *bool                     `json:"allow_comment,omitempty"`
	AllowShare   *bool                     `json:"allow_share,omitempty"`
}

// CreatePostPayload là payload duy nhất Client cần gửi lên
type CreatePostPayload struct {
	ID      string               `json:"id"` // ID do backend tạo ra, client không gửi lên
	UserID  string               // Trường này không lấy từ payload, sẽ được gán sau khi giải mã token
	Group   *string              `json:"group_id,omitempty"`       // Nếu có group_id thì sẽ là post nhóm, không có thì sẽ là post cá nhân
	Type    sharedEnums.PostType `json:"type" validate:"required"` // Bắt buộc
	Content string               `json:"content"`                  // Có thể rỗng nếu chỉ đăng ảnh
	Privacy PostPrivacyPayload   `json:"privacy" validate:"required"`

	// Optional (Pointer để phân biệt có gửi lên hay không)
	Context  *PostContextPayload `json:"context,omitempty"`
	Hashtags []string            `json:"hashtags,omitempty" validate:"omitempty,dive,max=50"`
	Mentions []string            `json:"mentions,omitempty" validate:"omitempty,dive,uuid"` // Kiểm tra chuẩn UUID

	// Các thành phần mở rộng
	Media     []MediaItemPayload `json:"media,omitempty" validate:"omitempty,max=10"` // Giới hạn max 10 ảnh/video
	Extension *ExtensionPayload  `json:"extension,omitempty"`
	Setting   *SettingPayload    `json:"setting,omitempty"`
}

// --- SUB-STRUCTS CHO CREATE ---

type PostPrivacyPayload struct {
	Scope        sharedEnums.PrivacyScope `json:"scope" validate:"required"`
	AllowComment *bool                    `json:"allow_comment" validate:"required"` // Pointer bool để tránh default false
	AllowShare   *bool                    `json:"allow_share" validate:"required"`
}

type PostContextPayload struct {
	Type     sharedEnums.ContextType `json:"type" validate:"required"`
	TargetID string                  `json:"target_id,omitempty"`
}

type MediaItemPayload struct {
	MediaType    sharedEnums.MediaType `json:"media_type" validate:"required"`
	URL          string                `json:"url" validate:"required,url"`
	ThumbnailURL string                `json:"thumbnail_url,omitempty" validate:"omitempty,url"`

	// Chỉ gửi thông tin cơ bản, backend/cloud tự lấy width/height chính xác sau
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
	Duration  float64 `json:"duration,omitempty"` // Cho video
	SizeBytes int64   `json:"size_bytes,omitempty"`
	MimeType  string  `json:"mime_type,omitempty"`

	TaggedUsers []TaggedUserPayload `json:"tagged_users,omitempty"`
}

type TaggedUserPayload struct {
	UserID string  `json:"user_id" validate:"required,uuid"`
	Name   string  `json:"name" validate:"required"`
	X      float64 `json:"x" validate:"min=0,max=1"` // Tọa độ phải từ 0->1
	Y      float64 `json:"y" validate:"min=0,max=1"`
}

type ExtensionPayload struct {
	// Client chỉ gửi 1 trong các cục này tùy tính năng người dùng chọn
	ShareData      *SharePayload      `json:"share_data,omitempty"`
	BackgroundData *BackgroundPayload `json:"background_data,omitempty"`
	QnAData        *QnAPayload        `json:"qna_data,omitempty"`
	ActivityData   *ActivityPayload   `json:"activity_data,omitempty"`
	LocationDetail *LocationPayload   `json:"location_detail,omitempty"`
}

type SharePayload struct {
	OriginalPostID string `json:"original_post_id" validate:"required,mongodb"`
}

type BackgroundPayload struct {
	ThemeID   string `json:"theme_id" validate:"required"`
	TextColor string `json:"text_color" validate:"required,hexcolor"`
}

type QnAPayload struct {
	Question   string `json:"question" validate:"required,max=255"`
	ButtonText string `json:"button_text" validate:"required,max=50"`
}

type ActivityPayload struct {
	Type       sharedEnums.ActivityType `json:"type" validate:"required"`
	ObjectID   string                   `json:"object_id,omitempty"`
	ObjectName string                   `json:"object_name" validate:"required"`
}

type LocationPayload struct {
	Longitude float64 `json:"longitude" validate:"min=-180,max=180"`
	Latitude  float64 `json:"latitude" validate:"min=-90,max=90"`
	Address   string  `json:"address" validate:"required"`
}

type SettingPayload struct {
	Schedule  *SchedulePayload  `json:"schedule,omitempty"`
	Targeting *TargetingPayload `json:"targeting,omitempty"`
}

type SchedulePayload struct {
	PublishTime time.Time `json:"publish_time" validate:"required,gt=now"` // Phải là giờ trong tương lai
}

type TargetingPayload struct {
	Locations []string `json:"locations,omitempty"`
	AgeMin    int      `json:"age_min,omitempty" validate:"omitempty,min=13"`
	AgeMax    int      `json:"age_max,omitempty" validate:"omitempty,gtefield=AgeMin"`
	Genders   []string `json:"genders,omitempty"`
	Languages []string `json:"languages,omitempty"`
	Interests []string `json:"interests,omitempty"`
}
