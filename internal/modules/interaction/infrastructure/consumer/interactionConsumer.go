package consumer

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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/utils"
	"github.com/gocql/gocql"
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
	type taskResult struct {
		targetID string
		userID   gocql.UUID
		err      error
	}
	worerCount := 2 // Số lượng worker trong pool
	resultChan := make(chan taskResult, 2)
	flag := true
	err := c.eventbus.Subscribe(ctx, constants.TopicReactComment.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data := event.Payload.(*interactionEvent.EntityReactionPayload)
		if data.TargetType != sharedEnums.ReactionTargetComment {
			flag = false
			return nil
		}
		switch event.Type {
		case constants.Created.String():
			dataComment, err := c.commentRepo.GetCommentByID(ctx, data.TargetID)
			if err != nil {
				return err
			}
			if dataComment == nil {
				return errors.New("comment not found for reaction")
			}
			for i := 0; i < worerCount; i++ {
				c.pool.Run(ctx, func() {
					safectx := context.WithoutCancel(ctx)
					switch i {
					case 0:
						entityReaction := &entity.EntityReaction{
							TargetID:     dataComment.ID.Hex(),
							UserID:       data.UserID,
							TargetType:   data.TargetType,
							ReactionCode: data.ReactionCode,
							CreatedAt:    data.CreatedAt,
						}
						err = c.entityReactionRepo.CreateReaction(safectx, entityReaction)
						if err != nil {
							return
						}
					case 1:
						convertedTargetID, err := gocql.ParseUUID(data.TargetID)
						if err != nil {
							resultChan <- taskResult{targetID: data.TargetID, userID: data.UserID, err: err}
							return
						}
						datahistory, err := c.reactionHistoryRepo.GetReactionHistoryByUserIDAndTargetID(ctx, data.UserID, convertedTargetID)
						if err != nil {
							resultChan <- taskResult{targetID: data.TargetID, userID: data.UserID, err: err}
							return
						}
						if datahistory != nil {
							datahistory.ReactionCode = data.ReactionCode
							datahistory.CreatedAt = data.CreatedAt
							datahistory.TargetType = data.TargetType
							err = c.reactionHistoryRepo.UpdateReactionHistory(ctx, datahistory)
							if err != nil {
								resultChan <- taskResult{targetID: data.TargetID, userID: data.UserID, err: err}
								return
							}
						} else {
							reactionHistory := &entity.UserReactionHistory{
								UserID:       data.UserID,
								TargetID:     data.TargetID,
								TargetType:   data.TargetType,
								ReactionCode: data.ReactionCode,
								CreatedAt:    data.CreatedAt,
							}
							err = c.reactionHistoryRepo.CreateReactionHistory(safectx, reactionHistory)
							if err != nil {
								resultChan <- taskResult{targetID: data.TargetID, userID: data.UserID, err: err}
								return
							}
						}
					}
				})
			}
		case constants.Deleted.String():
			dataComment, err := c.commentRepo.GetCommentByID(ctx, data.TargetID)
			if err != nil {
				return err
			}
			err = c.pool.Run(ctx, func() {
				safectx := context.WithoutCancel(ctx)
				utils.MappingDeleteReactioncode(dataComment, &data.ReactionCode)
				err = c.commentRepo.UpdateComment(safectx, dataComment)
				if err != nil {
					return
				}
				err = c.entityReactionRepo.DeleteReaction(safectx, data.TargetID, data.UserID)
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

func (c *InteractionConsumer) ConsumerEntityReaction(ctx context.Context) error {
	type resultTask struct {
		targetID string
		userID   gocql.UUID
		err      error
	}
	workerCount := 2 // Số lượng worker trong pool
	resultChan := make(chan resultTask, workerCount)
	// Implement the logic for consuming entity reaction events here
	err := c.eventbus.Subscribe(ctx, constants.TopicEntityReaction.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data := event.Payload.(*interactionEvent.EntityReactionPayload)
		// Xử lý sự kiện tương tác với entity (post hoặc comment)
		// Tùy thuộc vào event.Type (Created, Deleted, Updated), thực hiện các thao tác tương ứng
		switch event.Type {
		case constants.Created.String():
			// Handle creation of a reaction
			// Ví dụ: Cập nhật tổng số lượt thích của bài viết hoặc bình luận
			// Đồng thời lưu thông tin phản ứng vào Cassandra
			for i := 0; i < workerCount; i++ {
				c.pool.Run(ctx, func() {
					switch i {
					case 0:
						entityReaction := &entity.EntityReaction{
							TargetID:     data.TargetID,
							UserID:       data.UserID,
							TargetType:   data.TargetType,
							ReactionCode: data.ReactionCode,
							CreatedAt:    data.CreatedAt,
						}
						err := c.entityReactionRepo.CreateReaction(ctx, entityReaction)
						if err != nil {
							resultChan <- resultTask{targetID: data.TargetID, userID: data.UserID, err: err}
							return
						}
					case 1:
						convertedTargetID, err := gocql.ParseUUID(data.TargetID)
						datahistory, err := c.reactionHistoryRepo.GetReactionHistoryByUserIDAndTargetID(ctx, data.UserID, convertedTargetID)
						if err != nil {
							resultChan <- resultTask{targetID: data.TargetID, userID: data.UserID, err: err}
							return
						}
						if datahistory != nil {
							datahistory.ReactionCode = data.ReactionCode
							datahistory.CreatedAt = data.CreatedAt
							datahistory.TargetType = data.TargetType
							err = c.reactionHistoryRepo.UpdateReactionHistory(ctx, datahistory)
							if err != nil {
								resultChan <- resultTask{targetID: data.TargetID, userID: data.UserID, err: err}
								return
							}
						} else {
							reactionHistory := &entity.UserReactionHistory{
								UserID:       data.UserID,
								TargetID:     data.TargetID,
								TargetType:   data.TargetType,
								ReactionCode: data.ReactionCode,
								CreatedAt:    data.CreatedAt,
							}
							err = c.reactionHistoryRepo.CreateReactionHistory(ctx, reactionHistory)
							if err != nil {
								resultChan <- resultTask{targetID: data.TargetID, userID: data.UserID, err: err}
								return
							}
						}
					}
				})
			}
			c.pool.Wait()
		case constants.Deleted.String():
			// Handle deletion of a reaction
			// Ví dụ: Cập nhật tổng số lượt thích của bài viết hoặc bình luận
			// Đồng thời xóa thông tin phản ứng khỏi Cassandra
		default:
			// Handle unknown event types if necessary
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *InteractionConsumer) ConsumerFailedEntityReaction(ctx context.Context) error
