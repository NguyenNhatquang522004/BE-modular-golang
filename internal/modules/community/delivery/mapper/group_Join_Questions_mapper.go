package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------------------------------------------------
// MAPPER REQ TO ENTITY
// ---------------------------------------------------------

// ToEntityGroupJoinQuestion: Ánh xạ 100% từ CreateGroupJoinQuestionReq sang Entity GroupJoinQuestion
func ToEntityGroupJoinQuestion(r *req.CreateGroupJoinQuestionReq) *entity.GroupJoinQuestion {
	if r == nil {
		return nil
	}

	groupID, _ := primitive.ObjectIDFromHex(r.GroupID) // Cần validate ở handler trước khi vào mapper
	now := time.Now()

	e := &entity.GroupJoinQuestion{
		ID:         primitive.NewObjectID(),
		GroupID:    groupID,
		Content:    r.Content,
		Type:       r.Type,
		IsRequired: r.IsRequired,
		Order:      r.Order,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Ánh xạ Option
	if len(r.Options) > 0 {
		e.Options = make([]entity.QuestionOption, len(r.Options))
		for i, opt := range r.Options {
			e.Options[i] = entity.QuestionOption{
				Text:  opt.Text,
				Value: opt.Value,
			}
		}
	}

	return e
}

// UpdateToEntityGroupJoinQuestion: Cập nhật đè (Partial Update) từ UpdateGroupJoinQuestionReq vào Entity có sẵn
func UpdateToEntityGroupJoinQuestion(r *req.UpdateGroupJoinQuestionReq, e *entity.GroupJoinQuestion) {
	if r == nil || e == nil {
		return
	}

	if r.Content != nil {
		e.Content = *r.Content
	}
	if r.Type != nil {
		e.Type = *r.Type
	}
	if r.IsRequired != nil {
		e.IsRequired = *r.IsRequired
	}
	if r.Order != nil {
		e.Order = *r.Order
	}

	// Cập nhật lại toàn bộ mảng Option nếu có truyền lên
	if r.Options != nil {
		e.Options = make([]entity.QuestionOption, len(*r.Options))
		for i, opt := range *r.Options {
			e.Options[i] = entity.QuestionOption{
				Text:  opt.Text,
				Value: opt.Value,
			}
		}
	}

	// Tự động cập nhật thời gian sửa
	e.UpdatedAt = time.Now()
}

// ---------------------------------------------------------
// MAPPER ENTITY TO RES
// ---------------------------------------------------------

// ToResGroupJoinQuestion: Ánh xạ 100% từ Entity ra DTO Response
func ToResGroupJoinQuestion(e *entity.GroupJoinQuestion) *res.GroupJoinQuestionRes {
	if e == nil {
		return nil
	}

	response := &res.GroupJoinQuestionRes{
		ID:         e.ID.Hex(),
		GroupID:    e.GroupID.Hex(),
		Content:    e.Content,
		Type:       e.Type,
		IsRequired: e.IsRequired,
		Order:      e.Order,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}

	// Ánh xạ Option
	if len(e.Options) > 0 {
		response.Options = make([]res.ResQuestionOption, len(e.Options))
		for i, opt := range e.Options {
			response.Options[i] = res.ResQuestionOption{
				Text:  opt.Text,
				Value: opt.Value,
			}
		}
	}

	return response
}
