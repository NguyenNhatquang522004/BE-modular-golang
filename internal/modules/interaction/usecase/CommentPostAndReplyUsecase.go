package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/notificationEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type CommentPostAndReplyUsecase struct {
	commentRepo IRepositoryMongoDB.ICommentRepository
	pool        IRepositoryShare.IWorkerPool
	eventBus    events.EventBus
	contentGRPC pb.ContentServiceClient
}

func NewCommentPostAndReplyUsecase(commentRepo IRepositoryMongoDB.ICommentRepository, pool IRepositoryShare.IWorkerPool, eventBus events.EventBus, contentGRPC pb.ContentServiceClient) *CommentPostAndReplyUsecase {
	return &CommentPostAndReplyUsecase{
		commentRepo: commentRepo,
		pool:        pool,
		eventBus:    eventBus,
		contentGRPC: contentGRPC,
	}
}

func (u *CommentPostAndReplyUsecase) Execute(ctx context.Context, req *req.CommentPostAndReplyRequest) (*response.Response, error) {
	// Implement the logic for posting a comment or replying to a comment here
	entity, err := mapper.ToEntityComment(req.CreateCommentReq)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	err = u.commentRepo.CreateComment(ctx, entity)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	resa, err := u.contentGRPC.GetOwnerPost(ctx, &pb.GetOwnerPostRequest{PostId: entity.PostID.Hex()})
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	log.Printf("Owner of the post: %s", resa.OwnerId)
	notificationPayloadSendMention := &notificationEvent.NotificationPayload{
		UserID:      entity.UserID,
		SendUser:    req.Mentions,
		PreviewData: "người " + entity.UserID + " đã nhắc đến bạn trong một bình luận",
	}
	err = u.eventBus.Publish(ctx, string(constants.TopicSendNotificationType), "", string(constants.None), notificationPayloadSendMention)
	if err != nil {
		log.Printf("Failed to publish mention notification: %v", err)
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	notificationPayloadSendOWner := &notificationEvent.NotificationPayload{
		UserID:      entity.UserID,
		SendUser:    []string{resa.OwnerId},
		PreviewData: "người " + entity.UserID + " đã bình luận về bài viết của bạn",
	}
	err = u.eventBus.Publish(ctx, string(constants.TopicSendNotificationType), "", string(constants.None), notificationPayloadSendOWner)
	if err != nil {
		log.Printf("Failed to publish owner notification: %v", err)
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	conterpost := &interactionEvent.CounterPostPayload{
		PostID: entity.PostID.String(),
		UserID: entity.UserID,
	}
	err = u.eventBus.Publish(ctx, constants.TopicCounterPost.String(), entity.UserID, constants.Created.String(), conterpost)
	if err != nil {
		log.Printf("Failed to publish counter post event: %v", err)
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()), response.WithStatus(http.StatusBadRequest)), err
	}
	return response.NewResponse(response.WithData(""),
		response.WithMessage("Comment posted successfully"), response.WithStatus(http.StatusOK)), nil
}
