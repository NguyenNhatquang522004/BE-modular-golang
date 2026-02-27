package IRepostitoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type IMediaAssetsRepository interface {
	// Define methods for MediaAssetsRepository here
	CreateMediaAsset(ctx context.Context, asset *entity.MediaAsset) error
	CreateBulkMediaAssets(ctx context.Context, assets []*entity.MediaAsset) (int64, []*mongodbErrors.BulkError, error)
	GetMediaAssetByID(ctx context.Context, id string) (*entity.MediaAsset, error)
	GetMediaAssetsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetListMediaAssetsByAlbumID(ctx context.Context, albumID string) ([]*entity.MediaAsset, error)
	GetMediaAssetsByAlbumID(ctx context.Context, albumID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetMediaAssetsByPostID(ctx context.Context, postID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetMediaAssetsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateMediaAsset(ctx context.Context, asset *entity.MediaAsset) error
	UpdateBulkMediaAssets(ctx context.Context, assets []*entity.MediaAsset) (int64, []*mongodbErrors.BulkError, error)
	DeleteMediaAsset(ctx context.Context, id string) error
	DeleteBulkMediaAssets(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
	DeleteMediaAssetsByPostIDAndUserID(ctx context.Context, postID string, userID string) error
}
