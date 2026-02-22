package postgresErrors

type BulkError struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}
type AdsBulkError struct {
	ID           string `json:"id"`
	CampaignID   string `json:"campaign_id"`
	TargetPostID string `json:"target_post_id"`
	Reason       string `json:"reason"`
}

type AdCampaignsBulkError struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Reason    string `json:"reason"`
}
type AdAccountsBulkError struct {
	ID          string `json:"id"`
	OwnerUserID string `json:"owner_user_id"`
	Reason      string `json:"reason"`
}
