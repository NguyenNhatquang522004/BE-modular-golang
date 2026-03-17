package IRepositoryShare

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/constants"
)

// SetNX: Nếu thiếu cái này, bạn không thể làm Distributed Lock (khóa để ngăn 2 server cùng xử lý 1 đơn hàng).

// MGet: Nếu thiếu, khi cần lấy thông tin 50 user cache, bạn phải gọi Redis 50 lần (rất chậm). MGet chỉ gọi 1 lần.

// HGetAll: Cực kỳ cần thiết để map dữ liệu Redis Hash thẳng vào struct Go (như User Profile).

// BLPop: Nếu làm Worker xử lý background job đơn giản, đây là hàm bắt buộc (chờ có job mới chạy).

// Scan: Tuyệt đối không dùng hàm Keys trong Production (sẽ treo Redis). Phải dùng Scan.
type IRedis interface {
	// ---------------------------s
	// 1. BASIC KEYS & EXPIRATION (Quản lý chung)
	// ---------------------------
	Ping(ctx context.Context) (*response.Response, error)
	Exists(ctx context.Context, keys ...string) (*response.Response, error)
	Del(ctx context.Context, keys ...string) (*response.Response, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (*response.Response, error)
	TTL(ctx context.Context, key string) (*response.Response, error)

	// ---------------------------
	// 2. STRING (Caching đơn giản, JSON, Token)
	// ---------------------------
	Set(ctx context.Context, key string, value any, expiration time.Duration) (*response.Response, error)
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) (*response.Response, error) // Quan trọng cho Distributed Lock
	Get(ctx context.Context, key string) (*response.Response, error)
	MGet(ctx context.Context, keys ...string) (*response.Response, error) // Lấy nhiều key 1 lúc (tối ưu hiệu năng)

	// ---------------------------
	// 3. COUNTER (Đếm lượt view, limit request)
	// ---------------------------
	Incr(ctx context.Context, key string) (*response.Response, error)
	IncrBy(ctx context.Context, key string, value int64) (*response.Response, error)
	Decr(ctx context.Context, key string) (*response.Response, error)

	// ---------------------------
	// 4. HASH (Lưu Object, Shopping Cart, User Session phức tạp)
	// ---------------------------
	HSet(ctx context.Context, key string, values ...any) (*response.Response, error)
	HGet(ctx context.Context, key, field string) (*response.Response, error)
	HGetAll(ctx context.Context, key string) (*response.Response, error)
	HDel(ctx context.Context, key string, fields ...string) (*response.Response, error)
	HIncrBy(ctx context.Context, key, field string, incr int64) (*response.Response, error)

	// ---------------------------
	// 5. LIST (Message Queue đơn giản, Stack)
	// ---------------------------
	RPush(ctx context.Context, key string, values ...any) (*response.Response, error)             // Thêm vào đuôi
	LPop(ctx context.Context, key string) (*response.Response, error)                             // Lấy từ đầu
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) (*response.Response, error) // Block lấy (cho Worker)
	LLen(ctx context.Context, key string) (*response.Response, error)

	// ---------------------------
	// 6. SET (Lưu Unique ID, Tags, Friend List - không trùng lặp)
	// ---------------------------
	SAdd(ctx context.Context, key string, members ...any) (*response.Response, error)
	SMembers(ctx context.Context, key string) (*response.Response, error)
	SIsMember(ctx context.Context, key string, member any) (*response.Response, error)
	SRem(ctx context.Context, key string, members ...any) (*response.Response, error)

	// ---------------------------
	// 7. SORTED SET (Leaderboard, Ranking, Priority Queue)
	// ---------------------------
	ZAdd(ctx context.Context, key string, members ...Z) (*response.Response, error) // Cần định nghĩa struct Z
	ZRange(ctx context.Context, key string, start, stop int64) (*response.Response, error)
	ZRevRange(ctx context.Context, key string, start, stop int64) (*response.Response, error) // Top cao nhất
	ZRem(ctx context.Context, key string, members ...any) (*response.Response, error)
	ZScore(ctx context.Context, key, member string) (*response.Response, error)

	// ---------------------------
	// 8. UTILITIES (Scan keys - dùng cẩn thận)
	// ---------------------------
	Scan(ctx context.Context, match string) (*response.Response, error) // Tìm kiếm key theo pattern an toàn hơn Keys
	// ---------------------------
	// 9. CUSTOMIZE (Tùy biến nâng cao)
	// ---------------------------
	CustomizeSetCache(ctx context.Context, items map[string]any) error
	//any, string, bool, int, error : data , nextcursor , hasnext , limit
	CustomizeGetCache(ctx context.Context, items []string) (any, string, bool, int, error)

	// Các hàm tùy chỉnh khác (nếu cần)
	Lock(ctx context.Context, eventID string) (constants.ProcessStatus, bool, error) // Hàm khóa để đảm bảo chỉ 1 worker xử lý 1 eventID nhất định (Distributed Lock)
	Unlock(ctx context.Context, eventID string) error                                // Hàm mở khóa sau khi xử lý xong
	MarkCompleted(ctx context.Context, eventID string) error                         // Hàm đánh dấu event đã xử lý xong (để các worker khác biết mà không xử lý lại)
	GetViewedPosts(ctx context.Context, userID string) (map[string]bool, error)
	MarkPostsAsViewed(ctx context.Context, userID string, postIDs ...string) error
}
type Z struct {
	Score  float64
	Member any
}
