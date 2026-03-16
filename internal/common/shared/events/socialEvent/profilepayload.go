package socialEvent

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type ProfilePayload struct {
	UserID      string             `json:"user_id" binding:"required"`
	FirstName   string             `json:"first_name" binding:"required"`
	LastName    string             `json:"last_name" binding:"required"`
	Bio         string             `json:"bio"`
	DateOfBirth string             `json:"date_of_birth"`
	Gender      sharedEnums.Gender `json:"gender"` // Sử dụng Enum như yêu cầu

	// Các thông tin bổ sung (Optional)
	Avatar      *AvatarPayload      `json:"avatar"`
	CoverPhoto  *CoverPhotoPayload  `json:"cover_photo"`
	Address     *AddressPayload     `json:"address"`
	PhoneNumber string              `json:"phone_number"`
	SocialLinks *SocialLinksPayload `json:"social_links"`

	// Hồ sơ năng lực
	CVDocument     *CVDocumentPayload       `json:"cv_document"`
	WorkExperience []*WorkExperiencePayload `json:"work_experience"`
	Education      []*EducationPayload      `json:"education"`

	// Settings
	Settings *ProfileSettingsPayload `json:"settings"`
}

// --- SUB-DTOs ---

type AvatarPayload struct {
	ID  string `json:"id" binding:"required"` // Nhận string ID từ FE
	URL string `json:"url" binding:"required"`
}

type CoverPhotoPayload struct {
	ID        string  `json:"id" binding:"required"` // Nhận string ID từ FE
	URL       string  `json:"url" binding:"required"`
	PositionY float64 `json:"position_y"`
}

type AddressPayload struct {
	Street      string    `json:"street"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Coordinates []float64 `json:"coordinates"` // [Long, Lat]
}

type SocialLinksPayload struct {
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	Linkedin  string `json:"linkedin"`
	Instagram string `json:"instagram"`
	Github    string `json:"github"`
	Website   string `json:"website"`
}

type CVDocumentPayload struct {
	FileID   string `json:"file_id" binding:"required"` // Nhận string ID từ FE
	Filename string `json:"filename"`
}

type WorkExperiencePayload struct {
	ID          string     `json:"id"`
	Company     string     `json:"company" binding:"required"`
	Position    string     `json:"position" binding:"required"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"` // Null nếu đang làm việc
	IsCurrent   bool       `json:"is_current"`
	Description string     `json:"description"`
}

type EducationPayload struct {
	ID           string     `json:"id"`
	Institution  string     `json:"institution" binding:"required"`
	Degree       string     `json:"degree"`
	FieldOfStudy string     `json:"field_of_study"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
}

type ProfileSettingsPayload struct {
	IsPrivate         bool `json:"is_private"`
	AllowSearchEngine bool `json:"allow_search_engine"`
}
