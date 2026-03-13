package communityEvent

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

type GroupStatsPayload struct {
	GroupID            string              `json:"group_id"`
	MemberCount        int                 `json:"member_count"`
	PostCount          int                 `json:"post_count"`
	PendingMemberCount int                 `json:"pending_member_count"`
	PendingPostCount   int                 `json:"pending_post_count"`
	ReportedPostCount  int                 `json:"reported_post_count"`
	EventType          constants.EventType `json:"event_type"`
}

type DownloadGroupFilePayload struct {
	GroupID   string              `json:"group_id"`
	FileID    string              `json:"file_id"`
	Count     int                 `json:"count"`
	UserID    string              `json:"user_id"`
	EventType constants.EventType `json:"event_type"`
}
type DeleteCommunityRelationTargetPayload struct {
	TargetID string `json:"target_id"`
}

type CreateGroupPayload struct {
	GroupID string `json:"group_id"` // Dùng để update hoặc tracking, có thể là UUID hoặc ObjectID dưới dạng string
	// BẮT BUỘC
	Name       string                   `json:"name" validate:"required,min=3,max=100"`
	Privacy    sharedEnums.PrivacyScope `json:"privacy" validate:"required"`
	CategoryID string                   `json:"category_id" validate:"required,uuid4"` // Giả sử category_id là UUID

	// KHÔNG BẮT BUỘC (Có validation độ dài)
	Description string   `json:"description" validate:"omitempty,max=1000"`
	Tags        []string `json:"tags" validate:"omitempty,dive,max=30"` // dive: validate từng phần tử trong slice

	// MEDIA
	Cover  *CreateGroupCoverPayload  `json:"cover,omitempty" validate:"omitempty"`
	Avatar *CreateGroupAvatarPayload `json:"avatar,omitempty" validate:"omitempty"`

	// NỘI QUY & CÀI ĐẶT
	Rules    []CreateGroupRulePayload    `json:"rules,omitempty" validate:"omitempty,dive"`
	Settings *CreateGroupSettingsPayload `json:"settings,omitempty" validate:"omitempty"`

	// CẤU HÌNH TÍNH NĂNG
	CommunityChats      *CreateFeatureFlagPayload `json:"community_chats,omitempty" validate:"omitempty"`
	MembershipQuestions *CreateFeatureFlagPayload `json:"membership_questions,omitempty" validate:"omitempty"`
}

// --- Các struct phụ trợ cho Create ---

type CreateGroupCoverPayload struct {
	URL       string  `json:"url" validate:"required,url"`
	PositionY float64 `json:"position_y" validate:"min=0,max=100"`
}

type CreateGroupAvatarPayload struct {
	URL string `json:"url" validate:"required,url"`
}

type CreateGroupRulePayload struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required,max=2000"`
}

type CreateGroupSettingsPayload struct {
	RequireApprovalToJoin bool                    `json:"require_approval_to_join"`
	RequireApprovalToPost bool                    `json:"require_approval_to_post"`
	AllowMemberPosting    bool                    `json:"allow_member_posting"`
	WhoCanApproveMember   []*sharedEnums.RoleType `json:"who_can_approve_member" validate:"omitempty,dive"`
}

type CreateFeatureFlagPayload struct {
	IsEnabled bool `json:"is_enabled"`
}
type UpdateGroupPayload struct {
	GroupID     string                    `json:"group_id"` // Dùng để update hoặc tracking, có thể là UUID hoặc ObjectID dưới dạng string
	Name        *string                   `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string                   `json:"description,omitempty" validate:"omitempty,max=1000"`
	Tags        *[]string                 `json:"tags,omitempty" validate:"omitempty,dive,max=30"`
	Privacy     *sharedEnums.PrivacyScope `json:"privacy,omitempty" validate:"omitempty"`
	CategoryID  *string                   `json:"category_id,omitempty" validate:"omitempty,uuid4"`

	Cover  *UpdateGroupCoverPayload  `json:"cover,omitempty" validate:"omitempty"`
	Avatar *UpdateGroupAvatarPayload `json:"avatar,omitempty" validate:"omitempty"`

	Rules    *[]UpdateGroupRulePayload   `json:"rules,omitempty" validate:"omitempty,dive"`
	Settings *UpdateGroupSettingsPayload `json:"settings,omitempty" validate:"omitempty"`

	CommunityChats      *UpdateFeatureFlagPayload `json:"community_chats,omitempty" validate:"omitempty"`
	MembershipQuestions *UpdateFeatureFlagPayload `json:"membership_questions,omitempty" validate:"omitempty"`
}

// --- Các struct phụ trợ cho Update (Dùng Pointer) ---

type UpdateGroupCoverPayload struct {
	URL       *string  `json:"url,omitempty" validate:"omitempty,url"`
	PositionY *float64 `json:"position_y,omitempty" validate:"omitempty,min=0,max=100"`
}

type UpdateGroupAvatarPayload struct {
	URL *string `json:"url,omitempty" validate:"omitempty,url"`
}

type UpdateGroupRulePayload struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,max=255"`
	Content *string `json:"content,omitempty" validate:"omitempty,max=2000"`
}

type UpdateGroupSettingsPayload struct {
	RequireApprovalToJoin *bool                    `json:"require_approval_to_join,omitempty"`
	RequireApprovalToPost *bool                    `json:"require_approval_to_post,omitempty"`
	AllowMemberPosting    *bool                    `json:"allow_member_posting,omitempty"`
	WhoCanApproveMember   *[]*sharedEnums.RoleType `json:"who_can_approve_member,omitempty" validate:"omitempty,dive"`
}

type UpdateFeatureFlagPayload struct {
	IsEnabled *bool `json:"is_enabled,omitempty"`
}
