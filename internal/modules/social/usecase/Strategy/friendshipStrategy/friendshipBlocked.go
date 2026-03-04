package friendshipstrategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type FriendshipBlocked struct {
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
	friendshipProducer IGraph.IFriendshipMessage
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
}

func NewFriendshipBlocked(friendshipRepo IRepositoryPostgres.IFriendshipsRepository, followRepo IRepositoryPostgres.IFollowersRepository, followProducer IGraph.IFollowMessage, friendshipProducer IGraph.IFriendshipMessage, blockRepo IRepositoryPostgres.IBlockRepository, blockProducer IGraph.IBlockMessage) *FriendshipBlocked {
	return &FriendshipBlocked{
		friendshipRepo:     friendshipRepo,
		followRepo:         followRepo,
		followProducer:     followProducer,
		friendshipProducer: friendshipProducer,
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
	}
}

// execute(req *req.FriendShipUseCaseRequest) (*response.Response, error)
func (r *FriendshipBlocked) Execute(ctx context.Context, req *req.FriendShipUseCaseRequest) (*response.Response, error) {
	initdata, err := r.friendshipRepo.GetFriendshipTableByTableId(ctx, uuid.MustParse(req.FriendshipID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	if initdata == nil {
		return response.NewResponse(response.WithData(""), response.WithMessage("Friendship not found"), response.WithStatus("")), nil
	}
	initdata.Status = enum.StatusFriendship_Blocked
	err = r.friendshipRepo.UpdateFriendshipStatus(ctx, initdata)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	err = r.followRepo.DeleteBatchSoftFollowUser(ctx, initdata.Requester_ID, initdata.Recipient_ID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	payloadFollow := &socialEvent.FollowDeletePayload{
		Follower_UserID: initdata.Requester_ID.String(),
		Followed_UserID: initdata.Recipient_ID.String(),
	}
	payloadFollow2 := &socialEvent.FollowDeletePayload{
		Follower_UserID: initdata.Recipient_ID.String(),
		Followed_UserID: initdata.Requester_ID.String(),
	}
	err = r.followProducer.PublishFollowDeleteMessage(ctx, payloadFollow, payloadFollow2)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	var finalid string
	if initdata.Requester_ID == uuid.MustParse(req.UserID) {
		finalid = initdata.Recipient_ID.String()
	}
	if initdata.Recipient_ID == uuid.MustParse(req.UserID) {
		finalid = initdata.Requester_ID.String()
	}

	err = r.blockRepo.CreateBlockUser(ctx, uuid.MustParse(req.UserID), uuid.MustParse(finalid), enum.Type_Block_Profile)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	payloadBlock := &socialEvent.BlockUpdatePayload{
		Blocker_UserID: req.UserID,
		Blocked_UserID: finalid,
		Status:         enum.Type_Block_Profile,
	}
	err = r.blockProducer.PublishBlockUpdateMessage(ctx, payloadBlock)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Friendship blocked executed successfully"), response.WithStatus("200")), nil
}

// Type() enum.StatusFriendship // Để nhận diện strategy này dùng cho Enum nào
func (h *FriendshipBlocked) Type() enum.StatusFriendship {
	return enum.StatusFriendship_Blocked
}
