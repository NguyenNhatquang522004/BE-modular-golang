package identity

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/keycloak"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/postgres"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	postgres.NewUserRepository,
	keycloak.NewKeycloakRepository,
	wire.Bind(new(IRepositoryPostgres.IUserRepository), new(*postgres.UserRepository)),
	wire.Bind(new(keycloak.IKeycloakRepository), new(*keycloak.KeycloakRepository)),
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
