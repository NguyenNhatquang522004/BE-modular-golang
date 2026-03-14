package IRepositoryConsumer

import "context"

type ICommunityStats interface {
	ConsumerCommunityStats(ctx context.Context) error
	ConsumerFailedCommunityStats(ctx context.Context) error
}
