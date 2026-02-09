package client

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/google/wire"
)

var providerDatabase = wire.NewSet(
	database.NewMongoDatabase,
	// database.NewRedisClient,
	database.NewPostgresDB,
	// database.NewElasticClient,
	// database.NewCassandraSession,
	// database.NewNeo4jDriver,
)

var providerGRPC = wire.NewSet(

	ProvideIdentityClient,
	ProvideGRPCConnection,

)
// var providerResilience = wire.NewSet(
//
//	resilience.NewBreakerProvider,
//
// )
var ProviderSet = wire.NewSet(
	providerDatabase,
	providerGRPC,
)
