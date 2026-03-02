package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"

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
type AcceptGroupJoinRequest struct {
	GroupID      string `json:"group_id" binding:"required"`
	UserID       string `json:"user_id" binding:"required"`        // ID của người được chấp nhận vào nhóm
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
	Accept       bool   `json:"accept"`                            // true nếu chấp nhận, false nếu từ chối
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
	UserActionID string `json:"user_action_id" binding:"required"` // ID của người thực hiện hành động (có thể là admin hoặc chính user đó)
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

type StatsGroupEventRequest struct {
	GroupID         string              `json:"group_id"`
	EventID         string              `json:"event_id"`
	GoingCount      int                 `json:"going_count"`
	InterestedCount int                 `json:"interested_count"`
	EventType       constants.EventType `json:"event_type"`
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
