package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	 "github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) IRepositoryShare.IRedis {
	return &RedisAdapter{
		client: client,
	}
}
func (r *RedisAdapter) Ping(ctx context.Context) (*response.Response, error) {
	result, err := r.client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Redis ping successful"),
		response.WithStatus("success"),
	), nil

}

func (r *RedisAdapter) Exists(ctx context.Context, keys ...string) (*response.Response, error) {
	result, err := r.client.Exists(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Keys existence checked successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) Del(ctx context.Context, keys ...string) (*response.Response, error) {
	result, err := r.client.Del(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Keys deleted successfully"),
		response.WithStatus("success"),
	), nil
}

func (r *RedisAdapter) Expire(ctx context.Context, key string, expiration time.Duration) (*response.Response, error) {
	result, err := r.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key expiration set successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) TTL(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key TTL retrieved successfully"),
		response.WithStatus("success"),
	), nil
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value any, expiration time.Duration) (*response.Response, error) {
	result, err := r.client.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key set successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (*response.Response, error) {
	result, err := r.client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key set successfully if not exists"),
		response.WithStatus("success"),
	), nil
} // Quan trọng cho Distributed Lock
func (r *RedisAdapter) Get(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key retrieved successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) MGet(ctx context.Context, keys ...string) (*response.Response, error) {
	result, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Multiple keys retrieved successfully"),
		response.WithStatus("success"),
	), nil
} // Lấy nhiều key 1 lúc (tối ưu hiệu năng)
func (r *RedisAdapter) Incr(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key incremented successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) IncrBy(ctx context.Context, key string, value int64) (*response.Response, error) {
	result, err := r.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key incremented by value successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) Decr(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Key decremented successfully"),
		response.WithStatus("success"),
	), nil
}

func (r *RedisAdapter) HSet(ctx context.Context, key string, values ...any) (*response.Response, error) {
	result, err := r.client.HSet(ctx, key, values...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Hash fields set successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) HGet(ctx context.Context, key, field string) (*response.Response, error) {
	result, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Hash field retrieved successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) HGetAll(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Hash fields retrieved successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) HDel(ctx context.Context, key string, fields ...string) (*response.Response, error) {
	result, err := r.client.HDel(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Hash fields deleted successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) HIncrBy(ctx context.Context, key, field string, incr int64) (*response.Response, error) {
	result, err := r.client.HIncrBy(ctx, key, field, incr).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Hash field incremented successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) RPush(ctx context.Context, key string, values ...any) (*response.Response, error) {
	result, err := r.client.RPush(ctx, key, values...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Values pushed to list successfully"),
		response.WithStatus("success"),
	), nil
} // Thêm vào đuôi
func (r *RedisAdapter) LPop(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.LPop(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Value popped from list successfully"),
		response.WithStatus("success"),
	), nil
} // Lấy từ đầu
func (r *RedisAdapter) LLen(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.LLen(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("List length retrieved successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) SAdd(ctx context.Context, key string, members ...any) (*response.Response, error) {
	result, err := r.client.SAdd(ctx, key, members...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Members added to set successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) SMembers(ctx context.Context, key string) (*response.Response, error) {
	result, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Set members retrieved successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) SIsMember(ctx context.Context, key string, member any) (*response.Response, error) {
	result, err := r.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Membership checked successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) BLPop(ctx context.Context, timeout time.Duration, keys ...string) (*response.Response, error) {
	result, err := r.client.BLPop(ctx, timeout, keys...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Blocking pop from list successful"),
		response.WithStatus("success"),
	), nil
} // Block lấy (cho Worker)
func (r *RedisAdapter) SRem(ctx context.Context, key string, members ...any) (*response.Response, error) {
	result, err := r.client.SRem(ctx, key, members...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Members removed from set successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) ZAdd(ctx context.Context, key string, members ...IRepositoryShare.Z) (*response.Response, error) {
	// 1. Create a slice of the type the redis library expects
	redisMembers := make([]redis.Z, len(members))

	// 2. Map the data from your shared struct to the library's struct
	for i, m := range members {
		redisMembers[i] = redis.Z{
			Score:  m.Score,
			Member: m.Member,
		}
	}

	// 3. Pass the new slice using the ellipsis (...) operator
	result, err := r.client.ZAdd(ctx, key, redisMembers...).Result()
	if err != nil {
		return nil, err
	}

	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Members added to sorted set successfully"),
		response.WithStatus("success"),
	), nil
} // Cần định nghĩa struct Z
func (r *RedisAdapter) ZRange(ctx context.Context, key string, start, stop int64) (*response.Response, error) {
	result, err := r.client.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Sorted set range retrieved successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) ZRevRange(ctx context.Context, key string, start, stop int64) (*response.Response, error) {
	result, err := r.client.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Sorted set reverse range retrieved successfully"),
		response.WithStatus("success"),
	), nil
} // Top cao nhất
func (r *RedisAdapter) ZRem(ctx context.Context, key string, members ...any) (*response.Response, error) {
	result, err := r.client.ZRem(ctx, key, members...).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Members removed from sorted set successfully"),
		response.WithStatus("success"),
	), nil
}
func (r *RedisAdapter) ZScore(ctx context.Context, key, member string) (*response.Response, error) {
	result, err := r.client.ZScore(ctx, key, member).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Member score retrieved successfully"),
		response.WithStatus("success"),
	), nil
}

func (r *RedisAdapter) Scan(ctx context.Context, match string) (*response.Response, error) {
	result, _, err := r.client.Scan(ctx, 0, match, 0).Result()
	if err != nil {
		return nil, err
	}
	return response.NewResponse(
		response.WithData(result),
		response.WithMessage("Keys scanned successfully"),
		response.WithStatus("success"),
	), nil
} // Tìm kiếm key theo pattern an toàn hơn Keys
func (r *RedisAdapter) CustomizeSetCache(ctx context.Context, items map[string]any) error {
	ttl := 30 * time.Minute // Thời gian hết hạn chung
	for key, value := range items {
		_, err := r.Set(ctx, key, value, ttl)
		if err != nil {
			// Nếu lỗi, trả về lỗi ngay lập tức (hoặc có thể log lại rồi continue tùy logic)
			return fmt.Errorf("failed to set cache for key %s: %w", key, err)
		}
	}

	return nil
}
func (r *RedisAdapter) CustomizeGetCache(ctx context.Context, items []string) (any, string, bool, int, error) {
	var data any
	var cursor string
	var hasNext bool
	var limit int
	for _, item := range items {
		getdata, err := r.Get(ctx, item)
		if err != nil {
			return nil, "", false, 0, err
		}
		switch {
		case strings.Contains(item, "user"):
			data = getdata.Data
		case strings.Contains(item, "cursor"):
			cursor, _ = getdata.Data.(string)
		case strings.Contains(item, "hasnext"):
			hasNext, _ = getdata.Data.(bool)
		case strings.Contains(item, "limit"):
			limit, _ = getdata.Data.(int)
		}
	}
	return data, cursor, hasNext, limit, nil
}
