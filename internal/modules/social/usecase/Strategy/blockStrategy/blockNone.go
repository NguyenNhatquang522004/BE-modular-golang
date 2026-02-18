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

type BlockNone struct {
	blockRepo          IRepositoryPostgres.IBlockRepository
	blockProducer      IGraph.IBlockMessage
	friendshipRepo     IRepositoryPostgres.IFriendshipsRepository
	friendShipProducer IGraph.IFriendshipMessage
	followRepo         IRepositoryPostgres.IFollowersRepository
	followProducer     IGraph.IFollowMessage
}

func NewBlockNone(blockRepo IRepositoryPostgres.IBlockRepository, blockProducer IGraph.IBlockMessage, friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer IGraph.IFriendshipMessage,
	followRepo IRepositoryPostgres.IFollowersRepository, followProducer IGraph.IFollowMessage) *BlockNone {
	return &BlockNone{
		blockRepo:          blockRepo,
		blockProducer:      blockProducer,
		friendshipRepo:     friendshipRepo,
		friendShipProducer: friendshipProducer,
		followRepo:         followRepo,
		followProducer:     followProducer,
	}
}

func (h *BlockNone) Execute(ctx context.Context, req *req.BlockCreateRequest) (*response.Response, error) {
	init, errinit := h.blockRepo.GetBlockIndiscriminate(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if errinit != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(errinit.Error()), response.WithStatus("")), errinit
	}
	init.Type_Block = req.Status
	err := h.blockRepo.DeleteBlockUser(uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	payloadBlock1 := &socialEvent.BlockDeletePayload{
		Blocker_UserID: req.BlockerUserID,
		Blocked_UserID: req.BlockedUserID,
	}
	err = h.blockProducer.PublishBlockDeleteMessage(ctx, payloadBlock1)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("")), err
	}
	return response.NewResponse(response.WithData(""), response.WithMessage("No block action executed"), response.WithStatus("200")), nil
}

func (h *BlockNone) Type() enum.Type_Block {
	return enum.Type_Block_None
}
