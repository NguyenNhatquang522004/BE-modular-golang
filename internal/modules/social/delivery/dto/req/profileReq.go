package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateProfileReq dùng cho việc tạo mới Profile
type CreateProfileReq struct {
	UserID      string      `json:"user_id" binding:"required"`
	FirstName   string      `json:"first_name" binding:"required"`
	LastName    string      `json:"last_name" binding:"required"`
	Bio         string      `json:"bio"`
	DateOfBirth time.Time   `json:"date_of_birth"`
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
	Notifications *ProfileNotificationsReq `json:"notifications"`
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
	Company     string     `json:"company" binding:"required"`
	Position    string     `json:"position" binding:"required"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"` // Null nếu đang làm việc
	IsCurrent   bool       `json:"is_current"`
	Description string     `json:"description"`
}

type EducationReq struct {
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

type ProfileNotificationsReq struct {
	EmailFrequency string   `json:"email_frequency"`
	PushTypes      []string `json:"push_types"`
}

// ToEntity chuyển đổi từ DTO Request sang Entity để lưu DB
func (req *CreateProfileReq) ToEntity() *entity.Profiles {
	now := time.Now()

	// 1. Khởi tạo object chính
	profile := &entity.Profiles{
		ID:          primitive.NewObjectID(),
		UserID:      req.UserID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		FullName:    req.FirstName + " " + req.LastName, // Logic Denormalization
		Bio:         req.Bio,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender, // Map trực tiếp Enum
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Xử lý Slug (Cần hàm util riêng để convert tên thành slug, ví dụ tạm thời)
	// profile.Slug = utils.GenerateSlug(profile.FullName)

	// 2. Mapping Avatar
	if req.Avatar != nil {
		profile.Avatar = &entity.ProfileAvatar{
			URL:       req.Avatar.URL,
			UpdatedAt: now,
		}
	}

	// 3. Mapping Cover Photo
	if req.CoverPhoto != nil {
		profile.CoverPhoto = &entity.ProfileCover{
			URL:       req.CoverPhoto.URL,
			PositionY: req.CoverPhoto.PositionY,
		}
	}

	// 4. Mapping Address
	if req.Address != nil {
		profile.Address = &entity.ProfileAddress{
			Street:      req.Address.Street,
			City:        req.Address.City,
			Country:     req.Address.Country,
			Coordinates: req.Address.Coordinates,
		}
	}

	// 5. Mapping Social Links
	if req.SocialLinks != nil {
		profile.SocialLinks = &entity.SocialLinks{
			Facebook:  req.SocialLinks.Facebook,
			Twitter:   req.SocialLinks.Twitter,
			Linkedin:  req.SocialLinks.Linkedin,
			Instagram: req.SocialLinks.Instagram,
			Github:    req.SocialLinks.Github,
			Website:   req.SocialLinks.Website,
		}
	}

	// 6. Mapping CV Document (Convert String ID -> ObjectID)
	if req.CVDocument != nil {
		if fileID, err := primitive.ObjectIDFromHex(req.CVDocument.FileID); err == nil {
			profile.CVDocument = &entity.CVDocument{
				FileID:     fileID,
				Filename:   req.CVDocument.Filename,
				UploadedAt: now,
			}
		}
	}

	// 7. Mapping Work Experience (Tự tạo ObjectID cho sub-doc)
	if len(req.WorkExperience) > 0 {
		var works []*entity.WorkExperience
		for _, w := range req.WorkExperience {
			works = append(works, &entity.WorkExperience{
				ID:          primitive.NewObjectID(), // Tạo ID mới
				Company:     w.Company,
				Position:    w.Position,
				StartDate:   w.StartDate,
				EndDate:     w.EndDate,
				IsCurrent:   w.IsCurrent,
				Description: w.Description,
			})
		}
		profile.WorkExperience = works
	}

	// 8. Mapping Education (Tự tạo ObjectID cho sub-doc)
	if len(req.Education) > 0 {
		var edus []*entity.Education
		for _, e := range req.Education {
			edus = append(edus, &entity.Education{
				ID:           primitive.NewObjectID(), // Tạo ID mới
				Institution:  e.Institution,
				Degree:       e.Degree,
				FieldOfStudy: e.FieldOfStudy,
				StartDate:    e.StartDate,
				EndDate:      e.EndDate,
			})
		}
		profile.Education = edus
	}

	// 9. Mapping Settings (Set default nếu nil hoặc map từ request)
	if req.Settings != nil {
		profile.Settings = entity.ProfileSettings{
			IsPrivate:         req.Settings.IsPrivate,
			AllowSearchEngine: req.Settings.AllowSearchEngine,
		}
	} else {
		// Default settings
		profile.Settings = entity.ProfileSettings{
			IsPrivate:         false,
			AllowSearchEngine: true,
		}
	}

	// 10. Mapping Notifications
	if req.Notifications != nil {
		profile.Notifications = entity.ProfileNotifications{
			EmailFrequency: req.Notifications.EmailFrequency,
			PushTypes:      req.Notifications.PushTypes,
		}
	}

	return profile
}
