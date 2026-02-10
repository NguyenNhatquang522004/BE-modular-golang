package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	UserID uuid.UUID `json:"user_id" binding:"required,uuid"`
}
type UserSessionReq struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// 2. Foreign Key: Viết hoa UserID
	UserID uuid.UUID `gorm:"type:uuid;not null"`

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

// Helper để map từ Entity sang Response DTO
func (r *UserSessionReq) ToUserSessionResp() *UserSessionReq {
	if r == nil {
		return nil
	}
	return &UserSessionReq{
		ID:              r.ID,
		DeviceName:      r.DeviceName,
		OsVersion:       r.OsVersion,
		Browser:         r.Browser,
		IpAddress:       r.IpAddress,
		LocationCity:    r.LocationCity,
		LocationCountry: r.LocationCountry,
		IsActive:        r.IsActive,
		LastActiveAt:    r.LastActiveAt,
		CreatedAt:       r.CreatedAt,
	}
}

// // Helper để map một danh sách (List)
//
//	func ToUserSessionRespList(sessions []entity.UserSession) []UserSessionReq {
//		result := make([]UserSessionReq, len(sessions))
//		for i, s := range sessions {
//			result[i] = *ToUserSessionResp(&s)
//		}
//		return result
//	}
func (req *UserSessionReq) ToEntity() *entity.UserSession {
	if req == nil {
		return nil
	}
	return &entity.UserSession{
		// ID: Postgres tự sinh (gen_random_uuid)
		UserID:          req.UserID,
		DeviceName:      req.DeviceName,
		OsVersion:       req.OsVersion,
		Browser:         req.Browser,
		IpAddress:       req.IpAddress,       // Lấy từ Context request
		LocationCity:    req.LocationCity,    // Lấy từ GeoIP service
		LocationCountry: req.LocationCountry, // Lấy từ GeoIP service
		IsActive:        true,
		LastActiveAt:    time.Now(),
	}
}
