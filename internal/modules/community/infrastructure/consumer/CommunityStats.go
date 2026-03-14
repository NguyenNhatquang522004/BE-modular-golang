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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type CommunityStats struct {
	events          events.EventBus
	pool            IRepositoryShare.IWorkerPool
	redisRepo       IRepositoryShare.IRedis
	groupEventRepo  IRepositoryMongodb.IGroupEventsRepository
	groupMemberRepo IRepositoryMongodb.IGroupMembersRepository
	groupFile       IRepositoryMongodb.IGroupFilesRepository
}

func NewCommunityStats(events events.EventBus, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis, groupEventRepo IRepositoryMongodb.IGroupEventsRepository, groupMemberRepo IRepositoryMongodb.IGroupMembersRepository, groupFile IRepositoryMongodb.IGroupFilesRepository) *CommunityStats {
	return &CommunityStats{
		events:          events,
		pool:            pool,
		redisRepo:       redisRepo,
		groupEventRepo:  groupEventRepo,
		groupMemberRepo: groupMemberRepo,
		groupFile:       groupFile,
	}

}

func (c *CommunityStats) ConsumerCommunityStats(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicCommunityStats.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
				// Process the event here
				switch event.Type {
				case constants.Created.String():
					processErr = c.handleCreatedEvent(ctx, event)
					// Handle created event
				case constants.Updated.String():
					processErr = c.handleUpdatedEvent(ctx, event)
					// Handle updated event
				case constants.Deleted.String():
					processErr = c.handleDeletedEvent(ctx, event)
					// Handle deleted event
				default:
					processErr = errors.New("unknown event type: " + event.Type)
				}
			})
			if err != nil {
				errchan <- errors.New("failed to run worker for event " + event.ID + ": " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done()
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				c.redisRepo.MarkCompleted(ctx, event.ID)
				errchan <- nil
			}
		}
		wg.Wait()
		close(errchan)
		var finalErr error
		for err := range errchan {
			if finalErr == nil {
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
func (c *CommunityStats) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement logic to handle created event
	data, err := utils.ParsePayload[communityEvent.CommunityStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.LastActiveAt != nil {
		datamember, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, data.UserID, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group member: %w", err)
		}
		datamember.LastActiveAt = *data.LastActiveAt
		err = c.groupMemberRepo.UpdateGroupMember(ctx, datamember)
		if err != nil {
			return fmt.Errorf("failed to update group member: %w", err)
		}
		// Use datamember as needed
	}
	if data.Going != nil || data.Interested != nil {
		dataEvent, err := c.groupEventRepo.GetGroupEventByID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group event: %w", err)
		}
		if data.Going != nil {
			dataEvent.AttendeesCount.Going = dataEvent.AttendeesCount.Going + *data.Going
		}
		if data.Interested != nil {
			dataEvent.AttendeesCount.Interested = dataEvent.AttendeesCount.Interested + *data.Interested
		}
		err = c.groupEventRepo.UpdateGroupEvent(ctx, dataEvent)
		if err != nil {
			return fmt.Errorf("failed to update group event: %w", err)
		}
	}
	if data.DownloadCount != nil {
		dataFile, err := c.groupFile.GetGroupFileByID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group file: %w", err)
		}
		dataFile.DownloadCount = dataFile.DownloadCount + *data.DownloadCount
		err = c.groupFile.UpdateGroupFile(ctx, dataFile)
		if err != nil {
			return fmt.Errorf("failed to update group file: %w", err)
		}
	}
	// Use the parsed data
	return nil
}

func (c *CommunityStats) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement logic to handle updated event
	// Implement logic to handle created event
	data, err := utils.ParsePayload[communityEvent.CommunityStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.LastActiveAt != nil {
		datamember, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, data.UserID, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group member: %w", err)
		}
		datamember.LastActiveAt = *data.LastActiveAt
		err = c.groupMemberRepo.UpdateGroupMember(ctx, datamember)
		if err != nil {
			return fmt.Errorf("failed to update group member: %w", err)
		}
		// Use datamember as needed
	}
	if data.Going != nil || data.Interested != nil {
		dataEvent, err := c.groupEventRepo.GetGroupEventByID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group event: %w", err)
		}
		if data.Going != nil {
			dataEvent.AttendeesCount.Going = dataEvent.AttendeesCount.Going + *data.Going
		}
		if data.Interested != nil {
			dataEvent.AttendeesCount.Interested = dataEvent.AttendeesCount.Interested + *data.Interested
		}
		err = c.groupEventRepo.UpdateGroupEvent(ctx, dataEvent)
		if err != nil {
			return fmt.Errorf("failed to update group event: %w", err)
		}
	}
	if data.DownloadCount != nil {
		dataFile, err := c.groupFile.GetGroupFileByID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group file: %w", err)
		}
		dataFile.DownloadCount = dataFile.DownloadCount + *data.DownloadCount
		err = c.groupFile.UpdateGroupFile(ctx, dataFile)
		if err != nil {
			return fmt.Errorf("failed to update group file: %w", err)
		}
	}
	// Use the parsed data
	return nil
}

func (c *CommunityStats) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	// Implement logic to handle deleted event
	// Implement logic to handle created event
	data, err := utils.ParsePayload[communityEvent.CommunityStatsPayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf(""))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	if data.LastActiveAt != nil {
		datamember, err := c.groupMemberRepo.GetGroupMemberByUserIDAndGroupID(ctx, data.UserID, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group member: %w", err)
		}
		datamember.LastActiveAt = *data.LastActiveAt
		err = c.groupMemberRepo.UpdateGroupMember(ctx, datamember)
		if err != nil {
			return fmt.Errorf("failed to update group member: %w", err)
		}
		// Use datamember as needed
	}
	if data.Going != nil || data.Interested != nil {
		dataEvent, err := c.groupEventRepo.GetGroupEventByID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group event: %w", err)
		}
		if data.Going != nil {
			dataEvent.AttendeesCount.Going = dataEvent.AttendeesCount.Going - *data.Going
		}
		if data.Interested != nil {
			dataEvent.AttendeesCount.Interested = dataEvent.AttendeesCount.Interested - *data.Interested
		}
		err = c.groupEventRepo.UpdateGroupEvent(ctx, dataEvent)
		if err != nil {
			return fmt.Errorf("failed to update group event: %w", err)
		}
	}
	if data.DownloadCount != nil {
		dataFile, err := c.groupFile.GetGroupFileByID(ctx, data.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get group file: %w", err)
		}
		dataFile.DownloadCount = dataFile.DownloadCount - *data.DownloadCount
		err = c.groupFile.UpdateGroupFile(ctx, dataFile)
		if err != nil {
			return fmt.Errorf("failed to update group file: %w", err)
		}
	}
	// Use the parsed data
	return nil
}
func (c *CommunityStats) ConsumerFailedCommunityStats(ctx context.Context) error {
	return nil
}
