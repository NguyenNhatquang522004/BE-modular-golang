package friendshipstrategy

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type FriendshipDeclined struct {
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
	friendshipProducer IGraph.IFriendshipMessage
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
}

func NewFriendshipDeclined(friendshipRepo IRepositoryPostgres.IFriendshipsRepository, followRepo IRepositoryPostgres.IFollowersRepository, followProducer IGraph.IFollowMessage, friendshipProducer IGraph.IFriendshipMessage, blockRepo IRepositoryPostgres.IBlockRepository, blockProducer IGraph.IBlockMessage) *FriendshipDeclined {
	return &FriendshipDeclined{
		friendshipRepo:     friendshipRepo,
		followRepo:         followRepo,
		followProducer:     followProducer,
		friendshipProducer: friendshipProducer,
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
	}
}

// execute(req *req.FriendShipUseCaseRequest) (*response.Response, error)
func (r *FriendshipDeclined) Execute(ctx context.Context, req *req.FriendShipUseCaseRequest) (*response.Response, error) {
	initdata, err := r.friendshipRepo.GetFriendshipTableByTableId(ctx, uuid.MustParse(req.FriendshipID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	if initdata == nil {
		return response.NewResponse(response.WithData(""), response.WithMessage("Friendship not found"), response.WithStatus("")), nil
	}
	err = r.friendshipRepo.DeleteSoftFriendship(ctx, initdata.ID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Friendship declined executed successfully"), response.WithStatus("200")), nil
}

// Type() enum.StatusFriendship // Để nhận diện strategy này dùng cho Enum nào
func (r *FriendshipDeclined) Type() enum.StatusFriendship {
	return enum.StatusFriendship_Declined
}
