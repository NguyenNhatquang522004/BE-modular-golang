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

type CreateGroupEventPayload struct {
	// GroupID có thể lấy từ URL param (param: group_id) hoặc Body.
	// Nếu lấy từ URL thì bỏ trường này đi để tránh trùng lặp.
	GroupID     string `json:"group_id" validate:"required,mongodb"`
	CreatorID   string `json:"creator_id" validate:"required,uuid"` // Thêm trường creatorID để biết ai là người tạo sự kiện
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"required,max=5000"`
	CoverURL    string `json:"cover_url" validate:"omitempty,url"`

	// StartTime phải ở tương lai, EndTime phải lớn hơn StartTime
	StartTime time.Time `json:"start_time" validate:"required,gt=now"`
	EndTime   time.Time `json:"end_time" validate:"required,gtfield=StartTime"`

	Location EventLocationPayload `json:"location" validate:"required"`
}
type DeleteGroupEventPayload struct {
	EventID   string `json:"event_id" validate:"required,mongodb"`
	GroupID   string `json:"group_id" validate:"required,mongodb"`
	DeleteAll bool   `json:"delete_all"` // Nếu true, xóa tất cả bài viết và tương tác của user trong sự kiện, không chỉ xóa event record
}

// ==============================================================================
// 2. UPDATE PAYLOAD
// Dùng cho API: PATCH /api/v1/events/:event_id
// Dùng con trỏ (*) để hỗ trợ Partial Update (cập nhật 1 phần).
// ==============================================================================

type UpdateGroupEventPayload struct {
	ID          string  `json:"id" validate:"required,mongodb"` // Dùng để update hoặc tracking, có thể là UUID hoặc ObjectID dưới dạng string
	Title       *string `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=5000"`
	CoverURL    *string `json:"cover_url,omitempty" validate:"omitempty,url"`

	// Khi update thời gian, bạn sẽ cần logic code để check StartTime < EndTime ở tầng Service
	StartTime *time.Time `json:"start_time,omitempty" validate:"omitempty,gt=now"`
	EndTime   *time.Time `json:"end_time,omitempty" validate:"omitempty"`

	// Location thường sẽ được update toàn bộ object thay vì từng field nhỏ bên trong
	Location *EventLocationPayload `json:"location,omitempty" validate:"omitempty"`
}

// ==============================================================================
// SHARED PAYLOADS
// ==============================================================================

// EventLocationPayload chứa thông tin địa điểm từ request của user
type EventLocationPayload struct {
	Type sharedEnums.EventLocationType `json:"type" validate:"required"` // Cần custom validator cho enum này nếu cần (vd: oneof=online offline)

	// Có thể custom validate: nếu type là 'offline' thì bắt buộc có Coordinates,
	// nếu 'online' thì Address phải là URL hợp lệ.
	Address string `json:"address" validate:"required"`

	// GeoJSON: [Kinh độ (Longitude), Vĩ độ (Latitude)]
	// Validate mảng phải có đúng 2 phần tử (nếu được truyền lên)
	Coordinates []float64 `json:"coordinates,omitempty" validate:"omitempty,len=2"`
}

// QuestionOptionPayload tách biệt với Entity để gắn tag validate (ví dụ: max length)
type QuestionOptionPayload struct {
	Text  string `json:"text" validate:"required,max=255"`
	Value string `json:"value" validate:"required,max=255"`
}

// =============================================================================
// 1. CREATE PAYLOAD (Dùng cho method POST)
// =============================================================================

type CreateGroupJoinQuestionPayload struct {
	// GroupID: Best practice là truyền qua URL Params (VD: POST /groups/:group_id/questions)
	// Nhưng nếu cấu trúc API của bạn yêu cầu nhận từ Body, ta để kiểu string để validate
	// mã Hex của ObjectID hợp lệ trước khi parse thành primitive.ObjectID ở tầng Service.
	GroupID string `json:"group_id" validate:"required,mongodb"`

	// Validate min/max độ dài để tránh rác DB
	Content string `json:"content" validate:"required,min=3,max=1000"`

	// Validate type là bắt buộc (Nên kết hợp với tag `oneof=text checkbox...` nếu có string cụ thể)
	Type sharedEnums.QuestionType `json:"type" validate:"required"`

	// Sử dụng tag `dive` để yêu cầu validator đi sâu vào từng phần tử trong mảng
	Options []QuestionOptionPayload `json:"options,omitempty" validate:"omitempty,dive"`

	// Dùng con trỏ (*bool, *int) cho các field required để bắt buộc client PHẢI gửi lên giá trị.
	// Nếu không dùng con trỏ, client bỏ trống field này Go sẽ tự gán false/0 và pass qua required.
	IsRequired *bool `json:"is_required" validate:"required"`
	Order      *int  `json:"order" validate:"required,min=1"`
}

// =============================================================================
// 2. UPDATE PAYLOAD (Dùng cho method PATCH / PUT)
// =============================================================================

type UpdateGroupJoinQuestionPayload struct {
	ID string `json:"id" validate:"required,mongodb"` // Dùng để update hoặc tracking, có thể là UUID hoặc ObjectID dưới dạng string
	// LƯU Ý: ID của câu hỏi và GroupID KHÔNG NẰM Ở ĐÂY.
	// Best practice: Lấy từ URL Path (VD: PATCH /groups/:group_id/questions/:question_id)

	// TẤT CẢ các field đều dùng con trỏ (Pointer) để hỗ trợ Partial Update.
	// Lợi ích:
	// - Nếu client không gửi `content` -> JSON unmarshal ra `nil` -> Service bỏ qua không update DB.
	// - Nếu client gửi `content: ""` -> Lỗi validate vì min=3.
	Content *string                   `json:"content" validate:"omitempty,min=3,max=1000"`
	Type    *sharedEnums.QuestionType `json:"type" validate:"omitempty"`

	// Với mảng Options, bản chất Update thường là ghi đè (replace) toàn bộ mảng.
	// Nếu slice là `nil`, tức là client không gửi lên.
	// Nếu slice là `[]` (len=0), tức là client muốn xóa hết options.
	Options []QuestionOptionPayload `json:"options" validate:"omitempty,dive"`

	// Nhờ dùng con trỏ, ta có thể dễ dàng update IsRequired từ true thành false
	// mà không sợ nhầm lẫn với việc client không gửi gì.
	IsRequired *bool `json:"is_required" validate:"omitempty"`
	Order      *int  `json:"order" validate:"omitempty,min=1"`
}

type DeleteGroupJoinQuestionPayload struct {
	ID        string `json:"id" validate:"required,mongodb"` // Dùng để update hoặc tracking, có thể là UUID hoặc ObjectID dưới dạng string
	GroupID   string `json:"group_id" validate:"required,mongodb"`
	DeleteAll bool   `json:"delete_all"` // Nếu true, xóa tất cả bài viết và tương tác của user trong nhóm, không chỉ xóa member record
}
