package IRepostitoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type IMusicLibraryRepository interface {
	// Define methods for MusicLibraryRepository here
	CreateMusicLibrary(ctx context.Context, musicLibrary *entity.MusicLibrary) (*entity.MusicLibrary, error)
	CreateBulkMusicLibraries(ctx context.Context, musicLibraries []*entity.MusicLibrary) (int64, []*dto.BulkError, error)
	GetMusicLibraryByID(ctx context.Context, id string) (*entity.MusicLibrary, error)
	GetMusicLibraries(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error)
	GetMusicLibrariesByArtist(ctx context.Context, artist string, cursor string, limit int) (*dto.PaginationRes, error)
	UpdateMusicLibrary(ctx context.Context, musicLibrary *entity.MusicLibrary) (*entity.MusicLibrary, error)
	UpdateBulkMusicLibraries(ctx context.Context, musicLibraries []*entity.MusicLibrary) (int64, []*dto.BulkError, error)
	DeleteMusicLibrary(ctx context.Context, id string) error
	DeleteBulkMusicLibraries(ctx context.Context, ids []string) (int64, []*dto.BulkError, error)
}
