package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(cfg *configs.Config) (*mongo.Database, func(), error) {

	// 1. Setup Context với Timeout (chỉ dùng lúc connect thôi)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Config Client
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.MONGO_URI)

	// 3. Connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	// 4. Ping check (Bắt buộc)
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	log.Println("✅ Connected to MongoDB successfully!")

	// 5. Chọn Database luôn (Repository chỉ cần cái này)
	db := client.Database(cfg.MongoDB.MONGO_DB_NAME)

	// 6. Tạo hàm Cleanup (Wire sẽ gọi hàm này ở main.go)
	cleanup := func() {
		log.Println("⚠️ Closing MongoDB connection...")
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting mongodb: %v", err)
		}
	}

	return db, cleanup, nil
}
