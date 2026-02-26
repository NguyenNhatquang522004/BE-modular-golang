package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

// ToEntityEntityReactions chuyển đổi từ Request DTO sang Entity
func ToEntityEntityReactions(r *req.EntityReactionReq) *entity.EntityReaction {
	if r == nil {
		return nil
	}

	// Best Practice: Nếu client không gửi CreatedAt, tự động gán thời gian hiện tại
	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	return &entity.EntityReaction{
		TargetID:     r.TargetID,
		UserID:       r.UserID,
		TargetType:   r.TargetType,
		ReactionCode: r.ReactionCode,
		CreatedAt:    createdAt,
	}
}

// UpdateToEntityEntityReactions cập nhật dữ liệu từ Request DTO vào Entity hiện có
func UpdateToEntityEntityReactions(r *req.EntityReactionReq, e *entity.EntityReaction) {
	if r == nil || e == nil {
		return
	}

	// [Best Practice Note]: Cassandra không cho phép UPDATE Primary Key (TargetID, UserID)
	// Nếu bạn thay đổi TargetID/UserID, nó sẽ thực hiện lệnh UPSERT (tạo ra một record mới).
	// Code dưới đây vẫn ánh xạ 100% trường theo đúng yêu cầu của bạn.
	e.TargetID = r.TargetID
	e.UserID = r.UserID
	e.TargetType = r.TargetType
	e.ReactionCode = r.ReactionCode

	if !r.CreatedAt.IsZero() {
		e.CreatedAt = r.CreatedAt
	}
}

// ==============================
// RES MAPPER
// ==============================

// FromReqToRes chuyển đổi trực tiếp từ Request DTO sang Response DTO (Theo đúng yêu cầu của bạn)
func FromReqToResEntityReactions(r *req.EntityReactionReq) *res.EntityReactionRes {
	if r == nil {
		return nil
	}

	return &res.EntityReactionRes{
		TargetID:     r.TargetID,
		UserID:       r.UserID,
		TargetType:   r.TargetType,
		ReactionCode: r.ReactionCode,
		CreatedAt:    r.CreatedAt,
	}
}

// [Best Practice Bổ sung] Hàm này thường dùng nhất: Chuyển từ Entity sang Res
func FromEntityToResEntityReactions(e *entity.EntityReaction) *res.EntityReactionRes {
	if e == nil {
		return nil
	}

	return &res.EntityReactionRes{
		TargetID:     e.TargetID,
		UserID:       e.UserID,
		TargetType:   e.TargetType,
		ReactionCode: e.ReactionCode,
		CreatedAt:    e.CreatedAt,
	}
}
