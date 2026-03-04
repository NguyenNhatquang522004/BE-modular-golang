package res

type FailedPageResponse struct {
	PageID       string `json:"page_id"`
	UserActionID string `json:"user_action_id"` // ID của hành động người dùng, dùng để trackingsss
	ErrorMessage string `json:"error_message"`
}

type FailedPageRoleResponse struct {
	PageID       string `json:"page_id"`
	UserID       string `json:"user_id"`
	ErrorMessage string `json:"error_message"`
}

type FailedPageFollowerResponse struct {
	PageID       string `json:"page_id"`
	UserID       string `json:"user_id"`
	ErrorMessage string `json:"error_message"`
}
type FailedAdAccountResponse struct {
	AccountID    string `json:"account_id"`
	UserActionID string `json:"user_action_id"` // ID của hành động người dùng, dùng để tracking
	ErrorMessage string `json:"error_message"`
}

type FailedAdCampainResponse struct {
	CampaignID   string `json:"campaign_id"`
	UserActionID string `json:"user_action_id"` // ID của hành động người dùng, dùng để tracking
	AccountID    string `json:"account_id"`
	ErrorMessage string `json:"error_message"`
}
type FailedAdResponse struct {
	AdID         string `json:"ad_id"`
	UserActionID string `json:"user_action_id"` // ID của hành động người dùng, dùng để tracking
	CampaignID   string `json:"campaign_id"`
	ErrorMessage string `json:"error_message"`
}
