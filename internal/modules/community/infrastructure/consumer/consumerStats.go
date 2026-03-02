package consumer

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/communityEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
)

type ConsumerStats struct {
	events    events.EventBus
	groupRepo IRepositoryMongodb.IGroupRepository
	eventRepo IRepositoryMongodb.IGroupEventsRepository
	pool      IRepositoryShare.IWorkerPool
}

func NewConsumerStats(events events.EventBus, groupRepo IRepositoryMongodb.IGroupRepository, eventRepo IRepositoryMongodb.IGroupEventsRepository, pool IRepositoryShare.IWorkerPool) *ConsumerStats {
	return &ConsumerStats{
		events:    events,
		groupRepo: groupRepo,
		eventRepo: eventRepo,
		pool:      pool,
	}
}

func (c *ConsumerStats) ConsumerGroupStats(ctx context.Context) {
	err := c.events.Subscribe(ctx, constants.TopicGroupStats.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*communityEvent.GroupStatsPayload)
		if !ok {
			// Log error
			return nil
		}
		// Process data and update group stats in the database
		switch data.EventType {
		case constants.Created:
			datagroup, err := c.groupRepo.GetGroupByID(ctx, data.GroupID)
			if err != nil {
				// Log error
				return nil
			}
			datagroup.Stats.MemberCount = datagroup.Stats.MemberCount + data.MemberCount
			datagroup.Stats.PostCount = datagroup.Stats.PostCount + data.PostCount
			datagroup.Stats.PendingMemberCount = datagroup.Stats.PendingMemberCount + data.PendingMemberCount
			datagroup.Stats.PendingPostCount = datagroup.Stats.PendingPostCount + data.PendingPostCount
			datagroup.Stats.ReportedPostCount = datagroup.Stats.ReportedPostCount + data.ReportedPostCount
			err = c.groupRepo.UpdateGroup(ctx, datagroup)
		case constants.Deleted:
		case constants.Updated:
		default:
			// Log unknown event type
			return nil
		}
		return nil
	})
	if err != nil {
		// Log error
	}
}
func (c *ConsumerStats) ConsumerFailedGroupStats(ctx context.Context)
func (c *ConsumerStats) ConsumerEventStats(ctx context.Context) {
	err := c.events.Subscribe(ctx, constants.TopicEventStats.String(), func(ctx context.Context, event events.IntegrationEvent) error {
		data, ok := event.Payload.(*communityEvent.GroupEventStatsPayload)
		if !ok {
			// Log error
			return nil
		}
		// Process data and update group stats in the database
		switch data.EventType {
		case constants.Created:
			datagroup, err := c.eventRepo.GetGroupEventByID(ctx, data.EventID)
			if err != nil {
				// Log error
				return nil
			}
			datagroup.AttendeesCount.Going = datagroup.AttendeesCount.Going + data.GoingCount
			datagroup.AttendeesCount.Interested = datagroup.AttendeesCount.Interested + data.InterestedCount
			err = c.eventRepo.UpdateGroupEvent(ctx, datagroup)
		case constants.Deleted:
		case constants.Updated:
		default:
			// Log unknown event type
			return nil
		}
		return nil
	})
	if err != nil {
		// Log error
	}
}
func (c *ConsumerStats) ConsumerFailedEventStats(ctx context.Context)
