package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type ContentGRPC struct {
	pb.UnimplementedContentServiceServer
	postRepo IRepositoryMongodb.IPostRepository
}

func NewContentGRPC(postRepo IRepositoryMongodb.IPostRepository) *ContentGRPC {
	return &ContentGRPC{
		postRepo: postRepo,
	}
}

// Implement gRPC methods for content-related operations here
func (h *ContentGRPC) GetOwnerPost(ctx context.Context, req *pb.GetOwnerPostRequest) (*pb.GetOwnerPostResponse, error) {
	data, err := h.postRepo.GetPostByID(ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	return &pb.GetOwnerPostResponse{
		OwnerId: data.UserID,
	}, nil
}
