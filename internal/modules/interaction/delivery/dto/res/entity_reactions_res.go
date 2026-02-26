package res

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/gocql/gocql"
)

// EntityReactionRes trả về dữ liệu cho client
type EntityReactionRes struct {
	TargetID     string                     `json:"target_id"`
	UserID       gocql.UUID                 `json:"user_id"`
	TargetType   sharedEnums.ReactionTarget `json:"target_type"`
	ReactionCode sharedEnums.ReactionCode   `json:"reaction_code"`
	CreatedAt    time.Time                  `json:"created_at"`
}
