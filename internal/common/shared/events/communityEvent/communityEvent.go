package communityEvent

import (
	"time"

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
	GroupID   string `json:"group_id"`   // Dùng để update hoặc tracking, có thể là UUID hoặc ObjectID dưới dạng
	CreatorID string `json:"creator_id"` // Thêm trường creatorID để biết ai là người tạo nhóm
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
type DeleteGroupPayload struct {
	UserAction string `json:"user_action"` // Enum: "delete" hoặc "leave"
	GroupID    string `json:"group_id"`
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

// =====================================================================
// 1. SUB-PAYLOADS (Dùng chung cho Create/Update nếu cần)
// =====================================================================

type JoinAnswerPayload struct {
	// Dùng string thay vì primitive.ObjectID để tách biệt tầng HTTP và DB
	QuestionID string `json:"question_id" binding:"required,mongodb"`
	Answer     string `json:"answer" binding:"required,max=1000"`
}

type DisciplineInfoPayload struct {
	Reason    string     `json:"reason" binding:"required,max=500"`
	BannedBy  string     `json:"banned_by" binding:"required,uuid"`               // Assuming Postgres UUID format
	UntilDate *time.Time `json:"until_date,omitempty" binding:"omitempty,gt=now"` // Bắt buộc phải là ngày trong tương lai
}

// =====================================================================
// 2. CREATE PAYLOAD (POST /api/v1/groups/{group_id}/members)
// =====================================================================

type CreateGroupMemberPayload struct {
	// Dùng binding:"required" để đảm bảo client bắt buộc phải gửi.
	// Nếu group_id được lấy từ URL params (Path) thì không cần thiết để trong body.
	// Tôi vẫn để ở đây để đảm bảo tính đầy đủ.
	GroupID string                       `json:"group_id" binding:"required,mongodb"`
	UserID  string                       `json:"user_id" binding:"required,uuid"` // Assuming Postgres UUID
	Role    sharedEnums.RoleType         `json:"role" binding:"required"`         // Nên validate custom hoặc oneof='admin' 'member'
	Status  sharedEnums.ProcessingStatus `json:"status" binding:"required"`       // Nên validate custom hoặc oneof='active' 'pending' 'invited'

	// Optional fields (Có thể rỗng)
	InviterID   string              `json:"inviter_id,omitempty" binding:"omitempty,uuid"`
	JoinAnswers []JoinAnswerPayload `json:"join_answers,omitempty" binding:"dive"` // "dive" giúp validate vào sâu từng phần tử trong mảng
}

// =====================================================================
// 3. UPDATE PAYLOAD (PATCH /api/v1/groups/{group_id}/members/{user_id})
// =====================================================================

type UpdateGroupMemberPayload struct {
	GroupID string `json:"group_id" binding:"required,mongodb"` // Dùng để update hoặc tracking, có thể là UUID hoặc ObjectID dưới dạng string
	UserID  string `json:"user_id" binding:"required,uuid"`     // Assuming Postgres UUID
	// BEST PRACTICE: Dùng Pointer (*) cho update payload.
	// Nếu client không gửi trường này (nil), ta bỏ qua.
	// Nếu client gửi, ta lấy giá trị thực để update.
	Role   *sharedEnums.RoleType         `json:"role,omitempty" binding:"omitempty"`
	Status *sharedEnums.ProcessingStatus `json:"status,omitempty" binding:"omitempty"`

	// Kỷ luật: Nếu client muốn ban/mute
	DisciplineInfo *DisciplineInfoPayload `json:"discipline_info,omitempty" binding:"omitempty"`

	// Cập nhật huy hiệu gamification
	Badges []sharedEnums.UserBadge `json:"badges,omitempty" binding:"omitempty"`

	// Chú ý: KHÔNG cho phép update GroupID, UserID vì đây là danh tính core của record.
	// KHÔNG cho update JoinedAt, CreatedAt, UpdatedAt (Hệ thống tự lo).
	// KHÔNG cho update LastActiveAt ở API này (Nên update thông qua middleware hoặc sự kiện hoạt động).
}
type DeleteGroupMemberPayload struct {
	GroupID   string `json:"group_id" binding:"required,mongodb"`
	UserID    string `json:"user_id" binding:"required,uuid"`
	DeleteALL bool   `json:"delete_all"` // Nếu true, xóa tất cả bài viết và tương tác của user trong nhóm, không chỉ xóa member record
}
