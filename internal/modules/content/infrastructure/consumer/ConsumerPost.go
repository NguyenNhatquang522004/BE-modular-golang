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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/contentEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/mediaEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/pkg/pb/v1"
)

type ConsumerPost struct {
	events        events.EventBus
	pool          IRepositoryShare.IWorkerPool
	redisRepo     IRepositoryShare.IRedis
	postRepo      IRepositoryMongodb.IPostRepository
	mediaRepo     IRepositoryMongodb.IPostMediaRepository
	extensionRepo IRepositoryMongodb.IPostExtensionRepository
	settingRepo   IRepositoryMongodb.IPostSettingRepository
	insightRepo   IRepositoryCassandra.IPostInsights
	pbcommunity   pb.CommunityServiceClient
	pbbusiness    pb.BusinessServiceClient
	pbSocial      pb.SocialServiceClient
}

func NewConsumerPost(events events.EventBus,
	pool IRepositoryShare.IWorkerPool,
	redisRepo IRepositoryShare.IRedis,
	postRepo IRepositoryMongodb.IPostRepository,
	mediaRepo IRepositoryMongodb.IPostMediaRepository,
	extensionRepo IRepositoryMongodb.IPostExtensionRepository,
	settingRepo IRepositoryMongodb.IPostSettingRepository,
	insightRepo IRepositoryCassandra.IPostInsights,
	pbcommunity pb.CommunityServiceClient,
	pbbusiness pb.BusinessServiceClient,
	pbSocial pb.SocialServiceClient) *ConsumerPost {
	return &ConsumerPost{
		events:        events,
		pool:          pool,
		redisRepo:     redisRepo,
		postRepo:      postRepo,
		mediaRepo:     mediaRepo,
		extensionRepo: extensionRepo,
		settingRepo:   settingRepo,
		insightRepo:   insightRepo,
		pbcommunity:   pbcommunity,
		pbbusiness:    pbbusiness,
		pbSocial:      pbSocial,
	}
}

func (c *ConsumerPost) ConsumerPost(ctx context.Context) error {
	// Xử lý logic tiêu thụ sự kiện liên quan đến bài viết
	err := c.events.SubscribeBatch(ctx, constants.TopicPost.String(), 100, time.Duration(5)*time.Second, func(ctx context.Context, events []events.IntegrationEvent) error {
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
			ev := event
			wg.Add(1)
			var processErr error
			err = c.pool.Run(ctx, func() {
				defer wg.Done()
				switch ev.Type {
				case constants.Created.String():
					processErr = c.handleCreatedPost(ctx, ev)
				case constants.Updated.String():
					processErr = c.handleUpdatedPost(ctx, ev)
				case constants.Deleted.String():
					processErr = c.handleDeletedPost(ctx, ev)
				default:
					processErr = errors.New("unknown event type: " + ev.Type)
				}
			})
			if processErr != nil {
				c.redisRepo.Unlock(ctx, ev.ID) // Thất bại -> Mở khoá để lần sau làm lại
				errchan <- errors.New("event " + ev.ID + " failed: " + processErr.Error())
			} else {
				// Thành công -> Cập nhật trạng thái đã xử lý trong Redis
				c.redisRepo.MarkCompleted(ctx, ev.ID)
				errchan <- nil
			}
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + ev.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, ev.ID) // Mở khoá nếu không chạy được worker
				wg.Done()
			}
		}
		wg.Wait()
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
func (c *ConsumerPost) handleCreatedPost(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi tiêu thụ sự kiện tạo bài viết, ví dụ: insert bài viết vào DB, xử lý media, extension/setting, v.v.
	data, err := utils.ParsePayload[contentEvent.CreatePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	var wg sync.WaitGroup
	workercount := 5
	type taskResult struct {
		datapost      *entity.Post
		datamedia     *entity.PostMedia
		dataextension *entity.PostExtension
		datasetting   *entity.PostSetting
		err           error
	}

	taskResultchan := make(chan taskResult, workercount)
	for i := 0; i < workercount; i++ {
		wg.Add(1)
		err := c.pool.Run(ctx, func() {
			defer wg.Done()
			switch i {
			case 0:
				entitypost := mapper.ToCreateEntityPostPayload(data)
				if entitypost == nil {
					taskResultchan <- taskResult{err: fmt.Errorf("failed to map payload to entity")}
					return
				}
				_, err := c.postRepo.CreatePost(ctx, entitypost)
				if err != nil {
					taskResultchan <- taskResult{err: fmt.Errorf("failed to create post in repository: %w", err)}
					return
				}
				if entitypost.Context.Type == sharedEnums.ContextTypeGroup {
					payload := &communityEvent.GroupStatsPayload{
						GroupID:            entitypost.Context.TargetID,
						MemberCount:        0,
						PostCount:          0,
						PendingMemberCount: 0,
						PendingPostCount:   1, // Tăng số lượng bài viết đang chờ duyệt lên 1
						ReportedPostCount:  0,
						EventType:          constants.Created,
					}
					err = c.events.Publish(ctx, constants.TopicGroupStats.String(), entitypost.Context.TargetID, constants.Created.String(), payload)
					if err != nil {
						taskResultchan <- taskResult{err: fmt.Errorf("failed to publish group stats event: %w", err)}
						return
					}
				}
				taskResultchan <- taskResult{datapost: entitypost, err: nil}
			case 1:
				entitymedia := mapper.ToCreateEntityPostMediaPayload(data.ID, data.Media) // PostID sẽ được gán sau khi tạo post
				if entitymedia == nil {
					taskResultchan <- taskResult{err: fmt.Errorf("failed to map media payload to entity")}
					return
				}
				err = c.mediaRepo.CreatePostMedia(ctx, entitymedia)
				if err != nil {
					taskResultchan <- taskResult{err: fmt.Errorf("failed to create post media in repository: %w", err)}
					return
				}
				var payloads []mediaEvent.CreateMediaAssetsPayload
				for index, item := range entitymedia.Items {
					itemmediaitempayload := mediaEvent.MediaItemPayload{
						PostID:       data.ID,     // PostID sẽ được gán sau khi tạo post
						AlbumID:      "",          // Chưa có album trong yêu cầu tạo post, để trống hoặc gán sau nếu có
						GroupID:      *data.Group, // Chưa có group trong yêu cầu tạo post, để trống hoặc gán sau nếu có
						MediaID:      entitymedia.Items[index].ID.Hex(),
						MediaType:    item.MediaType,
						URL:          item.URL,
						ThumbnailURL: item.ThumbnailURL,
						Metadata: mediaEvent.MetadataPayload{
							Width:     item.Metadata.Width,
							Height:    item.Metadata.Height,
							Duration:  item.Metadata.Duration,
							SizeBytes: item.Metadata.SizeBytes,
							MimeType:  item.Metadata.MimeType,
						},
						Order:    item.Order,
						Hashtags: data.Hashtags,
						TaggedUsers: func() []mediaEvent.TaggedUserPayload {
							var taggedUsers []mediaEvent.TaggedUserPayload
							for _, user := range item.TaggedUsers {
								taggedUsers = append(taggedUsers, mediaEvent.TaggedUserPayload{
									UserID: user.UserID,
									Name:   user.Name,
									X:      user.X,
									Y:      user.Y,
								})
							}
							return taggedUsers
						}(),
					}
					itempayload := mediaEvent.CreateMediaAssetsPayload{
						UserID: data.UserID,
						Items:  []mediaEvent.MediaItemPayload{itemmediaitempayload},
					}
					payloads = append(payloads, itempayload)
					err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.ID, constants.Created.String(), itempayload)
					if err != nil {
						taskResultchan <- taskResult{err: fmt.Errorf("failed to publish media asset event: %w", err)}
						return
					}
				}
				taskResultchan <- taskResult{datamedia: entitymedia, err: nil}
			case 2:
				dataprofile, err := c.pbSocial.GetInfoUserByID(ctx, &pb.UserSocialIDRequest{UserId: data.UserID})
				if err != nil {
					taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to call GetInfoUserByID: "+err.Error())}
					return
				}
				if dataprofile == nil {
					taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to get user info by ID: "+err.Error())}
					return
				}
				entityextension := mapper.ToCreateEntityPostExtensionPayload(data.ID, data.Extension, dataprofile.UserId, dataprofile.AuthorName, dataprofile.AuthorAvatar, data.Content, "") // PostID sẽ được gán sau khi tạo post
				if data.Extension != nil && entityextension == nil {
					taskResultchan <- taskResult{err: fmt.Errorf("failed to map extension payload to entity")}
					return
				}
				err = c.extensionRepo.CreatePostExtension(ctx, entityextension)
				if err != nil {
					taskResultchan <- taskResult{err: fmt.Errorf("failed to create post extension in repository: %w", err)}
					return
				}
				taskResultchan <- taskResult{dataextension: entityextension, err: nil}
			case 3:
				// PostSetting cần thông tin role của user để map đúng, nên phải xử lý sau khi có
				var entitysetting *entity.PostSetting
				switch data.Context.Type {
				case sharedEnums.ContextTypeGroup:
					res, err := c.pbcommunity.GetRoleUserInGroup(ctx, &pb.GetRoleUserInGroupRequest{
						GroupId: data.Context.TargetID,
						UserId:  data.UserID,
					})
					if err != nil {
						taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to call GetRoleUserInGroup: "+err.Error())}
						return
					}
					if res == nil {
						taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to get role user in group: "+err.Error())}
						return
					}
					convertrole, err := sharedEnums.RoleTypeString(res.Role)
					if err != nil {
						taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to convert role string to enum: "+err.Error())}
						return
					}

					entitysetting = mapper.ToCreateEntityPostSettingPayload(data.ID, data.UserID, convertrole, data.Setting)
				case sharedEnums.ContextTypeUserWall:
					entitysetting = mapper.ToCreateEntityPostSettingPayload(data.ID, data.UserID, sharedEnums.RoleTypeUser, data.Setting)
					if data.Setting != nil && entitysetting == nil {
						taskResultchan <- taskResult{err: fmt.Errorf("failed to map setting payload to entity")}
						return
					}
				case sharedEnums.ContextTypePage:
					res, err := c.pbbusiness.GetRoleUserInPage(ctx, &pb.GetRoleUserInPageRequest{
						PageId: data.Context.TargetID,
						UserId: data.UserID,
					})
					if err != nil {
						taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to call GetRoleUserInPage: "+err.Error())}
						return
					}
					if res == nil {
						taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to get role user in page: "+err.Error())}
						return
					}
					convertrole, err := sharedEnums.RoleTypeString(res.Role)
					if err != nil {
						taskResultchan <- taskResult{err: fmt.Errorf("%s", "failed to convert role string to enum: "+err.Error())}
						return
					}
					entitysetting = mapper.ToCreateEntityPostSettingPayload(data.ID, data.UserID, convertrole, data.Setting)
				default:
					taskResultchan <- taskResult{err: fmt.Errorf("unknown context type: %s", data.Context.Type.String())}
					return
				}
			}
		})
		if err != nil {
			wg.Done()
			taskResultchan <- taskResult{err: fmt.Errorf("failed to run worker: %w", err)}
			return fmt.Errorf("failed to run worker: %w", err)
		}
	}
	wg.Wait()
	close(taskResultchan)
	var finalErr error
	for result := range taskResultchan {
		if result.err != nil {
			finalErr = errors.Join(finalErr, result.err)
		}
	}
	return finalErr
}

func (c *ConsumerPost) handleUpdatedPost(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi tiêu thụ sự kiện cập nhật bài viết, ví dụ: cập nhật nội dung bài viết trong DB, xử lý thay đổi media/extension/setting, v.v.
	data, err := utils.ParsePayload[contentEvent.UpdatePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	workercount := 4
	type taskResult struct {
		datapost      *entity.Post
		datamedia     *entity.PostMedia
		dataextension *entity.PostExtension
		datasetting   *entity.PostSetting
		err           error
	}
	resultChan := make(chan taskResult, workercount)
	var wg sync.WaitGroup
	for i := 0; i < workercount; i++ {
		wg.Add(1)
		err := c.pool.Run(ctx, func() {
			defer wg.Done()
			switch i {
			case 0:
				datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
				if err != nil {
					resultChan <- taskResult{err: fmt.Errorf("failed to get post by ID: %w", err)}
					return
				}
				resultChan <- taskResult{datapost: datapost, err: nil}

			case 1:
				datamedia, err := c.mediaRepo.GetByPostID(ctx, data.PostID)
				if err != nil {
					resultChan <- taskResult{err: fmt.Errorf("failed to get post media by post ID: %w", err)}
					return
				}
				resultChan <- taskResult{datamedia: datamedia, err: nil}
			case 2:
				dataextension, err := c.extensionRepo.GetByPostID(ctx, data.PostID)
				if err != nil {
					resultChan <- taskResult{err: fmt.Errorf("failed to get post extension by post ID: %w", err)}
					return
				}
				resultChan <- taskResult{dataextension: dataextension, err: nil}
			case 3:
				datasetting, err := c.settingRepo.GetPostSettingByPostID(ctx, data.PostID)
				if err != nil {
					resultChan <- taskResult{err: fmt.Errorf("failed to get post setting by post ID: %w", err)}
					return
				}
				resultChan <- taskResult{datasetting: datasetting, err: nil}
			default:
				resultChan <- taskResult{err: nil} // Các worker còn lại không làm gì
			}
		})
		if err != nil {
			wg.Done()
			resultChan <- taskResult{err: fmt.Errorf("failed to run worker: %w", err)}
			return fmt.Errorf("failed to run worker: %w", err)

		}
	}
	wg.Wait()
	close(resultChan)
	var datapost *entity.Post
	var datamedia *entity.PostMedia
	var dataextension *entity.PostExtension
	var datasetting *entity.PostSetting
	for result := range resultChan {
		if result.err != nil {
			return result.err // Nếu có lỗi từ bất kỳ worker nào, trả về lỗi đó
		}
		if result.datapost != nil {
			datapost = result.datapost
		}
		if result.datamedia != nil {
			datamedia = result.datamedia
		}
		if result.dataextension != nil {
			dataextension = result.dataextension
		}
		if result.datasetting != nil {
			datasetting = result.datasetting
		}
	}
	errchan := make(chan error, workercount)
	for i := 0; i < workercount; i++ {
		wg.Add(1)
		err := c.pool.Run(ctx, func() {
			defer wg.Done()
			switch i {
			case 0:
				if datapost == nil {
					errchan <- fmt.Errorf("post not found with ID: %s", data.PostID)
					return
				}
				mapper.UpdateEntityPostFromPayload(datapost, data)
				_, err = c.postRepo.UpdatePost(ctx, datapost)
				if err != nil {
					errchan <- fmt.Errorf("failed to update post in repository: %w", err)
					return
				}

				errchan <- nil
			case 1:
				if data.Media != nil {
					if datamedia == nil {
						errchan <- fmt.Errorf("post media not found with post ID: %s", data.PostID)
						return
					}
					mapper.UpdateEntityPostMediaFromPayload(datamedia, data)
					err = c.mediaRepo.UpdatePostMedia(ctx, datamedia)
					if err != nil {
						errchan <- fmt.Errorf("failed to update post media in repository: %w", err)
						return
					}
					errchan <- nil
					if data.Deletemedia != nil {
						_, _, err = c.mediaRepo.DeleteBulkByIDs(ctx, *data.Deletemedia)
						if err != nil {
							errchan <- fmt.Errorf("failed to delete post media in repository: %w", err)
							return
						}
						for _, mediaID := range *data.Deletemedia {
							payload := mediaEvent.DeleteMediaAssetsPayload{
								MediaID:   mediaID,
								MessageID: "",
							}
							err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.PostID, constants.Deleted.String(), payload)
							if err != nil {
								errchan <- fmt.Errorf("failed to publish media asset deletion event: %w", err)
								return
							}
						}
						if err != nil {
							errchan <- fmt.Errorf("failed to publish media asset deletion event: %w", err)
							return
						}
					}
					errchan <- nil
				}
			case 2:
				if data.Extension != nil {
					if dataextension == nil {
						errchan <- fmt.Errorf("post extension not found with post ID: %s", data.PostID)
						return
					}
					mapper.UpdateEntityPostExtensionFromPayload(dataextension, data)
					err = c.extensionRepo.UpdatePostExtension(ctx, dataextension)
					if err != nil {
						errchan <- fmt.Errorf("failed to update post extension in repository: %w", err)
						return
					}
					errchan <- nil
				}
			case 3:
				if data.Setting != nil {
					if datasetting == nil {
						errchan <- fmt.Errorf("post setting not found with post ID: %s", data.PostID)
						return
					}
					mapper.UpdateEntityPostSettingFromPayload(datasetting, data)
					_, err = c.settingRepo.UpdatePostSetting(ctx, datasetting.PostID.Hex(), datasetting)
					if err != nil {
						errchan <- fmt.Errorf("failed to update post setting in repository: %w", err)
						return
					}
					errchan <- nil
				}
			default:
				errchan <- nil // Các worker còn lại không làm gì
			}
		})
		if err != nil {
			wg.Done()
			errchan <- fmt.Errorf("failed to run worker: %w", err)
			return fmt.Errorf("failed to run worker: %w", err)
		}
	}
	if err != nil {
		return fmt.Errorf("failed to run worker: %w", err)
	}
	wg.Wait()
	if datapost.Context.Type == sharedEnums.ContextTypeGroup && datapost.Status == sharedEnums.ProcessingActive {
		payload := &communityEvent.GroupStatsPayload{
			GroupID:            datapost.Context.TargetID,
			MemberCount:        0,
			PostCount:          1,
			PendingMemberCount: 0,
			PendingPostCount:   0, // Cập nhật lại số lượng bài viết đang chờ duyệt nếu có thay đổi về status
			ReportedPostCount:  0,
			EventType:          constants.Updated,
		}
		err = c.events.Publish(ctx, constants.TopicGroupStats.String(), datapost.Context.TargetID, constants.Updated.String(), payload)
		if err != nil {
			errchan <- fmt.Errorf("failed to publish group stats event: %w", err)
		}
	}
	if datapost.Context.Type == sharedEnums.ContextTypeGroup && datapost.Status == sharedEnums.ProcessingFailed {
		payload := &communityEvent.GroupStatsPayload{
			GroupID:            datapost.Context.TargetID,
			MemberCount:        0,
			PostCount:          0,
			PendingMemberCount: 0,
			PendingPostCount:   1, // Cập nhật lại số lượng bài viết đang chờ duyệt nếu có thay đổi về status
			ReportedPostCount:  0,
			EventType:          constants.Deleted,
		}
		err = c.events.Publish(ctx, constants.TopicGroupStats.String(), datapost.Context.TargetID, constants.Deleted.String(), payload)
		if err != nil {
			errchan <- fmt.Errorf("failed to publish group stats event: %w", err)
		}
		payloadDelete := &contentEvent.DeletePostPayload{
			PostID: data.PostID,
			Reason: "Post status changed to failed, treated as deleted in group stats",
		}
		err = c.events.Publish(ctx, constants.TopicPost.String(), data.PostID, constants.Deleted.String(), payloadDelete)
		if err != nil {
			errchan <- fmt.Errorf("failed to publish post deletion event for group stats update: %w", err)
		}
	}
	close(errchan)
	var finalErr error
	for err := range errchan {
		if err != nil {
			finalErr = errors.Join(finalErr, err)
		}
	}
	return finalErr
}
func (c *ConsumerPost) handleDeletedPost(ctx context.Context, event events.IntegrationEvent) error {
	// Xử lý logic khi tiêu thụ sự kiện xoá bài viết, ví dụ: xoá media liên quan, xoá extension/setting, v.v.
	data, err := utils.ParsePayload[contentEvent.DeletePostPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.PostID != "" && data.GroupID == "" && data.PageID == "" && data.ReelID == "" && data.UserID == "" {
		datapost, err := c.postRepo.GetPostByID(ctx, data.PostID)
		if err != nil {
			return fmt.Errorf("failed to get post by ID %s: %w", data.PostID, err)
		}
		if datapost == nil {
			return fmt.Errorf("post with ID %s not found", data.PostID)
		}
		var datamedia *entity.PostMedia
		if datapost.IsShared == true {
			datamedia = nil
		} else {
			datamedia, err = c.mediaRepo.GetByPostID(ctx, data.PostID)
			if err != nil {
				return fmt.Errorf("failed to get post media by post ID %s: %w", data.PostID, err)
			}
		}
		var wg sync.WaitGroup
		workercout := 5
		errchan := make(chan error, workercout)
		for i := 0; i < workercout; i++ {
			wg.Add(1)
			err := c.pool.Run(ctx, func() {
				defer wg.Done()
				switch i {
				case 0:
					err = c.postRepo.DeletePost(ctx, data.PostID)
					if err != nil {
						errchan <- fmt.Errorf("failed to delete post in repository: %w", err)
						return
					}
				case 1:

					err = c.mediaRepo.DeleteByPostID(ctx, data.PostID)
					if err != nil {
						errchan <- fmt.Errorf("failed to delete post media by post ID in repository: %w", err)
						return
					}
				case 2:
					err = c.extensionRepo.DeleteByPostID(ctx, data.PostID)
					if err != nil {
						errchan <- fmt.Errorf("failed to delete post extension by post ID in repository: %w", err)
						return
					}
				case 3:
					err = c.settingRepo.DeletePostSetting(ctx, data.PostID)
					if err != nil {
						errchan <- fmt.Errorf("failed to delete post setting by post ID in repository: %w", err)
						return
					}
				case 5:
					if datamedia != nil {
						for _, media := range datamedia.Items {
							payload := mediaEvent.DeleteMediaAssetsPayload{
								MediaID:   media.ID.Hex(),
								MessageID: "",
							}
							err := c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.PostID, constants.Deleted.String(), payload)
							if err != nil {
								errchan <- fmt.Errorf("failed to publish media asset deletion event: %w", err)
								return
							}
						}
					}

				default:
					errchan <- nil // Các worker còn lại không làm gì
				}
			})
			if err != nil {
				wg.Done()
				errchan <- fmt.Errorf("failed to run worker: %w", err)
				return fmt.Errorf("failed to run worker: %w", err)
			}
		}
		wg.Wait()
		close(errchan)
		var finalErr error
		for err := range errchan {
			if err != nil {
				finalErr = errors.Join(finalErr, err)
			}
		}
		return finalErr
	} else {
		var list []string
		if data.GroupID != "" {
			datagroup, err := c.postRepo.GetAllPostByGroupID(ctx, data.GroupID)
			if err != nil {
				return fmt.Errorf("failed to get posts by group ID %s: %w", data.GroupID, err)
			}
			for _, post := range datagroup {
				list = append(list, post.ID.Hex())
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.PostID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
				MediaID:   "", // Xoá tất cả media liên quan đến các bài viết trong group
				MessageID: "",
				GroupID:   data.GroupID,
				PageID:    "",
				ReelID:    "",
				StoryID:   "",
				CommentID: "",
			})
			if err != nil {
				return fmt.Errorf("failed to publish media asset deletion event for group posts: %w", err)
			}
		}
		if data.PageID != "" {
			datapage, err := c.postRepo.GetAllPostByPageID(ctx, data.PageID)
			if err != nil {
				return fmt.Errorf("failed to get posts by page ID %s: %w", data.PageID, err)
			}
			for _, post := range datapage {
				list = append(list, post.ID.Hex())
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.PostID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
				MediaID:   "", // Xoá tất cả media liên quan đến các bài viết trong group
				MessageID: "",
				GroupID:   "",
				PageID:    data.PageID,
				ReelID:    "",
				StoryID:   "",
				CommentID: "",
			})
			if err != nil {
				return fmt.Errorf("failed to publish media asset deletion event for group posts: %w", err)
			}
		}
		if data.ReelID != "" {
			datareel, err := c.postRepo.GetAllPostByReelID(ctx, data.ReelID)
			if err != nil {
				return fmt.Errorf("failed to get posts by reel ID %s: %w", data.ReelID, err)
			}
			for _, post := range datareel {
				list = append(list, post.ID.Hex())
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.PostID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
				MediaID:   "", // Xoá tất cả media liên quan đến các bài viết trong group
				MessageID: "",
				GroupID:   "",
				PageID:    "",
				ReelID:    data.ReelID,
				StoryID:   "",
				CommentID: "",
			})
			if err != nil {
				return fmt.Errorf("failed to publish media asset deletion event for group posts: %w", err)
			}
		}
		if data.UserID != "" {
			datauser, err := c.postRepo.GetAllPostByUserID(ctx, data.UserID)
			if err != nil {
				return fmt.Errorf("failed to get posts by user ID %s: %w", data.UserID, err)
			}
			for _, post := range datauser {
				list = append(list, post.ID.Hex())
			}
			err = c.events.Publish(ctx, constants.TopicMediaAsset.String(), data.PostID, constants.Deleted.String(), &mediaEvent.DeleteMediaAssetsPayload{
				MediaID:   "", // Xoá tất cả media liên quan đến các bài viết trong group
				MessageID: "",
				GroupID:   "",
				PageID:    "",
				ReelID:    "",
				StoryID:   "",
				CommentID: "",
				UserID:    data.UserID,
			})
			if err != nil {
				return fmt.Errorf("failed to publish media asset deletion event for group posts: %w", err)
			}
		}

		var wg sync.WaitGroup
		workercout := 4
		errchan := make(chan error, workercout)
		for i := 0; i < workercout; i++ {
			wg.Add(1)
			err := c.pool.Run(ctx, func() {
				defer wg.Done()
				switch i {
				case 0:
					_, _, err = c.postRepo.DeleteBulkPosts(ctx, list)
					if err != nil {
						errchan <- fmt.Errorf("failed to delete post in repository: %w", err)
						return
					}
					errchan <- nil
				case 1:
					_, _, err = c.mediaRepo.DeleteBulkByPostIDs(ctx, list)
					if err != nil {
						errchan <- fmt.Errorf("failed to publish media asset deletion event: %w", err)
						return
					}
				case 2:
					_, _, err = c.extensionRepo.DeleteBulkByPostIDs(ctx, list)
					if err != nil {
						errchan <- fmt.Errorf("failed to delete post extension by post ID in repository: %w", err)
						return
					}
				case 3:
					_, _, err = c.settingRepo.DeleteBulkPostSettings(ctx, list)
					if err != nil {
						errchan <- fmt.Errorf("failed to delete post setting by post ID in repository: %w", err)
						return
					}
				default:
					errchan <- nil // Các worker còn lại không làm gì
				}
			})
			if err != nil {
				wg.Done()
				errchan <- fmt.Errorf("failed to run worker: %w", err)
				return fmt.Errorf("failed to run worker: %w", err)
			}
		}
		wg.Wait()
		close(errchan)
		var finalErr error
		for err := range errchan {
			if err != nil {
				finalErr = errors.Join(finalErr, err)
			}
		}
		return finalErr
	}
}
func (c *ConsumerPost) ConsumerFailedPost(ctx context.Context) error {
	// Xử lý logic khi tiêu thụ sự kiện thất bại, ví dụ: ghi log, retry, v.v.
	return nil
}
