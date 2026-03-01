package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// CreateGroupMemberReq - Ánh xạ 100% các trường đầu vào cần thiết
type CreateGroupMemberReq struct {
	GroupID        string                       `json:"group_id" binding:"required"`
	UserID         string                       `json:"user_id" binding:"required"`
	Role           sharedEnums.RoleType         `json:"role" binding:"required"`
	Status         sharedEnums.ProcessingStatus `json:"status" binding:"required"`
	InviterID      string                       `json:"inviter_id,omitempty"`
	JoinAnswers    []ReqJoinAnswer              `json:"join_answers,omitempty"`
	DisciplineInfo *ReqDisciplineInfo           `json:"discipline_info,omitempty"` // Sử dụng pointer vì có thể null
	Badges         []sharedEnums.UserBadge      `json:"badges,omitempty"`
}

// UpdateGroupMemberReq - Dùng pointer 100% để hỗ trợ Partial Update
type UpdateGroupMemberReq struct {
	GroupID        *string                       `json:"group_id,omitempty"`
	UserID         *string                       `json:"user_id,omitempty"`
	Role           *sharedEnums.RoleType         `json:"role,omitempty"`
	Status         *sharedEnums.ProcessingStatus `json:"status,omitempty"`
	InviterID      *string                       `json:"inviter_id,omitempty"`
	JoinAnswers    *[]ReqJoinAnswer              `json:"join_answers,omitempty"`
	DisciplineInfo *ReqDisciplineInfo            `json:"discipline_info,omitempty"`
	Badges         *[]sharedEnums.UserBadge      `json:"badges,omitempty"`
}

// --- Nested Structs ---

type ReqJoinAnswer struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

type ReqDisciplineInfo struct {
	Reason    string     `json:"reason" binding:"required"`
	BannedBy  string     `json:"banned_by" binding:"required"`
	UntilDate *time.Time `json:"until_date,omitempty"` // Pointer vì có thể bị ban vĩnh viễn (null)
}
