package client

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/middleware"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/concurrency"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/email"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/redis"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/seaweedfs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/resilience"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/socket"
	"github.com/google/wire"
)

var ProviderDatabase = wire.NewSet(
	database.NewMongoDatabase,
	database.NewRedisClient,
	database.NewPostgresDB,
	database.NewElasticClient,
	// database.NewCassandraSession,
	database.NewNeo4jDriver,
	database.NewSeaweedFSClient,
)

var ProviderGRPC = wire.NewSet(

	ProvideIdentityClient,
	ProvideGRPCConnection,
)
var ProviderMiddleware = wire.NewSet(
	middleware.NewAuthMiddleware,
)
var ProviderSocket = wire.NewSet(
	socket.NewHub,
	wire.Bind(new(socket.Manager), new(*socket.Hub)),
)
var ProviderSeaweedfs = wire.NewSet(
	seaweedfs.NewSeaweedfsAdapter,
)
var ProviderKafka = wire.NewSet(
	kafka.NewKafkaEventBus,
	wire.Bind(new(events.EventBus), new(*kafka.KafkaEventBus)),
)
var ProviderResilience = wire.NewSet(

	resilience.NewBreakerProvider,
)
var ProviderCache = wire.NewSet(
	redis.NewRedisAdapter,
)
var ProviderLifecycle = wire.NewSet(
	concurrency.NewManager,
	// wire.Bind(new(concurrency.Service), new(*concurrency.Manager)),
	ProvideLifecycleManager,
)
var WorkerPoolProvider = wire.NewSet(
	concurrency.NewWorkerPool,
	wire.Bind(new(IRepositoryShare.IWorkerPool), new(IRepositoryShare.IWorkerPool)),
)
var EmailProvider = wire.NewSet(
	email.NewEmailAdapter,
	wire.Bind(new(IRepositoryShare.IEmail), new(*email.EmailAdapter)),
)
var ProviderSet = wire.NewSet(
	ProviderDatabase,
	ProviderGRPC,
	ProviderMiddleware,
	ProviderSocket,
	ProviderSeaweedfs,
	ProviderResilience,
	ProviderCache,
	ProviderKafka,
	ProviderLifecycle,
	WorkerPoolProvider,
	EmailProvider,
)
