package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSessionReq struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// 2. Foreign Key: Viết hoa UserID
	UserID string

	// 3. Data Fields: Viết hoa chữ cái đầu (PascalCase)
	DeviceName      string
	OsVersion       string
	Browser         string
	IpAddress       string
	LocationCity    string
	LocationCountry string
	RefreshToken    string

	IsActive     bool
	LastActiveAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Dùng gorm.DeletedAt chuẩn hơn *time.Time
}

func ToUserSessionResp(e *entity.UserSession) *UserSessionReq {
	if e == nil {
		return nil
	}
	return &UserSessionReq{
		ID:              e.ID,
		DeviceName:      e.DeviceName,
		OsVersion:       e.OsVersion,
		Browser:         e.Browser,
		IpAddress:       e.IpAddress,
		LocationCity:    e.LocationCity,
		LocationCountry: e.LocationCountry,
		IsActive:        e.IsActive,
		LastActiveAt:    e.LastActiveAt,
		CreatedAt:       e.CreatedAt,
	}
}

// Helper để map một danh sách (List)
func ToUserSessionRespList(sessions []entity.UserSession) []UserSessionReq {
	result := make([]UserSessionReq, len(sessions))
	for i, s := range sessions {
		result[i] = *ToUserSessionResp(&s)
	}
	return result
}
func (req *UserSessionReq) ToEntity(userID uuid.UUID, refreshToken string, ip string, city string, country string) *entity.UserSession {
	return &entity.UserSession{
		// ID: Postgres tự sinh (gen_random_uuid)
		UserID:          userID,
		DeviceName:      req.DeviceName,
		OsVersion:       req.OsVersion,
		Browser:         req.Browser,
		IpAddress:       ip,      // Lấy từ Context request
		LocationCity:    city,    // Lấy từ GeoIP service
		LocationCountry: country, // Lấy từ GeoIP service
		IsActive:        true,
		LastActiveAt:    time.Now(),
	}
}
