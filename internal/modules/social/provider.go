package social

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IKafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/infrastructure/kafka"
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
var KafkaRepositorySet = wire.NewSet(
	kafka.NewBlockMessage,
	kafka.NewFollowMessage,
	kafka.NewFriendshipMessage,
	wire.Bind(new(IKafka.IBlockMessage), new(*kafka.BlockMessage)),
	wire.Bind(new(IKafka.IFollowMessage), new(*kafka.FollowMessage)),
	wire.Bind(new(IKafka.IFriendshipMessage), new(*kafka.FriendshipMessage)),
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

var HandlerSet = wire.NewSet()

var ModuleSocialSet = wire.NewSet(
	NewModuleSocial,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
	KafkaRepositorySet,
)
