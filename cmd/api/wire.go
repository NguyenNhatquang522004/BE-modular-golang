//go:build wireinject
// +build wireinject

package main

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/google/wire"
)

func InitializeApp(config *configs.Config) (*App, error) {
	wire.Build(
		database.NewCassandraConnection,
		database.NewElasticConnection,
		database.NewNeo4jConnection,
		database.NewMongodbConnection,
		database.NewRedisConnection,
		database.NewPostgresConnection,
		database.ProvideMongoDatabase,
		database.ProvideRedisClient,
		database.ProvideCassandraSession,
		database.ProvideElasticClient,
		database.ProvideNeo4jDriver,
		database.ProvidePostgresGormDB,
		NewGinServer,
		NewApp,
	)
	return &App{}, nil
}
