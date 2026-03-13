package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
)

// GroupJoinQuestionRes - Phản hồi 100% dữ liệu từ Entity
type GroupJoinQuestionRes struct {
	ID         string                   `json:"id"`
	GroupID    string                   `json:"group_id"`
	Content    string                   `json:"content"`
	Type       sharedEnums.QuestionType `json:"type"`
	Options    []ResQuestionOption      `json:"options,omitempty"`
	IsRequired bool                     `json:"is_required"`
	Order      int                      `json:"order"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

// --- Nested Structs ---

type ResQuestionOption struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}
