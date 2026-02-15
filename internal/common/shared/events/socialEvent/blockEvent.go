package socialEvent

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"

type BlockCreatePayload struct {
	Blocker_UserID string          `json:"blocker_id"`
	Blocked_UserID string          `json:"blocked_id"`
	Status         enum.Type_Block `json:"type_block"`
}
type BlockDeletePayload struct {
	Blocker_UserID string `json:"blocker_id"`
	Blocked_UserID string `json:"blocked_id"`
}

type BlockUpdatePayload struct {
	Blocker_UserID string          `json:"blocker_id"`
	Blocked_UserID string          `json:"blocked_id"`
	Status         enum.Type_Block `json:"type_block"`
}
