package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

type BlockerUserIDRequest struct {
	BlockerUserID string `json:"blocker_user_id" validate:"required,uuid4"`
}
type BlockedUserIDRequest struct {
	BlockedUserID string `json:"blocked_user_id" validate:"required,uuid4"`
}

type BlockCreateRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
	Status enum.Type_Block `json:"status_block" validate:"required,oneof=0 1"`
}
type BlockDeleteRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
}

type BlockUpdateRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
	Status enum.Type_Block `json:"status_block" validate:"required,oneof=0 1"`
}

type BlockIsBlockedRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
}

type BlockGetBlockedUsersRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
}

type BlockPaginationTypeBlockRequest struct {
	BlockerUserID string             `form:"blocker_user_id" validate:"required,uuid4"`
	Metadata      *dto.PaginationReq `form:"metadata"`
}
