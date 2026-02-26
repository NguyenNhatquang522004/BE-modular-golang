package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type SharePostToGroupUseCase struct {
	postRepo          IRepositoryMongodb.IPostRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	postinsightRepo   IRepositoryCassandra.IPostInsights
	postsetting       IRepositoryMongodb.IPostSettingRepository
}

func NewSharePostToGroupUseCase() *SharePostToGroupUseCase {
	return &SharePostToGroupUseCase{}
}
func (s *SharePostToGroupUseCase) Execute(ctx context.Context, req *req.SharePostToGroupRequest) (*response.Response, error) {
	return nil, nil
}
