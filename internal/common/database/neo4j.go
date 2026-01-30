package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// 1. Wrapper Struct
type Neo4jConnection struct {
	neo4jDB neo4j.DriverWithContext
}

func NewNeo4jConnection(cfg *configs.Config) (*Neo4jConnection, error) {
	db := &Neo4jConnection{}
	err := db.ConnectNeo4j(&cfg.Neo4jDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}
// 2. Hàm khởi tạo kết nối
// Receiver đặt là 'n' (viết tắt của Neo4j)
func (n *Neo4jConnection) ConnectNeo4j(cfg *configs.Neo4jConfig) error {
	// URI thường là "bolt://localhost:7687" hoặc "neo4j://localhost:7687"
	uri := cfg.NEO4J_URI
	user := cfg.NEO4J_USER
	pass := cfg.NEO4J_PASSWORD

	// Tạo Auth Token
	authToken := neo4j.BasicAuth(user, pass, "")

	// Khởi tạo Driver (Sử dụng Context-aware Driver của bản v5)
	driver, err := neo4j.NewDriverWithContext(uri, authToken)
	if err != nil {
		return fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	// Kiểm tra kết nối (Connectivity Check)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to neo4j: %w", err)
	}

	log.Println("Connected to Neo4j successfully!")

	// 3. Gán driver vào struct
	n.neo4jDB = driver

	return nil
}

// 3. Hàm Getter
func (n *Neo4jConnection) GetDriver() neo4j.DriverWithContext {
	return n.neo4jDB
}

// 4. Hàm đóng kết nối (Quan trọng với Neo4j)
func (n *Neo4jConnection) Close(ctx context.Context) {
	if n.neo4jDB != nil {
		n.neo4jDB.Close(ctx)
	}
}
