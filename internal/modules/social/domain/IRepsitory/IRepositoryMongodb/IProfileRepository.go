package IRepositoryMongodb

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
)

type IProfileRepositoryMongodb interface {
	// Define methods for profile repository
	CreateProfile(profileData *req.ProfileReq) error
	GetProfileByID(profileID string) (*entity.Profiles, error)
	UpdateProfile(profileID string, updateData *req.ProfileReq) error
	GetProfileByUserID(userID string) (*entity.Profiles, error)
}
