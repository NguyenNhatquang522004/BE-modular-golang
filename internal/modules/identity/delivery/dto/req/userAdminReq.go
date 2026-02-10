package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"

type RoleReq struct {
	RoleName []enum.RoleType `json:"role_name" binding:"required"`
	RoleID   []string        `json:"role_id" binding:"required"`
}
type IDFusionRoleReq struct {
	UserID string `json:"user_id" binding:"required"`
	RoleReq
}

type DeleteUserReq struct {
	UserID string `json:"user_id" binding:"required"`
}
type AdminCreateUserReq struct {
	CreateUserReq
	RoleReq
}
