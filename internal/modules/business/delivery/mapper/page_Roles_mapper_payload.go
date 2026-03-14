package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToPageRoleEntity(payload *businessEvent.CreatePageRolePayload) (*entity.PageRole, error) {
	// 1. Chuyển đổi an toàn từ String sang ObjectID
	pageObjectID, err := primitive.ObjectIDFromHex(payload.PageID)
	if err != nil {
		return nil, err // Trả về lỗi nếu ID không hợp lệ dù đã qua validator
	}

	now := time.Now()

	// 2. Khởi tạo và trả về Entity
	return &entity.PageRole{
		ID:                primitive.NewObjectID(), // Tự động sinh ID mới
		PageID:            pageObjectID,
		UserID:            payload.UserID,
		Role:              payload.Role,
		CustomPermissions: payload.CustomPermissions,
		CreatedAt:         now,
		UpdatedAt:         now,
		AssignedBy:        payload.AssignedBy, // Lấy từ payload (đã được set trong Service layer)
	}, nil
}

func MapUpdateToPageRoleEntity(existingEntity *entity.PageRole, payload *businessEvent.UpdatePageRolePayload) {
	// Dùng cờ (flag) để kiểm tra xem thực sự có thay đổi dữ liệu không
	hasChanged := false

	// 1. Chỉ cập nhật Role nếu client thực sự truyền lên (khác nil)
	if payload.Role != nil && *payload.Role != existingEntity.Role {
		existingEntity.Role = *payload.Role
		hasChanged = true
	}

	// 2. Cập nhật CustomPermissions nếu client có truyền mảng mới
	// Lưu ý: Nếu client truyền mảng rỗng [], nó sẽ ghi đè xóa hết quyền (đúng logic PATCH)
	if payload.CustomPermissions != nil {
		existingEntity.CustomPermissions = payload.CustomPermissions
		hasChanged = true
	}

	// 3. Chỉ cập nhật UpdatedAt nếu thực sự có field bị thay đổi
	if hasChanged {
		existingEntity.UpdatedAt = time.Now()
	}
}
