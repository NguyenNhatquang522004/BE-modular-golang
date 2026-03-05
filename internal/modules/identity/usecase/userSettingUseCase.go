package usecase

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryMongodb"
)

type UserSettingUseCase struct {
	userSettingRepo IRepositoryMongodb.IUserSettingRepository
}

func NewUserSettingUseCase(userSettingRepo IRepositoryMongodb.IUserSettingRepository) *UserSettingUseCase {
	return &UserSettingUseCase{
		userSettingRepo: userSettingRepo,
	}
}
func (k *UserSettingUseCase) CreateUserSetting(ctx context.Context, userID string) (*response.Response, error) {
	user, err := k.userSettingRepo.CreateUserSettingDefault(ctx, userID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), errors.New(err.Error())
	}
	return response.NewResponse(response.WithData(user),
		response.WithMessage("Create user setting successfully"),
		response.WithStatus("200")), nil

}
func (k *UserSettingUseCase) UpdateUserSettings(ctx context.Context, userID string, settings *req.UserSettingReq) (*response.Response, error) {
	data := settings.ToEntity(userID)
	user, err := k.userSettingRepo.UpdateUserSettings(ctx, userID, data)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), errors.New(err.Error())
	}
	return response.NewResponse(response.WithData(user),
		response.WithMessage("Update user settings successfully"),
		response.WithStatus("200")), nil
}
