package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type UpdateAlbumUseCase struct {
	ablumRepo       IRepositoryMongodb.IAlbumsRepository
	mediaAssetsRepo IRepositoryMongodb.IMediaAssetsRepository
	pool            IRepositoryShare.IWorkerPool
}

func NewUpdateAlbumUseCase(ablumRepo IRepositoryMongodb.IAlbumsRepository, mediaAssetsRepo IRepositoryMongodb.IMediaAssetsRepository, pool IRepositoryShare.IWorkerPool) *UpdateAlbumUseCase {
	return &UpdateAlbumUseCase{
		ablumRepo:       ablumRepo,
		mediaAssetsRepo: mediaAssetsRepo,
		pool:            pool,
	}
}

func (uc *UpdateAlbumUseCase) Execute(ctx context.Context, req *req.UpdateAlbumRequest) (*res.UpdateAlbumResponse, error) {
	// Thực hiện logic cập nhật album ở đây
	dataAlbum, err := uc.ablumRepo.GetAlbumByID(ctx, req.AlbumID)
	if err != nil {
		return &res.UpdateAlbumResponse{
			AlbumID:            req.AlbumID,
			AlbumsErrorMessage: err.Error(),
		}, err
	}
	mapper.UpdateToEntityAlbum(req.UpdateAlbumReq, dataAlbum)
	err = uc.ablumRepo.UpdateAlbum(ctx, dataAlbum)
	if err != nil {
		return &res.UpdateAlbumResponse{
			AlbumID:            req.AlbumID,
			AlbumsErrorMessage: err.Error(),
		}, err
	}
	return &res.UpdateAlbumResponse{
		AlbumID:            req.AlbumID,
		AlbumsErrorMessage: "",
	}, nil
}
