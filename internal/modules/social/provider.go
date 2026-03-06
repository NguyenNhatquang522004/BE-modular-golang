package social

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/infrastructure/mongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/infrastructure/postgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/usecase"
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
	usecase.NewUseCase,
)
var ProducerSet = wire.NewSet()
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
