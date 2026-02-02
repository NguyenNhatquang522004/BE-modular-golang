package domain

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gocql/gocql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type InitEntityModule interface {
	InitPostgres(db *gorm.DB) error
	InitMongo(db *mongo.Database) error
	InitCassandra(session *gocql.Session) error
	InitElastic(client *elasticsearch.Client) error
	InitNeo4j(driver neo4j.DriverWithContext) error
}
