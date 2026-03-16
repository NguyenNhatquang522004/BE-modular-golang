package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerProfile struct {
	profileRepo IRepositoryMongodb.IProfileRepositoryMongodb
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
	redisRepo   IRepositoryShare.IRedis
	pb          pb.IdentityServiceClient
}

func NewConsumerProfile(profileRepo IRepositoryMongodb.IProfileRepositoryMongodb, events events.EventBus, pbClient pb.IdentityServiceClient) *ConsumerProfile {
	return &ConsumerProfile{
		profileRepo: profileRepo,
		events:      events,
		pb:          pbClient,
	}
}

func (c *ConsumerProfile) ConsumerProfile(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicProfile.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {

			// 1. Thử khóa event này trong Redis để đảm bảo chỉ 1 worker xử lý 1 eventID nhất định (Distributed Lock)
			status, acquired, err := c.redisRepo.Lock(ctx, event.ID)
			if err != nil {
				errchan <- fmt.Errorf("failed to acquire lock for event %s: %w", event.ID, err)
				continue
			}
			if !acquired {
				if status == constants.StatusProcessing {
					log.Printf("Event %s is currently being processed by another worker. Skipping.\n", event.ID)
					errchan <- errors.New("event is being processed by another worker")
				} else {
					log.Printf("Event %s has already been processed with status %s. Skipping.\n", event.ID, status)
				}
				continue
			}
			ev := event
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch ev.Type {
				case constants.Created.String():
					processErr = c.handleCreatedProfile(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedProfile(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedProfile(ctx, ev)
				default:
					errchan <- fmt.Errorf("unsupported event type %s for event ID %s. Marking as failed", ev.Type, ev.ID)
				}

			})
			if processErr != nil {
				c.redisRepo.Unlock(ctx, ev.ID) // Thất bại -> Mở khoá để lần sau làm lại
				errchan <- fmt.Errorf("event %s failed: %w", ev.ID, processErr)
			} else {
				// CỰC KỲ QUAN TRỌNG: Thành công -> Đánh dấu Vĩnh viễn (Hoặc 24h)
				c.redisRepo.MarkCompleted(ctx, ev.ID)
				errchan <- nil
			}
			if err != nil {
				wg.Done()
				errchan <- fmt.Errorf("failed to run event %s in worker pool: %w", ev.ID, err)
				c.redisRepo.Unlock(ctx, ev.ID) // Mở khóa ngay nếu có lỗi khi chạy goroutine
			}
		}
		wg.Wait()
		close(errchan)
		var batchErr error
		for err := range errchan {
			if err != nil {
				batchErr = errors.Join(batchErr, err)
			}
		}
		return batchErr
	})
	if err != nil {
		return err
	}
	return err
}

func (c *ConsumerProfile) handleCreatedProfile(ctx context.Context, event events.IntegrationEvent) error {
	// Gọi hàm parse chuẩn
	data, err := utils.ParsePayload[socialEvent.ProfilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	entity, err := mapper.ToEntityProfilePayload(data)
	if err != nil {
		return kafka.NewNonRetryableError(err) // Lỗi mapper thường là lỗi data sai, cũng ném vào DLQ
	}
	datausersetting, err := c.pb.GetUserSettingByID(ctx, &pb.UserSettingIDRequest{UserId: data.UserID})
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	if entity.Avatar != nil {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.Avatar.ID.Hex(),
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          entity.Avatar.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	if entity.CoverPhoto != nil {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.CoverPhoto.ID.Hex(),
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          entity.CoverPhoto.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	if entity.CVDocument != nil {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.CVDocument.FileID.Hex(),
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypePDF,
			URL:          entity.CVDocument.Filename,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	entity.Settings.AllowSearchEngine = datausersetting.AllowSearchEngineIndexing
	entity.Settings.IsPrivate = datausersetting.Allow_Profile_View_From
	// Lỗi DB thì cứ trả về bình thường để Retry
	return c.profileRepo.CreateProfile(ctx, entity)
}
func (c *ConsumerProfile) handleUpdatedProfile(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[socialEvent.ProfilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	dataprofile, err := c.profileRepo.GetProfileByID(ctx, data.UserID)
	if err != nil {
		return err // Lỗi kết nối DB -> Retry
	}
	datausersetting, err := c.pb.GetUserSettingByID(ctx, &pb.UserSettingIDRequest{UserId: data.UserID})
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}

	dataprofile.Settings.AllowSearchEngine = datausersetting.AllowSearchEngineIndexing
	dataprofile.Settings.IsPrivate = datausersetting.Allow_Profile_View_From
	if dataprofile == nil {
		// CẬP NHẬT: Tuỳ vào logic nghiệp vụ của bạn.
		// Lệnh Update mà user không tồn tại thì không bao giờ thành công được -> NonRetryableError
		return kafka.NewNonRetryableError(errors.New("profile not found for update"))
	}
	if data.Avatar != nil && data.Avatar.ID != dataprofile.Avatar.ID.Hex() {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   dataprofile.Avatar.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for old avatar: %w", err))
		}
	}
	if data.CVDocument != nil && data.CVDocument.FileID != dataprofile.CVDocument.FileID.Hex() {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   dataprofile.CVDocument.FileID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for old CV document: %w", err))
		}
	}
	if data.CoverPhoto != nil && data.CoverPhoto.ID != dataprofile.CoverPhoto.ID.Hex() {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   dataprofile.CoverPhoto.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for old cover photo: %w", err))
		}
	}
	entity, err := mapper.ToEntityUpdateProfilePayload(dataprofile, data)
	if entity.Avatar != nil {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.Avatar.ID.Hex(),
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          entity.Avatar.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	if entity.CoverPhoto != nil {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.CoverPhoto.ID.Hex(),
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypeImage,
			URL:          entity.CoverPhoto.URL,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	if entity.CVDocument != nil {
		item := &mediaEvent.CreateMediaAssetsPayload{
			UserID: data.UserID,
		}
		item.Items = append(item.Items, mediaEvent.MediaItemPayload{
			MediaID:      entity.CVDocument.FileID.Hex(),
			AlbumID:      "",
			GroupID:      "",
			PageID:       "",
			PostID:       "",
			CommentID:    "",
			StoryID:      "",
			ReelID:       "",
			MessageID:    "",
			MediaType:    sharedEnums.MediaTypePDF,
			URL:          entity.CVDocument.Filename,
			ThumbnailURL: "",
			Metadata: mediaEvent.MetadataPayload{
				Width:     0, // Cần bổ sung nếu có thông tin
				Height:    0, // Cần bổ sung nếu có thông tin
				Duration:  0,
				SizeBytes: 0,
				MimeType:  "",
			},
			Order:       0,
			Hashtags:    []string{},
			TaggedUsers: []mediaEvent.TaggedUserPayload{},
		})
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), item)
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset event: %w", err))
		}
	}
	err = c.profileRepo.UpdateProfile(ctx, entity)
	if err != nil {
		return err
	}
	return c.profileRepo.UpdateProfile(ctx, entity)
}
func (c *ConsumerProfile) handleDeletedProfile(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[socialEvent.ProfilePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("invalid payload format: %w", err))
	}
	dataprofile, err := c.profileRepo.GetProfileByID(ctx, data.UserID)
	if err != nil {
		return err // Lỗi kết nối DB -> Retry
	}
	if dataprofile == nil {
		// Lệnh Delete mà user không tồn tại thì coi như đã thành công (Idempotent)
		return nil
	}
	if dataprofile.Avatar.ID != primitive.NilObjectID {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Created.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   dataprofile.Avatar.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for avatar: %w", err))
		}
	}
	if dataprofile.CVDocument.FileID != primitive.NilObjectID {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   dataprofile.CVDocument.FileID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for CV document: %w", err))
		}
	}
	if dataprofile.CoverPhoto.ID != primitive.NilObjectID {
		err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.UserID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
			MediaID:   dataprofile.CoverPhoto.ID.Hex(),
			MessageID: "",
			GroupID:   "",
			PageID:    "",
			ReelID:    "",
			StoryID:   "",
			CommentID: "",
		})
		if err != nil {
			return kafka.NewNonRetryableError(fmt.Errorf("failed to publish media asset deletion event for cover photo: %w", err))
		}
	}
	err = c.profileRepo.DeleteProfile(ctx, data.UserID)
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerProfile) ConsumerFailedProfile(ctx context.Context) error {
	return nil
}
