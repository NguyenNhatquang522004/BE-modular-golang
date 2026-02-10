package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/google/uuid"
)

type GetUserByIDReq struct {
	UserID string `json:"user_id" binding:"required"`
}

type GetUserByEmailReq struct {
	Email string `json:"email" binding:"required,email"`
}

type PanigationUsersReq struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit" binding:"required,min=1,max=100"`
}

type UpdateUserReq struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Username    string    `json:"username" binding:"omitempty,min=3,max=50"`
	Email       string    `json:"email" binding:"omitempty,email"`
	PhoneNumber string    `json:"phone_number" binding:"omitempty,min=10,max=15"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

func (r *UpdateUserReq) ToEntity() *entity.User {
	return &entity.User{
		ID:          r.ID,
		Username:    r.Username,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
		IsActive:    r.IsActive,
	}
}

type CreateUserReq struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,min=10,max=15"`
	Password    string `json:"password" binding:"required,min=6"`
}

func (r *CreateUserReq) ToEntity() *entity.User {
	return &entity.User{
		Username:    r.Username,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
	}
}
