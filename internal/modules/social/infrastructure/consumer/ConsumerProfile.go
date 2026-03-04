package consumer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/socialEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
)

type ConsumerProfile struct {
	profileRepo IRepositoryMongodb.IProfileRepositoryMongodb
	events      events.EventBus
	pool        IRepositoryShare.IWorkerPool
}

func NewConsumerProfile(profileRepo IRepositoryMongodb.IProfileRepositoryMongodb, events events.EventBus) *ConsumerProfile {
	return &ConsumerProfile{
		profileRepo: profileRepo,
		events:      events,
	}
}

func (c *ConsumerProfile) ConsumerProfile(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicProfile.String(), 100, time.Duration(5)*time.Second, func(ctx context.Context, events []events.IntegrationEvent) error {
		var wg sync.WaitGroup
		errchan := make(chan error, len(events))
		for _, event := range events {
			ev := event
			wg.Add(1)
			err := c.pool.Run(ctx, func() {
				switch ev.Type {
				case constants.Created.String():
					err := c.handleCreatedProfile(ctx, ev)
					if err != nil {
						errchan <- err
					}
					errchan <- nil
				case constants.Updated.String():
					err := c.handleUpdatedProfile(ctx, ev)
					if err != nil {
						errchan <- err
					}
					errchan <- nil
				case constants.Deleted.String():
					err := c.handleDeletedProfile(ctx, ev)
					if err != nil {
						errchan <- err
					}
					errchan <- nil
				}
			})
			if err != nil {
				errchan <- err
				wg.Done()
			}
		}
		wg.Wait()
		close(errchan)
		for err := range errchan {
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerProfile) handleCreatedProfile(ctx context.Context, event events.IntegrationEvent) error {
	data, ok := event.Payload.(*socialEvent.ProfilePayload)
	if !ok {
		return kafka.NewNonRetryableError(errors.New("invalid event payload for profile event"))
	}
	entity, err := mapper.ToEntityProfilePayload(data)
	if err != nil {
		return err
	}
	err = c.profileRepo.CreateProfile(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerProfile) handleUpdatedProfile(ctx context.Context, event events.IntegrationEvent) error {
	data, ok := event.Payload.(*socialEvent.ProfilePayload)
	if !ok {
		kafka.NewNonRetryableError(errors.New("invalid event payload for profile event"))
		return kafka.NewNonRetryableError(errors.New("invalid event payload for profile event"))
	}
	dataprofile, err := c.profileRepo.GetProfileByID(ctx, data.UserID)
	if err != nil {
		return err
	}
	if dataprofile == nil {
		return errors.New("profile not found for update")
	}
	entity, err := mapper.ToEntityUpdateProfilePayload(dataprofile, data)
	if err != nil {
		return err
	}
	err = c.profileRepo.UpdateProfile(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerProfile) handleDeletedProfile(ctx context.Context, event events.IntegrationEvent) error {
	data, ok := event.Payload.(*socialEvent.ProfilePayload)
	if !ok {
		return kafka.NewNonRetryableError(errors.New("invalid event payload for profile event"))
	}
	err := c.profileRepo.DeleteProfile(ctx, data.UserID)
	if err != nil {
		return err
	}
	return nil
}
func (c *ConsumerProfile) ConsumerFailedProfile(ctx context.Context) error {
	return nil
}
