package res

type FailedGroup struct {
	GroupID      string `json:"group_id"`
	UserID       string `json:"user_id,omitempty"`
	ErrorMessage error  `json:"error_message"`
}

type FailedMember struct {
	GroupID      string `json:"group_id"`
	UserID       string `json:"user_id"`
	ErrorMessage error  `json:"error_message"`
}
