package usecase

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communicationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/delivery/dto/res"
)

type DeletePrivateConversationGroupUsecase struct {
	events events.EventBus
}

func NewDeletePrivateConversationGroupUsecase(events events.EventBus) *DeletePrivateConversationGroupUsecase {
	return &DeletePrivateConversationGroupUsecase{
		events: events,
	}
}
func (u *DeletePrivateConversationGroupUsecase) Execute(ctx context.Context, req *req.DeletePrivateConversationGroupRequest) (*res.FailedPrivateConversationResponse, error) {
	payload := &communicationEvent.DeletePrivateConversationGroupPayload{
		TargetID: req.GroupID,
	}
	err := u.events.Publish(ctx, constants.TopicDeleteRelationTarget.String(), req.GroupID, constants.Deleted.String(), payload)
	if err != nil {
		return nil, err
	}
	return &res.FailedPrivateConversationResponse{
		ConversationID:      req.GroupID,
		UserIDSenderFirst:   "",
		UserIDReceiverFirst: "",
		ErrorMessage:        nil,
	}, nil
}
