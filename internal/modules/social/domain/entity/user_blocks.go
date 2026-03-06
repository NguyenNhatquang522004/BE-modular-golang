package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBlock struct {
	ID             uuid.UUID                `gorm:"type:uuid;primaryKey;not null;"`                                        // ID của mối quan hệ bạn bè
	Blocker_UserID uuid.UUID                `gorm:"type:uuid;index:idx_user_block,unique;not null" json:"blocker_user_id"` // Người chặn
	Blocked_UserID uuid.UUID                `gorm:"type:uuid;index:idx_user_block,unique;not null" json:"blocked_user_id"` // Người bị chặn                                    // Lý do chặn (tùy chọn)
	Type_Block     []sharedEnums.Type_Block `gorm:"type:varchar(20);default:'full';" json:"type"`                          // 'full' (chặn hoàn toàn), 'partial' (chặn một phần)
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
