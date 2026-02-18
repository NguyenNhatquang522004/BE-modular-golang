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
	UserID string `json:"user_id" validate:"required,uuid4"`
	*FriendshipIDRequest
	*FriendStatusRequest
}
type CreateFriendshipRequest struct {
	*FriendshipIDRequest
	*RequesterIDRequest
	*RecipientIDRequest
	*FriendStatusRequest
}
type DeleteSoftFriendshipRequest struct {
	*FriendshipIDRequest
}

type DeleteHardFriendshipRequest struct {
	*FriendshipIDRequest
}
type PaginationFriendshipRequest struct {
	UserID string `form:"user_id" validate:"required,uuid4"`
	*FriendStatusRequest
	Metadata *dto.PaginationReq `form:"metadata"`
}
