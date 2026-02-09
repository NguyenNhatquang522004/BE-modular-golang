package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *configs.Config) (*redis.Client, func(), error) {
	// 1. Tạo Context timeout cho việc kết nối ban đầu
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Format địa chỉ
	addr := fmt.Sprintf("%s:%s", cfg.RedisDB.Host, cfg.RedisDB.Port)

	// 3. Cấu hình Options
	opts := &redis.Options{
		Addr:         addr,
		Password:     cfg.RedisDB.Password,
		DB:           cfg.RedisDB.DB,
		PoolSize:     10, // Nên đưa vào config nếu cần tune sau này
		MinIdleConns: 5,
	}

	// 4. Tạo Client (Lưu ý: NewClient chưa kết nối ngay lập tức, nó chỉ tạo struct)
	client := redis.NewClient(opts)

	// 5. Ping kiểm tra kết nối (Quan trọng)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("✅ Connected to Redis successfully!")

	// 6. Tạo hàm Cleanup cho Wire
	cleanup := func() {
		log.Println("⚠️ Closing Redis connection...")
		if err := client.Close(); err != nil {
			log.Printf("Error closing redis: %v", err)
		}
	}

	// 7. Trả về Client trực tiếp
	return client, cleanup, nil
}
