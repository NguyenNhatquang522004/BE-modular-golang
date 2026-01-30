package database

import (
	"context"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	RedisDB *redis.Client
}
func NewRedisConnection(cfg *configs.Config) (*RedisConnection, error) {
	db := &RedisConnection{}
	err := db.ConnectRedis(&cfg.RedisDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func (r *RedisConnection) ConnectRedis(cfg *configs.RedisConfig) error {
	// Tạo context có timeout cho việc Ping kiểm tra
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Format địa chỉ "Host:Port"
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Cấu hình Options
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password, // Mật khẩu (để trống nếu không có)
		DB:       cfg.DB,       // Database index (mặc định là 0)

		// Một số cấu hình Pool mặc định (có thể tùy chỉnh nếu cần)
		PoolSize:     10, // Số lượng kết nối tối đa trong pool
		MinIdleConns: 5,  // Số lượng kết nối nhàn rỗi tối thiểu giữ lại
	})

	// 3. Ping thử để đảm bảo Redis server đang sống
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	fmt.Println("Connected to Redis successfully!")

	// 4. Gán client vào struct (QUAN TRỌNG)
	r.RedisDB = rdb

	return nil
}

// 3. Hàm Getter để lấy Client ra dùng
func (r *RedisConnection) GetClient() *redis.Client {
	return r.RedisDB
}
