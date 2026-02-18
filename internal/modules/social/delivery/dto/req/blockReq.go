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
type StatusBlockRequest struct {
	Status enum.Type_Block `json:"status_block" validate:"required,oneof=0 1"`
}
type BlockCreateRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
	*StatusBlockRequest
}
type BlockDeleteRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
	*StatusBlockRequest
}

type BlockUpdateRequest struct {
	*BlockerUserIDRequest
	*BlockedUserIDRequest
	StatusBlockRequest
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
	BlockType     enum.Type_Block    `form:"block_type" validate:"required,oneof=0 1"`
	Metadata      *dto.PaginationReq `form:"metadata"`
}
