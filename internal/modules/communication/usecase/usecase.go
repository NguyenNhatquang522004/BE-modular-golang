package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
)

type ICreatePrivateConversationUsecase interface {
	Execute(ctx context.Context, req *req.CreatePrivateConversationRequest) (*res.FailedPrivateConversationResponse, error)
}
type IUpdatePrivateConversationUsecase interface {
	Execute(ctx context.Context, req *req.UpdatePrivateConversationRequest) (*res.FailedPrivateConversationResponse, error)
}
type IDeletePrivateConversationUsecase interface {
	Execute(ctx context.Context, req *req.DeletePrivateConversationRequest) (*res.FailedPrivateConversationResponse, error)
}
type ICreateChannelGroupConversationUsecase interface {
	Execute(ctx context.Context, req *req.CreateChannelGroupConversationRequest) (*res.FailedPrivateConversationResponse, error)
}
type IUpdateChannelGroupConversationUsecase interface {
	Execute(ctx context.Context, req *req.UpdateChannelConversationRequest) (*res.FailedChannelGroupConversationResponse, error)
}
type IDeleteChannelGroupConversationUsecase interface {
	Execute(ctx context.Context, req *req.DeleteChannelConversationRequest) (*res.FailedChannelGroupConversationResponse, error)
}
type ICreateGroupToChannelConversationUsecase interface {
	Execute(ctx context.Context, req *req.CreateGroupToChannelConversationRequest) (*res.FailedChannelGroupConversationResponse, error)
}
type ICreateGroupConversationUsecase interface {
	Execute(ctx context.Context, req *req.CreateGroupConversationRequest) (*res.FailedChannelGroupConversationResponse, error)
}
type IUpdateGroupConversationUsecase interface {
	Execute(ctx context.Context, req *req.UpdateGroupConversationRequest) (*res.FailedChannelGroupConversationResponse, error)
}
type IDeleteGroupConversationUsecase interface {
	Execute(ctx context.Context, req *req.DeleteGroupConversationRequest) (*res.FailedChannelGroupConversationResponse, error)
}
type IAddChannelParticipantUsecase interface {
	Execute(ctx context.Context, req *req.AddChannelParticipantRequest) (*res.FailedParticipantResponse, error)
}
type IAddGroupToChannelParticipantUsecase interface {
	Execute(ctx context.Context, req *req.AddGroupToChannelParticipantRequest) (*res.FailedParticipantResponse, error)
}
type IAddGroupParticipantUsecase interface {
	Execute(ctx context.Context, req *req.AddGroupParticipantRequest) (*res.FailedParticipantResponse, error)
}
type IRemoveGroupParticipantUsecase interface {
	Execute(ctx context.Context, req *req.RemoveGroupParticipantRequest) (*res.FailedParticipantResponse, error)
}
type IUpdateGroupParticipantUsecase interface {
	Execute(ctx context.Context, req *req.UpdateGroupParticipantRequest) (*res.FailedParticipantResponse, error)
}
type ICreateMessageUsecase interface {
	Execute(ctx context.Context) error
}
type IDeleteMessageUsecase interface {
	Execute(ctx context.Context) error
}
type IUpdateMessageUsecase interface {
	Execute(ctx context.Context) error
}
type IReactMessageUsecase interface {
	Execute(ctx context.Context) error
}
type IReadMessageUsecase interface {
	Execute(ctx context.Context) error
}
type IUnreadMessageUsecase interface {
	Execute(ctx context.Context) error
}
type IMessageReplyStoryUsecase interface {
	Execute(ctx context.Context) error
}
type IGetCoversationListUsecase interface {
	Execute(ctx context.Context) error
}

type Usecase struct {
}

func NewUsecase() *Usecase {
	return &Usecase{}
}
