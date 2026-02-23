package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type GetPostDetailUseCase struct {
	postRepo          IRepositoryMongodb.IPostRepository
	postMediaRepo     IRepositoryMongodb.IPostMediaRepository
	postSetting       IRepositoryMongodb.IPostSettingRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	postInsightRepo   IRepositoryCassandra.IPostInsights 
}

func NewGetPostDetailUseCase(postRepo IRepositoryMongodb.IPostRepository, postMediaRepo IRepositoryMongodb.IPostMediaRepository, postSetting IRepositoryMongodb.IPostSettingRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, postInsightRepo IRepositoryCassandra.IPostInsights) *GetPostDetailUseCase {
	return &GetPostDetailUseCase{
		postRepo:          postRepo,
		postMediaRepo:     postMediaRepo,
		postSetting:       postSetting,
		postExtensionRepo: postExtensionRepo,
		postInsightRepo:   postInsightRepo,
	}
}
