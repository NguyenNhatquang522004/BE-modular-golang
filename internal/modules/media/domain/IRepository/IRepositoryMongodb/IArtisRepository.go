package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
)

type IArtistRepository interface {
	CreatedArtist(ctx context.Context, artist *entity.Artist) error
	UpdatedArtist(ctx context.Context, artist *entity.Artist) error
	DeletedArtist(ctx context.Context, artistID string) error
	GetArtistByID(ctx context.Context, artistID string) (*entity.Artist, error)
}
