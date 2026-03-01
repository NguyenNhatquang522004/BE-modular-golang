package grpc

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type MediaGRPCHandler struct {
	pb.UnimplementedMediaServiceServer
	storyRepo IRepositoryMongodb.IStoryRepository
	reelRepo  IRepositoryMongodb.IReelRepository
}

func NewMediaGRPCHandler(storyRepo IRepositoryMongodb.IStoryRepository, reelRepo IRepositoryMongodb.IReelRepository) *MediaGRPCHandler {
	return &MediaGRPCHandler{
		storyRepo: storyRepo,
		reelRepo:  reelRepo,
	}
}

func (h *MediaGRPCHandler) GetUserIDbyStoryID(ctx context.Context, req *pb.GetUserIDbyStoryIDRequest) (*pb.GetUserIDbyStoryIDResponse, error) {
	story, err := h.storyRepo.GetStoryByID(ctx, req.StoryId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserIDbyStoryIDResponse{
		UserId: story.UserID,
	}, nil
}
