package main

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/database"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/grpc"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server         *gin.Engine
	Mongo          *database.MongodbConnection
	Redis          *database.RedisConnection
	Postgres       *database.PostgresConnection
	Elastic        *database.ElasticConnection
	Cassandra      *database.CassandraConnection
	Neo4j          *database.Neo4jConnection
	GRPCServer     *grpc.GRPCServer
	IdentityModule *identity.ModuleIdentity
}

func NewApp(
	server *gin.Engine,
	// Wire sẽ tự động bơm các biến này vào đây
	mongo *database.MongodbConnection,
	redis *database.RedisConnection,
	postgres *database.PostgresConnection,
	elastic *database.ElasticConnection,
	cassandra *database.CassandraConnection,
	neo4j *database.Neo4jConnection,
	GRPCServer *grpc.GRPCServer,
	identityModule *identity.ModuleIdentity,

) *App {

	return &App{Server: server, Mongo: mongo, Redis: redis, Postgres: postgres, Elastic: elastic, Cassandra: cassandra,
		Neo4j: neo4j, GRPCServer: GRPCServer, IdentityModule: identityModule}
}

func NewGinServer() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	return r
}
