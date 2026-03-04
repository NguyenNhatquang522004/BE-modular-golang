package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IStrategy"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/enum"
	"github.com/google/uuid"
)

type BlockUseCase struct {
	blockRepo            IRepositoryPostgres.IBlockRepository
	blockProducer        IGraph.IBlockMessage
	handlerblockStrategy map[enum.Type_Block]IStrategy.IBlockStrategy
}

func NewBlockUseCase(blockRepo IRepositoryPostgres.IBlockRepository,
	blockProducer IGraph.IBlockMessage,
	friendshipRepo IRepositoryPostgres.IFriendshipsRepository, friendshipProducer IGraph.IFriendshipMessage,
	followRepo IRepositoryPostgres.IFollowersRepository,
	followProducer IGraph.IFollowMessage, handlerblockStrategy []IStrategy.IBlockStrategy) *BlockUseCase {
	hMap := make(map[enum.Type_Block]IStrategy.IBlockStrategy)
	for _, handler := range handlerblockStrategy {
		hMap[handler.Type()] = handler
	}
	return &BlockUseCase{
		blockRepo:            blockRepo,
		blockProducer:        blockProducer,
		handlerblockStrategy: hMap,
	}
}

func (uc *BlockUseCase) UseCaseBlockUser(ctx context.Context, req *req.BlockCreateRequest) (*response.Response, error) {
	handler, ok := uc.handlerblockStrategy[req.Status]
	if !ok {
		return response.NewResponse(response.WithData(""), response.WithMessage("Invalid block type"), response.WithStatus("400")), nil
	}
	data, err := handler.Execute(ctx, req)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), err
	}

	return data, nil
}

func (uc *BlockUseCase) IsBlockedUseCase(ctx context.Context, req *req.BlockIsBlockedRequest) (*response.Response, error) {
	data, err := uc.blockRepo.IsBlocked(ctx, uuid.MustParse(req.BlockerUserID), uuid.MustParse(req.BlockedUserID))
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(data), response.WithMessage("Check block status successfully"), response.WithStatus("200")), nil
}
func (uc *BlockUseCase) GetPaginationTypeBlockUseCase(ctx context.Context, req *req.BlockPaginationTypeBlockRequest) (*response.Response, error) {
	data, err := uc.blockRepo.GetPaginationTypeBlock(ctx, uuid.MustParse(req.BlockerUserID), req.Metadata.Cursor, req.Metadata.Limit, req.BlockType)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), err
	}
	return response.NewResponse(response.WithData(data), response.WithMessage("Get pagination type block successfully"), response.WithStatus("200")), nil
}
