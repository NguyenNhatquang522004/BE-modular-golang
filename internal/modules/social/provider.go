package social

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/IRepsitory/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/infrastructure/postgres"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	postgres.NewFriendshipsRepository,
	postgres.NewFollowersRepository,
	postgres.NewBlockRepository,
	wire.Bind(new(IRepositoryPostgres.IFriendshipsRepository), new(*postgres.FriendshipsRepository)),
	wire.Bind(new(IRepositoryPostgres.IFollowersRepository), new(*postgres.FollowersRepository)),
	wire.Bind(new(IRepositoryPostgres.IBlockRepository), new(*postgres.BlockRepository)),
)

var UseCaseSet = wire.NewSet()

var HandlerSet = wire.NewSet()

var ModuleSocialSet = wire.NewSet(
	NewModuleSocial,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
)
