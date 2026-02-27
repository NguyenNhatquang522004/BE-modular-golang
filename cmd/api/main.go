package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.LifecycleManager.Start(ctx)

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
	srv := &http.Server{
		Addr:    httpPort,
		Handler: app.Server, // Gán Gin Engine làm Handler	
	}
	go func() {
		log.Printf("🚀 Starting HTTP server on %s", httpPort)
		// ListenAndServe luôn trả về error, trừ khi nó là ErrServerClosed (lỗi do mình tắt)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server failed: %v", err)
		}
	}()
	log.Printf("Starting HTTP server on %s", httpPort)

	// Hàm này sẽ chặn (block) tại đây để giữ app luôn chạy
	if err := app.Server.Run(httpPort); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}

	// --- 3. GRACEFUL SHUTDOWN ---
	quit := make(chan os.Signal, 1)	
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Chờ tín hiệu tắt
	log.Println("🛑 Shutting down system...")
	timeoutCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	// Bước A: Tắt HTTP Server (để ngừng nhận request mới)
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Printf("⚠️ HTTP Server Shutdown Error: %v", err)
	} else {
		log.Println("✅ HTTP Server stopped.")
	}

	// Bước B: Hủy context (báo cho Kafka Consumer ngừng fetch tin mới)
	cancel()

	// Bước C: TẮT TẤT CẢ WORKER (Chỉ 1 dòng duy nhất!)
	// Nó sẽ tự gọi Stop() cho OrderConsumer, NotifConsumer,... và chờ chúng xong
	app.LifecycleManager.Stop()

	log.Println("👋 Server exited properly")
}
