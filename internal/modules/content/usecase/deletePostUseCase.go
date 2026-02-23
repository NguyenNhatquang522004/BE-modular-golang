package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type DeletePostUseCase struct {
	postRepo          IRepositoryMongodb.IPostRepository
	postMediaRepo     IRepositoryMongodb.IPostMediaRepository
	postSetting       IRepositoryMongodb.IPostSettingRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	postInsightRepo   IRepositoryCassandra.IPostInsights
}

func NewDeletePostUseCase(postRepo IRepositoryMongodb.IPostRepository, postMediaRepo IRepositoryMongodb.IPostMediaRepository, postSetting IRepositoryMongodb.IPostSettingRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, postInsightRepo IRepositoryCassandra.IPostInsights) *DeletePostUseCase {
	return &DeletePostUseCase{
		postRepo:          postRepo,
		postMediaRepo:     postMediaRepo,
		postSetting:       postSetting,
		postExtensionRepo: postExtensionRepo,
		postInsightRepo:   postInsightRepo,
	}
}
