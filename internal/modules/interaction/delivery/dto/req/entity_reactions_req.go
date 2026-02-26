package req

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// EntityReactionReq ánh xạ 100% các trường của EntityReaction
// Sử dụng thêm tag `validate` (thường dùng chung với Gin/Echo) để kiểm tra dữ liệu đầu vào
type EntityReactionReq struct {
	TargetID     string                     `json:"target_id" validate:"required"`
	UserID       gocql.UUID                 `json:"user_id" validate:"required"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type" validate:"required"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code" validate:"required"`
	CreatedAt    time.Time                  `json:"created_at"`
}
