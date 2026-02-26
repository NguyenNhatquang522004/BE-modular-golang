package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/google/uuid"
)

type UserRole struct {
	// 1. Sửa uint thành uuid.UUID để khớp với type:uuid
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// 2. UserID và Role tạo thành cặp Unique để tránh trùng lặp
	// index:idx_user_role_unique,unique -> Đảm bảo 1 user chỉ có 1 role "admin" duy nhất
	Role        sharedEnums.RoleType `gorm:"type:varchar(20);not null;index:idx_user_role_unique,unique"`
	Description string               `gorm:"type:varchar(255);"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}
