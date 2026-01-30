package domain

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Followers struct {
	Follower_ID     uuid.UUID     `gorm:"type:uuid;primaryKey;"`
	Follower_UserID uuid.UUID     `gorm:"type:uuid;primaryKey" json:"follower_user_id"`
	Followed_UserID uuid.UUID     `gorm:"type:uuid;primaryKey" json:"followed_user_id"`
	Is_Muted        enum.Is_Muted `gorm:"type:varchar(10);default:'no';"` // 'yes', 'no'
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	// Người đi theo dõi (Fan)
	FollowerUser *domain.User `gorm:"foreignKey:Follower_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Người được theo dõi (Idol)
	FollowedUser *domain.User `gorm:"foreignKey:Followed_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
