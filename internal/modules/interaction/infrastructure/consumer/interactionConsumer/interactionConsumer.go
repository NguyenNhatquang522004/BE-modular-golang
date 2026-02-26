package interactionConsumer

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
)

type InteractionConsumer struct {
	pool                IRepositoryShare.IWorkerPool
	eventbus            events.EventBus
	commentRepo         IRepositoryMongoDB.ICommentRepository
	entityReactionRepo  IRepositoryCassandra.IReactionsRepository
	reactionHistoryRepo IRepositoryCassandra.IReactionHistoryRepository
}

func NewInteractionConsumer(pool IRepositoryShare.IWorkerPool, eventbus events.EventBus, commentRepo IRepositoryMongoDB.ICommentRepository, entityReactionRepo IRepositoryCassandra.IReactionsRepository, reactionHistoryRepo IRepositoryCassandra.IReactionHistoryRepository) *InteractionConsumer {
	return &InteractionConsumer{
		pool:                pool,
		eventbus:            eventbus,
		commentRepo:         commentRepo,
		entityReactionRepo:  entityReactionRepo,
		reactionHistoryRepo: reactionHistoryRepo,
	}
}

func (c *InteractionConsumer) ConsumerReactionComment(ctx context.Context) error {
	// Implement the logic for consuming reaction comment events here
	flag := true
	err := c.eventbus.Subscribe(ctx, constants.TopicReactComment.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data := event.Payload.(*interactionEvent.CommentReactionPayload)
		if data.TargetType != sharedEnums.ReactionTargetComment {
			flag = false
			return nil
		}
		switch event.Type {
		case constants.Created.String():
			dataComment, err := c.commentRepo.GetCommentByID(ctx, data.CommentID)
			if err != nil {
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				mappingCreateReactioncode(dataComment, &data.ReactionCode)
				err = c.commentRepo.UpdateComment(safectx, dataComment)
				if err != nil {
					return
				}
				entityReaction := &entity.EntityReaction{
					TargetID:     data.CommentID,
					UserID:       data.UserID,
					TargetType:   data.TargetType,
					ReactionCode: data.ReactionCode,
					CreatedAt:    data.CreatedAt,
				}
				err = c.entityReactionRepo.CreateReaction(safectx, entityReaction)
				if err != nil {
					return
				}
				reactionHistory := &entity.UserReactionHistory{
					UserID:       data.UserID,
					TargetID:     data.CommentID,
					TargetType:   data.TargetType,
					ReactionCode: data.ReactionCode,
					CreatedAt:    data.CreatedAt,
				}
				err = c.reactionHistoryRepo.CreateReactionHistory(safectx, reactionHistory)
				if err != nil {
					return
				}
			})

		case constants.Deleted.String():
			dataComment, err := c.commentRepo.GetCommentByID(ctx, data.CommentID)
			if err != nil {
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				mappingDeleteReactioncode(dataComment, &data.ReactionCode)
				err = c.commentRepo.UpdateComment(safectx, dataComment)
				if err != nil {
					return
				}
				err = c.entityReactionRepo.DeleteReaction(safectx, data.CommentID, data.UserID)
				if err != nil {
					return
				}
			})
		default:
			// Handle unknown event types if necessary
		}
		return nil
	})
	// Đảm bảo tất cả tác vụ đã hoàn thành trước khi trả về
	c.pool.Wait()
	if !flag {
		return errors.New("invalid event type for comment reaction")
	}
	if err != nil {
		return err
	}
	return nil
}

func mappingCreateReactioncode(datacomment *entity.Comment, react *sharedEnums.ReactionCode) {
	datacomment.Reactions.Total += 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Reactions.Like += 1
	case sharedEnums.ReactionLove:
		datacomment.Reactions.Love += 1
	case sharedEnums.ReactionHaha:
		datacomment.Reactions.Haha += 1
	case sharedEnums.ReactionSad:
		datacomment.Reactions.Sad += 1
	case sharedEnums.ReactionAngry:
		datacomment.Reactions.Angry += 1
	case sharedEnums.ReactionWow:
		datacomment.Reactions.Wow += 1
	}
}

func mappingDeleteReactioncode(datacomment *entity.Comment, react *sharedEnums.ReactionCode) {
	datacomment.Reactions.Total -= 1
	switch *react {
	case sharedEnums.ReactionLike:
		datacomment.Reactions.Like -= 1
	case sharedEnums.ReactionLove:
		datacomment.Reactions.Love -= 1
	case sharedEnums.ReactionHaha:
		datacomment.Reactions.Haha -= 1
	case sharedEnums.ReactionSad:
		datacomment.Reactions.Sad -= 1
	case sharedEnums.ReactionAngry:
		datacomment.Reactions.Angry -= 1
	case sharedEnums.ReactionWow:
		datacomment.Reactions.Wow -= 1
	}
}
