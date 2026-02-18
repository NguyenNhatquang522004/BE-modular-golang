package blockStrategy

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type BlockFull struct {
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	friendShipProducer IGraph.IFriendshipMessage
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
}

func NewBlockFull(blockRepo IRepositoryPostgres.IBlockRepository,
	blockProducer IGraph.IBlockMessage,
	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer IGraph.IFriendshipMessage,
	followRepo IRepositoryPostgres.IFollowersRepository,
	followProducer IGraph.IFollowMessage) *BlockFull {
	return &BlockFull{
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
		friendshipRepo:     friendshipRepo,
		friendShipProducer: friendshipProducer,
		followRepo:         followRepo,
		followProducer:     followProducer,
	}
}

// Execute(ctx context.Context, req *req.BlockCreateRequest) error
// Type() enum.Type_Block // Để nhận diện strategy này dùng cho Enum nào

func (h *BlockFull) Execute(ctx context.Context, req *req.BlockCreateRequest) (*response.Response, error) {
	init, errinit := h.blockRepo.GetBlockIndiscriminate(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if errinit != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(errinit.Error()), response.WithStatus("")), errinit
	}
	if init != nil {
		init.Type_Block = req.Status
		err := h.blockRepo.UpdateBlockUser(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID), req.Status)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
		}
	} else {
		err := h.blockRepo.CreateBlockUser(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID), req.Status)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
		}
	}
	follow, err := h.followRepo.GetFollowerIndiscriminate(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	if err == nil && follow != nil {
		err = h.followRepo.DeleteBatchSoftFollowUser(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
		}
	}
	for _, followlist := range follow {
		payloadFollow := &socialEvent.FollowDeletePayload{
			Follower_UserID: followlist.Follower_UserID.String(),
			Followed_UserID: followlist.Followed_UserID.String(),
		}
		err = h.followProducer.PublishFollowDeleteMessage(ctx, payloadFollow)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
		}
	}
	friendship, err := h.friendshipRepo.GetFriendshipIndiscriminate(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	if err == nil && friendship != nil {
		err = h.friendshipRepo.UpdateFriendshipStatus(friendship.ID, enum.StatusFriendship_Blocked)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
		}
	}
	payloadFriendship1 := &socialEvent.FriendshipsDeletePayload{
		Requester_ID: friendship.Requester_ID.String(),
		Recipient_ID: friendship.Recipient_ID.String(),
	}
	err = h.friendShipProducer.PublishFriendshipDeleteMessage(ctx, payloadFriendship1)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	payloadBlock1 := &socialEvent.BlockUpdatePayload{
		Blocker_UserID: req.BlockerUserID,
		Blocked_UserID: req.BlockedUserID,
		Status:         req.Status,
	}
	err = h.blockProducer.PublishBlockUpdateMessage(ctx, payloadBlock1)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Block full executed successfully"), response.WithStatus("200")), nil
}

func (h *BlockFull) Type() enum.Type_Block {
	return enum.Type_Block_Full
}
