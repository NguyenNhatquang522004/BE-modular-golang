package content

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/infrastructure/mongodb"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	mongodb.NewPostExtensionRepository,
	mongodb.NewPostMediaRepository,
	mongodb.NewPostSettingRepository,
	mongodb.NewPostEditLogsRepository,
	mongodb.NewPostRepository,
	wire.Bind(new(IRepositoryMongodb.IPostExtensionRepository), new(*mongodb.PostExtensionRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostMediaRepository), new(*mongodb.PostMediaRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostSettingRepository), new(*mongodb.PostSettingRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostEditLogsRepository), new(*mongodb.PostEditLogsRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostRepository), new(*mongodb.PostRepository)),
)
var UsecaseSet = wire.NewSet()
var HandlerSet = wire.NewSet()
var producerSet = wire.NewSet()
var consumerSet = wire.NewSet()

var ModuleContentSet = wire.NewSet(
	NewModuleContent,
	RepositorySet,
	UsecaseSet,
	HandlerSet,
	producerSet,
	consumerSet,
)
