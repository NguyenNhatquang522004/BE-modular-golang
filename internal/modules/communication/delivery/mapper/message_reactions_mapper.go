// File: mapper/message_reaction_mapper.go
package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

// ToEntity tạo mới Entity từ Request.
func ToEntityMessageReaction(r *req.MessageReactionReq) *entity.MessageReaction {
	if r == nil {
		return nil
	}
	return &entity.MessageReaction{
		ConversationID: r.ConversationID,
		MessageID:      r.MessageID,
		UserID:         r.UserID,
		ReactionCode:   r.ReactionCode,
		CreatedAt:      r.CreatedAt,
	}
}

// UpdateToEntity cập nhật Entity có sẵn từ Request.
func UpdateToEntityMessageReaction(r *req.MessageReactionReq, e *entity.MessageReaction) {
	if r == nil || e == nil {
		return
	}
	e.ConversationID = r.ConversationID
	e.MessageID = r.MessageID
	e.UserID = r.UserID
	e.ReactionCode = r.ReactionCode
	e.CreatedAt = r.CreatedAt
}

// =========================================================================
// MAPPER RES
// =========================================================================

// ReqToRes map trực tiếp từ Request sang Response.
func ToResFromReqMessageReaction(r *req.MessageReactionReq) *res.MessageReactionRes {
	if r == nil {
		return nil
	}
	return &res.MessageReactionRes{
		ConversationID: r.ConversationID,
		MessageID:      r.MessageID,
		UserID:         r.UserID,
		ReactionCode:   r.ReactionCode,
		CreatedAt:      r.CreatedAt,
	}
}

// EntityToRes map từ Entity sang Response (Best practice khi lấy dữ liệu từ DB).
func ToResFromEntityMessageReaction(e *entity.MessageReaction) *res.MessageReactionRes {
	if e == nil {
		return nil
	}
	return &res.MessageReactionRes{
		ConversationID: e.ConversationID,
		MessageID:      e.MessageID,
		UserID:         e.UserID,
		ReactionCode:   e.ReactionCode,
		CreatedAt:      e.CreatedAt,
	}
}
