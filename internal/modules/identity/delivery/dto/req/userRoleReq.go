package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"

type UpdateRoleofUserReq struct {
	UserID string   `json:"user_id" binding:"required,uuid"`
	RoleID []string `json:"role_id" binding:"required,dive,uuid"`
}
type GetAllRoleReq struct {
	RoleID      string `form:"role_id" binding:"omitempty,uuid"`
	Description string `form:"description" binding:"omitempty,max=255"`
}

type CreateRoleReq struct {
	Role        sharedEnums.RoleType `json:"role" binding:"required,max=100"`
	Description string               `json:"description" binding:"omitempty,max=255"`
}

type DeleteRoleReq struct {
	RoleID string `json:"role_id" binding:"required,uuid"`
}

type GetUserRolesReq struct {
	RoleName string `form:"role_name" binding:"omitempty,max=100"`
}
type AssignRoleToUserWithNameReq struct {
	RoleName string `json:"role_name" binding:"required,max=100"`
}
type AssignRoleToUserWithIDReq struct {
	RoleID string `json:"role_id" binding:"required,uuid"`
}
type UpdateRoleDescriptionReq struct {
	RoleID      string `json:"role_id" binding:"required,uuid"`
	Description string `json:"description" binding:"required,max=255"`
}
