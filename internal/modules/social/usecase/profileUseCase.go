package usecase

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileUseCase struct {
	profileRepo   IRepositoryMongodb.IProfileRepositoryMongodb
	client        *mongo.Database
	seaweedfsRepo IRepositoryShare.ISeaweedfs
}

func NewProfileUseCase(client *mongo.Database, profileRepo IRepositoryMongodb.IProfileRepositoryMongodb, seaweedfsRepo IRepositoryShare.ISeaweedfs) *ProfileUseCase {
	return &ProfileUseCase{
		client:        client,
		profileRepo:   profileRepo,
		seaweedfsRepo: seaweedfsRepo,
	}
}
func (r *ProfileUseCase) CreateProfileUseCase(req *req.CreateAndUpdateProfileRequest) (*response.Response, error) {
	err := r.profileRepo.CreateProfile(req.ProfileReq)
	if err != nil {
		return response.NewResponse(
			response.WithData(""),
			response.WithStatus("error"),
			response.WithMessage("Failed to create profile"),
		), err
	}
	return response.NewResponse(
		response.WithData(""),
		response.WithStatus("success"),
		response.WithMessage("Profile created successfully"),
	), nil
}
func (r *ProfileUseCase) GetProfileByIDUseCase(req *req.ProfileIDRequest) (*response.Response, error) {
	data, err := r.profileRepo.GetProfileByID(req.ProfileID)
	if data.Avatar != nil {
		data.Avatar.URL = r.seaweedfsRepo.GetPublicURL(data.Avatar.URL)
	}
	if data.CoverPhoto != nil {
		data.CoverPhoto.URL = r.seaweedfsRepo.GetPublicURL(data.CoverPhoto.URL)
	}
	if data.CVDocument != nil {
		data.CVDocument.Filename = r.seaweedfsRepo.GetPublicURL(data.CVDocument.Filename)
	}

	if err != nil {
		return response.NewResponse(
			response.WithData(""),
			response.WithStatus("error"),
			response.WithMessage("Failed to retrieve profile"),
		), err
	}
	return response.NewResponse(
		response.WithData(data),
		response.WithStatus("success"),
		response.WithMessage("Profile retrieved successfully"),
	), nil
}
func (r *ProfileUseCase) UpdateProfileUseCase(req *req.CreateAndUpdateProfileRequest) (*response.Response, error) {
	err := r.profileRepo.UpdateProfile(req.ProfileID, req.ProfileReq)
	if err != nil {
		return response.NewResponse(
			response.WithData(""),
			response.WithStatus("error"),
			response.WithMessage("Failed to update profile"),
		), err
	}
	return response.NewResponse(
		response.WithData(""),
		response.WithStatus("success"),
		response.WithMessage("Profile updated successfully"),
	), nil
}
func (r *ProfileUseCase) GetProfileByUserIDUseCase(req *req.ProfileIDRequest) (*response.Response, error) {
	data, err := r.profileRepo.GetProfileByUserID(req.ProfileID)
	if data.Avatar != nil {
		data.Avatar.URL = r.seaweedfsRepo.GetPublicURL(data.Avatar.URL)
	}
	if data.CoverPhoto != nil {
		data.CoverPhoto.URL = r.seaweedfsRepo.GetPublicURL(data.CoverPhoto.URL)
	}
	if data.CVDocument != nil {
		data.CVDocument.Filename = r.seaweedfsRepo.GetPublicURL(data.CVDocument.Filename)
	}
	if err != nil {
		return response.NewResponse(
			response.WithData(""),
			response.WithStatus("error"),
			response.WithMessage("Failed to retrieve profile"),
		), err
	}

	return response.NewResponse(
		response.WithData(data),
		response.WithStatus("success"),
		response.WithMessage("Profile retrieved successfully"),
	), nil
}
