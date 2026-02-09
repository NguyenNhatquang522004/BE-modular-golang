package database

import (
	"fmt"
	"log"
	"strconv" // Dùng để chuyển đổi Port từ string sang int

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/gocql/gocql"
)

// 1. Tạo Wrapper Struct
func NewCassandraSession(cfg *configs.Config) (*gocql.Session, func(), error) {
	conf := cfg.CassandraDB

	log.Printf("Connecting to Cassandra at %s:%s...", conf.Host, conf.Port)

	// 1. Config Cluster
	cluster := gocql.NewCluster(conf.Host)
	
	port, err := strconv.Atoi(conf.Port)
	if err != nil {
		port = 9042
	}
	cluster.Port = port
	cluster.Keyspace = conf.Keyspace
	cluster.Consistency = gocql.Quorum

	if conf.User != "" && conf.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: conf.User,
			Password: conf.Password,
		}
	}

	// 2. Create Session
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cassandra session: %w", err)
	}

	log.Println("✅ Connected to Cassandra successfully!")

	// 3. Cleanup Function
	cleanup := func() {
		log.Println("⚠️ Closing Cassandra session...")
		session.Close()
	}

	return session, cleanup, nil
}