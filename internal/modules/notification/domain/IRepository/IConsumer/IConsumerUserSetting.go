package IConsumer

import "context"

type IConsumerUserSetting interface {
	ConsumerUserNotificationSettings(ctx context.Context) error
	ConsumerFailedUserNotificationSettings(ctx context.Context) error
}
