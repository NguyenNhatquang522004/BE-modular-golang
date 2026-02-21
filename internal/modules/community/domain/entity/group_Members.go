package entity

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionGroupMembers = "GroupMembers"
)

type GroupMember struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. LINKING
	// Index: Compound { group_id: 1, user_id: 1 } (Unique) -> Đảm bảo user không join trùng
	GroupID primitive.ObjectID `bson:"group_id" json:"group_id"`

	// Index: { user_id: 1, joined_at: -1 } -> List các nhóm tôi đã tham gia
	UserID string `bson:"user_id" json:"user_id"`

	// 2. VAI TRÒ & TRẠNG THÁI
	Role   enum.MemberRole   `bson:"role" json:"role"`     // 'admin', 'member'
	Status enum.MemberStatus `bson:"status" json:"status"` // 'active', 'pending', 'invited'

	// 3. THÔNG TIN GIA NHẬP
	// Người mời (Postgres UUID -> String)
	InviterID string `bson:"inviter_id,omitempty" json:"inviter_id,omitempty"`

	// Câu trả lời khi xin vào nhóm (Chỉ có khi status = pending/active)
	JoinAnswers []JoinAnswer `bson:"join_answers,omitempty" json:"join_answers,omitempty"`

	// 4. KỶ LUẬT (Ban/Mute)
	// Pointer để null khi user trạng thái bình thường
	DisciplineInfo *DisciplineInfo `bson:"discipline_info,omitempty" json:"discipline_info,omitempty"`

	// 5. GAMIFICATION
	// Danh sách huy hiệu
	Badges []enum.MemberBadge `bson:"badges,omitempty" json:"badges,omitempty"`

	// 6. METADATA
	JoinedAt time.Time `bson:"joined_at" json:"joined_at"`

	// Index: { group_id: 1, last_active_at: 1 } -> Lọc thành viên "Tàu ngầm" (ít tương tác)
	LastActiveAt time.Time `bson:"last_active_at" json:"last_active_at"`

	// 8. TIMESTAMPS
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	// Soft Delete
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// --- JOIN ANSWERS ---
type JoinAnswer struct {
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
	Answer     string             `bson:"answer" json:"answer"`
}

// --- DISCIPLINE INFO (Kỷ luật) ---
type DisciplineInfo struct {
	Reason string `bson:"reason" json:"reason"`

	// UserID người thực hiện ban (Admin/Mod) -> Postgres UUID
	BannedBy string `bson:"banned_by" json:"banned_by"`

	// Thời điểm hết hạn. Null = Vĩnh viễn.
	UntilDate *time.Time `bson:"until_date,omitempty" json:"until_date,omitempty"`
}

func (GroupMember) CollectionName() string {
	return CollectionGroupMembers
}
