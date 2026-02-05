package database

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbConnection struct {
	mongodbDB *mongo.Client
}

func NewMongodbConnection(cfg *configs.Config) (*MongodbConnection, error) {
	db := &MongodbConnection{}
	_, err := db.ConnectMongodb(&cfg.MongoDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func (ctx *MongodbConnection) ConnectMongodb(configMongodb *configs.MongodbConfig) (*mongo.Client, error) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(configMongodb.MONGO_URI)
	client, err := mongo.Connect(c, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(c, nil)
	if err != nil {
		return nil, err
	}
	ctx.mongodbDB = client
	return client, nil
}
func (ctx *MongodbConnection) GetDatabase() *mongo.Client {
	return ctx.mongodbDB
}

func (ctx *MongodbConnection) GetMongoDatabase(dbName string) *mongo.Database {
	return ctx.mongodbDB.Database(dbName)
}
