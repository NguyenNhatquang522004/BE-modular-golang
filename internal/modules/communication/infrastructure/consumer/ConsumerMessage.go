package consumer

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
)

type ConsumerMessage struct {
	conversationRepo        IRepositoryMongodb.IConversationsRepository
	conversationParticipant IRepositoryMongodb.IConversationParticipantsRepository
	messageStateRepo        IRepositoryCassandra.IConversationReadStateRepository
	messageReact            IRepositoryCassandra.IMessageReactionsRepository
	messageRepo             IRepositoryCassandra.IMessageRepository
	events                  events.EventBus
	pool                    IRepositoryShare.IWorkerPool
}

func NewConsumerMessage(conversationRepo IRepositoryMongodb.IConversationsRepository, conversationParticipant IRepositoryMongodb.IConversationParticipantsRepository, messageStateRepo IRepositoryCassandra.IConversationReadStateRepository, messageReact IRepositoryCassandra.IMessageReactionsRepository, messageRepo IRepositoryCassandra.IMessageRepository, events events.EventBus, pool IRepositoryShare.IWorkerPool) *ConsumerMessage {
	return &ConsumerMessage{
		conversationRepo:        conversationRepo,
		conversationParticipant: conversationParticipant,
		messageStateRepo:        messageStateRepo,
		messageReact:            messageReact,
		messageRepo:             messageRepo,
		events:                  events,
		pool:                    pool,
	}
}
