package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
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

func (u *ReactionPostUsecase) Execute(ctx context.Context, req *contentEvent.PostStatsPayload) (*response.Response, error) {
	// Implement the logic for reacting to a post here
	payload := &contentEvent.PostStatsPayload{
		PostID:         req.PostID,
		UserID:         req.UserID,
		TotalReactions: 0,
		Comments:       0,
		Shares:         0,
		Views:          0,
		Like:           0,
		Love:           0,
		Haha:           0,
		Wow:            0,
		Sad:            0,
		Angry:          0,
		Type:           req.Type,
	}

	err := u.eventBus.Publish(ctx, constants.TopicPostStats.String(), req.UserID, req.Type.String(), payload)
	if err != nil {
		return nil, err
	}

	return response.NewResponse(response.WithData(req),
		response.WithMessage("Reaction added successfully"), response.WithStatus(http.StatusOK)), nil
}
