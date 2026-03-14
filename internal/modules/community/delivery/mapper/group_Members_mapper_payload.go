package mapper

import (
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =====================================================================
// 1. CREATE MAPPER (Payload -> Mới Entity)
// =====================================================================

// ToGroupMemberEntity biến đổi Create payload thành MongoDB entity.
// Trả về error nếu việc parse ObjectID bị lỗi.
func ToGroupMemberEntity(p *communityEvent.CreateGroupMemberPayload) (*entity.GroupMember, error) {
	// 1. Parse GroupID từ string sang ObjectID
	groupID, err := primitive.ObjectIDFromHex(p.GroupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group_id format: %w", err)
	}

	// 2. Setup thời gian hệ thống
	now := time.Now().UTC() // Best practice: Luôn lưu thời gian ở chuẩn UTC

	// 3. Khởi tạo Entity với các trường mặc định
	member := &entity.GroupMember{
		ID:           primitive.NewObjectID(), // Tự động sinh ID mới cho record
		GroupID:      groupID,
		UserID:       p.UserID,
		Role:         sharedEnums.RoleTypeMember,
		Status:       sharedEnums.ProcessingPending,
		InviterID:    p.InviterID,
		JoinedAt:     now,
		LastActiveAt: now,
		CreatedAt:    now,
		UpdatedAt:    now,
		Badges:       make([]sharedEnums.UserBadge, 0), // Khởi tạo mảng rỗng thay vì nil
	}

	// 4. Map mảng JoinAnswers (nếu có)
	if len(p.JoinAnswers) > 0 {
		answers := make([]entity.JoinAnswer, 0, len(p.JoinAnswers))
		for _, ansPayload := range p.JoinAnswers {
			// Parse QuestionID trong từng phần tử
			qID, err := primitive.ObjectIDFromHex(ansPayload.QuestionID)
			if err != nil {
				return nil, fmt.Errorf("invalid question_id '%s': %w", ansPayload.QuestionID, err)
			}

			answers = append(answers, entity.JoinAnswer{
				QuestionID: qID,
				Answer:     ansPayload.Answer,
			})
		}
		member.JoinAnswers = answers
	}

	return member, nil
}

// =====================================================================
// 2. UPDATE MAPPER (Existing Entity + Update Payload -> Updated Entity)
// =====================================================================

// ApplyUpdateToGroupMember ghi đè các trường từ payload lên entity hiện tại.
// Hàm này thay đổi trực tiếp (mutate) con trỏ existingEntity.
func ApplyUpdateToGroupMember(existingEntity *entity.GroupMember, p *communityEvent.UpdateGroupMemberPayload) *entity.GroupMember {
	// 1. Kiểm tra và cập nhật các trường con trỏ (Partial Update)
	if p.Role != nil {
		existingEntity.Role = *p.Role
	}

	if p.Status != nil {
		existingEntity.Status = *p.Status
	}

	// 2. Cập nhật mảng (Nếu client gửi mảng rỗng [], nó sẽ ghi đè thành mảng rỗng)
	if p.Badges != nil {
		existingEntity.Badges = p.Badges
	}

	// 3. Cập nhật thông tin kỷ luật (DisciplineInfo)
	if p.DisciplineInfo != nil {
		// Khởi tạo struct mới để tránh tham chiếu bộ nhớ từ payload
		existingEntity.DisciplineInfo = &entity.DisciplineInfo{
			Reason:    p.DisciplineInfo.Reason,
			BannedBy:  p.DisciplineInfo.BannedBy,
			UntilDate: p.DisciplineInfo.UntilDate,
		}
	}

	// 4. Cập nhật UpdatedAt (Bắt buộc cho mọi thao tác Update)
	existingEntity.UpdatedAt = time.Now().UTC()

	return existingEntity
}
