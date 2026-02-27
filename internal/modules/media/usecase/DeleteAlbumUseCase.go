package usecase

import (
	"context"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepostitoryMongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteAlbumUseCase struct {
	ablumRepo       IRepostitoryMongodb.IAlbumsRepository
	mediaAssetsRepo IRepostitoryMongodb.IMediaAssetsRepository
	pool            IRepositoryShare.IWorkerPool
}

func NewDeleteAlbumUseCase(ablumRepo IRepostitoryMongodb.IAlbumsRepository, mediaAssetsRepo IRepostitoryMongodb.IMediaAssetsRepository, pool IRepositoryShare.IWorkerPool) *DeleteAlbumUseCase {
	return &DeleteAlbumUseCase{
		ablumRepo:       ablumRepo,
		mediaAssetsRepo: mediaAssetsRepo,
		pool:            pool,
	}
}

func (uc *DeleteAlbumUseCase) Execute(ctx context.Context, req *req.DeleteAlbumRequest) (*res.DeleteAlbumResponse, error) {
	// Thực hiện logic xóa album ở đây
	var failedMediaAssets []*res.FailedMediaAssetsResponse
	dataAlbum, err := uc.ablumRepo.GetAlbumByID(ctx, req.AlbumID)
	if err != nil {
		return &res.DeleteAlbumResponse{
			AlbumID:            req.AlbumID,
			AlbumsErrorMessage: err.Error(),
			Fail:               failedMediaAssets,
		}, err
	}
	dataMediaAssets, err := uc.mediaAssetsRepo.GetListMediaAssetsByAlbumID(ctx, dataAlbum.ID.Hex())
	if err != nil {
		return &res.DeleteAlbumResponse{
			AlbumID:            req.AlbumID,
			AlbumsErrorMessage: err.Error(),
			Fail:               failedMediaAssets,
		}, err
	}
	for _, mediaAsset := range dataMediaAssets {
		err := uc.pool.Run(ctx, func() {
			mediaAsset.AlbumID = primitive.NilObjectID
			err = uc.mediaAssetsRepo.UpdateMediaAsset(ctx, mediaAsset)
			if err != nil {
				failedMediaAssets = append(failedMediaAssets, &res.FailedMediaAssetsResponse{
					MediaAssetsId: mediaAsset.ID.Hex(),
					ErrorMessage:  fmt.Sprintf("Failed to remove album association from media asset: %v", err),
				})
			}
		})
		if err != nil {
			failedMediaAssets = append(failedMediaAssets, &res.FailedMediaAssetsResponse{
				MediaAssetsId: mediaAsset.ID.Hex(),
				ErrorMessage:  fmt.Sprintf("Failed to process media asset in worker pool: %v", err),
			})
		}
	}
	err = uc.ablumRepo.DeleteAlbum(ctx, req.AlbumID)
	if err != nil {
		return &res.DeleteAlbumResponse{
			AlbumID:            req.AlbumID,
			AlbumsErrorMessage: err.Error(),
			Fail:               failedMediaAssets,
		}, err
	}
	return &res.DeleteAlbumResponse{
		AlbumID:            req.AlbumID,
		AlbumsErrorMessage: "",
		Fail:               failedMediaAssets,
	}, nil
	return nil, nil
}
