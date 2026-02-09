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
func NewNeo4jDriver(cfg *configs.Config) (neo4j.DriverWithContext, func(), error) {
	// 1. Auth Token
	authToken := neo4j.BasicAuth(cfg.Neo4jDB.NEO4J_USER, cfg.Neo4jDB.NEO4J_PASSWORD, "")

	// 2. Create Driver
	driver, err := neo4j.NewDriverWithContext(cfg.Neo4jDB.NEO4J_URI, authToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	// 3. Verify Connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to verify neo4j connectivity: %w", err)
	}

	log.Println("✅ Connected to Neo4j successfully!")

	// 4. Cleanup Function
	cleanup := func() {
		log.Println("⚠️ Closing Neo4j driver...")
		if err := driver.Close(context.Background()); err != nil {
			log.Printf("Error closing neo4j: %v", err)
		}
	}

	return driver, cleanup, nil
}