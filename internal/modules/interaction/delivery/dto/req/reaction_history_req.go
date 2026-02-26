package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// UserReactionHistoryReq ánh xạ đầy đủ các trường của Entity.
// Bổ sung các tag validate để phục vụ việc kiểm tra dữ liệu đầu vào từ API.
type UserReactionHistoryReq struct {
	UserID       gocql.UUID                 `json:"user_id" validate:"required"`
	CreatedAt    time.Time                  `json:"created_at"`
	TargetID     string                     `json:"target_id" validate:"required"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code" validate:"required"`
}
