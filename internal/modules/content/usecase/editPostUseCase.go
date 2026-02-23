package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type EditPostUseCase struct {
	postRepo          IRepositoryMongodb.IPostRepository
	postMediaRepo     IRepositoryMongodb.IPostMediaRepository
	postSetting       IRepositoryMongodb.IPostSettingRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	postInsightRepo   IRepositoryCassandra.IPostInsights
}

func NewEditPostUseCase(postRepo IRepositoryMongodb.IPostRepository, postMediaRepo IRepositoryMongodb.IPostMediaRepository, postSetting IRepositoryMongodb.IPostSettingRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, postInsightRepo IRepositoryCassandra.IPostInsights) *EditPostUseCase {
	return &EditPostUseCase{
		postRepo:          postRepo,
		postMediaRepo:     postMediaRepo,
		postSetting:       postSetting,
		postExtensionRepo: postExtensionRepo,
		postInsightRepo:   postInsightRepo,
	}
}
