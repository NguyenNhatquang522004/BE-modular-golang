package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionGroups = "Groups"
)

type Group struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// User ID người tạo (Postgres UUID -> String)
	// Index: { creator_id: 1 } -> Tìm "Nhóm tôi đã tạo"
	CreatorID string `bson:"creator_id" json:"creator_id"`

	// 1. ĐỊNH DANH & SEO
	// Index: Text Index { name: "text", description: "text", tags: "text" } -> Tìm kiếm nhóm
	Name string `bson:"name" json:"name"`

	// Index: Unique { slug: 1 } -> URL thân thiện
	Slug string `bson:"slug" json:"slug"`

	Description string   `bson:"description" json:"description"`
	Tags        []string `bson:"tags,omitempty" json:"tags,omitempty"`

	// 2. MEDIA
	Cover  GroupCover  `bson:"cover" json:"cover"`
	Avatar GroupAvatar `bson:"avatar" json:"avatar"`

	// 3. PHÂN LOẠI
	// Index: { privacy: 1, category_id: 1 } -> Filter nhóm
	Privacy    sharedEnums.PrivacyScope `bson:"privacy" json:"privacy"`
	CategoryID string                   `bson:"category_id" json:"category_id"`

	// 4. NỘI QUY (Embedded Array)
	Rules []GroupRule `bson:"rules,omitempty" json:"rules,omitempty"`

	// 5. CÀI ĐẶT QUẢN TRỊ
	Settings GroupSettings `bson:"settings" json:"settings"`

	// 6. CẤU HÌNH TÍNH NĂNG (Feature Flags)
	CommunityChats      FeatureFlag `bson:"community_chats" json:"community_chats"`
	MembershipQuestions FeatureFlag `bson:"membership_questions" json:"membership_questions"`

	// 7. METRICS (Denormalization)
	// Index: { "stats.member_count": -1 } -> Gợi ý "Nhóm phổ biến nhất"
	Stats GroupStats `bson:"stats" json:"stats"`

	// 8. TIMESTAMPS
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	// Soft Delete
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- MEDIA ---
type GroupCover struct {
	ID        primitive.ObjectID `bson:"id" json:"id"` // ID của media asset (PostgreSQL UUID -> String)
	URL       string             `bson:"url" json:"url"`
	PositionY float64            `bson:"position_y" json:"position_y"` // Để căn chỉnh ảnh bìa (0.0 - 100.0)
}

type GroupAvatar struct {
	ID  primitive.ObjectID `bson:"id" json:"id"` // ID của media asset (PostgreSQL UUID -> String)
	URL string             `bson:"url" json:"url"`
}

// --- RULES ---
type GroupRule struct {
	Title   string `bson:"title" json:"title"`
	Content string `bson:"content" json:"content"`
}

// --- SETTINGS ---
type GroupSettings struct {
	RequireApprovalToJoin bool `bson:"require_approval_to_join" json:"require_approval_to_join"`
	RequireApprovalToPost bool `bson:"require_approval_to_post" json:"require_approval_to_post"`
	AllowMemberPosting    bool `bson:"allow_member_posting" json:"allow_member_posting"`

	// Enum: Ai có quyền duyệt thành viên?
	WhoCanApproveMember []*sharedEnums.RoleType `bson:"who_can_approve_member" json:"who_can_approve_member"`
}

// --- FEATURE FLAGS (Cấu hình tính năng) ---
type FeatureFlag struct {
	IsEnabled bool `bson:"is_enabled" json:"is_enabled"`
}

// --- STATS (Metrics) ---
type GroupStats struct {
	MemberCount        int `bson:"member_count" json:"member_count"`
	PostCount          int `bson:"post_count" json:"post_count"`
	PendingMemberCount int `bson:"pending_member_count" json:"pending_member_count"`
	PendingPostCount   int `bson:"pending_post_count" json:"pending_post_count"`
	ReportedPostCount  int `bson:"reported_post_count" json:"reported_post_count"`
}

func (Group) CollectionName() string {
	return CollectionGroups
}
