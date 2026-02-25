package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/delivery/dto/req"
)

type ICreateDeleteUserNotificationSettingsUseCase interface {
	Execute(ctx context.Context) error
}
type IUpdateUserNotificationSettingsUseCase interface {
	Execute(ctx context.Context, req *req.UpdateUserNotificationSettingsRequest) (*response.Response, error)
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
	CreateNotificationSettings     ICreateDeleteUserNotificationSettingsUseCase
	UpdateUserNotificationSettings IUpdateUserNotificationSettingsUseCase
	CreateNotificationTemplate     ICreateDeleteNotificationTemplateUseCase
	UpdateNotificationTemplate     IUpdateNotificationTemplateUseCase
	SendNotificationType           ISendNotificationTypeUseCase
	ScheduleSend                   IScheduleSendUseCase
}

func NewUseCase(CreateNotificationSettings ICreateDeleteUserNotificationSettingsUseCase,
	UpdateUserNotificationSettings IUpdateUserNotificationSettingsUseCase,
	CreateNotificationTemplate ICreateDeleteNotificationTemplateUseCase,
	UpdateNotificationTemplate IUpdateNotificationTemplateUseCase,

	SendNotificationType ISendNotificationTypeUseCase) *UseCase {
	return &UseCase{
		CreateNotificationSettings:     CreateNotificationSettings,
		UpdateUserNotificationSettings: UpdateUserNotificationSettings,
		CreateNotificationTemplate:     CreateNotificationTemplate,
		UpdateNotificationTemplate:     UpdateNotificationTemplate,
		SendNotificationType:           SendNotificationType,
	}
}
