package notification

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/infrastructure/cassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/infrastructure/mongodb"
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	mongodb.NewNotificationTemplatesRepository,
	mongodb.NewUserNotificationSettingsRepository,
	cassandra.NewNotificationsRepository,
	wire.Bind(new(IRepositoryMongodb.INotificationTemplatesRepository), new(*mongodb.NotificationTemplatesRepository)),
	wire.Bind(new(IRepositoryMongodb.IUserNotificationSettingsRepository), new(*mongodb.UserNotificationSettingsRepository)),
	wire.Bind(new(IRepositoryCassandra.INotificationsRepository), new(*cassandra.NotificationsRepository)),
)
var UseCaseSet = wire.NewSet()
var HandlerSet = wire.NewSet()
var ProducerSet = wire.NewSet()
var ConsumerSet = wire.NewSet()

var ModuleNotificationSet = wire.NewSet(
	NewModuleNotification,
	RepositorySet,
	UseCaseSet,
	HandlerSet,
	ProducerSet,
	ConsumerSet,
)
