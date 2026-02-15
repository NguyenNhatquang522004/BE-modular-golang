package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/kafka"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org
// @host localhost:8081 // Replace with your actual host
// @BasePath /
func main() {
	// 1. Load Config
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}
	if len(cfg.Kafka.BROKERS) > 0 {
		kafkaBroker := cfg.Kafka.BROKERS[0]
		log.Println("⚡ Starting Kafka Topic Migration...")

		// Duyệt qua danh sách Topic đã định nghĩa ở constants
		for _, t := range constants.SocialTopics {
			// ReplicationFactor: Để 1 khi dev local, 3 khi Production
			err := kafka.EnsureTopicExists(kafkaBroker, t.Name.String(), t.Partitions, 1)
			if err != nil {
				// Tùy chọn: Panic nếu không tạo được topic quan trọng, hoặc chỉ log warning
				log.Fatalf("❌ Failed to ensure topic %s: %v", t.Name, err)
			}
		}
		log.Println("✅ Kafka Topic Migration Completed!")
	}
	// 2. Init App (Wire làm hết việc tạo object ở đây)
	app, cleanup, err := InitializeApp(cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	ctx := context.Background()
	go func() {
		log.Println("Starting Realtime Hub Manager...")
		app.Hub.Run(ctx) // <--- CHÍNH LÀ NÓ
	}()
	// 3. Chạy gRPC Server (Trong Goroutine riêng biệt)
	// Lý do: Để nó không chặn luồng chính
	go func() {
		grpcPort := cfg.GRPCServer.GRPC_SERVER_PORT // Ví dụ: ":50051"
		log.Printf("Starting gRPC server on %s", grpcPort)
		if err := app.GRPCServer.Run(grpcPort); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// 4. Chạy HTTP Server (Gin) (Ở luồng chính)
	httpPort := fmt.Sprintf(":%v", cfg.Server.Port)
	log.Printf("Starting HTTP server on %s", httpPort)

	// Hàm này sẽ chặn (block) tại đây để giữ app luôn chạy
	if err := app.Server.Run(httpPort); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
