package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendships struct {
	Requester_ID uuid.UUID             `gorm:"type:uuid;primaryKey;not null;"` // Người gửi lời mời kết bạn
	Recipient_ID uuid.UUID             `gorm:"type:uuid;primaryKey;not null;"` // Người nhận lời mời kết bạn
	Status       enum.StatusFriendship `gorm:"type:varchar(20);"`              // 'pending', 'accepted', 'blocked'
	Created_At   time.Time
	Updated_At   time.Time
	Deleted_At   gorm.DeletedAt `gorm:"index"`
	// Người gửi lời mời kết bạn
	RequesterUser *entity.User `gorm:"foreignKey:Requester_ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Người nhận lời mời kết bạn
	RecipientUser *entity.User `gorm:"foreignKey:Recipient_ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
