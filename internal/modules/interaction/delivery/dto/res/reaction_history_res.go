package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// UserReactionHistoryRes dùng để trả dữ liệu về cho client.
type UserReactionHistoryRes struct {
	UserID       gocql.UUID                 `json:"user_id"`
	CreatedAt    time.Time                  `json:"created_at"`
	TargetID     string                     `json:"target_id"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code"`
}
