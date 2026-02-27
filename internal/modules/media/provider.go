package media

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/IRepository/IRepositoryMongodb"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/infrastructure/cassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/infrastructure/mongodb"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	mongodb.NewMediaAssetsRepository,
	mongodb.NewLiveSessionRepository,
	mongodb.NewStoryRepository,
	mongodb.NewReelRepository,
	cassandra.NewLiveCommentsRepository,
	cassandra.NewStoryViewRepository,
	wire.Bind(new(IRepositoryMongodb.IMediaAssetsRepository), new(*mongodb.MediaAssetsRepository)),
	wire.Bind(new(IRepositoryMongodb.ILiveSessionRepository), new(*mongodb.LiveSessionRepository)),
	wire.Bind(new(IRepositoryMongodb.IStoryRepository), new(*mongodb.StoryRepository)),
	wire.Bind(new(IRepositoryMongodb.IReelRepository), new(*mongodb.ReelRepository)),
	wire.Bind(new(IRepositoryCassandra.ILiveCommentsRepository), new(*cassandra.LiveCommentsRepository)),
	wire.Bind(new(IRepositoryCassandra.IStoryViewRepository), new(*cassandra.StoryViewRepository)),
)
var UsecaseSet = wire.NewSet()
var HandlerSet = wire.NewSet()
var ProducerSet = wire.NewSet()
var ConsumerSet = wire.NewSet()
var ModuleMediaSet = wire.NewSet(
	NewModuleMedia,
	RepositorySet,
	UsecaseSet,
	HandlerSet,
	ProducerSet,
	ConsumerSet,
)
