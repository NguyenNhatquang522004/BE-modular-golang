package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	res "github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type SharePostUseCase struct {
	eventbus events.EventBus
}

func NewSharePostUseCase(eventbus events.EventBus, postRepo IRepositoryMongodb.IPostRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, pool IRepositoryShare.IWorkerPool) *SharePostUseCase {
	return &SharePostUseCase{
		eventbus: eventbus,
	}
}
func (s *SharePostUseCase) Execute(ctx context.Context, req *req.SharePostRequest) (*res.FailSharePostResponse, error) {
	err := s.eventbus.Publish(ctx, constants.TopicSharePost.String(), req.PostID, constants.Created.String(), req)
	return &res.FailSharePostResponse{
		UserID:  req.UserID,
		PostID:  req.PostID,
		Message: "Failed to share post",
	}, err
}
