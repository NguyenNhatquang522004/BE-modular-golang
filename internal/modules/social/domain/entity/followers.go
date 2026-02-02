package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Followers struct {
	Follower_UserID uuid.UUID `gorm:"type:uuid;primaryKey" json:"follower_user_id"`
	Followed_UserID uuid.UUID `gorm:"type:uuid;primaryKey" json:"followed_user_id"`
	IsMuted         bool      `gorm:"default:false" json:"is_muted"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	// Người đi theo dõi (Fan)
	FollowerUser *entity.User `gorm:"foreignKey:Follower_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Người được theo dõi (Idol)
	FollowedUser *entity.User `gorm:"foreignKey:Followed_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
