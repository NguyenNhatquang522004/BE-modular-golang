package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DLQ            DLQConfig
	WorkerPool     WorkerPoolConfig
	CircuitBreaker CircuitBreakerConfig
	GRPCServer     GRPCServerConfig
	Kafka          KafkaConfig
	KeyCloak       KeyCloakConfig
	EmailSMTP      EmailSMTPConfig
	GoogleAuth     GoogleAuthConfig
	JWT            JWTConfig
	SEAWEEDFS      SEAWEEDFSConfig
	Neo4jDB        Neo4jConfig
	MongoDB        MongodbConfig
	ElasticDB      ElasticsearchConfig
	CassandraDB    CassandraConfig
	RedisDB        RedisConfig
	PostgresDB     PostgresConfig
	Server         ServerConfig
	VNPay          VNPayConfig
}

// DLQ (Dead Letter Queue) Configuration
type DLQConfig struct {
	TOPIC            string
	MAX_RETRIES      int
	RETRY_BACKOFF_MS int
}

type VNPayConfig struct {
	TMN_CODE    string
	HASH_SECRET string
	PAYMENT_URL string
	API_URL     string
	RETURN_URL  string
	IPN_URL     string
	IS_SANDBOX  bool
}
type WorkerPoolConfig struct {
	WORKER_POOL_SIZE_MAX          int
	WORKER_POOL_POST_INSIGHTS_MAX int
}
type CircuitBreakerConfig struct {
	Name                   string
	MaxRequests            uint32
	Interval               int64
	FailureThreshold       float64
	RecoveryTimeoutSeconds int
	ExpectedResponseTimeMS int
}
type GRPCServerConfig struct {
	GRPC_SERVER_HOST string
	GRPC_SERVER_PORT string
}
type KafkaConfig struct {
	BROKERS        []string
	CONSUMER_GROUP string
}
type KeyCloakConfig struct {
	KEYCLOAK_SERVER_URL    string
	KEYCLOAK_REALM         string
	KEYCLOAK_CLIENT_ID     string
	KEYCLOAK_CLIENT_SECRET string
	REDIRECT_URI           string
	KEYCLOAK_AUTH_URL      string
	KEYCLOAK_TOKEN_URL     string
	KEYCLOAK_USERINFO_URL  string
}
type EmailSMTPConfig struct {
	SMTP_HOST        string
	SMTP_PORT        int
	SMTP_EMAIL       string
	SMTP_PASSWORD    string
	SMTP_SENDER_NAME string
}
type GoogleAuthConfig struct {
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL  string
}
type JWTConfig struct {
	JWT_SECRET_KEY       string
	JWT_EXPIRES_IN_HOURS int64
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
	Host         string
	Port         string
	Password     string
	DB           int
	LockTTL      int
	CompletedTTL int
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
		// --- DLQ Config ---
		DLQ: DLQConfig{
			TOPIC:            viper.GetString("KAFKA_TOPIC_DLQ"),
			MAX_RETRIES:      viper.GetInt("KAFKA_MAX_RETRIES"),
			RETRY_BACKOFF_MS: viper.GetInt("KAFKA_RETRY_BACKOFF_MS"),
		},
		// --- Worker Pool ---
		WorkerPool: WorkerPoolConfig{
			WORKER_POOL_SIZE_MAX:          int(viper.GetInt64("WORKER_POOL_SIZE_MAX")),
			WORKER_POOL_POST_INSIGHTS_MAX: int(viper.GetInt64("WORKER_POOL_POST_INSIGHTS_MAX")),
		},
		// --- Circuit Breaker ---
		CircuitBreaker: CircuitBreakerConfig{
			FailureThreshold:       viper.GetFloat64("CIRCUIT_BREAKER_FAILURE_THRESHOLD"),
			RecoveryTimeoutSeconds: int(viper.GetDuration("CIRCUIT_BREAKER_RECOVERY_TIMEOUT_SECONDS").Seconds()),
			ExpectedResponseTimeMS: viper.GetInt("CIRCUIT_BREAKER_EXPECTED_RESPONSE_TIME_MS"),
			Name:                   viper.GetString("CIRCUIT_BREAKER_NAME"),
			MaxRequests:            viper.GetUint32("CIRCUIT_BREAKER_MAX_REQUESTS"),
			Interval:               int64(viper.GetDuration("CIRCUIT_BREAKER_INTERVAL_MS").Milliseconds()),
		},
		// --- GRPC SERVER ---
		GRPCServer: GRPCServerConfig{
			GRPC_SERVER_HOST: viper.GetString("GRPC_SERVER_HOST"),
			GRPC_SERVER_PORT: viper.GetString("GRPC_SERVER_PORT"),
		},
		// --- Kafka ---
		Kafka: KafkaConfig{
			BROKERS:        viper.GetStringSlice("KAFKA_BROKERS"),
			CONSUMER_GROUP: viper.GetString("KAFKA_CONSUMER_GROUP"),
		},
		// --- KeyCloak ---
		KeyCloak: KeyCloakConfig{
			KEYCLOAK_SERVER_URL:    viper.GetString("KEYCLOAK_SERVER_URL"),
			KEYCLOAK_REALM:         viper.GetString("KEYCLOAK_REALM"),
			KEYCLOAK_CLIENT_ID:     viper.GetString("KEYCLOAK_CLIENT_ID"),
			KEYCLOAK_CLIENT_SECRET: viper.GetString("KEYCLOAK_CLIENT_SECRET"),
			REDIRECT_URI:           viper.GetString("KEYCLOAK_REDIRECT_URI"),
			KEYCLOAK_AUTH_URL:      viper.GetString("KEYCLOAK_AuthURL"),
			KEYCLOAK_TOKEN_URL:     viper.GetString("KEYCLOAK_TokenURL"),
			KEYCLOAK_USERINFO_URL:  viper.GetString("KEYCLOAK_UserInfoURL"),
		},
		// --- Email SMTP ---
		EmailSMTP: EmailSMTPConfig{
			SMTP_HOST:        viper.GetString("SMTP_HOST"),
			SMTP_PORT:        viper.GetInt("SMTP_PORT"),
			SMTP_EMAIL:       viper.GetString("SMTP_EMAIL"),
			SMTP_PASSWORD:    viper.GetString("SMTP_PASSWORD"),
			SMTP_SENDER_NAME: viper.GetString("SMTP_SENDER_NAME"),
		},
		// --- Google Auth ---
		GoogleAuth: GoogleAuthConfig{
			GOOGLE_CLIENT_ID:     viper.GetString("GOOGLE_CLIENT_ID"),
			GOOGLE_CLIENT_SECRET: viper.GetString("GOOGLE_CLIENT_SECRET"),
			GOOGLE_REDIRECT_URL:  viper.GetString("GOOGLE_REDIRECT_URL"),
		},
		// --- JWT ---
		JWT: JWTConfig{
			JWT_SECRET_KEY:       viper.GetString("JWT_SECRET_KEY"),
			JWT_EXPIRES_IN_HOURS: viper.GetInt64("JWT_EXPIRES_IN_HOURS"),
		},
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
			Host:         viper.GetString("REDIS_HOST"),
			Port:         viper.GetString("REDIS_PORT"),
			Password:     viper.GetString("REDIS_PASSWORD"),
			DB:           viper.GetInt("REDIS_DB"), // Lưu ý: Field này là int nên dùng GetInt
			LockTTL:      viper.GetInt("REDIS_LOCK_TTL"),
			CompletedTTL: viper.GetInt("REDIS_LOCK_COMPLETED_TTL"),
			// Thêm các cấu hình khác nếu cần (ví dụ: PoolSize, MinIdleConns...)
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

		// --- VNPay ---
		VNPay: VNPayConfig{
			TMN_CODE:    viper.GetString("VNPAY_TMN_CODE"),
			HASH_SECRET: viper.GetString("VNPAY_HASH_SECRET"),
			PAYMENT_URL: viper.GetString("VNPAY_PAYMENT_URL"),
			API_URL:     viper.GetString("VNPAY_API_URL"),
			RETURN_URL:  viper.GetString("VNPAY_RETURN_URL"),
			IPN_URL:     viper.GetString("VNPAY_IPN_URL"),
			IS_SANDBOX:  viper.GetBool("VNPAY_IS_SANDBOX"),
		},
	}

	return cfg, nil
}
