package consumer

import (
	"context"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
)

type ConsumerCounterComment struct {
	comment IRepositoryMongoDB.ICommentRepository
	events  events.EventBus
	pool    IRepositoryShare.IWorkerPool
}

func NewConsumerCounterComment(comment IRepositoryMongoDB.ICommentRepository, events events.EventBus, pool IRepositoryShare.IWorkerPool) *ConsumerCounterComment {
	return &ConsumerCounterComment{
		comment: comment,
		events:  events,
		pool:    pool,
	}
}

func (c *ConsumerCounterComment) ConsumeCounterComment(ctx context.Context) error {
	err := c.events.Subscribe(ctx, constants.TopicCounterComment.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*interactionEvent.CommentCountPayload)
		switch event.Type {
		case constants.Created.String():
			if !ok {
				log.Printf("Error casting event payload to CommentCountRequest")
				return nil
			}
			dataComment, err := c.comment.GetCommentByID(ctx, data.CommentID)
			if err != nil {
				log.Printf("Error getting comment by ID: %v", err)
				return nil
			}
			if dataComment == nil {
				log.Printf("Comment with ID %s not found", data.CommentID)
				return nil
			}
			dataComment.ReportCount = data.ReportCount + dataComment.ReportCount
			dataComment.MentionCount = data.MentionCount + dataComment.MentionCount
			dataComment.ReplyCount = data.ReplyCount + dataComment.ReplyCount
			err = c.comment.UpdateComment(ctx, dataComment)
			if err != nil {
				log.Printf("Error updating comment: %v", err)
				return nil
			}
		case constants.Updated.String():

		}

		return nil
	})
	if err != nil {
		log.Printf("Error subscribing to topic %s: %v", constants.TopicCounterComment.String(), err)
		return err
	}
	return nil
}

func (c *ConsumerCounterComment) ConsumerFailedCounterComment(ctx context.Context) error {
	return nil
}
