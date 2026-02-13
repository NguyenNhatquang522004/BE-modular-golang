package client

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/middleware"
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

var providerResilience = wire.NewSet(

	resilience.NewBreakerProvider,
)
var prodviderCache = wire.NewSet(
	redis.NewRedisAdapter,
)
var ProviderSet = wire.NewSet(
	providerDatabase,
	providerGRPC,
	prodviderMiddleware,
	providerSocket,
	providerSeaweedfs,
	providerResilience,
	prodviderCache,
)
