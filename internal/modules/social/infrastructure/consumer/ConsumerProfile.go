package consumer

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
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
		// var wg sync.WaitGroup

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsumerProfile) ConsumerFailedProfile(ctx context.Context) error {
	return nil
}
