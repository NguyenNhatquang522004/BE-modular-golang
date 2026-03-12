package IConsumer

type IConsumerMusicStats interface {
	ConsumerMusicStats() error
	ConsumerFailedMusicStats() error
}
