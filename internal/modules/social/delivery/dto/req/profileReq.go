package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

type ProfileIDRequest struct {
	ProfileID string `json:"profile_id" validate:"required,uuid4"`
}
type CreateAndUpdateProfileRequest struct {
	*ProfileIDRequest
	*ProfileReq
	Attachments []*dto.FileUploadInput
}

// CreateProfileReq dùng cho việc tạo mới Profile
type ProfileReq struct {
	UserID      string      `json:"user_id" binding:"required"`
	FirstName   string      `json:"first_name" binding:"required"`
	LastName    string      `json:"last_name" binding:"required"`
	Bio         string      `json:"bio"`
	DateOfBirth string      `json:"date_of_birth"`
	Gender      enum.Gender `json:"gender"` // Sử dụng Enum như yêu cầu

	// Các thông tin bổ sung (Optional)
	Avatar      *AvatarReq      `json:"avatar"`
	CoverPhoto  *CoverPhotoReq  `json:"cover_photo"`
	Address     *AddressReq     `json:"address"`
	PhoneNumber string          `json:"phone_number"`
	SocialLinks *SocialLinksReq `json:"social_links"`

	// Hồ sơ năng lực
	CVDocument     *CVDocumentReq       `json:"cv_document"`
	WorkExperience []*WorkExperienceReq `json:"work_experience"`
	Education      []*EducationReq      `json:"education"`

	// Settings
	Settings      *ProfileSettingsReq      `json:"settings"`
}

// --- SUB-DTOs ---

type AvatarReq struct {
	URL string `json:"url" binding:"required"`
}

type CoverPhotoReq struct {
	URL       string  `json:"url" binding:"required"`
	PositionY float64 `json:"position_y"`
}

type AddressReq struct {
	Street      string    `json:"street"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Coordinates []float64 `json:"coordinates"` // [Long, Lat]
}

type SocialLinksReq struct {
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	Linkedin  string `json:"linkedin"`
	Instagram string `json:"instagram"`
	Github    string `json:"github"`
	Website   string `json:"website"`
}

type CVDocumentReq struct {
	FileID   string `json:"file_id" binding:"required"` // Nhận string ID từ FE
	Filename string `json:"filename"`
}

type WorkExperienceReq struct {
	ID          string     `json:"id"`
	Company     string     `json:"company" binding:"required"`
	Position    string     `json:"position" binding:"required"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"` // Null nếu đang làm việc
	IsCurrent   bool       `json:"is_current"`
	Description string     `json:"description"`
}

type EducationReq struct {
	ID           string     `json:"id"`
	Institution  string     `json:"institution" binding:"required"`
	Degree       string     `json:"degree"`
	FieldOfStudy string     `json:"field_of_study"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
}

type ProfileSettingsReq struct {
	IsPrivate         bool `json:"is_private"`
	AllowSearchEngine bool `json:"allow_search_engine"`
}

