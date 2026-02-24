package content

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IStrategy"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/infrastructure/cassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/infrastructure/mongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/usecase"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/usecase/strategy"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	mongodb.NewPostExtensionRepository,
	mongodb.NewPostMediaRepository,
	mongodb.NewPostSettingRepository,
	mongodb.NewPostEditLogsRepository,
	mongodb.NewPostRepository,
	cassandra.NewPostInsightsRepository,
	wire.Bind(new(IRepositoryMongodb.IPostExtensionRepository), new(*mongodb.PostExtensionRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostMediaRepository), new(*mongodb.PostMediaRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostSettingRepository), new(*mongodb.PostSettingRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostEditLogsRepository), new(*mongodb.PostEditLogsRepository)),
	wire.Bind(new(IRepositoryMongodb.IPostRepository), new(*mongodb.PostRepository)),
	wire.Bind(new(IRepositoryCassandra.IPostInsights), new(*cassandra.PostInsightsRepository)),
)
var UsecaseSet = wire.NewSet(
	usecase.NewDeletePostUseCase,
	usecase.NewGetPostByUserIDUseCase,
	usecase.NewHideOrUnhidePostUseCase,
	usecase.NewEditPostUseCase,
	usecase.NewPublishPostUseCase,
	usecase.NewUsecaseContent,
	usecase.NewGetEnumUseCase,
	wire.Bind(new(usecase.IDeletePostUseCase), new(*usecase.DeletePostUseCase)),
	wire.Bind(new(usecase.IGetPostByUserIDUseCase), new(*usecase.GetPostByUserIDUseCase)),
	wire.Bind(new(usecase.IHideOrUnhidePostUseCase), new(*usecase.HideOrUnhidePostUseCase)),
	wire.Bind(new(usecase.IEditPostUseCase), new(*usecase.EditPostUseCase)),
	wire.Bind(new(usecase.IPublishPostUseCase), new(*usecase.PublishPostUseCase)),
	wire.Bind(new(usecase.IGetEnumUseCase), new(*usecase.GetEnumUseCase)),
)

var HandlerSet = wire.NewSet()
var producerSet = wire.NewSet()
var consumerSet = wire.NewSet()
var StrategyPublishPostSet = wire.NewSet(
	strategy.NewPublishPostStrategy,
	strategy.NewPostExtensionStrategy,
	strategy.NewPostInsightStrategy,
	strategy.NewMediaStrategy,
	strategy.NewSettingStrategy,
	strategy.NewPublishPostStrategy,
	strategy.NewPublishDeleteStrategy,
	wire.Bind(new(IStrategy.IPublishPostStrategy), new(*strategy.PostExtensionStrategy)),
	wire.Bind(new(IStrategy.IPublishPostStrategy), new(*strategy.PostInsightStrategy)),
	wire.Bind(new(IStrategy.IPublishPostStrategy), new(*strategy.MediaStrategy)),
	wire.Bind(new(IStrategy.IPublishPostStrategy), new(*strategy.SettingStrategy)),
	wire.Bind(new(IStrategy.IPublishPostStrategy), new(*strategy.PostInsightStrategy)),
	wire.Bind(new(IStrategy.IPublishDeleteStrategy), new(*strategy.PostExtensionStrategy)),
	wire.Bind(new(IStrategy.IPublishDeleteStrategy), new(*strategy.PostInsightStrategy)),
	wire.Bind(new(IStrategy.IPublishDeleteStrategy), new(*strategy.MediaStrategy)),
	wire.Bind(new(IStrategy.IPublishDeleteStrategy), new(*strategy.SettingStrategy)),
	wire.Bind(new(IStrategy.IPublishDeleteStrategy), new(*strategy.PostInsightStrategy)),
	wire.Bind(new(IStrategy.IPublishDeleteStrategy), new(*strategy.EditLogStrategy)),
)
var ModuleContentSet = wire.NewSet(
	NewModuleContent,
	RepositorySet,
	UsecaseSet,
	HandlerSet,
	producerSet,
	consumerSet,
)
