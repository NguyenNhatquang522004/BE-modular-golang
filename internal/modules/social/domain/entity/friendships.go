package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendships struct {
	ID           uuid.UUID             `gorm:"type:uuid;primaryKey;not null;"` // ID của mối quan hệ bạn bè
	Requester_ID uuid.UUID             `gorm:"type:uuid;index:idx_friendship,unique;not null"`
	Recipient_ID uuid.UUID             `gorm:"type:uuid;index:idx_friendship,unique;not null"`
	Status       enum.StatusFriendship `gorm:"type:varchar(20);"` // 'pending', 'accepted', 'blocked'
	Created_At   time.Time
	Updated_At   time.Time
	Deleted_At   gorm.DeletedAt `gorm:"index"`

}
