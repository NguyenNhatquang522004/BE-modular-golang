// File: mapper/conversation_read_state_mapper.go
package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
)

// =========================================================================
// MAPPER REQ
// =========================================================================

// ToEntity tạo mới Entity từ Request.
func ToEntityConversationReadState(r *req.ConversationReadStateReq) *entity.ConversationReadState {
	if r == nil {
		return nil
	}
	return &entity.ConversationReadState{
		ConversationID:    r.ConversationID,
		UserID:            r.UserID,
		LastReadMessageID: r.LastReadMessageID,
		LastReadAt:        r.LastReadAt,
	}
}

// UpdateToEntity cập nhật Entity có sẵn từ Request.
func UpdateToEntityConversationReadState(r *req.ConversationReadStateReq, e *entity.ConversationReadState) {
	if r == nil || e == nil {
		return
	}
	e.ConversationID = r.ConversationID
	e.UserID = r.UserID
	e.LastReadMessageID = r.LastReadMessageID
	e.LastReadAt = r.LastReadAt
}

// =========================================================================
// MAPPER RES
// =========================================================================

// ReqToRes map trực tiếp từ Request sang Response.
func ToResFromReqConversationReadState(r *req.ConversationReadStateReq) *res.ConversationReadStateRes {
	if r == nil {
		return nil
	}
	return &res.ConversationReadStateRes{
		ConversationID:    r.ConversationID,
		UserID:            r.UserID,
		LastReadMessageID: r.LastReadMessageID,
		LastReadAt:        r.LastReadAt,
	}
}

// EntityToRes map từ Entity sang Response (Dùng sau khi query Cassandra).
func ToResFromEntityConversationReadState(e *entity.ConversationReadState) *res.ConversationReadStateRes {
	if e == nil {
		return nil
	}
	return &res.ConversationReadStateRes{
		ConversationID:    e.ConversationID,
		UserID:            e.UserID,
		LastReadMessageID: e.LastReadMessageID,
		LastReadAt:        e.LastReadAt,
	}
}
