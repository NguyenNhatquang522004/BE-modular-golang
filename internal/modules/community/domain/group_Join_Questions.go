package domain

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =============================================================================
// SUB-STRUCT: OPTION
// =============================================================================

type QuestionOption struct {
	Text  string `bson:"text" json:"text"`   // VD: "Dưới 18 tuổi"
	Value string `bson:"value" json:"value"` // VD: "under_18"
}

// =============================================================================
// ROOT ENTITY: GROUP JOIN QUESTION
// =============================================================================

type GroupJoinQuestion struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 1. LINKING
	// Index: Compound { group_id: 1, order: 1 } -> Lấy list câu hỏi của nhóm theo thứ tự
	GroupID primitive.ObjectID `bson:"group_id" json:"group_id"`

	// 2. CONTENT
	Content string            `bson:"content" json:"content"` // "Bạn bao nhiêu tuổi?"
	Type    enum.QuestionType `bson:"type" json:"type"`       // 'text', 'checkbox'...

	// Danh sách lựa chọn (Chỉ dùng cho MultipleChoice/Checkbox)
	// Dùng omitempty để tiết kiệm nếu là câu hỏi Text
	Options []QuestionOption `bson:"options,omitempty" json:"options,omitempty"`

	// 3. SETTINGS
	IsRequired bool `bson:"is_required" json:"is_required"`
	Order      int  `bson:"order" json:"order"` // 1, 2, 3...

	// 4. TIMESTAMPS
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
