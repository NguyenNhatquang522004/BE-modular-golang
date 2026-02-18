package social

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IProducer/IGraph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/infrastructure/mongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/infrastructure/postgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/infrastructure/producer/graph"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/usecase"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/usecase/Strategy/blockStrategy"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	postgres.NewFriendshipsRepository,
	postgres.NewFollowersRepository,
	postgres.NewBlockRepository,
	mongodb.NewProfileRepository,
	wire.Bind(new(IRepositoryPostgres.IFriendshipsRepository), new(*postgres.FriendshipsRepository)),
	wire.Bind(new(IRepositoryPostgres.IFollowersRepository), new(*postgres.FollowersRepository)),
	wire.Bind(new(IRepositoryPostgres.IBlockRepository), new(*postgres.BlockRepository)),
	wire.Bind(new(IRepositoryMongodb.IProfileRepositoryMongodb), new(*mongodb.ProfileRepository)),
)
var UseCaseSet = wire.NewSet(
	usecase.NewAdminSocialUseCase,
	usecase.NewBlockUseCase,
	usecase.NewFollowUseCase,
	usecase.NewFriendshipUseCase,
	usecase.NewProfileUseCase,
	wire.Bind(new(usecase.IAdminSocialUseCase), new(*usecase.AdminSocialUseCase)),
	wire.Bind(new(usecase.IBlockUseCase), new(*usecase.BlockUseCase)),
	wire.Bind(new(usecase.IFollowUseCase), new(*usecase.FollowUseCase)),
	wire.Bind(new(usecase.IFriendshipUseCase), new(*usecase.FriendshipUseCase)),
	wire.Bind(new(usecase.IProfileUseCase), new(*usecase.ProfileUseCase)),
	usecase.NewUseCase,
)
var FriendshipStrategySet = wire.NewSet(
	blockStrategy.NewBlockFull,
	blockStrategy.NewBlockProfile,
	blockStrategy.NewBlockChat,
	blockStrategy.NewBlockNone,
	blockStrategy.NewProviderBlockStrategy,
)
var BlockStrategySet = wire.NewSet(
	blockStrategy.NewBlockFull,
	blockStrategy.NewBlockProfile,
	blockStrategy.NewBlockChat,
	blockStrategy.NewBlockNone,
	blockStrategy.NewProviderBlockStrategy,
)
var ProducerSet = wire.NewSet(
	graph.NewBlockMessage,
	graph.NewFollowMessage,
	graph.NewFriendshipMessage,
	wire.Bind(new(IGraph.IBlockMessage), new(*graph.BlockMessage)),
	wire.Bind(new(IGraph.IFollowMessage), new(*graph.FollowMessage)),
	wire.Bind(new(IGraph.IFriendshipMessage), new(*graph.FriendshipMessage)),
)
var ConsumerSet = wire.NewSet()
var HandlerSet = wire.NewSet()

var ModuleSocialSet = wire.NewSet(
	NewModuleSocial,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
	ProducerSet,
	ConsumerSet,
)
