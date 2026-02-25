package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
)

type UpdateUserNotificationSettingsUseCase struct {
	userNotificationRepo IRepositoryMongodb.IUserNotificationSettingsRepository
}

func NewUpdateUserNotificationSettingsUseCase(userNotificationRepo IRepositoryMongodb.IUserNotificationSettingsRepository) *UpdateUserNotificationSettingsUseCase {
	return &UpdateUserNotificationSettingsUseCase{
		userNotificationRepo: userNotificationRepo,
	}
}
func (uc *UpdateUserNotificationSettingsUseCase) Execute(ctx context.Context, req *req.UpdateUserNotificationSettingsRequest) (*response.Response, error) {
	// Implement the logic to update user notification settings based on the request data
	// You can interact with the repository layer to perform database operations if needed
	data, err := uc.userNotificationRepo.GetUserNotificationSettingsByUserID(ctx, req.UserID)
	mapper.UpdateToEntityUserNotificationSetting(req.UpdateUserNotificationSettingReq, data)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("Lấy thông tin cài đặt thông báo người dùng thất bại"),
			response.WithStatus(http.StatusBadRequest)), err
	}
	err = uc.userNotificationRepo.UpdateUserNotificationSettings(ctx, data)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("Cập nhật cài đặt thông báo người dùng thất bại"),
			response.WithStatus(http.StatusBadRequest)), err
	}
	// Update the user notification settings with the new data from the request
	// You can add your update logic here
	return response.NewResponse(response.WithData(data),
		response.WithMessage("Cập nhật cài đặt thông báo người dùng thành công"),
		response.WithStatus(http.StatusOK)), nil
}
