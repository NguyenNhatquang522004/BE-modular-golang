package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/enum"
)

type SearchHistoryReq struct {
	UserID     string                 `json:"user_id" validate:"required"`
	Keyword    string                 `json:"keyword" validate:"required,min=1"`
	TargetID   *string                `json:"target_id,omitempty"`
	TargetType *enum.SearchTargetType `json:"target_type,omitempty"`
}