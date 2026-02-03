package domain

import (
	"context"
	"time"
)

type RedisRepository interface {
	// --- Hệ thống & Kết nối ---
	Ping(ctx context.Context) error
	Close() error

	// --- Các thao tác Key chung ---
	Exists(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)

	// --- Strings (Làm việc với Value đơn giản hoặc JSON string) ---
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)

	// --- Hashes (Lưu Object phức tạp) ---
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// --- Lists (Hàng đợi/Queue) ---
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// --- Sets (Tập hợp duy nhất) ---
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SMembers(ctx context.Context, key string) ([]string, error)

	// --- Sorted Sets (Bảng xếp hạng/Ranking) ---
	ZAdd(ctx context.Context, key string, score float64, member interface{}) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// --- Transactions & Pub/Sub (Nâng cao) ---
	Subscribe(ctx context.Context, channel string) (interface{}, error)
	Publish(ctx context.Context, channel string, message interface{}) error
}
