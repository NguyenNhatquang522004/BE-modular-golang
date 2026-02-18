package req

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
)

type RequesterIDRequest struct {
	RequesterID string `json:"requester_id" validate:"required,uuid4"`
}

type RecipientIDRequest struct {
	RecipientID string `json:"recipient_id" validate:"required,uuid4"`
}

type FriendshipIDRequest struct {
	FriendshipID string `json:"friendship_id" validate:"required,uuid4"`
}
type FriendStatusRequest struct {
	Status enum.StatusFriendship `json:"status" validate:"required,oneof=accepted rejected"`
}
type FriendShipUseCaseRequest struct {
	*FriendshipIDRequest
	*RequesterIDRequest
	*RecipientIDRequest
	*FriendStatusRequest
}
type CreateFriendshipRequest struct {
	*FriendshipIDRequest
	*RequesterIDRequest
	*RecipientIDRequest
	*FriendStatusRequest
}
type UpdateFriendshipStatusRequest struct {
	*FriendshipIDRequest
	*FriendStatusRequest
}
type UpdateBlockFriendshipRequest struct {
	*FriendshipIDRequest
	Status enum.Type_Block `json:"status" validate:"required,oneof=blocked unblocked"`
}
type DeleteSoftFriendshipRequest struct {
	*FriendshipIDRequest
}

type DeleteHardFriendshipRequest struct {
	*FriendshipIDRequest
}

type PaginationFriendshipRequest struct {
	UserID   string             `form:"user_id" validate:"required,uuid4"`
	Metadata *dto.PaginationReq `form:"metadata"`
}
