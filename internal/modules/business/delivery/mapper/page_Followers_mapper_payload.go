package mapper

import (
	"errors"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToEntityFromCreatePayload(req *businessEvent.CreatePageFollowerPayload) (*entity.PageFollower, error) {
	if req == nil {
		return nil, errors.New("create payload is nil")
	}

	// 1. Chuyển string thành primitive.ObjectID an toàn
	pageObjectID, err := primitive.ObjectIDFromHex(req.PageID)
	if err != nil {
		return nil, errors.New("invalid page_id format")
	}

	now := time.Now()

	// 2. Khởi tạo Entity với ID mới và Timestamps
	ent := &entity.PageFollower{
		ID:        primitive.NewObjectID(),
		PageID:    pageObjectID,
		UserID:    req.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 3. Mapping Settings (Xử lý trường hợp client không gửi settings)
	if req.Settings != nil {
		ent.Settings = &entity.FollowerSettings{
			NotificationLevel: req.Settings.NotificationLevel,
			IsFavorite:        req.Settings.IsFavorite,
		}
	} else {
		// Best Practice: Gán giá trị mặc định nếu client không truyền
		ent.Settings = &entity.FollowerSettings{
			NotificationLevel: sharedEnums.NotifAll, // Giả sử All là mức mặc định
			IsFavorite:        false,
		}
	}

	return ent, nil
}

// ==========================================
// 2. MAPPER CHO UPDATE
// Nhận vào Entity hiện tại và Payload Update -> Trả ra Entity đã được cập nhật
// ==========================================
func UpdateEntityFromPayload(ent *entity.PageFollower, req *businessEvent.UpdatePageFollowerPayload) *entity.PageFollower {
	// 1. Safety check: Tránh panic nil pointer
	if ent == nil || req == nil || req.Settings == nil {
		return ent
	}

	// 2. Đảm bảo Entity có Settings để gán (phòng trường hợp DB đang lưu nil)
	if ent.Settings == nil {
		ent.Settings = &entity.FollowerSettings{}
	}

	isModified := false

	// 3. Gán giá trị có điều kiện (Chỉ gán nếu pointer KHÔNG nil)
	if req.Settings.NotificationLevel != nil {
		ent.Settings.NotificationLevel = *req.Settings.NotificationLevel
		isModified = true
	}

	if req.Settings.IsFavorite != nil {
		ent.Settings.IsFavorite = *req.Settings.IsFavorite
		isModified = true
	}

	// 4. Chỉ cập nhật thời gian nếu thực sự có dữ liệu thay đổi
	if isModified {
		ent.UpdatedAt = time.Now()
	}

	return ent
}
