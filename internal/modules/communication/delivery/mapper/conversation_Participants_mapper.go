package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ========================
// MAPPER REQUEST -> ENTITY
// ========================

// ToEntityConversationParticipant tạo mới một Entity ConversationParticipant từ Request.
func ToEntityConversationParticipant(r *req.ConversationParticipantReq, conversationId string) *entity.ConversationParticipant {
	if r == nil {
		return nil
	}

	// Xử lý chuyển đổi String sang ObjectID an toàn
	id, _ := primitive.ObjectIDFromHex(r.ID)
	if r.ID == "" {
		id = primitive.NewObjectID() // Auto sinh ID nếu Req không truyền lên
	}

	convID, _ := primitive.ObjectIDFromHex(conversationId)

	// Xử lý các trường thời gian mặc định
	now := time.Now()
	joinedAt := r.JoinedAt
	if joinedAt.IsZero() {
		joinedAt = now
	}
	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := r.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = now
	}

	return &entity.ConversationParticipant{
		ID:                id,
		ConversationID:    convID,
		UserID:            r.UserID,
		Role:              *r.Role,
		Nickname:          r.Nickname,
		LastSeenAt:        r.LastSeenAt,
		LastSeenMessageID: r.LastSeenMessageID,
		MuteUntil:         r.MuteUntil,
		IsArchived:        r.IsArchived,
		ClearHistoryAt:    r.ClearHistoryAt,
		JoinedAt:          joinedAt,
		AddedByUserID:     r.AddedByUserID,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		DeletedAt:         r.DeletedAt,
	}
}
func ToEntityBulkConversationParticipant(r []*req.ConversationParticipantReq, conversationId string) []*entity.ConversationParticipant {
	if r == nil {
		return nil
	}
	var entities []*entity.ConversationParticipant
	for _, item := range r {
		entity := ToEntityConversationParticipant(item, conversationId)
		if entity != nil {
			entities = append(entities, entity)
		}
	}
	return entities
	// Xử lý chuyển đổi String sang ObjectID an toàn

}

// UpdateToEntityConversationParticipant cập nhật Entity hiện có dựa trên data từ Request.
func UpdateToEntityConversationParticipant(r *req.ConversationParticipantReq, e *entity.ConversationParticipant) {
	if r == nil || e == nil {
		return
	}

	// Cập nhật an toàn Khóa ngoại
	if r.ConversationID != "" {
		if convID, err := primitive.ObjectIDFromHex(r.ConversationID); err == nil {
			e.ConversationID = convID
		}
	}

	// Cập nhật các trường thông tin cơ bản
	if r.UserID != "" {
		e.UserID = r.UserID
	}
	if r.Role != nil {
		e.Role = *r.Role
	}
	if r.Nickname != "" {
		e.Nickname = r.Nickname
	}
	if !r.LastSeenAt.IsZero() {
		e.LastSeenAt = r.LastSeenAt
	}
	if r.LastSeenMessageID != "" {
		e.LastSeenMessageID = r.LastSeenMessageID
	}

	// Các trường Cài đặt cá nhân
	e.MuteUntil = r.MuteUntil
	e.IsArchived = r.IsArchived
	e.ClearHistoryAt = r.ClearHistoryAt

	// Cập nhật meta
	if !r.JoinedAt.IsZero() {
		e.JoinedAt = r.JoinedAt
	}
	if r.AddedByUserID != "" {
		e.AddedByUserID = r.AddedByUserID
	}

	e.UpdatedAt = time.Now() // Luôn tự động cập nhật thời gian sửa đổi

	if r.DeletedAt != nil {
		e.DeletedAt = r.DeletedAt
	}
}

// ========================
// MAPPER -> RESPONSE
// ========================
// MAPPER -> RESPONSE
// ========================

// ToResFromReqConversationParticipant chuyển trực tiếp từ DTO Request sang DTO Response.
func ToResFromReqConversationParticipant(r *req.ConversationParticipantReq) *res.ConversationParticipantRes {
	if r == nil {
		return nil
	}

	return &res.ConversationParticipantRes{
		ID:                r.ID,
		ConversationID:    r.ConversationID,
		UserID:            r.UserID,
		Role:              *r.Role,
		Nickname:          r.Nickname,
		LastSeenAt:        r.LastSeenAt,
		LastSeenMessageID: r.LastSeenMessageID,
		MuteUntil:         r.MuteUntil,
		IsArchived:        r.IsArchived,
		ClearHistoryAt:    r.ClearHistoryAt,
		JoinedAt:          r.JoinedAt,
		AddedByUserID:     r.AddedByUserID,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
		DeletedAt:         r.DeletedAt,
	}
}

// ToResFromEntityConversationParticipant chuyển từ DB Entity sang DTO Response trả về client.
func ToResFromEntityConversationParticipant(e *entity.ConversationParticipant) *res.ConversationParticipantRes {
	if e == nil {
		return nil
	}

	return &res.ConversationParticipantRes{
		ID:                e.ID.Hex(),
		ConversationID:    e.ConversationID.Hex(),
		UserID:            e.UserID,
		Role:              e.Role,
		Nickname:          e.Nickname,
		LastSeenAt:        e.LastSeenAt,
		LastSeenMessageID: e.LastSeenMessageID,
		MuteUntil:         e.MuteUntil,
		IsArchived:        e.IsArchived,
		ClearHistoryAt:    e.ClearHistoryAt,
		JoinedAt:          e.JoinedAt,
		AddedByUserID:     e.AddedByUserID,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
		DeletedAt:         e.DeletedAt,
	}
}
