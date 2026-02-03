package database

import (
	"fmt"
	"log"
	"strconv" // Dùng để chuyển đổi Port từ string sang int

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/gocql/gocql"
)

// 1. Tạo Wrapper Struct
type CassandraConnection struct {
	CassandraDB *gocql.Session
}

func NewCassandraConnection(cfg *configs.Config) (*CassandraConnection, error) {
	db := &CassandraConnection{}
	err := db.ConnectCassandra(&cfg.CassandraDB)
	if err != nil {
		return nil, err
	}
	return db , nil
}

// 2. Hàm khởi tạo kết nối
// Receiver đặt là 'cas' (viết tắt của Cassandra)
func (cas *CassandraConnection) ConnectCassandra(cfg *configs.CassandraConfig) error {
	log.Printf("Connecting to Cassandra at %s:%s...", cfg.Host, cfg.Port)

	// Tạo cấu hình Cluster
	cluster := gocql.NewCluster(cfg.Host)
	// Chuyển đổi Port từ string (trong config) sang int (gocql yêu cầu)
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		// Nếu không parse được port, dùng port mặc định 9042
		port = 9042
	}
	cluster.Port = port

	// Cấu hình Keyspace
	cluster.Keyspace = cfg.Keyspace

	// Cấu hình Consistency (Thường dùng Quorum cho hệ thống phân tán)
	cluster.Consistency = gocql.Quorum

	// Cấu hình xác thực (Nếu có User/Pass)
	if cfg.User != "" && cfg.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.User,
			Password: cfg.Password,
		}
	}

	// Tạo Session (Thực hiện kết nối)
	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create cassandra session: %w", err)
	}

	log.Println("Connected to Cassandra successfully!")

	// 3. Gán session vào struct (Quan trọng)
	cas.CassandraDB = session

	return nil
}

// 3. Hàm Getter để lấy Session ra dùng
func (cas *CassandraConnection) GetSession() *gocql.Session {
	return cas.CassandraDB
}

// 4. Hàm đóng kết nối (Nên gọi khi tắt app)
func (cas *CassandraConnection) Close() {
	if cas.CassandraDB != nil {
		cas.CassandraDB.Close()
	}
}
