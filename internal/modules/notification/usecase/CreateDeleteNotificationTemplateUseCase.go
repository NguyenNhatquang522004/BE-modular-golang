package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
)

type CreateDeleteNotificationTemplateUseCase struct {
	notificationRepo IRepositoryMongodb.INotificationTemplatesRepository
}

func NewCreateDeleteNotificationTemplateUseCase(notificationRepo IRepositoryMongodb.INotificationTemplatesRepository) *CreateDeleteNotificationTemplateUseCase {
	return &CreateDeleteNotificationTemplateUseCase{
		notificationRepo: notificationRepo,
	}
}

func (uc *CreateDeleteNotificationTemplateUseCase) Execute(ctx context.Context, req *req.CreateDeleteNotificationTemplateRequest) (*response.Response, error) {
	switch req.Type {
	case 1:
		// Handle create notification template
		entity := mapper.ToEntityNotificationTemplate(req.CreateNotificationTemplateReq)
		err := uc.notificationRepo.CreateTemplate(ctx, entity)
		if err != nil {
			return response.NewResponse(response.WithData(""),
				response.WithMessage("Tạo notification template thất bại"),
				response.WithStatus(http.StatusBadRequest)), err
		}
		return response.NewResponse(response.WithData(""),
			response.WithMessage("Tạo notification template thành công"),
			response.WithStatus(http.StatusOK)), nil
	case 2:
		// Handle delete notification template
		err := uc.notificationRepo.DeleteTemplate(ctx, req.TemplateID)
		if err != nil {
			return response.NewResponse(response.WithData(""),
				response.WithMessage("Xóa notification template thất bại"),
				response.WithStatus(http.StatusBadRequest)), err
		}
		return response.NewResponse(response.WithData(""),
			response.WithMessage("Xóa notification template thành công"),
			response.WithStatus(http.StatusOK)), nil
	}
	return response.NewResponse(response.WithData(""),
		response.WithMessage("Loại yêu cầu không hợp lệ"),
		response.WithStatus(http.StatusBadRequest)), nil
}
