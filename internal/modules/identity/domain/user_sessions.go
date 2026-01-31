package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	// 1. Primary Key: Nên đặt là ID cho chuẩn GORM, hoặc SessionID nếu thích
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// 2. Foreign Key: Viết hoa UserID
	UserID uuid.UUID `gorm:"type:uuid;not null"`

	// 3. Data Fields: Viết hoa chữ cái đầu (PascalCase)
	DeviceName      string `gorm:"type:varchar(255);"`
	OsVersion       string `gorm:"type:varchar(50);"`
	Browser         string `gorm:"type:varchar(50);"`
	IpAddress       string `gorm:"type:varchar(45);"`
	LocationCity    string `gorm:"type:varchar(100);"`
	LocationCountry string `gorm:"type:varchar(100);"`
	RefreshToken    string `gorm:"type:text;"`

	IsActive     bool `gorm:"default:true;"`
	LastActiveAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Dùng gorm.DeletedAt chuẩn hơn *time.Time

	// 4. Quan hệ Belongs To
	// - Dùng con trỏ *User
	// - foreignKey trỏ vào field 'UserID' ở trên
	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
