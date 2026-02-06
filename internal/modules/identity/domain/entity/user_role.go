package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
	"github.com/google/uuid"
)

type UserRole struct {
	// 1. Sửa uint thành uuid.UUID để khớp với type:uuid
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// 2. UserID và Role tạo thành cặp Unique để tránh trùng lặp
	// index:idx_user_role_unique,unique -> Đảm bảo 1 user chỉ có 1 role "admin" duy nhất
	UserID uuid.UUID     `gorm:"type:uuid;not null;index:idx_user_role_unique,unique"`
	Role   enum.RoleType `gorm:"type:varchar(20);not null;index:idx_user_role_unique,unique"`

	// 3. Foreign Key Relation (Belongs To)
	// OnDelete:CASCADE -> Rất quan trọng: Xóa User là xóa luôn quyền, tránh rác DB
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}
