package req

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"

// CreateGroupJoinQuestionReq - Ánh xạ đầy đủ các trường cần thiết để tạo câu hỏi
type CreateGroupJoinQuestionReq struct {
	GroupID    string                   `json:"group_id" binding:"required"`
	Content    string                   `json:"content" binding:"required"`
	Type       sharedEnums.QuestionType `json:"type" binding:"required"`
	Options    []ReqQuestionOption      `json:"options,omitempty"` // Tùy chọn, dùng cho trắc nghiệm
	IsRequired bool                     `json:"is_required"`
	Order      int                      `json:"order" binding:"required,min=1"`
	// Bỏ qua ID, CreatedAt, UpdatedAt vì hệ thống sẽ tự động sinh (Best Practice)
}

// UpdateGroupJoinQuestionReq - Dùng con trỏ 100% để hỗ trợ Partial Update
type UpdateGroupJoinQuestionReq struct {
	Content    *string                   `json:"content,omitempty"`
	Type       *sharedEnums.QuestionType `json:"type,omitempty"`
	Options    *[]ReqQuestionOption      `json:"options,omitempty"`
	IsRequired *bool                     `json:"is_required,omitempty"`
	Order      *int                      `json:"order,omitempty"`
	// GroupID thường không được phép thay đổi sau khi tạo nên không đưa vào đây.
}

// --- Nested Structs ---

type ReqQuestionOption struct {
	Text  string `json:"text" binding:"required"`
	Value string `json:"value" binding:"required"`
}
