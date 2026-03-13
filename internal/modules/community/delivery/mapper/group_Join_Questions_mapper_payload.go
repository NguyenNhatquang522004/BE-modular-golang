package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =============================================================================
// HELPER MAPPER
// =============================================================================

// mapOptions chuyển đổi mảng Option Payload sang Option Entity
func mapOptions(payloadOptions []communityEvent.QuestionOptionPayload) []entity.QuestionOption {
	// Nếu slice nil, trả về nil để tương thích với cấu hình omitempty của MongoDB
	if payloadOptions == nil {
		return nil
	}

	// Cấp phát trước dung lượng slice (Best practice để tối ưu memory allocation)
	options := make([]entity.QuestionOption, len(payloadOptions))
	for i, opt := range payloadOptions {
		options[i] = entity.QuestionOption{
			Text:  opt.Text,
			Value: opt.Value,
		}
	}
	return options
}

// =============================================================================
// 1. CREATE MAPPER
// =============================================================================

// ToGroupJoinQuestionEntity biến đổi Create Payload thành Entity.
// Do Payload có chứa GroupID dạng string nên có thể parse lỗi, ta cần trả về error.
func ToGroupJoinQuestionEntity(req *communityEvent.CreateGroupJoinQuestionPayload) (*entity.GroupJoinQuestion, error) {
	// 1. Parse string sang ObjectID
	groupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, err // Trả về lỗi để Controller/Service báo 400 Bad Request
	}

	// 2. Xử lý an toàn cho các trường Pointer (Tránh Panic nil pointer dereference)
	// Dù đã validate required, nhưng việc code phòng vệ luôn là Best Practice trong Go
	isRequired := false
	if req.IsRequired != nil {
		isRequired = *req.IsRequired
	}

	order := 1 // Giá trị mặc định an toàn
	if req.Order != nil {
		order = *req.Order
	}

	// 3. Khởi tạo Timestamps và ID mới
	now := time.Now()

	return &entity.GroupJoinQuestion{
		ID:         primitive.NewObjectID(), // Chủ động tạo ID ngay tại đây
		GroupID:    groupID,
		Content:    req.Content,
		Type:       req.Type,
		Options:    mapOptions(req.Options),
		IsRequired: isRequired,
		Order:      order,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// =============================================================================
// 2. UPDATE MAPPER
// =============================================================================

// ApplyUpdateToEntity nhận vào entity hiện tại và payload update.
// Hàm này sẽ trực tiếp biến đổi (mutate) giá trị trên entity gốc và cập nhật UpdatedAt.
func ApplyUpdateToEntity(existingEntity *entity.GroupJoinQuestion, req *communityEvent.UpdateGroupJoinQuestionPayload) {
	// Dùng cờ hasChanges để chỉ cập nhật UpdatedAt nếu thực sự có field bị thay đổi
	hasChanges := false

	if req.Content != nil {
		existingEntity.Content = *req.Content
		hasChanges = true
	}

	if req.Type != nil {
		existingEntity.Type = *req.Type
		hasChanges = true
	}

	// Xử lý Options:
	// - req.Options != nil nghĩa là client CÓ truyền field "options" (có thể là mảng rỗng [] để xóa)
	if req.Options != nil {
		existingEntity.Options = mapOptions(req.Options)
		hasChanges = true
	}

	if req.IsRequired != nil {
		existingEntity.IsRequired = *req.IsRequired
		hasChanges = true
	}

	if req.Order != nil {
		existingEntity.Order = *req.Order
		hasChanges = true
	}

	// Cập nhật thời gian nếu có bất kỳ sự thay đổi nào
	if hasChanges {
		existingEntity.UpdatedAt = time.Now()
	}
}
