package communication

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/infrastructure/cassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/infrastructure/mongodb"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	mongodb.NewConversationsRepository,
	mongodb.NewConversationParticipantsRepository,
	mongodb.NewCallLogsRepository,
	cassandra.NewMessageRepository,
	cassandra.NewMessageReactionsRepository,
	cassandra.NewConversationReadStateRepository,
	wire.Bind(new(IRepositoryMongodb.IConversationsRepository), new(*mongodb.ConversationsRepository)),
	wire.Bind(new(IRepositoryMongodb.IConversationParticipantsRepository), new(*mongodb.ConversationParticipantsRepository)),
	wire.Bind(new(IRepositoryMongodb.ICallLogsRepository), new(*mongodb.CallLogsRepository)),
	wire.Bind(new(IRepositoryCassandra.IMessageRepository), new(*cassandra.MessageRepository)),
	wire.Bind(new(IRepositoryCassandra.IMessageReactionsRepository), new(*cassandra.MessageReactionsRepository)),
	wire.Bind(new(IRepositoryCassandra.IConversationReadStateRepository), new(*cassandra.ConversationReadStateRepository)),
)
var UseCaseSet = wire.NewSet()

var HandlerSet = wire.NewSet()
var ProducerSet = wire.NewSet()

var ConsumerSet = wire.NewSet()

var ModuleCommunicationSet = wire.NewSet(
	NewModuleCommunication,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
	ProducerSet,
	ConsumerSet,
)
