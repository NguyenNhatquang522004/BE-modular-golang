package usecase

import (
	"context"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
)

type CommentReactionUsecase struct {
	pool     IRepositoryShare.IWorkerPool
	eventBus events.EventBus
}

func NewCommentReactionUsecase() *CommentReactionUsecase {
	return &CommentReactionUsecase{}
}

func (u *CommentReactionUsecase) Execute(ctx context.Context, req *interactionEvent.EntityReactionPayload) (*response.Response, error) {
	// Implement the logic for reacting to a comment here
	// err := u.eventBus.Publish(ctx, constants.TopicReactComment.String(), req.UserID.String(), string(req.Topic), req)
	// if err != nil {
	// 	return nil, err
	// }

	return response.NewResponse(response.WithData(req),
		response.WithMessage("Reaction added successfully"), response.WithStatus(http.StatusOK)), nil
}
