package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
)

// CreateGroupReq - Dùng cho API tạo mới (Tất cả các trường cần thiết theo Entity)
type CreateGroupReq struct {
	CreatorID           string            `json:"creator_id" binding:"required"` // Thường lấy từ token, nhưng ánh xạ 100% theo yêu cầu
	Name                string            `json:"name" binding:"required,max=255"`
	Slug                string            `json:"slug" binding:"required"`
	Description         string            `json:"description"`
	Tags                []string          `json:"tags"`
	Cover               ReqGroupCover     `json:"cover"`
	Avatar              ReqGroupAvatar    `json:"avatar"`
	Privacy             enum.GroupPrivacy `json:"privacy" binding:"required"`
	CategoryID          string            `json:"category_id" binding:"required"`
	Rules               []ReqGroupRule    `json:"rules"`
	Settings            ReqGroupSettings  `json:"settings"`
	CommunityChats      ReqFeatureFlag    `json:"community_chats"`
	MembershipQuestions ReqFeatureFlag    `json:"membership_questions"`
	// Stats, CreatedAt, UpdatedAt, DeletedAt không đưa vào Create Req vì đây là các trường do Server/DB tự quản lý (Best Practice).
}

// UpdateGroupReq - Dùng cho API cập nhật (Dùng con trỏ để hỗ trợ Partial Update)
type UpdateGroupReq struct {
	Name                *string            `json:"name,omitempty"`
	Slug                *string            `json:"slug,omitempty"`
	Description         *string            `json:"description,omitempty"`
	Tags                *[]string          `json:"tags,omitempty"`
	Cover               *ReqGroupCover     `json:"cover,omitempty"`
	Avatar              *ReqGroupAvatar    `json:"avatar,omitempty"`
	Privacy             *enum.GroupPrivacy `json:"privacy,omitempty"`
	CategoryID          *string            `json:"category_id,omitempty"`
	Rules               *[]ReqGroupRule    `json:"rules,omitempty"`
	Settings            *ReqGroupSettings  `json:"settings,omitempty"`
	CommunityChats      *ReqFeatureFlag    `json:"community_chats,omitempty"`
	MembershipQuestions *ReqFeatureFlag    `json:"membership_questions,omitempty"`
}

// --- Nested Structs cho Request ---

type ReqGroupCover struct {
	URL       string  `json:"url" binding:"omitempty,url"`
	PositionY float64 `json:"position_y"`
}

type ReqGroupAvatar struct {
	URL string `json:"url" binding:"omitempty,url"`
}

type ReqGroupRule struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type ReqGroupSettings struct {
	RequireApprovalToJoin bool                    `json:"require_approval_to_join"`
	RequireApprovalToPost bool                    `json:"require_approval_to_post"`
	AllowMemberPosting    bool                    `json:"allow_member_posting"`
	WhoCanApproveMember   []*sharedEnums.RoleType `json:"who_can_approve_member"`
}

type ReqFeatureFlag struct {
	IsEnabled bool `json:"is_enabled"`
}
