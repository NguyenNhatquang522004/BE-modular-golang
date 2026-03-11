package consumer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/interactionEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type CommentLiveStream struct {
	events          events.EventBus
	pool            IRepositoryShare.IWorkerPool
	redisRepo       IRepositoryShare.IRedis
	livecommentRepo IRepositoryCassandra.ILiveCommentsRepository
	pb              pb.SocialServiceClient
}

func NewCommentLiveStream(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, livecommentRepo IRepositoryCassandra.ILiveCommentsRepository, pb pb.SocialServiceClient) *CommentLiveStream {
	return &CommentLiveStream{
		events:          events,
		pool:            pool,
		redisRepo:       redisRepo,
		livecommentRepo: livecommentRepo,
		pb:              pb,
	}
}
func (c *CommentLiveStream) ConsumerLiveComment(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicCommentLive.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- errors.New("failed to acquire lock for event " + event.ID + ": " + err.Error())
				continue
			}
			if !can {
				if status == constants.StatusProcessing {
					errchan <- errors.New("event " + event.ID + " is currently being processed by another worker. Skipping.")
				} else {
					errchan <- errors.New("event " + event.ID + " has already been processed with status " + status.String() + ". Skipping.")
				}
				continue
			}
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch event.Type {
				case constants.Created.String():
					processErr = c.handleCreatedEvent(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdatedEvent(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeletedEvent(ctx, event)

				default:
					processErr = errors.New("unsupported event type: " + event.Type)
				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit task for event " + event.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
			}
		}
		wg.Done()
		close(errchan)
		var finalErr error
		for err := range errchan {
			if err != nil {
				finalErr = errors.Join(finalErr, err)
			}
		}
		return finalErr
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *CommentLiveStream) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.LiveCommentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datainfo, err := c.pb.GetInfoUserByID(ctx, &pb.UserSocialIDRequest{UserId: data.UserID})
	if err != nil {
		return fmt.Errorf("failed to get user info by ID %s: %w", data.UserID, err)
	}
	data.UserNickname = datainfo.GetAuthorName()
	data.UserAvatarURL = datainfo.GetAuthorAvatar()
	entitylivecomment := mapper.ToEntityLiveCommentPayload(data)
	err = c.livecommentRepo.CreateLiveComment(ctx, entitylivecomment)
	if err != nil {
		return fmt.Errorf("failed to create live comment in repository: %w", err)
	}
	payload := &interactionEvent.CreatedCommentPayload{
		ID:              entitylivecomment.CommentID.String(),
		UserID:          data.UserID,
		TargetID:        data.StreamID,
		Content:         data.Content,
		CreatedAt:       data.CreatedAt,
		Media:           nil, // Live comment không có media, nếu sau này có thêm thì sẽ map ở đây
		Mentions:        nil, // Live comment không có mentions, nếu sau này có thêm thì sẽ map ở đây
		ParentCommentID: nil, // Live comment không có parent comment, nếu sau này có thêm thì sẽ map ở đây
		RootCommentID:   nil, // Live comment không có root comment, nếu sau này có thêm thì sẽ map ở đây
		Type:            constants.Created,
	}
	err = c.events.Publish(ctx, constants.TopicComment.String(), data.StreamID, constants.Created.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish created comment event: %w", err)
	}
	return nil
}
func (c *CommentLiveStream) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.LiveCommentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datacommet, err := c.livecommentRepo.GetLiveCommentByCommentID(ctx, data.StreamID, data.CommentID)
	if err != nil {
		return fmt.Errorf("failed to get live comment by comment ID %s: %w", data.CommentID, err)
	}
	if datacommet == nil {
		return fmt.Errorf("live comment with comment ID %s not found", data.CommentID)
	}
	mapper.UpdateToEntityLiveCommentPayload(data, datacommet)
	err = c.livecommentRepo.UpdateLiveComment(ctx, datacommet)
	if err != nil {
		return fmt.Errorf("failed to update live comment in repository: %w", err)
	}
	payload := &interactionEvent.UpdatedCommentPayload{
		CommentID: datacommet.CommentID.String(),
		UserID:    data.UserID,
		Content:   data.Content,
		Type:      constants.Updated,
	}
	err = c.events.Publish(ctx, constants.TopicComment.String(), data.StreamID, constants.Updated.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish updated comment event: %w", err)
	}
	return nil
}
func (c *CommentLiveStream) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[mediaEvent.LiveCommentPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse event payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datacommet, err := c.livecommentRepo.GetLiveCommentByCommentID(ctx, data.StreamID, data.CommentID)
	if err != nil {
		return fmt.Errorf("failed to get live comment by comment ID %s: %w", data.CommentID, err)
	}
	if datacommet == nil {
		return fmt.Errorf("live comment with comment ID %s not found", data.CommentID)
	}
	err = c.livecommentRepo.DeleteLiveComment(ctx, data.StreamID, data.CommentID, datacommet.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to delete live comment in repository: %w", err)
	}
	payload := &interactionEvent.DeleteCommentPayload{
		CommentID: datacommet.CommentID.String(),
	}
	err = c.events.Publish(ctx, constants.TopicComment.String(), data.StreamID, constants.Deleted.String(), payload)
	if err != nil {
		return fmt.Errorf("failed to publish deleted comment event: %w", err)
	}
	return nil
}
func (c *CommentLiveStream) ConsumerFailedLiveComment(ctx context.Context) error {

	return nil
}
