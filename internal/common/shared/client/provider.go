package client

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/middleware"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/concurrency"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/redis"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/seaweedfs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/resilience"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/socket"
	"github.com/google/wire"
)

var providerDatabase = wire.NewSet(
	database.NewMongoDatabase,
	database.NewRedisClient,
	database.NewPostgresDB,
	database.NewElasticClient,
	// database.NewCassandraSession,
	database.NewNeo4jDriver,
	database.NewSeaweedFSClient,
)

var providerGRPC = wire.NewSet(

	ProvideIdentityClient,
	ProvideGRPCConnection,
)
var prodviderMiddleware = wire.NewSet(
	middleware.NewAuthMiddleware,
)
var providerSocket = wire.NewSet(
	socket.NewHub,
	wire.Bind(new(socket.Manager), new(*socket.Hub)),
)
var providerSeaweedfs = wire.NewSet(
	seaweedfs.NewSeaweedfsAdapter,
)
var providerKafka = wire.NewSet(
	kafka.NewKafkaEventBus,
	wire.Bind(new(events.EventBus), new(*kafka.KafkaEventBus)),
)
var providerResilience = wire.NewSet(

	resilience.NewBreakerProvider,
)
var prodviderCache = wire.NewSet(
	redis.NewRedisAdapter,
)
var providerLifecycle = wire.NewSet(
	concurrency.NewManager,
	// wire.Bind(new(concurrency.Service), new(*concurrency.Manager)),
	ProvideLifecycleManager,
)
var ProviderSet = wire.NewSet(
	providerDatabase,
	providerGRPC,
	prodviderMiddleware,
	providerSocket,
	providerSeaweedfs,
	providerResilience,
	prodviderCache,
	providerKafka,
	providerLifecycle,
)
