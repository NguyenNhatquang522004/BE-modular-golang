package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User_Blocks struct {
	Blocker_UserID uuid.UUID       `gorm:"type:uuid;primaryKey;not null" json:"blocker_user_id"` // Người chặn
	Blocked_UserID uuid.UUID       `gorm:"type:uuid;primaryKey;not null" json:"blocked_user_id"` // Người bị chặn
	Reason         string          `gorm:"type:text;" json:"reason"`                             // Lý do chặn (tùy chọn)
	Type_Block     enum.Type_Block `gorm:"type:varchar(20);default:'full';" json:"type"`         // 'full' (chặn hoàn toàn), 'partial' (chặn một phần)
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	// Người chặn
	BlockerUser *entity.User `gorm:"foreignKey:Blocker_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Người bị chặn
	BlockedUser *entity.User `gorm:"foreignKey:Blocked_UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
