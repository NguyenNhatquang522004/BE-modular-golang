package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	SEAWEEDFS   SEAWEEDFSConfig
	Neo4jDB     Neo4jConfig
	MongoDB     MongodbConfig
	ElasticDB   ElasticsearchConfig
	CassandraDB CassandraConfig
	RedisDB     RedisConfig
	PostgresDB  PostgresConfig
	Server      ServerConfig
}
type SEAWEEDFSConfig struct {
	SEAWEEDFS_FILER_URL string
	SEAWEEDFS_S3_URL    string
	S3_REGION           string
	S3_ACCESS_KEY       string
	S3_SECRET_KEY       string
	S3_BUCKET_NAME      string
}
type Neo4jConfig struct {
	NEO4J_URI      string
	NEO4J_USER     string
	NEO4J_PASSWORD string
}
type MongodbConfig struct {
	MONGO_URI     string
	MONGO_DB_NAME string
}
type ElasticsearchConfig struct {
	ES_ADDRESS string
}
type CassandraConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Keyspace string
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type PostgresConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	DB_SSLMODE  string
	DB_TimeZone string
}
type ServerConfig struct {
	Port    string
	BaseURL string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 4. MANUAL MAPPING (Quan trọng)
	// Vì file .env của bạn đặt tên biến lộn xộn (Host, POSTGRES_USER...)
	// nên ta phải lấy từng cái bỏ vào đúng chỗ trong Struct.
	cfg := &Config{
		// --- Postgres ---
		PostgresDB: PostgresConfig{
			Host:        viper.GetString("Host"), // Chú ý: trong file env bạn đang để key là "Host"
			Port:        viper.GetString("PORT"),
			User:        viper.GetString("POSTGRES_USER"),
			Password:    viper.GetString("POSTGRES_PASSWORD"),
			DBName:      viper.GetString("POSTGRES_DB"),
			DB_SSLMODE:  viper.GetString("DB_SSLMODE"),
			DB_TimeZone: viper.GetString("DB_TIMEZONE"),
		},

		// --- Redis ---
		RedisDB: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"), // Lưu ý: Field này là int nên dùng GetInt
		},

		// --- MongoDB ---
		MongoDB: MongodbConfig{
			MONGO_URI:     viper.GetString("MONGO_URI"),
			MONGO_DB_NAME: viper.GetString("MONGO_DB_NAME"),
		},

		// --- Elasticsearch ---
		ElasticDB: ElasticsearchConfig{
			ES_ADDRESS: viper.GetString("ES_ADDRESS"),
		},

		// --- Neo4j ---
		Neo4jDB: Neo4jConfig{
			NEO4J_URI:      viper.GetString("NEO4J_URI"),
			NEO4J_USER:     viper.GetString("NEO4J_USER"),
			NEO4J_PASSWORD: viper.GetString("NEO4J_PASSWORD"),
		},

		// --- Cassandra ---
		CassandraDB: CassandraConfig{
			Host:     viper.GetString("CASSANDRA_HOST"),
			Port:     viper.GetString("CASSANDRA_PORT"),
			Keyspace: viper.GetString("CASSANDRA_KEYSPACE"),
			// Thêm User/Pass nếu sau này bạn bật authentication cho Cassandra
			User:     viper.GetString("CASSANDRA_USER"),
			Password: viper.GetString("CASSANDRA_PASSWORD"),
		},

		// --- SeaweedFS (S3 & Filer) ---
		SEAWEEDFS: SEAWEEDFSConfig{
			SEAWEEDFS_FILER_URL: viper.GetString("SEAWEEDFS_FILER_URL"),
			SEAWEEDFS_S3_URL:    viper.GetString("SEAWEEDFS_S3_URL"),
			S3_REGION:           viper.GetString("S3_REGION"),
			S3_ACCESS_KEY:       viper.GetString("S3_ACCESS_KEY"),
			S3_SECRET_KEY:       viper.GetString("S3_SECRET_KEY"),
			S3_BUCKET_NAME:      viper.GetString("S3_BUCKET_NAME"),
		},

		// --- Server ---
		Server: ServerConfig{
			Port:    viper.GetString("SERVER_PORT"),
			BaseURL: viper.GetString("Base_URL"),
		},
	}

	return cfg, nil
}
