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
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/graphEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/IRepository/neo4j"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/entity"
)

type GraphLocationConsumer struct {
	events    events.EventBus
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
	graphRepo neo4j.IGraphRepository
}

func NewGraphLocationConsumer() *GraphLocationConsumer {
	return &GraphLocationConsumer{}
}

func (c *GraphLocationConsumer) ConsumerGraphLocation(ctx context.Context) error {
	err := c.events.SubscribeBatch(ctx, constants.TopicGraphLocation.String(), 100, time.Duration(5)*time.Minute, func(ctx context.Context, events []events.IntegrationEvent) error {
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
					processErr = errors.New("unknown event type: " + event.Type)
					return
				}
			})
			if err != nil {
				errchan <- errors.New("failed to submit event " + event.ID + " to worker pool: " + err.Error())
				c.redisRepo.Unlock(ctx, event.ID)
				wg.Done() // Decrement the WaitGroup counter since the task won't be processed
			}
			if processErr != nil {
				errchan <- errors.New("failed to process event " + event.ID + ": " + processErr.Error())
				c.redisRepo.Unlock(ctx, event.ID)
			} else {
				errchan <- nil
				c.redisRepo.MarkCompleted(ctx, event.ID)
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
func (c *GraphLocationConsumer) handleCreatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[graphEvent.LocationNodePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.graphRepo.UpsertLocationNode(ctx, &entity.LocationNode{
		CityID:      data.CityID,
		CountryCode: data.CountryCode,
		GeoHash:     data.GeoHash,
	})
	return nil
}
func (c *GraphLocationConsumer) handleUpdatedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[graphEvent.LocationNodePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.graphRepo.UpsertLocationNode(ctx, &entity.LocationNode{
		CityID:      data.CityID,
		CountryCode: data.CountryCode,
		GeoHash:     data.GeoHash,
	})
	return nil
}
func (c *GraphLocationConsumer) handleDeletedEvent(ctx context.Context, event events.IntegrationEvent) error {
	data, err := utils.ParsePayload[graphEvent.DeleteLocationNodePayload](event.Payload)
	if err != nil {
		return kafka.NewNonRetryableError(fmt.Errorf("failed to parse payload: %w", err))
	}
	if data == nil {
		return kafka.NewNonRetryableError(fmt.Errorf("payload is nil"))
	}
	err = c.graphRepo.DeleteLocationNode(ctx, data.CityID)
	if err != nil {
		return fmt.Errorf("failed to delete location node: %w", err)
	}
	return nil
}
func (c *GraphLocationConsumer) ConsumerFailedGraphLocation(ctx context.Context) error {
	return nil
}
