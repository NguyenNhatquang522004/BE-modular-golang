package IProducerNotification

import "context"

type IProducerPost interface {
	PublishPostStats(ctx context.Context) error
	PublishPostNotificationFriend(ctx context.Context) error
	PublishPostNotificationFollower(ctx context.Context) error
	PublishPostNotificationTag(ctx context.Context) error
}
