package notification

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryCassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IRepositoryMongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/IRepository/IStrategy"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/infrastructure/cassandra"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/infrastructure/mongodb"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/usecase/strategy"
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
var StrategyNotificationTypeSet = wire.NewSet(
	strategy.NewNotifCommentReply,
	strategy.NewNotifFriendAccept,
	strategy.NewNotifFriendRequest,
	strategy.NewNotifGroupInvite,
	strategy.NewNotifMention,
	strategy.NewNotifPostLike,
	strategy.NewNotifSystemAlert,
	strategy.NewStrategyNotificationType,
	wire.Bind(new(IStrategy.IStrategyTypeNotificationType), new(*strategy.NotifCommentReply)),
	wire.Bind(new(IStrategy.IStrategyTypeNotificationType), new(*strategy.NotifFriendAccept)),
	wire.Bind(new(IStrategy.IStrategyTypeNotificationType), new(*strategy.NotifFriendRequest)),
	wire.Bind(new(IStrategy.IStrategyTypeNotificationType), new(*strategy.NotifGroupInvite)),
	wire.Bind(new(IStrategy.IStrategyTypeNotificationType), new(*strategy.NotifMention)),
	wire.Bind(new(IStrategy.IStrategyTypeNotificationType), new(*strategy.NotiPostLike)),
	wire.Bind(new(IStrategy.IStrategyTypeNotificationType), new(*strategy.NotiFSystemAlert)),
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
