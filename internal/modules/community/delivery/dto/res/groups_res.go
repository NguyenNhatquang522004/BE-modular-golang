package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
)

// GroupRes - Trả về 100% dữ liệu của Entity cho Client
type GroupRes struct {
	ID                  string            `json:"id"` // ObjectID chuyển thành string
	CreatorID           string            `json:"creator_id"`
	Name                string            `json:"name"`
	Slug                string            `json:"slug"`
	Description         string            `json:"description"`
	Tags                []string          `json:"tags,omitempty"`
	Cover               ResGroupCover     `json:"cover"`
	Avatar              ResGroupAvatar    `json:"avatar"`
	Privacy             enum.GroupPrivacy `json:"privacy"`
	CategoryID          string            `json:"category_id"`
	Rules               []ResGroupRule    `json:"rules,omitempty"`
	Settings            ResGroupSettings  `json:"settings"`
	CommunityChats      ResFeatureFlag    `json:"community_chats"`
	MembershipQuestions ResFeatureFlag    `json:"membership_questions"`
	Stats               ResGroupStats     `json:"stats"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	DeletedAt           *time.Time        `json:"deleted_at,omitempty"`
}

// --- Nested Structs cho Response ---

type ResGroupCover struct {
	URL       string  `json:"url"`
	PositionY float64 `json:"position_y"`
}

type ResGroupAvatar struct {
	URL string `json:"url"`
}

type ResGroupRule struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ResGroupSettings struct {
	RequireApprovalToJoin bool                    `json:"require_approval_to_join"`
	RequireApprovalToPost bool                    `json:"require_approval_to_post"`
	AllowMemberPosting    bool                    `json:"allow_member_posting"`
	WhoCanApproveMember   []*sharedEnums.RoleType `json:"who_can_approve_member"`
}

type ResFeatureFlag struct {
	IsEnabled bool `json:"is_enabled"`
}

type ResGroupStats struct {
	MemberCount        int `json:"member_count"`
	PostCount          int `json:"post_count"`
	PendingMemberCount int `json:"pending_member_count"`
	PendingPostCount   int `json:"pending_post_count"`
	ReportedPostCount  int `json:"reported_post_count"`
}
