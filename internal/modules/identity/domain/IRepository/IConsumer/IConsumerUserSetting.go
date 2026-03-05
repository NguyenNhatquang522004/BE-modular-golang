package IConsumer

import "context"

type IConsumerUserSetting interface {
	ConsumerUserSettingEvent(ctx context.Context) error
	ConsumerFailedUserSettingEvent(ctx context.Context) error
}
