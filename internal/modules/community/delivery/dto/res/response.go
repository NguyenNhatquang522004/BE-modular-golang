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
type FailedGroupQA struct {
	GroupID      string `json:"group_id"`
	ErrorMessage error  `json:"error_message"`
}

type FailedGroupEvent struct {
	GroupID      string `json:"group_id"`
	EventID      string `json:"event_id"`
	ErrorMessage error  `json:"error_message"`
}

type FailGroupFile struct {
	GroupID      string `json:"group_id"`
	UserID       string `json:"user_id"`
	UserActionID string `json:"user_action_id"`
	FileID       string `json:"file_id"`
	Data		 any `json:"data,omitempty"`
	ErrorMessage error  `json:"error_message"`
}
