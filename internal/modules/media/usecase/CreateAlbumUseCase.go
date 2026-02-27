package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAlbumUseCase struct {
	ablumRepo       IRepositoryMongodb.IAlbumsRepository
	mediaAssetsRepo IRepositoryMongodb.IMediaAssetsRepository
	pool            IRepositoryShare.IWorkerPool
}

func NewCreateAlbumUseCase(ablumRepo IRepositoryMongodb.IAlbumsRepository, mediaAssetsRepo IRepositoryMongodb.IMediaAssetsRepository, pool IRepositoryShare.IWorkerPool) *CreateAlbumUseCase {
	return &CreateAlbumUseCase{
		ablumRepo:       ablumRepo,
		mediaAssetsRepo: mediaAssetsRepo,
		pool:            pool,
	}
}

func (uc *CreateAlbumUseCase) Execute(ctx context.Context, req *req.CreateAlbumRequest) (*res.CreateAlbumResponse, error) {
	// Thực hiện logic tạo album ở đây
	entityAlbum, err := mapper.ToEntityAlbum(req.AlbumReq)
	if err != nil {
		return &res.CreateAlbumResponse{
			Fail:               nil,
			AlbumID:            entityAlbum.ID.Hex(),
			AlbumsErrorMessage: err.Error(),
		}, err
	}
	type taskResult struct {
		id  string
		err error
	}
	var failedResponses []*res.FailedMediaAssetsResponse
	resultChan := make(chan taskResult, len(req.ItemMediaAssetsIDs))
	if len(req.ItemMediaAssetsIDs) > 0 {
		for _, assetID := range req.ItemMediaAssetsIDs {
			// Thực hiện logic xử lý từng assetID ở đây
			err := uc.pool.Run(ctx, func() {
				dataMedia, err := uc.mediaAssetsRepo.GetMediaAssetByID(ctx, assetID)
				if err != nil {
					failedResponses = append(failedResponses, &res.FailedMediaAssetsResponse{
						MediaAssetsId: assetID,
						ErrorMessage:  fmt.Sprintf("Failed to get media asset by ID: %v", err),
					})
					resultChan <- taskResult{"", errors.New(fmt.Sprintf("Failed to get media asset by ID %s: %v", assetID, err))}
					return
				}
				if dataMedia.AlbumID != primitive.NilObjectID {
					failedResponses = append(failedResponses, &res.FailedMediaAssetsResponse{
						MediaAssetsId: assetID,
						ErrorMessage:  "Media asset already belongs to an album",
					})
					resultChan <- taskResult{dataMedia.ID.Hex(), errors.New("Media asset already belongs to an album")}
					return

				}	
				dataMedia.AlbumID = entityAlbum.ID
				err = uc.mediaAssetsRepo.UpdateMediaAsset(ctx, dataMedia)
				if err != nil {
					failedResponses = append(failedResponses, &res.FailedMediaAssetsResponse{
						MediaAssetsId: assetID,
						ErrorMessage:  fmt.Sprintf("Failed to update media asset with album ID association: %v", err),
					})
					resultChan <- taskResult{"", errors.New(fmt.Sprintf("Failed to update media asset with album ID association %s: %v", assetID, err))}
					return
				}
				resultChan <- taskResult{assetID, nil}
			})
			if err != nil {
				return nil, err
			}
		}
		uc.pool.Wait() // Chờ tất cả các task hoàn thành
		close(resultChan)
	}
	return &res.CreateAlbumResponse{
		Fail:    failedResponses,
		AlbumID: entityAlbum.ID.Hex(),
	}, nil
}
