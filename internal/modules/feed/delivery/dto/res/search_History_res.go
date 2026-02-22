package res

import (
	"time"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/enum"
)

type SearchHistoryRes struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	Keyword    string                 `json:"keyword"`
	TargetID   *string                `json:"target_id,omitempty"`
	TargetType *enum.SearchTargetType `json:"target_type,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}