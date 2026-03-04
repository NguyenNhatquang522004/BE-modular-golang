package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/mapper"
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
func (r *ProfileUseCase) CreateProfileUseCase(ctx context.Context, req *req.CreateAndUpdateProfileRequest) (*response.Response, error) {
	entity, err := mapper.ToEntityProfile(req.ProfileReq)
	err = r.profileRepo.CreateProfile(ctx, entity)
	if err != nil {
		return response.NewResponse(
			response.WithData(""),
			response.WithStatus("error"),
			response.WithMessage("Failed to map profile data"),
		), err
	}
	err = r.profileRepo.CreateProfile(ctx, entity)
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
func (r *ProfileUseCase) GetProfileByIDUseCase(ctx context.Context, req *req.ProfileIDRequest) (*response.Response, error) {
	data, err := r.profileRepo.GetProfileByID(ctx, req.ProfileID)
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
func (r *ProfileUseCase) UpdateProfileUseCase(ctx context.Context, req *req.CreateAndUpdateProfileRequest) (*response.Response, error) {
	data, err := r.profileRepo.GetProfileByID(ctx, req.ProfileID)
	if err != nil {
		return response.NewResponse(
			response.WithData(""),
			response.WithStatus("error"),
			response.WithMessage("Failed to retrieve profile for update"),
		), err
	}
	// Map dữ liệu cập nhật từ req vào entity hiện tại
	mapper.ToEntityUpdateProfile(data, req.ProfileReq)
	err = r.profileRepo.UpdateProfile(ctx, data)
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
func (r *ProfileUseCase) GetProfileByUserIDUseCase(ctx context.Context, req *req.ProfileIDRequest) (*response.Response, error) {
	data, err := r.profileRepo.GetProfileByUserID(ctx, req.ProfileID)
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
