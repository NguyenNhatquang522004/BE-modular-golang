package IRepositoryMongodb

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
)

type IProfileRepositoryMongodb interface {
	// Define methods for profile repository
	CreateProfile(ctx context.Context, profileData *entity.Profiles) error
	GetProfileByID(ctx context.Context, profileID string) (*entity.Profiles, error)
	UpdateProfile(ctx context.Context, profileData *entity.Profiles) error
	GetProfileByUserID(ctx context.Context, userID string) (*entity.Profiles, error)
	DeleteProfile(ctx context.Context, profileID string) error
}
