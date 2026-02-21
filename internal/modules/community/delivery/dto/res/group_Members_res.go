package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
)

// GroupMemberRes - Phản hồi đầy đủ 100% các trường của Entity ra Client
type GroupMemberRes struct {
	ID             string             `json:"id"`
	GroupID        string             `json:"group_id"`
	UserID         string             `json:"user_id"`
	Role           enum.MemberRole    `json:"role"`
	Status         enum.MemberStatus  `json:"status"`
	InviterID      string             `json:"inviter_id,omitempty"`
	JoinAnswers    []ResJoinAnswer    `json:"join_answers,omitempty"`
	DisciplineInfo *ResDisciplineInfo `json:"discipline_info,omitempty"`
	Badges         []enum.MemberBadge `json:"badges,omitempty"`
	JoinedAt       time.Time          `json:"joined_at"`
	LastActiveAt   time.Time          `json:"last_active_at"`
}

// --- Nested Structs ---

type ResJoinAnswer struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type ResDisciplineInfo struct {
	Reason    string     `json:"reason"`
	BannedBy  string     `json:"banned_by"`
	UntilDate *time.Time `json:"until_date,omitempty"`
}