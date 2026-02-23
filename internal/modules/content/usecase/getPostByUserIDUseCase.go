package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
)

type GetPostByUserIDUseCase struct {
	postRepo          IRepositoryMongodb.IPostRepository
	postMediaRepo     IRepositoryMongodb.IPostMediaRepository
	postSetting       IRepositoryMongodb.IPostSettingRepository
	postExtensionRepo IRepositoryMongodb.IPostExtensionRepository
	postInsightRepo   IRepositoryCassandra.IPostInsights
}

func NewGetPostByUserIDUseCase(postRepo IRepositoryMongodb.IPostRepository, postMediaRepo IRepositoryMongodb.IPostMediaRepository, postSetting IRepositoryMongodb.IPostSettingRepository, postExtensionRepo IRepositoryMongodb.IPostExtensionRepository, postInsightRepo IRepositoryCassandra.IPostInsights) *GetPostByUserIDUseCase {
	return &GetPostByUserIDUseCase{
		postRepo:          postRepo,
		postMediaRepo:     postMediaRepo,
		postSetting:       postSetting,
		postExtensionRepo: postExtensionRepo,
		postInsightRepo:   postInsightRepo,
	}
}
