package domain

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email       string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PhoneNumber string    `gorm:"type:varchar(20);uniqueIndex;"`
	Password    string    `gorm:"type:varchar(255);not null"`
	Username    string    `gorm:"type:varchar(50);uniqueIndex;"`

	TypeLogin enum.LoginType `gorm:"type:varchar(20);"`

	OTPVerified string    `gorm:"type:varchar(20);"`
	OTPCode     string    `gorm:"type:varchar(10);"`
	OTPExpiry   time.Time `gorm:""`
	OTPAttempts int       `gorm:"default:0;"`

	IsActive             bool `gorm:"default:true;"`
	SettingsActiveStatus bool `gorm:"default:true;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Sessions []*UserSession `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
