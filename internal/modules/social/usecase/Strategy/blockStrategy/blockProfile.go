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

type BlockProfile struct {
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	friendShipProducer IGraph.IFriendshipMessage
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
}

func NewBlockProfile(blockRepo IRepositoryPostgres.IBlockRepository,
	blockProducer IGraph.IBlockMessage,
	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer IGraph.IFriendshipMessage,
	followRepo IRepositoryPostgres.IFollowersRepository,
	followProducer IGraph.IFollowMessage) *BlockProfile {
	return &BlockProfile{
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
		friendshipRepo:     friendshipRepo,
		friendShipProducer: friendshipProducer,
		followRepo:         followRepo,
		followProducer:     followProducer,
	}
}

// Execute(ctx context.Context, req *req.BlockCreateRequest) error
func (h *BlockProfile) Execute(ctx context.Context, req *req.BlockCreateRequest) (*response.Response, error) {
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
	payloadBlock1 := &socialEvent.BlockUpdatePayload{
		Blocker_UserID: req.BlockerUserID,
		Blocked_UserID: req.BlockedUserID,
		Status:         req.Status,
	}
	err := h.blockProducer.PublishBlockUpdateMessage(ctx, payloadBlock1)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("Block profile executed successfully"), response.WithStatus("200")), nil
}
func (h *BlockProfile) Type() enum.Type_Block {
	return enum.Type_Block_Profile
}

// Type() enum.Type_Block // Để nhận diện strategy này dùng cho Enum nào
