package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type IAlbumsRepository interface {
	CreateAlbum(ctx context.Context, album *entity.Album) error
	CreateBulkAlbums(ctx context.Context, albums []*entity.Album) (int64, []*mongodbErrors.BulkError, error)
	GetAlbumByID(ctx context.Context, id string) (*entity.Album, error)
	GetAlbumsByIDs(ctx context.Context, ids []string) ([]*entity.Album, error)
	GetAlbumsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error)
	GetAlbumsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateAlbum(ctx context.Context, album *entity.Album) error
	UpdateBulkAlbums(ctx context.Context, albums []*entity.Album) (int64, []*mongodbErrors.BulkError, error)
	DeleteAlbum(ctx context.Context, id string) error
	DeleteBulkAlbums(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error)
}
