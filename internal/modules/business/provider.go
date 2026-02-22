package business

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/IRepository/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/infrastructure/cassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/infrastructure/mongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/infrastructure/postgres"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	postgres.NewAdAccountsRepository,
	postgres.NewAdCampaignsRepository,
	postgres.NewAdsRepository,
	mongodb.NewPageFollowersRepository,
	mongodb.NewPagesRolesRepository,
	mongodb.NewPagesRepository,
	cassandra.NewPageDailyMetricsRepository,
	wire.Bind(new(IRepositoryMongodb.IPageFollowersRepository), new(*mongodb.PageFollowersRepository)),
	wire.Bind(new(IRepositoryMongodb.IPageRolesRepository), new(*mongodb.PagesRolesRepository)),
	wire.Bind(new(IRepositoryMongodb.IPagesRepository), new(*mongodb.PagesRepository)),
	wire.Bind(new(IRepositoryCassandra.IPageDailyMetricsRepository), new(*cassandra.PageDailyMetricsRepository)),
	wire.Bind(new(IRepositoryPostgres.IAdAccountsRepository), new(*postgres.AdAccountsRepository)),
	wire.Bind(new(IRepositoryPostgres.IAdCampaignsRepository), new(*postgres.AdCampaignsRepository)),
	wire.Bind(new(IRepositoryPostgres.IAdsRepository), new(*postgres.AdsRepository)),

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
var ModuleBusinessSet = wire.NewSet(
	NewModuleBusiness,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
	ProducerSet,
	ConsumerSet,
)
