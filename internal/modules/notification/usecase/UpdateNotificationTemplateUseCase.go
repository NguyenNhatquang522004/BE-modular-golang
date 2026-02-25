package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
)

type UpdateNotificationTemplateUseCase struct {
	notificationRepo IRepositoryMongodb.INotificationTemplatesRepository
}

func NewUpdateNotificationTemplateUseCase(notificationRepo IRepositoryMongodb.INotificationTemplatesRepository) *UpdateNotificationTemplateUseCase {
	return &UpdateNotificationTemplateUseCase{
		notificationRepo: notificationRepo,
	}
}
func (uc *UpdateNotificationTemplateUseCase) Execute(ctx context.Context, req *req.UpdateNotificationTemplateRequest) (*response.Response, error) {
	// Implement the logic to update notification template based on the request data
	// You can interact with the repository layer to perform database operations if needed
	// Update the notification template with the new data from the request
	// You can add your update logic here
	return response.NewResponse(response.WithData(""),
		response.WithMessage("Cập nhật notification template thành công"),
		response.WithStatus(http.StatusOK)), nil
}
