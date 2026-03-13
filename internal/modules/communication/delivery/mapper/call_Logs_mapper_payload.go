package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1. Create Mapper: Từ Payload -> Entity Mới
// Trả về error vì việc convert chuỗi hex sang primitive.ObjectID có thể thất bại
func ToCallLogEntity(req *communicationEvent.CreateCallLogPayload) (*entity.CallLog, error) {
	// Parse ObjectID an toàn
	convID, err := primitive.ObjectIDFromHex(req.ConversationID)
	if err != nil {
		return nil, err
	}

	// Best Practice: Luôn dùng UTC khi lưu time vào MongoDB
	now := time.Now().UTC()

	return &entity.CallLog{
		ID:             primitive.NewObjectID(),
		ConversationID: convID,
		CallerID:       req.CallerID,
		Participants:   req.Participants,
		Type:           req.Type,
		// Giả sử bạn có enum CallStatusOngoing/Calling trong package enum
		Status:      sharedEnums.CallStatusMissed,
		IsGroupCall: req.IsGroupCall,
		StartedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
		// EndedAt và DeletedAt để trống (nil) ban đầu
	}, nil
}

// 2. Update Mapper: Cập nhật Entity có sẵn từ Payload Update
// Sử dụng con trỏ (*entity.CallLog) để mutate trực tiếp Entity hiện tại giúp tiết kiệm memory
func MapUpdateCallLogEntity(existingEntity *entity.CallLog, req *communicationEvent.UpdateCallLogPayload) *entity.CallLog {
	now := time.Now().UTC()

	// Update các trường thay đổi
	existingEntity.Status = req.Status
	existingEntity.DurationSeconds = req.DurationSeconds
	existingEntity.UpdatedAt = now

	// Vì req.EndedAt là time.Time, còn entity.EndedAt là *time.Time
	// Cần gán giá trị vào một biến local để lấy địa chỉ con trỏ hợp lệ
	endedAt := req.EndedAt.UTC()
	existingEntity.EndedAt = &endedAt

	return existingEntity
}
