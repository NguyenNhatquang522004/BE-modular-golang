package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type HandlerSocialGRPC struct {
	pb.UnimplementedSocialServiceServer
	blockRepo IRepositoryPostgres.IBlockRepository
}

func NewHandlerSocialGRPC(blockRepo IRepositoryPostgres.IBlockRepository) *HandlerSocialGRPC {
	return &HandlerSocialGRPC{
		blockRepo: blockRepo,
	}
}
func (h *HandlerSocialGRPC) GetListBlockByUserID(ctx context.Context, req *pb.UserblockIDRequest) (*pb.ListBlockByUserIDResponse, error) {
	data, err := h.blockRepo.GetBlockByUserIDs(ctx, req.UserId, req.UserId2)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return &pb.ListBlockByUserIDResponse{
			Typeblock: []string{},
		}, nil
	}
	var typeBlocks []string
	for _, block := range data.Type_Block {
		typeBlocks = append(typeBlocks, block.String())
	}
	return &pb.ListBlockByUserIDResponse{
		Typeblock: typeBlocks,
	}, nil
}
