package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
)

type ReactionPostUsecase struct {
	pool     IRepositoryShare.IWorkerPool
	eventBus events.EventBus
}

func NewReactionPostUsecase(pool IRepositoryShare.IWorkerPool, eventBus events.EventBus) *ReactionPostUsecase {
	return &ReactionPostUsecase{
		pool:     pool,
		eventBus: eventBus,
	}
}

func (u *ReactionPostUsecase) Execute(ctx context.Context, req *interactionEvent.EntityReactionPayload) (*response.Response, error) {
	// Implement the logic for reacting to a post here
	payload := interactionEvent.EntityReactionPayload{
		TargetID:     req.TargetID,
		UserID:       req.UserID,
		TargetType:   req.TargetType,
		ReactionCode: req.ReactionCode,
		CreatedAt:    req.CreatedAt,
		Topic:        req.Topic,
	}
	err := u.eventBus.Publish(ctx, constants.TopicReactPost.String(), req.UserID.String(), string(req.Topic), payload)
	if err != nil {
		return nil, err
	}

	return response.NewResponse(response.WithData(req),
		response.WithMessage("Reaction added successfully"), response.WithStatus(http.StatusOK)), nil
}
