package consumer

import "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"

type ConsumerUserSetting struct {
	events events.EventBus
	pool IRepositoryShare.IWorkerPool
	UserSetting IRepositoryShare.IUserSettingRepository
	redisRepo IRepostiro
}

func NewConsumerUserSetting() *ConsumerUserSetting {
	return &ConsumerUserSetting{}
}
