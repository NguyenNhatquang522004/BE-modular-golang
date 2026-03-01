package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
)

type DeleteChannelGroupConversationUsecase struct {
	
}

func NewDeleteChannelGroupConversationUsecase() *DeleteChannelGroupConversationUsecase {
	return &DeleteChannelGroupConversationUsecase{}
}

func (u *DeleteChannelGroupConversationUsecase) Execute(ctx context.Context, req *req.DeleteChannelConversationRequest) (*res.FailedChannelGroupConversationResponse, error) {

	return &res.FailedChannelGroupConversationResponse{
		ChannelID:    req.ChannelID,
		ErrorMessage: nil,
	}, nil
}
