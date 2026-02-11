package identity

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/http"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryKeyCloak"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepository/IRepositoryPostgres"
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
	postgres.NewUserRoleRepository,
	wire.Bind(new(IRepositoryPostgres.IUserRepository), new(*postgres.UserRepository)),
	wire.Bind(new(IRepositoryKeyCloak.IKeycloakRepository), new(*keycloak.KeycloakRepository)),
	wire.Bind(new(IRepositoryMongodb.IUserSettingRepository), new(*mongodb.UserSettingRepository)),
	wire.Bind(new(IRepositoryPostgres.IUserSessionRepository), new(*postgres.UserSessionRepository)),
	wire.Bind(new(IRepositoryPostgres.IUserRoleRepository), new(*postgres.UserRoleRepository)),
)

var UseCaseSet = wire.NewSet(
	usecase.NewUserSessionUseCase,
	usecase.NewUserSettingUseCase,
	usecase.NewUserAuthUseCase,
	usecase.NewGoogleAuthUseCase,
	usecase.NewUserRoleUseCase,
	usecase.NewUserUseCase,
	usecase.NewAdminUseCase,
	wire.Bind(new(usecase.IUserSessionService), new(*usecase.UserSessionUseCase)),
	wire.Bind(new(usecase.IUserSettingService), new(*usecase.UserSettingUseCase)),
	wire.Bind(new(usecase.IUserAuthService), new(*usecase.UserAuthUseCase)),
	wire.Bind(new(usecase.IGoogleAuthUseCase), new(*usecase.GoogleAuthUseCase)),
	wire.Bind(new(usecase.IUserRoleUseCase), new(*usecase.UserRoleUseCase)),
	wire.Bind(new(usecase.IUserService), new(*usecase.UserUseCase)),
	wire.Bind(new(usecase.IAdminUseCase), new(*usecase.AdminUseCase)),
	usecase.NewUsecase,
)

var HandlerSet = wire.NewSet(
	http.NewAuthHandler,
	http.NewUserHandler,
	http.NewRoleHandler,
	http.NewSessionHandler,
	http.NewAdminHandler,
	wire.Bind(new(http.IHandlerAuth), new(*http.AuthHandler)),
	wire.Bind(new(http.IHandlerUser), new(*http.UserHandler)),
	wire.Bind(new(http.IHandlerRole), new(*http.RoleHandler)),
	wire.Bind(new(http.IHandlerUserSession), new(*http.SessionHandler)),
	wire.Bind(new(http.IHandleAdmin), new(*http.AdminHandler)),
	wire.Bind(new(http.IHandlerUserSetting), new(*http.SettingHandler)),
	http.NewHandler,

// Add handler providers here
)

var ModuleIndentitySet = wire.NewSet(
	NewModuleIdentity,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
)
