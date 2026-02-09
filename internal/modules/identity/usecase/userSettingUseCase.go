package usecase

import (
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
)

type UserSettingUseCase struct {
	userSettingRepo IRepositoryMongodb.IUserSettingRepository
}

func NewUserSettingUseCase(userSettingRepo IRepositoryMongodb.IUserSettingRepository) *UserSettingUseCase {
	return &UserSettingUseCase{
		userSettingRepo: userSettingRepo,
	}
}
func (k *UserSettingUseCase) CreateUserSetting(userID string) (*response.Response, error) {
	user, err := k.userSettingRepo.CreateUserSettingDefault(userID)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), errors.New(err.Error())
	}
	return response.NewResponse(response.WithData(user),
		response.WithMessage("Create user setting successfully"),
		response.WithStatus("200")), nil

}
func (k *UserSettingUseCase) UpdateUserSettings(userID string, settings *entity.UserSetting) (*response.Response, error) {
	user, err := k.userSettingRepo.UpdateUserSettings(userID, settings)
	if err != nil {
		return response.NewResponse(response.WithData(""), response.WithMessage(err.Error()), response.WithStatus("500")), errors.New(err.Error())
	}
	return response.NewResponse(response.WithData(user),
		response.WithMessage("Update user settings successfully"),
		response.WithStatus("200")), nil
}
