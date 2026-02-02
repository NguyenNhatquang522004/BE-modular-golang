package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender int

const (
	GenderMale Gender = iota
	GenderFemale
	GenderOther
	GenderHidden
)
const (
	collectionProfiles = "Profilesss"
)

type Profiles struct {
	// 1. ĐỊNH DANH
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// UserID từ Postgres (UUID). Lưu String để query dễ dàng và tương thích JSON.
	// Index: Unique
	UserID string `bson:"user_id" json:"user_id"`

	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`

	// Denormalization: Nên gộp FirstName + LastName khi save để search text nhanh
	FullName string `bson:"full_name" json:"full_name"`

	// Index: Unique. Dùng để tạo URL đẹp: facebook.com/nguyen-van-a
	Slug string `bson:"slug" json:"slug"`

	Bio         string      `bson:"bio" json:"bio"`
	DateOfBirth time.Time   `bson:"date_of_birth" json:"date_of_birth"`
	Gender      enum.Gender `bson:"gender" json:"gender"` // Enum Int -> Lưu String trong DB

	// 2. MEDIA & LIÊN HỆ
	// Dùng pointer (*) cho các object con.
	// Lợi ích: Nếu user chưa set avatar, DB sẽ không lưu field này (tiết kiệm),
	// và JSON trả về null giúp Frontend biết là chưa có.
	Avatar      *ProfileAvatar  `bson:"avatar,omitempty" json:"avatar,omitempty"`
	CoverPhoto  *ProfileCover   `bson:"cover_photo,omitempty" json:"cover_photo,omitempty"`
	Address     *ProfileAddress `bson:"address,omitempty" json:"address,omitempty"`
	PhoneNumber string          `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	SocialLinks *SocialLinks    `bson:"social_links,omitempty" json:"social_links,omitempty"`

	// 3. HỒ SƠ NĂNG LỰC
	CVDocument     *CVDocument       `bson:"cv_document,omitempty" json:"cv_document,omitempty"`
	WorkExperience []*WorkExperience `bson:"work_experience,omitempty" json:"work_experience,omitempty"`
	Education      []*Education      `bson:"education,omitempty" json:"education,omitempty"`

	// 4. META & SETTINGS
	// Settings nên khởi tạo mặc định, không nên để nil pointer
	Settings      ProfileSettings      `bson:"settings" json:"settings"`
	Notifications ProfileNotifications `bson:"notifications" json:"notifications"`

	// 5. TIMESTAMPS
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // Pointer để check null
}
type ProfileAvatar struct {
	URL       string    `bson:"url" json:"url"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type ProfileCover struct {
	URL       string  `bson:"url" json:"url"`
	PositionY float64 `bson:"position_y" json:"position_y"` // 0.0 đến 100.0
}

// --- LOCATION ---

type ProfileAddress struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	Country string `bson:"country" json:"country"`
	// Lưu tọa độ dạng mảng: [Longitude, Latitude]
	//omitempty vì user có thể chưa set tọa độ
	Coordinates []float64 `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
}

// --- SOCIAL LINKS ---

type SocialLinks struct {
	Facebook  string `bson:"facebook,omitempty" json:"facebook,omitempty"`
	Twitter   string `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Linkedin  string `bson:"linkedin,omitempty" json:"linkedin,omitempty"`
	Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Github    string `bson:"github,omitempty" json:"github,omitempty"`
	Website   string `bson:"website,omitempty" json:"website,omitempty"`
}

// --- PROFESSIONAL INFO ---

type CVDocument struct {
	FileID     primitive.ObjectID `bson:"file_id" json:"file_id"` // Link sang GridFS
	Filename   string             `bson:"filename" json:"filename"`
	UploadedAt time.Time          `bson:"uploaded_at" json:"uploaded_at"`
}

type WorkExperience struct {
	// Tự tạo ID cho sub-document để dễ xóa/sửa chính xác item này trong mảng
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Company     string             `bson:"company" json:"company"`
	Position    string             `bson:"position" json:"position"`
	StartDate   time.Time          `bson:"start_date" json:"start_date"`
	EndDate     *time.Time         `bson:"end_date" json:"end_date"` // Pointer: Null = Đang làm việc
	IsCurrent   bool               `bson:"is_current" json:"is_current"`
	Description string             `bson:"description" json:"description"`
}

type Education struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Institution  string             `bson:"institution" json:"institution"`
	Degree       string             `bson:"degree" json:"degree"`
	FieldOfStudy string             `bson:"field_of_study" json:"field_of_study"`
	StartDate    time.Time          `bson:"start_date" json:"start_date"`
	EndDate      *time.Time         `bson:"end_date" json:"end_date"` // Pointer: Null = Chưa tốt nghiệp
}

// --- SETTINGS (Riêng cho Profile) ---

type ProfileSettings struct {
	IsPrivate         bool `bson:"is_private" json:"is_private"` // Không dùng omitempty cho bool
	AllowSearchEngine bool `bson:"allow_search_engine" json:"allow_search_engine"`
}

type ProfileNotifications struct {
	EmailFrequency string   `bson:"email_frequency" json:"email_frequency"` // 'weekly', 'daily'
	PushTypes      []string `bson:"push_types" json:"push_types"`           // ['comment', 'friend_request']
}

func (Profiles) CollectionNameProfiles() string {
	return collectionProfiles
}
