package mapper

import (
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ==============================================================================
// 1. CREATE MAPPER
// Chuyển đổi từ CreatePayload -> Entity để lưu vào DB.
// Nhận thêm creatorID từ Context/Token (vì user không gửi cái này qua body).
// Trả về error vì quá trình parse ObjectID có thể thất bại.
// ==============================================================================

func ToGroupEventEntity(req *communityEvent.CreateGroupEventPayload) (*entity.GroupEvent, error) {
	// 1. Parse GroupID từ string sang MongoDB ObjectID
	groupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group_id format: %w", err)
	}

	// 2. Map dữ liệu
	event := &entity.GroupEvent{
		ID:          primitive.NewObjectID(), // Tự động generate ID mới
		GroupID:     groupID,
		CreatorID:   req.CreatorID,
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,

		// Ánh xạ struct con Location
		Location: entity.EventLocation{
			Type:        req.Location.Type,
			Address:     req.Location.Address,
			Coordinates: req.Location.Coordinates,
		},

		// Khởi tạo Stats mặc định là 0
		AttendeesCount: entity.EventAttendeeStats{
			Going:      0,
			Interested: 0,
		},

		// Luôn dùng UTC cho thời gian lưu vào DB để tránh lỗi múi giờ
		CreatedAt: time.Now().UTC(),
	}

	return event, nil
}

// ==============================================================================
// 2. UPDATE MAPPER
// Chuyển dữ liệu từ UpdatePayload đè lên Entity Cũ.
// Nhận vào một con trỏ Entity có sẵn (lấy từ DB lên) và cập nhật trực tiếp (In-place update).
// ==============================================================================

func ApplyUpdateToGroupEvent(event *entity.GroupEvent, req *communityEvent.UpdateGroupEventPayload) {
	// Kiểm tra từng trường con trỏ, nếu khác nil (client có gửi lên) thì mới cập nhật
	if req.Title != nil {
		event.Title = *req.Title
	}

	if req.Description != nil {
		event.Description = *req.Description
	}

	if req.CoverURL != nil {
		event.CoverURL = *req.CoverURL
	}

	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		event.EndTime = *req.EndTime
	}

	// Với Location là 1 nested struct, khi user update thì thường sẽ update cả cụm
	if req.Location != nil {
		event.Location = entity.EventLocation{
			Type:        req.Location.Type,
			Address:     req.Location.Address,
			Coordinates: req.Location.Coordinates,
		}
	}

	// Lưu ý: Không cập nhật ID, GroupID, CreatorID, AttendeesCount hay CreatedAt ở đây
}
