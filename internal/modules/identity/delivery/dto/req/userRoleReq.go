package req

type UpdateRoleofUserReq struct {
	RoleID []string `json:"role_id" binding:"required,dive,uuid"`
}
type GetAllRoleReq struct {
	RoleID      string `form:"role_id" binding:"omitempty,uuid"`
	Description string `form:"description" binding:"omitempty,max=255"`
}
