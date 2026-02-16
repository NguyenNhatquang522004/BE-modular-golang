package usecase

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type BlockUseCase struct {
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IProducer.IBlockMessage
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	friendShipProducer IProducer.IFriendshipMessage
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IProducer.IFollowMessage
}

func NewBlockUseCase(blockRepo IRepositoryPostgres.IBlockRepository,
	blockProducer IProducer.IBlockMessage,
	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer IProducer.IFriendshipMessage,
	followRepo IRepositoryPostgres.IFollowersRepository,
	followProducer IProducer.IFollowMessage) *BlockUseCase {
	return &BlockUseCase{
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
		friendshipRepo:     friendshipRepo,
		friendShipProducer: friendshipProducer,
		followRepo:         followRepo,
		followProducer:     followProducer,
	}
}

func (uc *BlockUseCase) CreateBlockUserUseCase(req *req.BlockCreateRequest) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := uc.blockRepo.CreateBlockUser(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID), req.Status)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	follow, err := uc.followRepo.GetFollowerIndiscriminate(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	if err == nil && follow != nil {
		err = uc.followRepo.DeleteBatchSoftFollowUser(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
		}
	}
	payloadFollow1 := &socialEvent.FollowDeletePayload{
		Follower_UserID: req.BlockerUserID,
		Followed_UserID: req.BlockedUserID,
	}
	payloadFollow2 := &socialEvent.FollowDeletePayload{
		Follower_UserID: req.BlockedUserID,
		Followed_UserID: req.BlockerUserID,
	}
	err = uc.followProducer.PublishFollowDeleteMessage(ctx, payloadFollow1, payloadFollow2)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	friendship, err := uc.friendshipRepo.GetFriendshipIndiscriminate(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	if err == nil && friendship != nil {
		err = uc.friendshipRepo.UpdateFriendshipStatus(friendship.ID, enum.StatusFriendship_Blocked)
		if err != nil {
			return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
		}
	}
	payloadFriendship1 := &socialEvent.FriendshipsDeletePayload{
		Requester_ID: req.BlockerUserID,
		Recipient_ID: req.BlockedUserID,
	}
	payloadFriendship2 := &socialEvent.FriendshipsDeletePayload{
		Requester_ID: req.BlockedUserID,
		Recipient_ID: req.BlockerUserID,
	}
	err = uc.friendShipProducer.PublishFriendshipDeleteMessage(ctx, payloadFriendship1, payloadFriendship2)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	payloadBlock1 := &socialEvent.BlockUpdatePayload{
		Blocker_UserID: req.BlockerUserID,
		Blocked_UserID: req.BlockedUserID,
		Status:         req.Status,
	}
	payloadBlock2 := &socialEvent.BlockUpdatePayload{
		Blocker_UserID: req.BlockedUserID,
		Blocked_UserID: req.BlockerUserID,
		Status:         req.Status,
	}
	err = uc.blockProducer.PublishBlockUpdateMessage(ctx, payloadBlock1, payloadBlock2)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}

	return response.NewResponse(response.WithData(""), response.WithMessage("Success"), response.WithStatus("200")), nil
}
func (uc *BlockUseCase) DeleteBlockUserUseCase(req *req.BlockDeleteRequest) (*response.Response, error) {
	
	return nil, nil
}
func (uc *BlockUseCase) UpdateBlockUserUseCase(req *req.BlockUpdateRequest) (*response.Response, error) {
	return nil, nil
}
func (uc *BlockUseCase) IsBlockedUseCase(req *req.BlockIsBlockedRequest) (*response.Response, error) {
	return nil, nil
}
func (uc *BlockUseCase) GetBlockedUsersUseCase(req *req.BlockGetBlockedUsersRequest) (*response.Response, error) {
	return nil, nil
}
func (uc *BlockUseCase) GetPaginationTypeBlockUseCase(req *req.BlockPaginationTypeBlockRequest) (*response.Response, error) {
	return nil, nil
}
