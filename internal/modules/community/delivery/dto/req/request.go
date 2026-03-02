package req

type CreateGroupRequest struct {
	*CreateGroupReq
}

type UpdateGroupRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
	*UpdateGroupReq
}
type DeleteGroupRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

type AddMemberGroupRequest struct {
	*CreateGroupMemberReq
}

type RemoveMemberGroupRequest struct {
	GroupID      string `json:"group_id" binding:"required"`
	UserID       string `json:"user_id" binding:"required"`
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
}

type UpdateMemberGroupRequest struct {
	GroupID      string `json:"group_id" binding:"required"`
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
	*UpdateGroupMemberReq
}
type CreateGroupQARequest struct {
	// Các trường cần thiết để tạo Group QA
	*CreateGroupJoinQuestionReq
}
type UpdateGroupQARequest struct {
	// Các trường cần thiết để cập nhật Group QA
	GroupID      string `json:"group_id" binding:"required"`
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
	QAID         string `json:"qa_id" binding:"required"`          // ID của câu hỏi cần cập nhật
	*UpdateGroupJoinQuestionReq
}

type DeleteGroupQARequest struct {
	GroupID      string `json:"group_id" binding:"required"`
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
	QAID         string `json:"qa_id" binding:"required"`          // ID của câu hỏi cần xóa
}
type CreateGroupEventRequest struct {
	// Các trường cần thiết để tạo Group Event
	*CreateGroupEventReq
}

type UpdateGroupEventRequest struct {
	// Các trường cần thiết để cập nhật Group
	EventID string `json:"event_id" binding:"required"`
	*UpdateGroupEventReq
}
type DeleteGroupEventRequest struct {
	GroupID      string `json:"group_id" binding:"required"`
	EventID      string `json:"event_id" binding:"required"`
	CreatorID    string `json:"creator_id" binding:"required"`
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
}

type CreateGroupFileRequest struct {
	*CreateGroupFileReq
}
type DeleteGroupFileRequest struct {
	GroupID      string `json:"group_id" binding:"required"`
	FileID       string `json:"file_id" binding:"required"`
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
}
type DownloadGroupFileRequest struct {
	GroupID      string `json:"group_id" binding:"required"`
	FileID       string `json:"file_id" binding:"required"`
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
}

