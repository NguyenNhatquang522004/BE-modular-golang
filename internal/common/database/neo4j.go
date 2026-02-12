package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var schemaQueries = []string{
	// ==========================================
	// 1. UNIQUE CONSTRAINTS (Định danh - Identity)
	// ==========================================
	"CREATE CONSTRAINT user_id_unique IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS UNIQUE",
	"CREATE CONSTRAINT topic_name_unique IF NOT EXISTS FOR (t:Topic) REQUIRE t.name IS UNIQUE",
	"CREATE CONSTRAINT group_id_unique IF NOT EXISTS FOR (g:Group) REQUIRE g.group_id IS UNIQUE",
	"CREATE CONSTRAINT page_id_unique IF NOT EXISTS FOR (p:Page) REQUIRE p.page_id IS UNIQUE",
	"CREATE CONSTRAINT device_id_unique IF NOT EXISTS FOR (d:Device) REQUIRE d.device_id IS UNIQUE",
	"CREATE CONSTRAINT phone_hash_unique IF NOT EXISTS FOR (c:PhoneContact) REQUIRE c.phone_hash IS UNIQUE",

	// ==========================================
	// 2. NODE INDEXES (Hiệu năng tìm kiếm)
	// ==========================================
	"CREATE INDEX user_verified_idx IF NOT EXISTS FOR (u:User) ON (u.is_verified)",
	"CREATE INDEX user_created_at_idx IF NOT EXISTS FOR (u:User) ON (u.created_at)",
	"CREATE INDEX user_last_active_idx IF NOT EXISTS FOR (u:User) ON (u.last_active_at)",
	"CREATE INDEX user_risk_score_idx IF NOT EXISTS FOR (u:User) ON (u.risk_score)",

	"CREATE INDEX topic_trending_idx IF NOT EXISTS FOR (t:Topic) ON (t.trending_score)",
	"CREATE INDEX topic_category_idx IF NOT EXISTS FOR (t:Topic) ON (t.category)",

	"CREATE INDEX group_member_count_idx IF NOT EXISTS FOR (g:Group) ON (g.member_count)",
	"CREATE INDEX location_geohash_idx IF NOT EXISTS FOR (l:Location) ON (l.geo_hash)",

	// ==========================================
	// 3. RELATIONSHIP INDEXES (Ranking & Time Decay) - QUAN TRỌNG NHẤT
	// ==========================================
	"CREATE INDEX rel_interacted_affinity_idx IF NOT EXISTS FOR ()-[r:INTERACTED_WITH]-() ON (r.affinity_score)",
	"CREATE INDEX rel_interacted_time_idx IF NOT EXISTS FOR ()-[r:INTERACTED_WITH]-() ON (r.last_interaction_at)",
	"CREATE INDEX rel_interested_score_idx IF NOT EXISTS FOR ()-[r:INTERESTED_IN]-() ON (r.score)",
	"CREATE INDEX rel_friend_since_idx IF NOT EXISTS FOR ()-[r:FRIEND]-() ON (r.since)",
	"CREATE INDEX rel_used_device_time_idx IF NOT EXISTS FOR ()-[r:USED_DEVICE]-() ON (r.last_used_at)",

	// ==========================================
	// 4. VECTOR INDEX (AI / Machine Learning)
	// Lưu ý: Cú pháp này dành cho Neo4j 5.x trở lên
	// ==========================================
	`CREATE VECTOR INDEX user_embedding_idx IF NOT EXISTS
	FOR (u:User) ON (u.embedding)
	OPTIONS {indexConfig: {
	 'vector.dimensions': 128,
	 'vector.similarity_function': 'cosine'
	}}`,

	// ==========================================
	// 5. EXISTENCE CONSTRAINTS (Enterprise Only)
	// Bỏ comment nếu bạn dùng bản Enterprise
	// ==========================================
	// "CREATE CONSTRAINT user_id_exists IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS NOT NULL",

	"CREATE CONSTRAINT country_code_unique IF NOT EXISTS FOR (c:Country) REQUIRE c.code IS UNIQUE;",
	"CREATE INDEX rel_invited_timestamp_idx IF NOT EXISTS FOR ()-[r:INVITED]-() ON (r.timestamp);",
	"CREATE INDEX rel_topic_child_of_idx IF NOT EXISTS FOR ()-[r:CHILD_OF]-() ON (r.weight);",
}

// 1. Wrapper Struct
func NewNeo4jDriver(cfg *configs.Config) (neo4j.DriverWithContext, func(), error) {
	// 1. Auth Token
	authToken := neo4j.BasicAuth(cfg.Neo4jDB.NEO4J_USER, cfg.Neo4jDB.NEO4J_PASSWORD, "")

	// 2. Create Driver
	driver, err := neo4j.NewDriverWithContext(cfg.Neo4jDB.NEO4J_URI, authToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	// 3. Verify Connectivity & Run Migration
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to verify neo4j connectivity: %w", err)
	}

	// --- TÍCH HỢP MIGRATION TẠI ĐÂY ---
	log.Println("⏳ Running Neo4j Schema Migration...")
	if err := applySchema(ctx, driver); err != nil {
		return nil, nil, fmt.Errorf("failed to apply neo4j schema: %w", err)
	}

	log.Println("✅ Connected to Neo4j and Schema is up-to-date!")

	// 4. Cleanup Function
	cleanup := func() {
		log.Println("⚠️ Closing Neo4j driver...")
		if err := driver.Close(context.Background()); err != nil {
			log.Printf("Error closing neo4j: %v", err)
		}
	}

	return driver, cleanup, nil
}
func applySchema(ctx context.Context, driver neo4j.DriverWithContext) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for _, query := range schemaQueries {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, query, nil)
			if err != nil {
				return nil, err
			}
			return result.Consume(ctx)
		})
		if err != nil {
			return fmt.Errorf("query failed [%s]: %w", query, err)
		}
	}
	return nil
}
