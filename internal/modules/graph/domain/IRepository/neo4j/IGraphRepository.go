package neo4j

import (
	"context"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/graphEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/entity"
)

// IGraphRepository bao quát toàn bộ vòng đời dữ liệu trên Neo4j
type IGraphRepository interface {

	// ==================================================
	// NHÓM 1: ĐỒNG BỘ DỮ LIỆU CƠ BẢN (WRITE - Dành cho Kafka Consumer)
	// ==================================================

	// Upsert Đỉnh (Nodes)
	UpsertUserNode(ctx context.Context, user *entity.UserNode) error
	UpsertPostNode(ctx context.Context, post *entity.PostNode) error
	UpsertTopicNode(ctx context.Context, topic *entity.TopicNode) error
	UpsertTopicAndConnectNeighbors(ctx context.Context, topicName string, embedding []float32) error
	UpsertGroupNode(ctx context.Context, group *entity.GroupNode) error
	UpsertPageNode(ctx context.Context, page *entity.PageNode) error
	UpsertLocationNode(ctx context.Context, loc *entity.LocationNode) error

	// Tạo Cạnh Tạo Nội Dung (Content Creation Edges)
	LinkAuthorToPost(ctx context.Context, userID string, postID string, createdAt int64) error
	LinkPageToPost(ctx context.Context, pageID string, postID string, createdAt int64) error
	LinkPostToTopic(ctx context.Context, postID string, topicName string, confidenceScore float64) error
	LinkPostToGroup(ctx context.Context, postID string, groupID string, createdAt int64) error

	// ==================================================
	// NHÓM 2: MẠNG LƯỚI XÃ HỘI & TĂNG TRƯỞNG (SOCIAL GRAPHS)
	// ==================================================

	// User kết nối User/Page/Group
	CreateFriendship(ctx context.Context, userA string, userB string, since int64) error
	CreateBlock(ctx context.Context, sourceUserID string, targetUserID string, since int64) error
	CreateFollow(ctx context.Context, followerID string, targetID string, since int64) error
	JoinGroup(ctx context.Context, userID string, groupID string, role string, joinedAt int64) error
	LikePage(ctx context.Context, userID string, pageID string, since int64) error

	// Đồng bộ danh bạ (Cho tính năng PYMK)
	SyncPhoneContact(ctx context.Context, userID string, phoneHash string, uploadedAt int64) error

	// ==================================================
	// NHÓM 3: TƯƠNG TÁC & CÁ NHÂN HÓA (INTERACTIONS & SCORING)
	// ==================================================

	// Tương tác ngắn hạn (Dành cho Trending)
	RecordRecentInteraction(ctx context.Context, userID string, postID string, interactionType string, weight float64, timestamp int64) error

	// Tương tác dài hạn (Tính điểm EdgeRank/Affinity)
	IncrementAffinityScore(ctx context.Context, sourceUserID string, targetID string, weight float64) error
	IncrementInteraction(
		ctx context.Context,
		userID string,
		targetID string,
		targetType string,
		like, comment, share, message, view int, // Số lượng thay đổi (thường là 1)
		weight float64, // Trọng số thay đổi (có thể dương hoặc âm)
		flag bool, // true = cộng, false = trừ
	) error
	IncrementInteractionByPost(
		ctx context.Context,
		userID string,
		postID string,
		like, comment, share, message, view int,
		flag bool,
	) (map[string]float64, error)
	// Tích lũy sở thích

	IncrementTopicInterest(ctx context.Context, userID string, limitK int) error
	CreateSharePost(ctx context.Context, userID string, newSharePostID string, originalPostID string, createdAt int64) error
	// ==================================================
	// NHÓM 4: TRUY VẤN DỮ LIỆU (READ - Dành cho REST/gRPC API)
	// ==================================================

	// 4.1. News Feed & Recommendation
	// Lấy Bảng tin cá nhân hóa (Mix giữa bạn bè, group, page và sở thích)
	GetSocialFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error)
	GetInterestFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error)
	GetSemanticDiscoveryFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error)
	GetGlobalTrendingFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error) // 4.2. People You May Know (PYMK)
	GetLocationFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error)
	GetContactSyncFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error)
	GetPeopleYouMayKnow(ctx context.Context, userID string, limit int) ([]graphEvent.SuggestedUser, error)
	GetPersonalizedNewsFeed(ctx context.Context, userID string, limit int, offset int) ([]string, error)

	// Gợi ý bạn bè dựa trên "Bạn chung" (Triadic Closure) và Danh bạ
	GetSuggestedFriends(ctx context.Context, userID string, limit int) ([]string, error)

	// 4.3. Real-time Discovery
	// Lấy danh sách Post đang hot nhất dựa trên InteractedRecentlyRel
	GetTrendingPosts(ctx context.Context, limit int) ([]string, error)
	// Lấy danh sách Topic đang hot
	GetTrendingTopics(ctx context.Context, limit int) ([]string, error)

	// 4.4. Graph Analytics & Maintenance (Thường gọi qua Cronjob)
	// Xóa các Post đã hết hạn TTL khỏi Graph
	CleanupExpiredPosts(ctx context.Context, currentTime int64) (int64, error)
}
