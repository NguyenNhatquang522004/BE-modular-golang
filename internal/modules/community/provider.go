package community

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/infrastructure/mongodb"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	mongodb.NewGroupRepository,
	mongodb.NewGroupMembersRepository,
	mongodb.NewGroupJoinQuestionsRepository,
	mongodb.NewGroupFilesRepository,
	mongodb.NewGroupEventsRepository,
	wire.Bind(new(IRepositoryMongodb.IGroupRepository), new(*mongodb.GroupRepository)),
	wire.Bind(new(IRepositoryMongodb.IGroupMembersRepository), new(*mongodb.GroupMembersRepository)),
	wire.Bind(new(IRepositoryMongodb.IGroupjoinQuestionsRepository), new(*mongodb.GroupJoinQuestionsRepository)),
	wire.Bind(new(IRepositoryMongodb.IGroupFilesRepository), new(*mongodb.GroupFilesRepository)),
	wire.Bind(new(IRepositoryMongodb.IGroupEventsRepository), new(*mongodb.GroupEventsRepository)),

// Repositories...
)

var UseCaseSet = wire.NewSet(
// UseCases...
)
var HandlerSet = wire.NewSet(
// Handlers...
)
var ProducerSet = wire.NewSet(
// Producers...
)
var ConsumerSet = wire.NewSet(
// Consumers...
)

var ModuleCommunitySet = wire.NewSet(
	NewModuleCommunity,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
	ProducerSet,
	ConsumerSet,
)
