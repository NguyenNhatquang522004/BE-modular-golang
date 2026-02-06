package identity

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryKeyCloak"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/keycloak"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/mongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/postgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/usecase"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	postgres.NewUserRepository,
	keycloak.NewKeycloakRepository,
	mongodb.NewUserSettingRepository,
	postgres.NewUserSessionRepository,
	wire.Bind(new(IRepositoryPostgres.IUserRepository), new(*postgres.UserRepository)),
	wire.Bind(new(IRepositoryKeyCloak.IKeycloakRepository), new(*keycloak.KeycloakRepository)),
	wire.Bind(new(IRepositoryMongodb.IUserSettingRepository), new(*mongodb.UserSettingRepository)),
	wire.Bind(new(IRepositoryPostgres.IUserSessionRepository), new(*postgres.UserSessionRepository)),
)

var UseCaseSet = wire.NewSet(
	usecase.NewUserSessionUseCase,
	usecase.NewUserSettingUseCase,
	usecase.NewUserAuthUseCase,
	usecase.NewGoogleAuthUseCase,
	wire.Bind(new(usecase.IUserSessionService), new(*usecase.UserSessionUseCase)),
	wire.Bind(new(usecase.IUserSettingService), new(*usecase.UseSettingUseCase)),
	wire.Bind(new(usecase.IUserAuthService), new(*usecase.UserAuthUseCase)),
	wire.Bind(new(usecase.IGoogleAuthUseCase), new(*usecase.GoogleAuthUseCase)),
)

var HandlerSet = wire.NewSet(
// Add handler providers here
)

var ModuleIndentitySet = wire.NewSet(
	RepositorySet,
	UseCaseSet,
	HandlerSet,
)
