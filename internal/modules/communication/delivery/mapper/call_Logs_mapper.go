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

// ToEntity tạo mới một Entity từ Request.
func ToEntityCallLogs(r *req.CallLogReq) *entity.CallLog {
	if r == nil {
		return nil
	}

	// Xử lý chuyển đổi String sang ObjectID an toàn
	id, _ := primitive.ObjectIDFromHex(r.ID)
	if r.ID == "" {
		id = primitive.NewObjectID() // Auto sinh ID nếu Req không có
	}

	convID, _ := primitive.ObjectIDFromHex(r.ConversationID)

	// Set default time nếu chưa được cung cấp
	now := time.Now()
	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := r.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = now
	}

	return &entity.CallLog{
		ID:              id,
		ConversationID:  convID,
		CallerID:        r.CallerID,
		Participants:    r.Participants,
		Type:            *r.Type,
		Status:          *r.Status,
		StartedAt:       r.StartedAt,
		EndedAt:         r.EndedAt,
		DurationSeconds: r.DurationSeconds,
		IsGroupCall:     r.IsGroupCall,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		DeletedAt:       r.DeletedAt,
	}
}

// UpdateToEntityCallLogs cập nhật Entity hiện có dựa trên data từ Request.
func UpdateToEntityCallLogs(r *req.CallLogReq, e *entity.CallLog) {
	if r == nil || e == nil {
		return
	}

	// Update ObjectID (chỉ cập nhật nếu Hex hợp lệ)
	if r.ConversationID != "" {
		if convID, err := primitive.ObjectIDFromHex(r.ConversationID); err == nil {
			e.ConversationID = convID
		}
	}

	if r.CallerID != "" {
		e.CallerID = r.CallerID
	}
	if len(r.Participants) > 0 {
		e.Participants = r.Participants
	}
	if r.Type != nil {
		e.Type = *r.Type
	}
	if r.Status != nil {
		e.Status = *r.Status
	}
	if !r.StartedAt.IsZero() {
		e.StartedAt = r.StartedAt
	}

	e.EndedAt = r.EndedAt
	e.DurationSeconds = r.DurationSeconds
	e.IsGroupCall = r.IsGroupCall
	e.UpdatedAt = time.Now() // Luôn update lại timestamp khi có thay đổi

	if r.DeletedAt != nil {
		e.DeletedAt = r.DeletedAt
	}
}

// ========================
// MAPPER -> RESPONSE
// ========================

// ToResFromEntityCallLogs chuyển từ DB Entity sang DTO Response trả về client.
func ToResFromEntityCallLogs(e *entity.CallLog) *res.CallLogRes {
	if e == nil {
		return nil
	}

	return &res.CallLogRes{
		ID:              e.ID.Hex(),
		ConversationID:  e.ConversationID.Hex(),
		CallerID:        e.CallerID,
		Participants:    e.Participants,
		Type:            e.Type,
		Status:          e.Status,
		StartedAt:       e.StartedAt,
		EndedAt:         e.EndedAt,
		DurationSeconds: e.DurationSeconds,
		IsGroupCall:     e.IsGroupCall,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		DeletedAt:       e.DeletedAt,
	}
}

// ToResFromReq chuyển trực tiếp từ DTO Request sang DTO Response (theo đúng yêu cầu của bạn).
func ToResFromReq(r *req.CallLogReq) *res.CallLogRes {
	if r == nil {
		return nil
	}

	return &res.CallLogRes{
		ID:              r.ID,
		ConversationID:  r.ConversationID,
		CallerID:        r.CallerID,
		Participants:    r.Participants,
		Type:            *r.Type,
		Status:          *r.Status,
		StartedAt:       r.StartedAt,
		EndedAt:         r.EndedAt,
		DurationSeconds: r.DurationSeconds,
		IsGroupCall:     r.IsGroupCall,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
		DeletedAt:       r.DeletedAt,
	}
}
