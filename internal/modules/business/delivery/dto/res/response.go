package res

type FailedPageResponse struct {
	PageID       string `json:"page_id"`
	UserActionID string `json:"user_action_id"` // ID của hành động người dùng, dùng để trackingsss
	ErrorMessage string `json:"error_message"`
}
