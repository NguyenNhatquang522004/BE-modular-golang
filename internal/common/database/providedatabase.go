package database

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"gorm.io/gorm"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gocql/gocql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProvideMongoDatabase(conn *MongodbConnection, cfg *configs.Config) *mongo.Database {
	// Tại đây, bạn chỉ rõ: "Lấy cái tên DB từ config, nhét vào hàm GetMongoDatabase"
	return conn.GetMongoDatabase(cfg.MongoDB.MONGO_DB_NAME)
}

// --- 2. REDIS ---
func ProvideRedisClient(conn *RedisConnection) *redis.Client {
	return conn.GetClient()
}

// --- 3. CASSANDRA ---
func ProvideCassandraSession(conn *CassandraConnection) *gocql.Session {
	return conn.GetSession()
}

// --- 4. ELASTICSEARCH ---
func ProvideElasticClient(conn *ElasticConnection) *elasticsearch.Client {
	return conn.GetClient()
}

// --- 5. NEO4J ---
func ProvideNeo4jDriver(conn *Neo4jConnection) neo4j.DriverWithContext {
	return conn.GetDriver()
}
func ProvidePostgresGormDB(conn *PostgresConnection) *gorm.DB {
	return conn.GetDB()
}
