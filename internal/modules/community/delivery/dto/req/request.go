package req

type CreateGroupRequest struct {
	*CreateGroupReq
}

type UpdateGroupRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"` // Thường lấy từ token, nhưng ánh xạ 100% theo yêu cầu
	*UpdateGroupReq
}
type DeleteGroupRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"` // Thường lấy từ token, nhưng ánh xạ 100% theo yêu cầu
}

type AddMemberGroupRequest struct {
	*CreateGroupMemberReq
}
