package interaction

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/IRepository/IRepositoryMongoDB"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/infrastructure/cassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/infrastructure/mongodb"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	mongodb.NewCommentRepository,
	mongodb.NewCommentEditLogsRepository,
	mongodb.NewSavedItemsRepository,
	cassandra.NewReactionsRepository,
	cassandra.NewReactionsHistoryRepository,
	wire.Bind(new(IRepositoryCassandra.IReactionsRepository), new(*cassandra.ReactionsRepository)),
	wire.Bind(new(IRepositoryCassandra.IReactionHistoryRepository), new(*cassandra.ReactionsHistoryRepository)),
	wire.Bind(new(IRepositoryMongoDB.ICommentRepository), new(*mongodb.CommentRepository)),
	wire.Bind(new(IRepositoryMongoDB.ICommentEditLogsRepository), new(*mongodb.CommentEditLogsRepository)),
	wire.Bind(new(IRepositoryMongoDB.ISavedItemsRepository), new(*mongodb.SavedItemsRepository)),
)
var UseCaseSet = wire.NewSet()
var HandlerSet = wire.NewSet()
var ProducerSet = wire.NewSet()
var ConsumerSet = wire.NewSet()
var ModuleInteractionSet = wire.NewSet(
	NewModuleInteraction,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
	ProducerSet,
	ConsumerSet,
)
