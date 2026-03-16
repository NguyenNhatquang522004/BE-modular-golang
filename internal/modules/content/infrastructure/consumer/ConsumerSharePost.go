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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerSharePost struct {
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	postRepo      IRepositoryMongodb.IPostRepository
	mediaRepo     IRepositoryMongodb.IPostMediaRepository
	extensionRepo IRepositoryMongodb.IPostExtensionRepository
	settingRepo   IRepositoryMongodb.IPostSettingRepository
	insightRepo   IRepositoryCassandra.IPostInsights
	pbcommunity   pb.CommunityServiceClient
}

func NewConsumerSharePost(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, postRepo IRepositoryMongodb.IPostRepository, mediaRepo IRepositoryMongodb.IPostMediaRepository, extensionRepo IRepositoryMongodb.IPostExtensionRepository, settingRepo IRepositoryMongodb.IPostSettingRepository, insightRepo IRepositoryCassandra.IPostInsights, pbcommunity pb.CommunityServiceClient) *ConsumerSharePost {
	return &ConsumerSharePost{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		postRepo:      postRepo,
		mediaRepo:     mediaRepo,
		extensionRepo: extensionRepo,
		settingRepo:   settingRepo,
		insightRepo:   insightRepo,
		pbcommunity:   pbcommunity,
	}
}

func (c *ConsumerSharePost) ConsumerSharePost(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicSharePost.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			status, can, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- err
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
				switch event.Type {
				case constants.Created.String():
					processErr = c.handleCreateSharePost(ctx, event)
				case constants.Updated.String():
					processErr = c.handleUpdateSharePost(ctx, event)
				case constants.Deleted.String():
					processErr = c.handleDeleteSharePost(ctx, event)
				default:
					processErr = errors.New("unsupported event type: " + event.Type)
					processErr = fmt.Errorf("event %s has unsupported type %s: %w", event.ID, event.Type, processErr)
				}
				if processErr != nil {
					errchan <- processErr
					c.redisRepo.Unlock(ctx, event.ID)
				} else {
					errchan <- nil
					c.redisRepo.MarkCompleted(ctx, event.ID)

				}
				if err != nil {
					errchan <- errors.New("failed to run worker for event " + event.ID + ": " + err.Error())
					c.redisRepo.Unlock(ctx, event.ID) // Mở khoá nếu không chạy được worker
					wg.Done()
				}
			})

		}
		var finalError error
		for err := range errchan {
			if err != nil {
				finalError = errors.Join(finalError, err)
			}
		}
		return finalError
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerSharePost) handleCreateSharePost(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[contentEvent.SharePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to get post by ID %s: %w", data.PostID, err)
	}
	if datapost == nil {
		return fmt.Errorf("post with ID %s not found", data.PostID)
	}
	if datapost.Privacy.AllowShare == false {
		return fmt.Errorf("post with ID %s does not allow sharing", data.PostID)
	}
	if datapost.Context.Type == sharedEnums.ContextTypeGroup {
		datagroup, err := c.pbcommunity.GetGroupInfo(ctx, &pb.GetGroupInfoRequest{GroupId: datapost.Context.TargetID})
		if err != nil {
			return fmt.Errorf("failed to get group info for group ID %s: %w", datapost.Context.TargetID, err)
		}
		if datagroup == nil {
			return fmt.Errorf("group with ID %s not found", datapost.Context.TargetID)
		}
		if datagroup.Protected == sharedEnums.ScopePrivate.String() {
			return fmt.Errorf("post with ID %s is in a private group and cannot be shared", data.PostID)
		}
	}
	convertID, err := primitive.ObjectIDFromHex(data.PostID)
	if err != nil {
		return fmt.Errorf("invalid post ID %s: %w", data.PostID, err)
	}
	timestamp := time.Now()
	workercount := 4 // Giới hạn số lượng worker đồng thời để tránh quá tải hệ thống
	var wg sync.WaitGroup
	errchan := make(chan error, workercount)
	for i := 0; i < workercount; i++ {
		wg.Add(1)
		err := c.pool.Run(ctx, func() {
			defer wg.Done()
			switch i {
			case 0:

				sharepost := datapost
				sharepost.ID = convertID
				sharepost.IsShared = true
				if data.GroupID != nil {
					sharepost.Context = &entity.PostContext{
						Type:     sharedEnums.ContextTypeGroup,
						TargetID: *data.GroupID,
					}
				}
				if data.PageID != nil {
					sharepost.Context = &entity.PostContext{
						Type:     sharedEnums.ContextTypePage,
						TargetID: *data.PageID,
					}
				} else {
					sharepost.Context = &entity.PostContext{
						Type:     sharedEnums.ContextTypeUserWall,
						TargetID: data.UserID,
					}
				}
				sharepost.Privacy.AllowComment = true
				sharepost.Privacy.AllowShare = true
				sharepost.UserID = data.UserID
				sharepost.PublishedAt = &timestamp
				sharepost.CreatedAt = timestamp
				sharepost.UpdatedAt = timestamp
				sharepost.Stats = entity.PostStats{
					TotalReactions: 0,
					Shares:         0,
					Comments:       0,
					Views:          0,
					Like:           0,
					Love:           0,
					Haha:           0,
					Wow:            0,
					Sad:            0,
					Angry:          0,
				}
				_, err = c.postRepo.CreatePost(ctx, sharepost)
				if err != nil {
					errchan <- fmt.Errorf("failed to create share post with ID %s: %w", sharepost.ID.Hex(), err)
					return
				}
				errchan <- nil
				datapost.Stats.Shares += 1
				_, err = c.postRepo.UpdatePost(ctx, datapost)
				if err != nil {
					errchan <- fmt.Errorf("failed to update share count for original post with ID %s: %w", datapost.ID.Hex(), err)
					return
				}
				errchan <- nil
			case 1:
				datamedia, err := c.mediaRepo.GetByPostID(ctx, data.PostID)
				if err != nil {
					errchan <- fmt.Errorf("failed to get media by post ID %s: %w", data.PostID, err)
					return
				}
				errchan <- nil
				datamedianew := datamedia
				datamedianew.ID = primitive.NewObjectID()
				datamedianew.PostID = convertID
				err = c.mediaRepo.CreatePostMedia(ctx, datamedianew)
				if err != nil {
					errchan <- fmt.Errorf("failed to create media for share post with ID %s: %w", convertID.Hex(), err)
					return
				}
				errchan <- nil
			case 2:
				dataextension, err := c.extensionRepo.GetByPostID(ctx, data.PostID)
				if err != nil {
					errchan <- fmt.Errorf("failed to get post extension by post ID %s: %w", data.PostID, err)
					return
				}
				errchan <- nil
				dataextensionnew := dataextension
				dataextensionnew.ID = primitive.NewObjectID()
				dataextensionnew.PostID = convertID
				dataextensionnew.ShareData.ParentPostID = dataextension.PostID
				dataextensionnew.ShareData.OriginalPostID = dataextension.ShareData.OriginalPostID
				err = c.extensionRepo.CreatePostExtension(ctx, dataextensionnew)
				if err != nil {
					errchan <- fmt.Errorf("failed to create post extension for share post with ID %s: %w", convertID.Hex(), err)
					return
				}
				errchan <- nil
			case 3:
				datasetting, err := c.settingRepo.GetPostSettingByPostID(ctx, data.PostID)
				if err != nil {
					errchan <- fmt.Errorf("failed to get post setting by post ID %s: %w", data.PostID, err)
					return
				}
				errchan <- nil
				datasettingnew := datasetting
				datasettingnew.ID = primitive.NewObjectID()
				datasettingnew.PostID = convertID
				_, err = c.settingRepo.CreatePostSetting(ctx, datasettingnew)
				if err != nil {
					errchan <- fmt.Errorf("failed to create post setting for share post with ID %s: %w", convertID.Hex(), err)
					return
				}
				errchan <- nil
			}
		})
		if err != nil {
			err = fmt.Errorf("failed to run worker for event %s: %w", event.ID, err)
			errchan <- err
			wg.Done() // Đảm bảo giảm wg nếu không chạy được worker
		}
	}
	wg.Wait()
	close(errchan)
	var finalError error
	for err := range errchan {
		if err != nil {
			finalError = errors.Join(finalError, err)
		}
	}
	return finalError
}
func (c *ConsumerSharePost) handleUpdateSharePost(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[contentEvent.SharePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}

	return nil
}

func (c *ConsumerSharePost) handleDeleteSharePost(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[contentEvent.SharePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(" failed to parse payload for event %s: %w", event.ID, err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.postRepo.DeletePost(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to delete share post with ID %s: %w", data.PostID, err)
	}
	err = c.mediaRepo.DeleteByPostID(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to delete media for share post with ID %s: %w", data.PostID, err)
	}
	err = c.extensionRepo.DeleteByPostID(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to delete post extension for share post with ID %s: %w", data.PostID, err)
	}
	err = c.settingRepo.DeletePostSetting(ctx, data.PostID)
	if err != nil {
		return fmt.Errorf("failed to delete post setting for share post with ID %s: %w", data.PostID, err)
	}
	return nil
}

func (c *ConsumerSharePost) ConsumerFailedSharePost(ctx context.Context) error {
	return nil
}
