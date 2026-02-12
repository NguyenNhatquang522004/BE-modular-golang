package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Followers struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;not null;"` // ID của mối quan hệ bạn bè
	Follower_UserID uuid.UUID `gorm:"type:uuid;index:idx_follower,unique;not null" json:"follower_user_id"`
	Followed_UserID uuid.UUID `gorm:"type:uuid;index:idx_follower,unique;not null" json:"followed_user_id"`
	IsMuted         bool      `gorm:"default:false" json:"is_muted"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	// Người đi theo dõi (Fan)
	// FollowerUser *entity.User `gorm:"foreignKey:Follower_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// // Người được theo dõi (Idol)
	// FollowedUser *entity.User `gorm:"foreignKey:Followed_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
