package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MapCreatePayloadToEntity chuyển đổi Create Payload thành MongoDB Entity
// Best Practice: Luôn trả về error vì có thao tác parse ObjectID từ string
func MapCreatePayloadToEntity(payload *communicationEvent.CreateParticipantPayload) (*entity.ConversationParticipant, error) {
	// 1. Parse string sang primitive.ObjectID của MongoDB
	convID, err := primitive.ObjectIDFromHex(payload.ConversationID)
	if err != nil {
		return nil, err // Trả về lỗi nếu string không phải là 24-byte hex hợp lệ
	}

	// 2. Lấy thời gian hiện tại chuẩn UTC
	now := time.Now().UTC()

	// 3. Khởi tạo Entity với các giá trị mặc định chuẩn
	participant := &entity.ConversationParticipant{
		ID:             primitive.NewObjectID(), // Tự động sinh _id mới cho Document
		ConversationID: convID,
		UserID:         payload.UserID,
		Role:           payload.Role,
		Nickname:       payload.Nickname,
		AddedByUserID:  payload.AddedByUserID,

		// Metadata & Timestamps
		JoinedAt:  now,
		CreatedAt: now,
		UpdatedAt: now,

		// Các trường trạng thái ban đầu (Khởi tạo rõ ràng để clean code)
		IsArchived:        false,
		LastSeenAt:        time.Time{}, // Zero value
		LastSeenMessageID: "",
		MuteUntil:         nil,
		ClearHistoryAt:    nil,
		DeletedAt:         nil,
	}

	return participant, nil
}

// ApplyUpdatePayloadToEntity áp dụng các thay đổi từ Update Payload vào Entity có sẵn.
// Best Practice: Nhận vào con trỏ (*entity) để sửa trực tiếp, không cần return Entity.
func ApplyUpdatePayloadToEntity(existingEntity *entity.ConversationParticipant, payload *communicationEvent.UpdateParticipantPayload) {
	isUpdated := false

	// Cập nhật Nickname nếu client có gửi lên
	if payload.Nickname != nil {
		existingEntity.Nickname = *payload.Nickname
		isUpdated = true
	}

	// Cập nhật Role nếu client có gửi lên
	if payload.Role != nil {
		existingEntity.Role = *payload.Role
		isUpdated = true
	}

	// Cập nhật trạng thái Archived nếu client có gửi lên
	if payload.IsArchived != nil {
		existingEntity.IsArchived = *payload.IsArchived
		isUpdated = true
	}

	// Cập nhật MuteUntil
	// ⚠️ Xử lý Best Practice cho việc UNMUTE trong Golang:
	// Do JSON 'null' và trường 'bị bỏ qua' đều làm MuteUntil == nil trong Golang.
	// Cách tốt nhất: Nếu client muốn UNMUTE, họ phải gửi một thời gian "Zero" (ví dụ: "0001-01-01T00:00:00Z").
	if payload.MuteUntil != nil {
		if payload.MuteUntil.IsZero() {
			// Nếu là thời gian rỗng -> Hủy Mute (Unmute)
			existingEntity.MuteUntil = nil
		} else {
			// Nếu có thời gian cụ thể -> Gán thời gian Mute (Chuyển về UTC)
			utcTime := payload.MuteUntil.UTC()
			existingEntity.MuteUntil = &utcTime
		}
		isUpdated = true
	}

	// Tự động cập nhật UpdatedAt nếu có bất kỳ trường nào bị thay đổi
	if isUpdated {
		existingEntity.UpdatedAt = time.Now().UTC()
	}
}
