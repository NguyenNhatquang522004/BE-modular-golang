package friendshipstrategy

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

type FriendShipAccepted struct {
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
	friendshipProducer IGraph.IFriendshipMessage
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
}

func NewFriendShipAccepted(friendshipRepo IRepositoryPostgres.IFriendshipsRepository, followRepo IRepositoryPostgres.IFollowersRepository, followProducer IGraph.IFollowMessage, friendshipProducer IGraph.IFriendshipMessage, blockRepo IRepositoryPostgres.IBlockRepository, blockProducer IGraph.IBlockMessage) *FriendShipAccepted {
	return &FriendShipAccepted{
		friendshipRepo:     friendshipRepo,
		followRepo:         followRepo,
		followProducer:     followProducer,
		friendshipProducer: friendshipProducer,
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
	}
}

// execute(req *req.FriendShipUseCaseRequest) (*response.Response, error)
func (r *FriendShipAccepted) Execute(req *req.FriendShipUseCaseRequest) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	datainit, err := r.friendshipRepo.GetFriendshipTableByTableId(uuid.MustParse(req.FriendshipID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	datainit.Status = enum.StatusFriendship_Accepted
	err = r.friendshipRepo.UpdateFriendshipStatus(datainit)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	GetFollow, err := r.followRepo.GetFollowerIndiscriminate(datainit.Recipient_ID, datainit.Requester_ID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	if len(GetFollow) == 1 {
		if GetFollow[0].Followed_UserID == uuid.MustParse(req.UserID) {
			err = r.followRepo.CreateFollowUser(uuid.MustParse(req.UserID), GetFollow[0].Follower_UserID)
			if err != nil {
				return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
			}
			payloadFollow := &socialEvent.FollowUpdatePayload{
				Follower_UserID: req.UserID,
				Followed_UserID: GetFollow[0].Follower_UserID.String(),
				IsMuted:         false,
			}
			err = r.followProducer.PublishFollowUpdateMessage(ctx, payloadFollow)
			if err != nil {
				return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
			}

		} else if GetFollow[0].Follower_UserID == uuid.MustParse(req.UserID) {
			err = r.followRepo.CreateFollowUser(GetFollow[0].Followed_UserID, uuid.MustParse(req.UserID))
			if err != nil {
				return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
			}
			payloadFollow := &socialEvent.FollowUpdatePayload{
				Follower_UserID: GetFollow[0].Followed_UserID.String(),
				Followed_UserID: req.UserID,
				IsMuted:         false,
			}
			err = r.followProducer.PublishFollowUpdateMessage(ctx, payloadFollow)
			if err != nil {
				return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
			}
		}
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Friendship accepted executed successfully"), response.WithStatus("200")), nil
}

// Type() enum.StatusFriendship // Để nhận diện strategy này dùng cho Enum nào
func (h *FriendShipAccepted) Type() enum.StatusFriendship {
	return enum.StatusFriendship_Accepted
}
