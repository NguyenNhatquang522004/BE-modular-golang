package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type HandlerSocialGRPC struct {
	pb.UnimplementedSocialServiceServer
	blockRepo   IRepositoryPostgres.IBlockRepository
	profileRepo IRepositoryMongodb.IProfileRepositoryMongodb
}

func NewHandlerSocialGRPC(blockRepo IRepositoryPostgres.IBlockRepository, profileRepo IRepositoryMongodb.IProfileRepositoryMongodb) *HandlerSocialGRPC {
	return &HandlerSocialGRPC{
		blockRepo:   blockRepo,
		profileRepo: profileRepo,
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

func (h *HandlerSocialGRPC) GetInfoUserByID(ctx context.Context, req *pb.UserSocialIDRequest) (*pb.InfoUserByIDResponse, error) {
	// Giả sử bạn có một phương thức trong blockRepo để lấy thông tin người dùng
	userInfo, err := h.profileRepo.GetProfileByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return &pb.InfoUserByIDResponse{
			UserId:       "",
			AuthorName:   "",
			AuthorAvatar: "",
		}, nil
	}
	return &pb.InfoUserByIDResponse{
		UserId:       userInfo.UserID,
		AuthorName:   userInfo.FullName,
		AuthorAvatar: userInfo.Avatar.URL,
	}, nil
}
