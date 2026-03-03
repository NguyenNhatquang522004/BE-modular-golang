package req

type CreateAdAccountRequest struct {
	*PageReq
}
type UpdateAdAccountRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	*PageReq
}
type DeletePageRequest struct {
	UserActionID string `json:"user_action_id" validate:"required"` // ID của hành động người dùng, dùng để tracking
	PageID       string `json:"page_id" validate:"required"`
}
