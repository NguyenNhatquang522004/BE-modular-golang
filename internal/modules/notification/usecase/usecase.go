package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
)

type IUseSettingNotificationUseCase interface {
	Execute(ctx context.Context, req *req.UseSettingNotificationRequest) error
}

type ICreateDeleteNotificationTemplateUseCase interface {
	Execute(ctx context.Context, req *req.CreateDeleteNotificationTemplateRequest) (*response.Response, error)
}
type IUpdateNotificationTemplateUseCase interface {
	Execute(ctx context.Context, req *req.UpdateNotificationTemplateRequest) (*response.Response, error)
}
type ISendNotificationTypeUseCase interface {
	Execute(ctx context.Context) error
}

type IScheduleSendUseCase interface {
	Execute(ctx context.Context) error
}

type UseCase struct {
	CreateNotificationTemplate ICreateDeleteNotificationTemplateUseCase
	UpdateNotificationTemplate IUpdateNotificationTemplateUseCase
	SendNotificationType       ISendNotificationTypeUseCase
	ScheduleSend               IScheduleSendUseCase
}

func NewUseCase(
	CreateNotificationTemplate ICreateDeleteNotificationTemplateUseCase,
	UpdateNotificationTemplate IUpdateNotificationTemplateUseCase,

	SendNotificationType ISendNotificationTypeUseCase) *UseCase {
	return &UseCase{
		CreateNotificationTemplate: CreateNotificationTemplate,
		UpdateNotificationTemplate: UpdateNotificationTemplate,
		SendNotificationType:       SendNotificationType,
	}
}
