package mapper

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

// ==============================
// REQ MAPPER
// ==============================

// ToEntity chuyển đổi từ Request DTO sang Entity.
func ToEntityUserReactionHistory(r *req.UserReactionHistoryReq) *entity.UserReactionHistory {
	if r == nil {
		return nil
	}

	// Best Practice: Đảm bảo Clustering Key (CreatedAt) luôn có giá trị để Cassandra sắp xếp đúng
	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	return &entity.UserReactionHistory{
		UserID:       r.UserID,
		CreatedAt:    createdAt,
		TargetID:     r.TargetID,
		TargetType:   r.TargetType,
		ReactionCode: r.ReactionCode,
	}
}

// UpdateToEntity cập nhật dữ liệu từ Request vào Entity hiện có.
func UpdateToEntityUserReactionHistory(r *req.UserReactionHistoryReq, e *entity.UserReactionHistory) {
	if r == nil || e == nil {
		return
	}

	// [LƯU Ý CASSANDRA]: Trong Cassandra, UserID (Partition Key) và CreatedAt (Clustering Key)
	// hợp thành Primary Key. Nếu bạn thay đổi 2 trường này, bản chất Cassandra sẽ tạo ra một record mới (Upsert)
	// thay vì update record cũ.
	e.UserID = r.UserID
	e.TargetID = r.TargetID
	e.TargetType = r.TargetType
	e.ReactionCode = r.ReactionCode

	if !r.CreatedAt.IsZero() {
		e.CreatedAt = r.CreatedAt
	}
}

// ==============================
// RES MAPPER
// ==============================

// FromReqToRes chuyển đổi trực tiếp từ Request DTO sang Response DTO theo yêu cầu.
func FromReqToResUserReactionHistory(r *req.UserReactionHistoryReq) *res.UserReactionHistoryRes {
	if r == nil {
		return nil
	}

	return &res.UserReactionHistoryRes{
		UserID:       r.UserID,
		CreatedAt:    r.CreatedAt,
		TargetID:     r.TargetID,
		TargetType:   r.TargetType,
		ReactionCode: r.ReactionCode,
	}
}

// FromEntityToRes (Bổ sung Best Practice): Thường dùng để map dữ liệu từ DB trả về cho API.
func FromEntityToResUserReactionHistory(e *entity.UserReactionHistory) *res.UserReactionHistoryRes {
	if e == nil {
		return nil
	}

	return &res.UserReactionHistoryRes{
		UserID:       e.UserID,
		CreatedAt:    e.CreatedAt,
		TargetID:     e.TargetID,
		TargetType:   e.TargetType,
		ReactionCode: e.ReactionCode,
	}
}
