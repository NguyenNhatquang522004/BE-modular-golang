package identity

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/repository_postgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/postgres"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	postgres.NewUserRepository,
	wire.Bind(new(repository_postgres.IUserRepository), new(*postgres.UserRepository)),
)

var UseCaseSet = wire.NewSet(
// Add use case providers here
)

var HandlerSet = wire.NewSet(
// Add handler providers here
)

var ModuleIndentitySet = wire.NewSet(
	RepositorySet,
	UseCaseSet,
	HandlerSet,
)
