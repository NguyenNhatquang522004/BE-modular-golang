package consumer

import (
	"context"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"
)

type ConsumerSheduler struct {
	cron            IRepositoryShare.IScheduler
	storyRepository IRepositoryMongodb.IStoryRepository
}

func NewConsumerSheduler(cron IRepositoryShare.IScheduler, storyRepository IRepositoryMongodb.IStoryRepository) *ConsumerSheduler {
	return &ConsumerSheduler{
		cron:            cron,
		storyRepository: storyRepository,
	}
}

func (c *ConsumerSheduler) ConsumerDeleteStoryExpiresAt(ctx context.Context) {
	// Implement the logic to consume messages related to deleting stories that have expired
	// This could involve connecting to a message queue, processing messages, and performing necessary actions
	err := c.cron.ScheduleJob("0 0 * * *", func() {
		// Logic to delete stories that have expired
		data, err := c.storyRepository.DeleteExpiredStories(ctx)
		if err != nil {
			// Handle error, possibly log it
			return
		}
		// Optionally, log the number of deleted stories or perform other actions based on the result
		log.Printf("Deleted %d expired stories", data)
	})
	if err != nil {
		// Handle error, possibly log it
		return
	}
}
