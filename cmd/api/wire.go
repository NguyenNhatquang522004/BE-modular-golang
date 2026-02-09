//go:build wireinject
// +build wireinject

package main

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/grpc"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/client"
	"github.com/google/wire"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity"
)

func InitializeApp(config *configs.Config) (*App, error) {
	wire.Build(
		database.NewCassandraConnection,
		database.NewElasticConnection,
		database.NewNeo4jConnection,
		database.NewMongodbConnection,
		database.NewRedisConnection,
		database.NewPostgresConnection,
		grpc.NewGRPCServer,
		client.ProvideMongoDatabase,
		// client.ProvideRedisClient,
		// client.ProvideCassandraSession,
		// client.ProvideElasticClient,
		// client.ProvideNeo4jDriver,
		client.ProvidePostgresGormDB,
		identity.ModuleIndentitySet,
		NewGinServer,
		NewApp,
	)
	return &App{}, nil
}
