package neo4j

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/graphEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/sync/errgroup"
)

type GraphRepository struct {
	// Có thể thêm các trường như Neo4j Driver, Logger, Config nếu cần
	driver    neo4j.DriverWithContext
	aiRepo    IRepositoryShare.IAI
	cfg       *configs.Config
	redisRepo IRepositoryShare.IRedis
}

func NewGraphRepository(driver neo4j.DriverWithContext, aiRepo IRepositoryShare.IAI, cfg *configs.Config, redisRepo IRepositoryShare.IRedis) *GraphRepository {
	return &GraphRepository{
		driver:    driver,
		aiRepo:    aiRepo,
		cfg:       cfg,
		redisRepo: redisRepo,
	}
}
func (r *GraphRepository) UpsertUserNode(ctx context.Context, user *entity.UserNode) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Dùng MERGE để đảm bảo không tạo User trùng lặp.
	// ON CREATE: Set lúc mới tạo | ON MATCH: Cập nhật lúc đã có
	cypher := `
		MERGE (u:User {user_id: $user_id})
		ON CREATE SET 
			u.created_at = $created_at,
			u.last_active_at = $last_active_at,
			u.is_verified = $is_verified
		ON MATCH SET 
			u.last_active_at = $last_active_at
	`
	params := map[string]any{
		"user_id":        user.UserID,
		"created_at":     user.CreatedAt,
		"last_active_at": user.LastActiveAt,
		"is_verified":    user.IsVerified,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) UpsertPostNode(ctx context.Context, post *entity.PostNode) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MERGE (p:Post {post_id: $post_id})
		ON CREATE SET 
			p.created_at = $created_at,
			p.ttl = $ttl
	`
	params := map[string]any{
		"post_id":    post.PostID,
		"created_at": post.CreatedAt,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}

func (r *GraphRepository) UpsertTopicNode(ctx context.Context, topic *entity.TopicNode) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := ` 
		MERGE (t:Topic {name: $name})
		ON CREATE SET
		t.name = $name,
		t.trending_score = $trending_score
		`
	params := map[string]any{
		"name":           topic.Name,
		"trending_score": topic.TrendingScore,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}

// UpsertTopicAndConnectNeighbors trả về Tên Topic chuẩn (Canonical Name)
// Ví dụ: Nhập "Trí tuệ nhân tạo", phát hiện giống 99% với "AI", nó sẽ không tạo mới mà trả về chữ "AI".
// UpsertTopicAndConnectNeighbors nhận vào một chuỗi topicName thô, tự động gọi AI lấy Vector,
// khử trùng lặp và vẽ mạng lưới ngữ nghĩa. Trả về Canonical Name.
func (r *GraphRepository) UpsertTopicAndConnectNeighbors(ctx context.Context, topicName string) (string, error) {
	// ==========================================
	// BƯỚC 1: GỌI AI LLM ĐỂ LẤY VECTOR (Nằm ngoài DB Transaction)
	// ==========================================
	embedding, err := r.aiRepo.GenerateBgeM3Embedding(topicName)
	if err != nil {
		return "", fmt.Errorf("failed to generate embedding for topic '%s': %w", topicName, err)
	}

	// ==========================================
	// BƯỚC 2: TƯƠNG TÁC VỚI NEO4J
	// ==========================================
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)
	topKNeighbors := r.cfg.AIAnalysisMediaAsset.TopKNeighbors
	// Các cấu hình ngưỡng (Thresholds)
	exactMatchThreshold := 0.99
	similarityThreshold := 0.85

	canonicalTopicName := topicName // Mặc định là tên gốc nếu không bị trùng

	// Bắt đầu Transaction (Càng nhanh càng tốt)
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		// ------------------------------------------
		// PHASE 1: DÒ TÌM TRÙNG LẶP NGỮ NGHĨA (DEDUP)
		// ------------------------------------------
		searchCypher := `
			CALL db.index.vector.queryNodes('topic_embeddings', 1, $embedding)
			YIELD node AS existingTopic, score
			WHERE score >= $exactMatchThreshold
			RETURN existingTopic.name AS matchedName
		`

		searchResult, err := tx.Run(ctx, searchCypher, map[string]any{
			"embedding":           embedding,
			"exactMatchThreshold": exactMatchThreshold,
		})
		if err != nil {
			return nil, err
		}

		// Nếu tìm thấy Topic giống 99% -> Ghi nhận tên của nó và ABORT việc tạo mới
		if searchResult.Next(ctx) {
			canonicalTopicName = searchResult.Record().Values[0].(string)
			// Phát hiện trùng lặp, bỏ qua Phase 2 và kết thúc Transaction an toàn.
			return nil, nil
		}

		// ------------------------------------------
		// PHASE 2: TẠO MỚI & VẼ MẠNG LƯỚI TƯƠNG ĐỒNG
		// ------------------------------------------
		upsertCypher := `
			// 1. Tạo Topic Mới
			MERGE (newTopic:Topic {name: $topicName})
			SET newTopic.embedding = $embedding
			WITH newTopic

			// 2. Quét hàng xóm để vẽ lưới
			CALL db.index.vector.queryNodes('topic_embeddings', $topKNeighbors + 1, $embedding)
			YIELD node AS neighborTopic, score AS similarityScore

			// 3. Lọc bỏ chính nó và lọc theo ngưỡng họ hàng (0.85)
			WHERE neighborTopic.name <> newTopic.name 
			  AND similarityScore >= $similarityThreshold

			// 4. Vẽ Cạnh Tương Đồng vô hướng (-)
			MERGE (newTopic)-[rel:RELATED_TO]-(neighborTopic)
			SET rel.similarity_score = similarityScore
		`

		_, err = tx.Run(ctx, upsertCypher, map[string]any{
			"topicName":           topicName,
			"embedding":           embedding,
			"similarityThreshold": similarityThreshold,
			"topKNeighbors":       topKNeighbors,
		})

		return nil, err
	})

	if err != nil {
		return "", err
	}

	// Trả về Topic Name chuẩn (Canonical Name) cho các luồng xử lý tiếp theo
	return canonicalTopicName, nil
}
func (r *GraphRepository) UpsertGroupNode(ctx context.Context, group *entity.GroupNode) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MERGE (g:Group {group_id: $group_id})
		ON CREATE SET 
			g.privacy = $privacy,
			g.member_count = $member_count
		ON MATCH SET
			g.privacy = $privacy,
			g.member_count = $member_count
	`
	params := map[string]any{
		"group_id":     group.GroupID,
		"privacy":      group.Privacy,
		"member_count": group.MemberCount,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) UpsertPageNode(ctx context.Context, page *entity.PageNode) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MERGE (p:Page {page_id: $page_id})
		ON CREATE SET 
			p.category_id = $category_id,
			p.rating = $rating
		ON MATCH SET
			p.category_id = $category_id,
			p.rating = $rating
	`
	params := map[string]any{
		"page_id":     page.PageID,
		"category_id": page.CategoryID,
		"rating":      page.Rating,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) UpsertLocationNode(ctx context.Context, loc *entity.LocationNode) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MERGE (l:Location {city_id: $city_id, country_code: $country_code})
		ON CREATE SET 
			l.geo_hash = $geo_hash
		ON MATCH SET
			l.geo_hash = $geo_hash
	`
	params := map[string]any{
		"city_id":      loc.CityID,
		"country_code": loc.CountryCode,
		"geo_hash":     loc.GeoHash,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) LinkAuthorToPost(ctx context.Context, userID string, postID string, createdAt int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (u:User {user_id: $user_id}), (p:Post {post_id: $post_id})
		MERGE (u)-[r:AUTHORED]->(p)
		ON CREATE SET r.created_at = $created_at
	`
	params := map[string]any{
		"user_id":    userID,
		"post_id":    postID,
		"created_at": createdAt,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) LinkPageToPost(ctx context.Context, pageID string, postID string, createdAt int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (pg:Page {page_id: $page_id}), (p:Post {post_id: $post_id})
		MERGE (pg)-[r:PUBLISHED]->(p)
		ON CREATE SET r.created_at = $created_at
	`
	params := map[string]any{
		"page_id":    pageID,
		"post_id":    postID,
		"created_at": createdAt,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) LinkPostToTopic(ctx context.Context, postID string, topicName string, confidenceScore float64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)
	cypher := `
		MATCH (p:Post {post_id: $post_id}), (t:Topic {name: $topic_name})
		MERGE (p)-[r:HAS_TOPIC]->(t)
		ON CREATE SET r.confidence_score = $confidence_score
		ON MATCH SET r.confidence_score = $confidence_score
	`
	params := map[string]any{
		"post_id":          postID,
		"topic_name":       topicName,
		"confidence_score": confidenceScore,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) LinkPostToGroup(ctx context.Context, postID string, groupID string, createdAt int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (p:Post {post_id: $post_id}), (g:Group {group_id: $group_id})
		MERGE (p)-[r:POSTED_IN]->(g)
		ON CREATE SET r.created_at = $created_at
	`
	params := map[string]any{
		"post_id":    postID,
		"group_id":   groupID,
		"created_at": createdAt,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) CreateFriendship(ctx context.Context, userA string, userB string, friendshipType string, interactionType int, since int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (a:User {user_id: $userA}), (b:User {user_id: $userB})
		MERGE (a)-[r:FRIEND]->(b)
		ON CREATE SET r.since = $since, r.type = $type, r.interaction_frequency = $interaction_frequency
		ON MATCH SET r.type = $type, r.interaction_frequency = $interaction_frequency
		MERGE (b)-[r2:FRIEND]->(a)
		ON CREATE SET r2.since = $since, r2.type = $type, r2.interaction_frequency = $interaction_frequency
		ON MATCH SET r2.type = $type, r2.interaction_frequency = $interaction_frequency
	
	`
	params := map[string]any{
		"userA":                 userA,
		"userB":                 userB,
		"since":                 since,
		"type":                  friendshipType,
		"interaction_frequency": interactionType,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) CreateBlock(ctx context.Context, sourceUserID string, targetUserID string, since int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (source:User {user_id: $sourceUserID}), (target:User {user_id: $targetUserID})
		MERGE (source)-[r:BLOCK]->(target)
		ON CREATE SET r.since = $since
	`
	params := map[string]any{
		"sourceUserID": sourceUserID,
		"targetUserID": targetUserID,
		"since":        since,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) CreateFollow(ctx context.Context, followerID string, targetID string, targetType sharedEnums.ContextType, since int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)
	var targetMatch string
	switch targetType {
	case sharedEnums.ContextTypeUserWall:
		targetMatch = "(target:User {user_id: $targetID})"
	case sharedEnums.ContextTypeGroup:
		targetMatch = "(target:Group {group_id: $targetID})"
	case sharedEnums.ContextTypePage:
		targetMatch = "(target:Page {page_id: $targetID})"
	default:
		return fmt.Errorf("invalid target type: %v", targetType)
	}
	cypher := fmt.Sprintf(`
		MATCH (follower:User {user_id: $followerID})
		MATCH %s
		MERGE (follower)-[r:FOLLOWS]->(target)
		ON CREATE SET r.since = $since
	`, targetMatch)

	params := map[string]any{
		"followerID": followerID,
		"targetID":   targetID,
		"since":      since,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) JoinGroup(ctx context.Context, userID string, groupID string, role string, joinedAt int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (u:User {user_id: $user_id}), (g:Group {group_id: $group_id})
		MERGE (u)-[r:MEMBER_OF]->(g)
		ON CREATE SET r.role = $role, r.joined_at = $joined_at
		ON MATCH SET r.role = $role
	`
	params := map[string]any{
		"user_id":   userID,
		"group_id":  groupID,
		"role":      role,
		"joined_at": joinedAt,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) LikePage(ctx context.Context, userID string, pageID string, since int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (u:User {user_id: $user_id}), (p:Page {page_id: $page_id})
		MERGE (u)-[r:LIKES]->(p)
		ON CREATE SET r.since = $since
	`
	params := map[string]any{
		"user_id": userID,
		"page_id": pageID,
		"since":   since,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return err
}
func (r *GraphRepository) SyncPhoneContact(ctx context.Context, userID string, phoneHash string, uploadedAt int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (u:User {user_id: $user_id})
		MERGE (c:PhoneContact {phone_hash: $phone_hash})
		ON CREATE SET c.created_at = $uploaded_at
		MERGE (u)-[r:HAS_CONTACT]->(c)
		ON CREATE SET r.uploaded_at = $uploaded_at
	`
	params := map[string]any{
		"user_id":     userID,
		"phone_hash":  phoneHash,
		"uploaded_at": uploadedAt,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}

// Tương tác ngắn hạn (Dành cho Trending)
func (r *GraphRepository) RecordRecentInteraction(ctx context.Context, userID string, postID string, interactionType string, weight float64, timestamp int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (u:User {user_id: $user_id}), (p:Post {post_id: $post_id})
		MERGE (u)-[r:INTERACTED_RECENTLY]->(p)
		ON CREATE SET r.type = $interaction_type, r.weight = $weight, r.timestamp = $timestamp
		ON MATCH SET r.type = $interaction_type, r.weight = $weight, r.timestamp = $timestamp
	`
	params := map[string]any{
		"user_id":          userID,
		"post_id":          postID,
		"interaction_type": interactionType,
		"weight":           weight,
		"timestamp":        timestamp,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}

// Tương tác dài hạn (Tính điểm EdgeRank/Affinity)
func (r *GraphRepository) IncrementAffinityScore(ctx context.Context, sourceUserID string, targetID string, weight float64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		MATCH (source:User {user_id: $sourceUserID})-[r]->(target {user_id: $targetID})
		WHERE type(r) IN ['FRIEND', 'FOLLOWS', 'MEMBER_OF']
		SET r.affinity_score = coalesce(r.affinity_score, 0) + $weight
	`
	params := map[string]any{
		"sourceUserID": sourceUserID,
		"targetID":     targetID,
		"weight":       weight,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})

	return err
}
func (r *GraphRepository) IncrementInteraction(
	ctx context.Context,
	userID string,
	targetID string,
	targetType sharedEnums.ContextType,
	like, comment, share, message, view int,
	weight float64, // Bạn có thể bỏ qua nếu muốn dùng weight cố định bên dưới
	flag bool,
) (float64, error) { // Trả về AffinityScore mới sau khi cập nhật
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// 1. Xác định Label dựa trên targetType (Sử dụng whitelist để bảo mật fmt.Sprintf)
	var targetMatch string
	switch targetType {
	case sharedEnums.ContextTypeUserWall:
		targetMatch = "(target:User {user_id: $targetID})"
	case sharedEnums.ContextTypePage:
		targetMatch = "(target:Page {page_id: $targetID})"
	case sharedEnums.ContextTypeGroup:
		targetMatch = "(target:Group {group_id: $targetID})"
	default:
		return 0, fmt.Errorf("invalid target type: %v", targetType)
	}

	// 2. Chuẩn bị hệ số multiplier
	multiplier := 1.0
	if !flag {
		multiplier = -1.0
	}

	// 3. Câu lệnh Cypher tối ưu
	// - Sử dụng WITH để tính toán biến trung gian, giúp câu lệnh sạch hơn.
	// - Sử dụng logic tính điểm Ranking chuẩn ngay trong DB.
	cypher := fmt.Sprintf(`
		MATCH (u:User {user_id: $userID})
		MATCH %s
		MERGE (u)-[r:INTERACTED_WITH]->(target)
		ON CREATE SET 
			r.like_count = 0, r.comment_count = 0, r.share_count = 0, 
			r.message_count = 0, r.profile_view_count = 0, r.affinity_score = 0.0

		// Tính toán giá trị dự kiến
		WITH r, $m AS m,
		     (r.like_count + ($like * m)) AS nLike,
		     (r.comment_count + ($comment * m)) AS nComment,
		     (r.share_count + ($share * m)) AS nShare,
		     (r.message_count + ($message * m)) AS nMessage,
		     (r.profile_view_count + ($view * m)) AS nView

		// Cập nhật và đảm bảo không âm (Floor at 0)
		SET r.like_count = CASE WHEN nLike > 0 THEN nLike ELSE 0 END,
		    r.comment_count = CASE WHEN nComment > 0 THEN nComment ELSE 0 END,
		    r.share_count = CASE WHEN nShare > 0 THEN nShare ELSE 0 END,
		    r.message_count = CASE WHEN nMessage > 0 THEN nMessage ELSE 0 END,
		    r.profile_view_count = CASE WHEN nView > 0 THEN nView ELSE 0 END,
		    r.last_interaction_at = $now

		// TÍNH TOÁN LẠI AFFINITY SCORE (Ranking Algorithm)
		// Trọng số: View=1, Like=3, Comment=5, Message=8, Share=10
		// Bạn có thể tùy chỉnh các hằng số này theo chiến lược của mình
		WITH r
		SET r.affinity_score = (r.profile_view_count * 1.0) + 
		                       (r.like_count * 3.0) + 
		                       (r.comment_count * 5.0) + 
		                       (r.message_count * 8.0) + 
		                       (r.share_count * 10.0)

		RETURN r.affinity_score AS newScore
	`, targetMatch)

	params := map[string]any{
		"userID":   userID,
		"targetID": targetID,
		"like":     like,
		"comment":  comment,
		"share":    share,
		"message":  message,
		"view":     view,
		"m":        multiplier,
		"now":      time.Now().Unix(),
	}

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return 0.0, err
		}
		if res.Next(ctx) {
			return res.Record().Values[0].(float64), nil
		}
		return 0.0, fmt.Errorf("failed to retrieve updated score")
	})

	if err != nil {
		return 0, err
	}

	return result.(float64), nil
}
func (r *GraphRepository) IncrementInteractionByPost(
	ctx context.Context,
	userID string,
	postID string,
	like, comment, share, message, view int,
	flag bool,
) (map[string]float64, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	multiplier := 1.0
	if !flag {
		multiplier = -1.0
	}

	// CÂU LỆNH CYPHER THẦN THÁNH:
	// 1. Tìm User tương tác và Bài Post.
	// 2. Tìm TẤT CẢ các Đỉnh có nối với Post qua cạnh AUTHORED, PUBLISHED, hoặc POSTED_IN.
	//    (Dùng dấu "-" không định hướng để quét cả 2 chiều).
	// 3. Loại trừ trường hợp tự tương tác với chính mình.
	// 4. Cập nhật đồng loạt cho tất cả các "target" tìm được.
	cypher := `
		MATCH (u:User {user_id: $userID})
		MATCH (p:Post {post_id: $postID})
		
		// Quét tìm Chủ nhân (User/Page) hoặc Nơi chứa (Group)
		MATCH (target)-[:AUTHORED|PUBLISHED|POSTED_IN]-(p)
		
		// Đảm bảo không tự tăng điểm thân thiết với chính bản thân mình
		WHERE elementId(u) <> elementId(target)

		// Bắt đầu cập nhật cho TỪNG target tìm được
		MERGE (u)-[r:INTERACTED_WITH]->(target)
		ON CREATE SET 
			r.like_count = 0, r.comment_count = 0, r.share_count = 0, 
			r.message_count = 0, r.profile_view_count = 0, r.affinity_score = 0.0

		// Tính toán biến tạm
		WITH u, target, r, $m AS m,
		     (r.like_count + ($like * m)) AS nLike,
		     (r.comment_count + ($comment * m)) AS nComment,
		     (r.share_count + ($share * m)) AS nShare,
		     (r.message_count + ($message * m)) AS nMessage,
		     (r.profile_view_count + ($view * m)) AS nView

		// Cập nhật và chặn số âm
		SET r.like_count = CASE WHEN nLike > 0 THEN nLike ELSE 0 END,
		    r.comment_count = CASE WHEN nComment > 0 THEN nComment ELSE 0 END,
		    r.share_count = CASE WHEN nShare > 0 THEN nShare ELSE 0 END,
		    r.message_count = CASE WHEN nMessage > 0 THEN nMessage ELSE 0 END,
		    r.profile_view_count = CASE WHEN nView > 0 THEN nView ELSE 0 END,
		    r.last_interaction_at = $now

		// Tính Affinity Score (Ranking)
		WITH target, r
		SET r.affinity_score = (r.profile_view_count * 1.0) + 
		                       (r.like_count * 3.0) + 
		                       (r.comment_count * 5.0) + 
		                       (r.message_count * 8.0) + 
		                       (r.share_count * 10.0)

		// Thu thập ID của target (vì target có thể là User, Page hoặc Group nên ta dùng coalesce để rút gọn)
		RETURN coalesce(target.user_id, target.page_id, target.group_id) AS targetID, 
		       r.affinity_score AS newScore
	`

	params := map[string]any{
		"userID":  userID,
		"postID":  postID,
		"like":    like,
		"comment": comment,
		"share":   share,
		"message": message,
		"view":    view,
		"m":       multiplier,
		"now":     time.Now().Unix(),
	}

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		// Tạo Map để hứng kết quả
		updatedScores := make(map[string]float64)
		for res.Next(ctx) {
			record := res.Record()
			tID := record.Values[0].(string)
			score := record.Values[1].(float64)
			updatedScores[tID] = score
		}
		return updatedScores, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(map[string]float64), nil
}
func (r *GraphRepository) UpdateInterestGraph(ctx context.Context, userID string, postID string, actionWeight float64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// LEARNING_RATE: Xác định việc User thay đổi sở thích nhanh hay chậm.
	learningRate := 0.05

	cypher := `
		MATCH (u:User {user_id: $userID})
		MATCH (p:Post {post_id: $postID})
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)-[ht:HAS_TOPIC]->(t:Topic)
		// 1. Tìm các Topic của Bài Post
		MATCH (p)-[ht:HAS_TOPIC]->(t:Topic)

		// 2. Kết nối User với Topic (MERGE chuẩn 100% theo Struct)
		MERGE (u)-[r:INTERESTED_IN]->(t)
		ON CREATE SET 
		    r.score = 0.0, 
		    r.short_term_score = 0.0, // [BỔ SUNG 1] Đồng bộ với Struct mới
		    r.last_engaged_at = 0

		// 3. Tính toán Lực tác động tuyệt đối (Absolute Alpha)
		WITH u, t, r, ht, (abs($actionWeight) * coalesce(ht.confidence_score, 0.5) * $learningRate) AS alpha, $actionWeight AS rawWeight
		
		// 4. CÔNG THỨC TOÁN HỌC ĐỐI XỨNG [BỔ SUNG 2]
		WITH r, coalesce(r.score, 0.0) AS currentScore, alpha, rawWeight
		WITH r, currentScore, 
		     CASE 
		         // Nếu hành động Dương (Like/Share) -> Tịnh tiến lên 1.0
		         WHEN rawWeight > 0 THEN currentScore + (alpha * (1.0 - currentScore))
		         // Nếu hành động Âm (Unlike/Hide) -> Tịnh tiến về 0.0 (Điểm càng cao trừ càng mạnh)
		         WHEN rawWeight < 0 THEN currentScore - (alpha * currentScore)
		         // Nếu Weight = 0 thì giữ nguyên
		         ELSE currentScore 
		     END AS predictedScore
		
		// 5. ÉP KIỂU (CLAMPING): Lớp khiên cuối cùng đảm bảo không bao giờ lố [0.0, 1.0] do sai số Float
		SET r.score = CASE 
			WHEN predictedScore > 1.0 THEN 1.0
			WHEN predictedScore < 0.0 THEN 0.0
			ELSE predictedScore 
		END,
		r.last_engaged_at = $now
	`

	params := map[string]any{
		"userID":       userID,
		"postID":       postID,
		"actionWeight": actionWeight, // Truyền số âm vào thoải mái!
		"learningRate": learningRate,
		"now":          time.Now().Unix(),
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	return err
}

// Tích lũy sở thích
func (r *GraphRepository) IncrementTopicInterest(ctx context.Context, userID string, limitK int) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// lambda: Hệ số suy giảm thời gian.
	// Ví dụ: 0.000008 làm cho bài viết cách đây 1 ngày bị giảm khoảng 50% giá trị.
	lambda := 0.000008

	cypher := `
		MATCH (u:User {user_id: $userID})-[r_recent:INTERACTED_RECENTLY]->(p:Post)
		
		// 1. CẮT KHUNG CỬA SỔ (SLIDING WINDOW)
		// Lấy Top K bài viết tương tác gần nhất của User này
		WITH u, r_recent, p
		ORDER BY r_recent.timestamp DESC
		LIMIT $limitK
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)-[ht:HAS_TOPIC]->(t:Topic)

		// 2. LẤY TOPIC TỪ CÁC BÀI VIẾT NÀY
		MATCH (p)-[ht:HAS_TOPIC]->(t:Topic)

		// 3. TÍNH TOÁN ĐIỂM SỐ BẰNG TOÁN HỌC (TIME DECAY)
		// Tính khoảng cách thời gian (Delta T) tính bằng giây
		WITH u, t, r_recent, ht, ($now - r_recent.timestamp) AS deltaT
		
		// Áp dụng công thức: Trọng số hành động * Độ tự tin AI * Phân rã thời gian
		WITH u, t, (r_recent.weight * ht.confidence_score * exp(-$lambda * deltaT)) AS interactionScore
		
		// 4. GOM NHÓM (AGGREGATION) VÀ TÍNH TỔNG ĐIỂM CHO TỪNG TOPIC
		WITH u, t, SUM(interactionScore) AS totalTopicScore
		
		// 5. CHUẨN HÓA VÀ LƯU TRỮ (Ghi đè hoặc tạo mới cạnh INTERESTED_IN)
		// Điểm này thể hiện độ "cuồng" của user với topic đó trong thời gian gần đây
		MERGE (u)-[rel:INTERESTED_IN]->(t)
		ON CREATE SET rel.score = 0.0
		SET rel.short_term_score = totalTopicScore,
		    rel.last_engaged_at = $now
	`

	params := map[string]any{
		"userID": userID,
		"limitK": limitK, // Ví dụ: 50
		"lambda": lambda,
		"now":    time.Now().Unix(),
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	return err
}

// CreateSharePost tạo một bài Share (Wrapper Post) trỏ về bài gốc
// CreateSharePost tạo một bài Share (Wrapper Post) trỏ về bài gốc và bài trung gian
func (r *GraphRepository) CreateSharePost(ctx context.Context, userID string, newSharePostID string, originalPostID string, parentPostID string, createdAt int64) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	cypher := `
		// 1. Tìm User Share, Bài Gốc và Bài Trung Gian (Parent)
		MATCH (u:User {user_id: $userID})
		MATCH (original:Post {post_id: $originalPostID})
		MATCH (parent:Post {post_id: $parentPostID})
		
		// 2. Tạo Bài Post mới (Bài Share)
		MERGE (sharePost:Post {post_id: $newSharePostID})
		ON CREATE SET sharePost.created_at = $createdAt
		
		// 3. Nối User với Bài Share
		MERGE (u)-[auth:AUTHORED]->(sharePost)
		ON CREATE SET auth.created_at = $createdAt
		
		// 4. Mũi tên "Nội Dung": Nối thẳng về Bài Gốc (để load hình ảnh/video)
		MERGE (sharePost)-[rep:REPOSTS]->(original)
		ON CREATE SET rep.created_at = $createdAt

		// 5. Mũi tên "Lan Truyền": Nối về Bài Trung Gian (để ghi công Spreader)
		// Nếu parent == original, Neo4j sẽ dán thêm cạnh SHARED_VIA vào cùng 1 node, rất an toàn và hợp lý.
		MERGE (sharePost)-[via:SHARED_VIA]->(parent)
		ON CREATE SET via.created_at = $createdAt
	`

	params := map[string]any{
		"userID":         userID,
		"newSharePostID": newSharePostID,
		"originalPostID": originalPostID,
		"parentPostID":   parentPostID,
		"createdAt":      createdAt,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	return err
}

// GetSocialFeed (G1) - Lấy Top K bài viết từ Vòng tròn xã hội (Bạn bè, Group, Page)
// Áp dụng EdgeRank (Affinity) + Time Decay + Anti-Block + Xuyên thấu Share
func (r *GraphRepository) GetSocialFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error) {
	// Dùng AccessModeRead để báo cho Neo4j biết đây là tác vụ Read-Only (Giúp tối ưu routing trong Cluster)
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	cypher := `
		// 1. CHỌN MỤC TIÊU: Nhảy từ User sang các nguồn nội dung (Bạn bè, Group đã tham gia, Page đã Like)
		MATCH (u:User {user_id: $userID})-[:FRIEND|FOLLOWS|MEMBER_OF|LIKES]->(target)
		
		// 2. LẤY BÀI VIẾT: Tìm các bài Post được Đăng/Share bởi các mục tiêu này
		MATCH (target)-[:AUTHORED|PUBLISHED|POSTED_IN]->(p:Post)
		
		// 3. XUYÊN THẤU SHARE & TÌM TÁC GIẢ THẬT SỰ (True Author)
		// Dùng *0..1 để bắt cả bài gốc lẫn bài share.
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)<-[:AUTHORED|PUBLISHED|POSTED_IN]-(trueAuthor)
		
		// 4. BỘ LỌC AN TOÀN (FILTERING)
		// - Kiểm tra Block 2 chiều vô hướng (KHÔNG dùng mũi tên -> mà dùng -)
		// - Kiểm tra TTL (Hạn sử dụng của bài viết). Dùng coalesce để xử lý bài không có TTL.
		WHERE NOT (u)-[:BLOCK]-(trueAuthor)
		  AND coalesce(p.ttl, 9999999999999) > $now
		
		// 5. LẤY ĐIỂM THÂN THIẾT (AFFINITY SCORE)
		// Dùng OPTIONAL MATCH vì có thể User vừa kết bạn, chưa từng tương tác (chưa có cạnh INTERACTED_WITH)
		OPTIONAL MATCH (u)-[interact:INTERACTED_WITH]->(target)
		
		// 6. TÍNH TOÁN ĐIỂM SỐ CƠ BẢN (BASE SCORING)
		WITH p, 
		     coalesce(interact.affinity_score, 1.0) AS affinity, // Nếu chưa tương tác, mặc định điểm thân thiết là 1.0
		     p.created_at AS recency,
		     $now AS now
		
		// Công thức EdgeRank + Time Decay (Lực hấp dẫn của thời gian)
		// Đổi giây ra giờ (/ 3600.0) để đồ thị suy giảm mượt mà hơn
		WITH p, (affinity * 1.0) / (((now - recency) / 3600.0) + 2.0)^1.5 AS baseScore
		
		// 7. SẮP XẾP VÀ TRẢ VỀ
		ORDER BY baseScore DESC
		LIMIT $limitK
		
		RETURN p.post_id AS postID, baseScore AS score
	`

	params := map[string]any{
		"userID": userID,
		"now":    time.Now().Unix(),
		"limitK": limit,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var posts []graphEvent.ScoredPost
		for res.Next(ctx) {
			record := res.Record()
			postID := record.Values[0].(string)
			score := record.Values[1].(float64)

			posts = append(posts, graphEvent.ScoredPost{
				PostID: postID,
				Score:  score,
			})
		}
		return posts, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch G1 social feed: %w", err)
	}

	return result.([]graphEvent.ScoredPost), nil
}

// GetInterestFeed (G2) - Lấy Top K bài viết dựa trên Topic User quan tâm (Content-Based)
// Áp dụng Trộn điểm (Short-term + Long-term) + Trọng số AI + Time Decay + Anti-Block
func (r *GraphRepository) GetInterestFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	cypher := `
		// 1. TÌM SỞ THÍCH CỦA USER VÀ TRỘN ĐIỂM (BLENDING SCORE)
		MATCH (u:User {user_id: $userID})-[interest:INTERESTED_IN]->(t:Topic)
		
		// Trộn Sở thích Ngắn hạn (Ví dụ: 70% trọng số) và Dài hạn (30% trọng số)
		WITH u, t, 
		     (coalesce(interest.short_term_score, 0.0) * 0.7 + coalesce(interest.score, 0.0) * 0.3) AS topicWeight
		
		// 2. CẮT TỈA (PRUNING) - Tối ưu hiệu năng cực đại
		// Chỉ lấy Top 20 chủ đề mà User đang quan tâm nhất hiện tại để đi tìm bài viết
		ORDER BY topicWeight DESC
		LIMIT 20
		
		// 3. TÌM BÀI VIẾT TỪ TOP 20 TOPICS NÀY
		// ht.confidence_score là độ tự tin của con AI khi gắn tag cho bài viết
		MATCH (t)<-[ht:HAS_TOPIC]-(p:Post)
		
		// 4. XUYÊN THẤU SHARE VÀ TÌM TÁC GIẢ THẬT SỰ
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)<-[:AUTHORED|PUBLISHED|POSTED_IN]-(trueAuthor)
		
		// 5. BỘ LỌC AN TOÀN & CÁ NHÂN HÓA
		// - Chặn Block 2 chiều
		// - Bài viết còn hạn (TTL)
		// - KHÔNG lấy bài do chính User đang xem đăng (Self-exclusion)
		WHERE NOT (u)-[:BLOCK]-(trueAuthor)
		  AND coalesce(p.ttl, 9999999999999) > $now
		  AND elementId(u) <> elementId(trueAuthor)
		
		// 6. CỘNG DỒN ĐIỂM CHÉO (MULTI-TOPIC AGGREGATION)
		// Nếu 1 bài Post có nhiều Topic khớp với sở thích, điểm phải được cộng dồn (SUM)
		WITH p, 
		     $now AS now,
		     SUM(topicWeight * coalesce(ht.confidence_score, 0.5)) AS totalRelevance
		
		// 7. CÔNG THỨC TIME DECAY (Chuẩn hóa cùng Scale với G1)
		// Điểm = Độ liên quan / Lực hấp dẫn thời gian
		WITH p, totalRelevance / (((now - p.created_at) / 3600.0) + 2.0)^1.5 AS finalScore
		
		// 8. SẮP XẾP VÀ TRẢ VỀ
		ORDER BY finalScore DESC
		LIMIT $limitK
		
		RETURN p.post_id AS postID, finalScore AS score
	`

	params := map[string]any{
		"userID": userID,
		"now":    time.Now().Unix(),
		"limitK": limit,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var posts []graphEvent.ScoredPost
		for res.Next(ctx) {
			record := res.Record()
			postID := record.Values[0].(string)
			score := record.Values[1].(float64)

			posts = append(posts, graphEvent.ScoredPost{
				PostID: postID,
				Score:  score,
			})
		}
		return posts, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch G2 interest feed: %w", err)
	}

	return result.([]graphEvent.ScoredPost), nil
}

// GetSemanticDiscoveryFeed (G3) - Lấy Top K bài viết từ các Chủ đề Lân cận (Semantic Discovery)
// Phá vỡ Filter Bubble bằng cách nhảy qua cạnh RELATED_TO, kết hợp Time Decay và Anti-Block
func (r *GraphRepository) GetSemanticDiscoveryFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	cypher := `
		// ==========================================
		// BƯỚC 1: XÁC ĐỊNH "TÂM CHẤN" SỞ THÍCH CỦA USER
		// ==========================================
		MATCH (u:User {user_id: $userID})-[interest:INTERESTED_IN]->(t:Topic)
		WITH u, t, 
		     (coalesce(interest.short_term_score, 0.0) * 0.7 + coalesce(interest.score, 0.0) * 0.3) AS baseWeight
		ORDER BY baseWeight DESC
		LIMIT 15 // Chỉ lấy Top 15 chủ đề mà User đang thích nhất làm "Bàn đạp"

		// ==========================================
		// BƯỚC 2: RANDOM WALK - NHẢY SANG CÁC CHỦ ĐỀ HÀNG XÓM
		// ==========================================
		MATCH (t)-[rel:RELATED_TO]-(neighborTopic:Topic)
		
		// [TUYỆT CHIÊU 100% BEST PRACTICE]: ĐẢM BẢO SỰ "MỚI MẺ"
		// Bắt buộc User CHƯA TỪNG thích chủ đề hàng xóm này (Chống gợi ý trùng lặp với G2)
		// Chỉ lấy những họ hàng có độ tương đồng đủ cao (>= 0.8)
		WHERE NOT (u)-[:INTERESTED_IN]->(neighborTopic)
		  AND rel.similarity_score >= $similarityThreshold

		// Tính "Điểm Khám phá" (Discovery Weight) = Điểm chủ đề gốc x Độ tương đồng
		WITH u, neighborTopic, SUM(baseWeight * rel.similarity_score) AS discoveryWeight
		ORDER BY discoveryWeight DESC
		LIMIT 20 // Lọc ra Top 20 chủ đề "Hàng xóm" tiềm năng nhất

		// ==========================================
		// BƯỚC 3: TÌM BÀI VIẾT TỪ CHỦ ĐỀ KHÁM PHÁ VÀ XUYÊN THẤU SHARE
		// ==========================================
		MATCH (neighborTopic)<-[ht:HAS_TOPIC]-(p:Post)
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)<-[:AUTHORED|PUBLISHED|POSTED_IN]-(trueAuthor)

		// ==========================================
		// BƯỚC 4: BỘ LỌC AN TOÀN (FILTERING)
		// ==========================================
		WHERE NOT (u)-[:BLOCK]-(trueAuthor)
		  AND coalesce(p.ttl, 9999999999999) > $now
		  AND elementId(u) <> elementId(trueAuthor)

		// ==========================================
		// BƯỚC 5: CỘNG DỒN ĐIỂM CHÉO VÀ TIME DECAY
		// ==========================================
		WITH p, 
		     $now AS now, 
		     SUM(discoveryWeight * coalesce(ht.confidence_score, 0.5)) AS postRelevance
		
		WITH p, postRelevance / (((now - p.created_at) / 3600.0) + 2.0)^1.5 AS finalScore
		
		ORDER BY finalScore DESC
		LIMIT $limitK
		
		RETURN p.post_id AS postID, finalScore AS score
	`

	params := map[string]any{
		"userID":              userID,
		"now":                 time.Now().Unix(),
		"limitK":              limit,
		"similarityThreshold": 0.80, // Có thể chỉnh threshold này tùy theo mức độ "dũng cảm" của thuật toán
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var posts []graphEvent.ScoredPost
		for res.Next(ctx) {
			record := res.Record()
			postID := record.Values[0].(string)
			score := record.Values[1].(float64)

			posts = append(posts, graphEvent.ScoredPost{
				PostID: postID,
				Score:  score,
			})
		}
		return posts, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch G3 semantic discovery feed: %w", err)
	}

	return result.([]graphEvent.ScoredPost), nil
}

// GetGlobalTrendingFeed (G4) - Lấy Top K bài viết đang Trending toàn cầu
// Dựa trên Điểm Trending của Topic, Áp dụng Time Decay, Xuyên thấu Share và Anti-Block
func (r *GraphRepository) GetGlobalTrendingFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	cypher := `
		// ==========================================
		// BƯỚC 1: LẤY USER VÀ TÌM TOP TRENDING TOPICS
		// ==========================================
		MATCH (u:User {user_id: $userID})
		
		// Quét Index cực nhanh để lấy ra 20 Chủ đề đang Hot nhất toàn mạng lưới
		MATCH (t:Topic)
		WHERE coalesce(t.trending_score, 0.0) > 0
		WITH u, t
		ORDER BY t.trending_score DESC
		LIMIT 20 // Cắt tỉa (Pruning) để bảo vệ hiệu năng Database

		// ==========================================
		// BƯỚC 2: TÌM BÀI VIẾT TỪ CÁC CHỦ ĐỀ HOT NÀY
		// ==========================================
		MATCH (t)<-[ht:HAS_TOPIC]-(p:Post)
		
		// Xuyên thấu bài Share để tìm gốc rễ
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)<-[:AUTHORED|PUBLISHED|POSTED_IN]-(trueAuthor)

		// ==========================================
		// BƯỚC 3: KHIÊN BẢO VỆ TÍCH HỢP (FILTERING)
		// ==========================================
		WHERE NOT (u)-[:BLOCK]-(trueAuthor)
		  AND coalesce(p.ttl, 9999999999999) > $now
		  AND elementId(u) <> elementId(trueAuthor) // Không bơm bài của chính mình vào mục Trending

		// ==========================================
		// BƯỚC 4: TÍNH ĐIỂM TRENDING CỦA BÀI VIẾT
		// ==========================================
		// Nếu 1 bài viết gắn cả 2 hashtag đang top 1 và top 2 trending -> Điểm cộng dồn cực to
		WITH p, 
		     $now AS now,
		     SUM(t.trending_score * coalesce(ht.confidence_score, 0.5)) AS totalTrendWeight
		
		// Áp dụng Time Decay: Trend thì trend, nhưng bài từ 3 ngày trước thì phải nhường chỗ cho bài vừa đăng 1 tiếng trước
		WITH p, totalTrendWeight / (((now - p.created_at) / 3600.0) + 2.0)^1.5 AS finalScore
		
		// ==========================================
		// BƯỚC 5: SẮP XẾP VÀ TRẢ VỀ CHO GO
		// ==========================================
		ORDER BY finalScore DESC
		LIMIT $limitK
		
		RETURN p.post_id AS postID, finalScore AS score
	`

	params := map[string]any{
		"userID": userID,
		"now":    time.Now().Unix(),
		"limitK": limit,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var posts []graphEvent.ScoredPost
		for res.Next(ctx) {
			record := res.Record()
			postID := record.Values[0].(string)
			score := record.Values[1].(float64)

			posts = append(posts, graphEvent.ScoredPost{
				PostID: postID,
				Score:  score,
			})
		}
		return posts, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch G4 global trending feed: %w", err)
	}

	return result.([]graphEvent.ScoredPost), nil
}

// GetLocationFeed (G5) - Lấy Top K bài viết từ những người lạ ở gần khu vực User (Location-based)
// Áp dụng Cắt tỉa siêu đô thị (Megacity Pruning), Anti-Friend-Overlap, Xuyên thấu Share và Anti-Block
func (r *GraphRepository) GetLocationFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	cypher := `
		// ==========================================
		// BƯỚC 1: XÁC ĐỊNH "TỌA ĐỘ TÂM CHẤN" CỦA USER
		// ==========================================
		MATCH (u:User {user_id: $userID})-[:LOCATED_IN]->(loc:Location)

		// ==========================================
		// BƯỚC 2: KHÁM PHÁ CỘNG ĐỒNG ĐỊA PHƯƠNG & CẮT TỈA (PRUNING)
		// ==========================================
		// Quét ngược lại để tìm những thực thể (User, Page, Group) ở cùng Location
		MATCH (loc)<-[:LOCATED_IN]-(localSource)
		
		// [TUYỆT CHIÊU ANTI-OVERLAP]: Đảm bảo tính "Khám phá người lạ"
		// Bỏ qua chính bản thân User
		// Bỏ qua những người đã là Bạn Bè hoặc đang Follow (Vì luồng G1 đã lấy rồi, không lấy trùng lại ở đây)
		WHERE elementId(localSource) <> elementId(u)
		  AND NOT (u)-[:FRIEND|FOLLOWS]-(localSource)
		
		// [TUYỆT CHIÊU ANTI-EXPLOSION]: Chặn bùng nổ dữ liệu ở thành phố lớn
		WITH u, localSource
		LIMIT 300 // Chỉ lấy mẫu ngẫu nhiên 300 nguồn địa phương để đảm bảo truy vấn siêu tốc (O(1))

		// ==========================================
		// BƯỚC 3: LẤY BÀI VIẾT & XUYÊN THẤU SHARE
		// ==========================================
		MATCH (localSource)-[:AUTHORED|PUBLISHED|POSTED_IN]->(p:Post)
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)<-[:AUTHORED|PUBLISHED|POSTED_IN]-(trueAuthor)

		// ==========================================
		// BƯỚC 4: KHIÊN BẢO VỆ TÍCH HỢP (FILTERING)
		// ==========================================
		WHERE NOT (u)-[:BLOCK]-(trueAuthor)
		  AND coalesce(p.ttl, 9999999999999) > $now
		  AND elementId(u) <> elementId(trueAuthor) // Không bơm bài của chính mình

		// ==========================================
		// BƯỚC 5: CHẤM ĐIỂM (LOCAL RECENCY SCORING)
		// ==========================================
		// Đối với luồng địa phương, tính "Thời sự" (Mới đăng) là quan trọng nhất
		WITH p, $now AS now
		WITH p, 1.0 / (((now - p.created_at) / 3600.0) + 2.0)^1.5 AS finalScore

		// ==========================================
		// BƯỚC 6: SẮP XẾP VÀ TRẢ VỀ CHO GO
		// ==========================================
		ORDER BY finalScore DESC
		LIMIT $limitK

		RETURN p.post_id AS postID, finalScore AS score
	`

	params := map[string]any{
		"userID": userID,
		"now":    time.Now().Unix(),
		"limitK": limit,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var posts []graphEvent.ScoredPost
		for res.Next(ctx) {
			record := res.Record()
			postID := record.Values[0].(string)
			score := record.Values[1].(float64)

			posts = append(posts, graphEvent.ScoredPost{
				PostID: postID,
				Score:  score,
			})
		}
		return posts, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch G5 location feed: %w", err)
	}

	return result.([]graphEvent.ScoredPost), nil
}

// GetContactSyncFeed (G6) - Lấy Top K bài viết từ những người có liên kết trong Danh bạ
// Dựa trên Phone Hash (ẩn danh), Áp dụng Xuyên thấu Share, Anti-Block và Lọc trùng bạn bè
func (r *GraphRepository) GetContactSyncFeed(ctx context.Context, userID string, limit int) ([]graphEvent.ScoredPost, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	cypher := `
		// ==========================================
		// BƯỚC 1: XÁC ĐỊNH MẠNG LƯỚI DANH BẠ (CONTACT SYNC)
		// ==========================================
		MATCH (u:User {user_id: $userID})
		
		// Tìm những người có chung băm số điện thoại thông qua Node PhoneContact
		// Dấu '-' không mũi tên giúp quét cả 2 chiều: A lưu số B, hoặc B lưu số A đều dính.
		MATCH (u)-[:HAS_CONTACT]-(c:PhoneContact)-[:HAS_CONTACT]-(realLifeContact:User)

		// ==========================================
		// BƯỚC 2: CẮT TỈA (PRUNING) VÀ LỌC TRÙNG LẶP
		// ==========================================
		// Không lấy lại những người đã kết bạn / theo dõi (Vì G1 đã xử lý)
		// Không tự lấy chính mình
		WHERE elementId(u) <> elementId(realLifeContact)
		  AND NOT (u)-[:FRIEND|FOLLOWS]-(realLifeContact)
		
		// Cắt tỉa: Chỉ lấy ngẫu nhiên một nhóm danh bạ để đảm bảo truy vấn nhanh
		WITH u, realLifeContact
		LIMIT 50

		// ==========================================
		// BƯỚC 3: TÌM BÀI VIẾT VÀ XUYÊN THẤU SHARE
		// ==========================================
		MATCH (realLifeContact)-[:AUTHORED|PUBLISHED|POSTED_IN]->(p:Post)
		
		// Bắt cả bài gốc lẫn bài Share
		MATCH (p)-[:REPOSTS*0..1]->(basePost:Post)<-[:AUTHORED|PUBLISHED|POSTED_IN]-(trueAuthor)

		// ==========================================
		// BƯỚC 4: KHIÊN BẢO VỆ (FILTERING)
		// ==========================================
		// Chặn Block 2 chiều với tác giả thực sự
		WHERE NOT (u)-[:BLOCK]-(trueAuthor)
		  AND coalesce(p.ttl, 9999999999999) > $now
		  AND elementId(u) <> elementId(trueAuthor)

		// ==========================================
		// BƯỚC 5: CHẤM ĐIỂM (CONTACT SCORING)
		// ==========================================
		// Vì độ tin cậy ngoài đời đã rất cao, ta lấy điểm mặc định là 1.5 (Cao hơn điểm bạn bè chưa tương tác)
		// Chỉ phân rã theo thời gian để lấy bài mới
		WITH p, 
		     1.5 AS baseContactAffinity, 
		     $now AS now
		     
		WITH p, (baseContactAffinity) / (((now - p.created_at) / 3600.0) + 2.0)^1.5 AS finalScore

		// ==========================================
		// BƯỚC 6: SẮP XẾP VÀ TRẢ VỀ
		// ==========================================
		ORDER BY finalScore DESC
		LIMIT $limitK

		RETURN p.post_id AS postID, finalScore AS score
	`

	params := map[string]any{
		"userID": userID,
		"now":    time.Now().Unix(),
		"limitK": limit,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var posts []graphEvent.ScoredPost
		for res.Next(ctx) {
			record := res.Record()
			postID := record.Values[0].(string)
			score := record.Values[1].(float64)

			posts = append(posts, graphEvent.ScoredPost{
				PostID: postID,
				Score:  score,
			})
		}
		return posts, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch G6 contact sync feed: %w", err)
	}

	return result.([]graphEvent.ScoredPost), nil
}

// GetPeopleYouMayKnow (G7) - Trả về Top K người dùng tiềm năng để gợi ý kết bạn
// Dựa trên Thuật toán Friend-of-Friend (Triadic Closure), kết hợp điểm Danh bạ và Anti-Block
func (r *GraphRepository) GetPeopleYouMayKnow(ctx context.Context, userID string, limit int) ([]graphEvent.SuggestedUser, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	cypher := `
		// ==========================================
		// BƯỚC 1: QUÉT BẠN CỦA BẠN (FRIEND OF FRIEND - FOAF)
		// ==========================================
		MATCH (u:User {user_id: $userID})
		
		// Tìm những người là bạn của bạn bè (nhảy 2 bước vô hướng)
		MATCH (u)-[:FRIEND]-(mutualFriend:User)-[:FRIEND]-(candidate:User)

		// ==========================================
		// BƯỚC 2: BỘ LỌC AN TOÀN & CHỐNG TRÙNG LẶP
		// ==========================================
		// 1. Không gợi ý chính bản thân mình
		// 2. Không gợi ý những người ĐÃ LÀ BẠN hoặc ĐANG FOLLOW
		// 3. Khắc chế Block 2 chiều (Bạn block họ, hoặc họ block bạn đều không hiện)
		WHERE elementId(u) <> elementId(candidate)
		  AND NOT (u)-[:FRIEND|FOLLOWS]-(candidate)
		  AND NOT (u)-[:BLOCK]-(candidate)

		// ==========================================
		// BƯỚC 3: ĐẾM BẠN CHUNG (MUTUAL FRIENDS)
		// ==========================================
		// Dùng WITH + COUNT để nhóm ứng viên lại và đếm xem có bao nhiêu nhánh mutualFriend nối u và candidate
		WITH u, candidate, COUNT(mutualFriend) AS mutualFriendsCount

		// ==========================================
		// BƯỚC 4: TĂNG CƯỜNG ĐIỂM SỐ BẰNG DANH BẠ (CONTACT BOOST)
		// ==========================================
		// Kiểm tra xem ứng viên này có nằm trong danh bạ (Phone Hash) của User hay không
		OPTIONAL MATCH (u)-[:HAS_CONTACT]-(c:PhoneContact)-[:HAS_CONTACT]-(candidate)
		
		// count(c) > 0 trả về true nếu có trong danh bạ. Nếu true -> thưởng 5.0 điểm, false -> 0 điểm.
		WITH candidate, 
		     mutualFriendsCount, 
		     CASE WHEN count(c) > 0 THEN 5.0 ELSE 0.0 END AS contactBoostScore

		// ==========================================
		// BƯỚC 5: TÍNH TỔNG ĐIỂM VÀ SẮP XẾP
		// ==========================================
		// Điểm = (Số bạn chung x 1.0) + Điểm thưởng danh bạ
		WITH candidate, mutualFriendsCount, (mutualFriendsCount * 1.0) + contactBoostScore AS pymkScore
		
		ORDER BY pymkScore DESC, mutualFriendsCount DESC
		LIMIT $limitK

		RETURN candidate.user_id AS candidateID, mutualFriendsCount, pymkScore
	`

	params := map[string]any{
		"userID": userID,
		"limitK": limit,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		var suggestions []graphEvent.SuggestedUser
		for res.Next(ctx) {
			record := res.Record()
			candidateID := record.Values[0].(string)
			mutualFriends := int(record.Values[1].(int64))
			score := record.Values[2].(float64)

			suggestions = append(suggestions, graphEvent.SuggestedUser{
				UserID:             candidateID,
				MutualFriendsCount: mutualFriends,
				Score:              score,
			})
		}
		return suggestions, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch G7 pymk feed: %w", err)
	}

	return result.([]graphEvent.SuggestedUser), nil
}

func (r *GraphRepository) GetPersonalizedNewsFeed(ctx context.Context, userID string) (*graphEvent.NewsFeedResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second) // Giới hạn API 2 giây
	defer cancel()

	eg, egCtx := errgroup.WithContext(timeoutCtx)

	// Các biến chứa dữ liệu từ Tầng 1
	var g1, g2, g3, g4, g5, g6 []graphEvent.ScoredPost
	var g7_pymk []graphEvent.SuggestedUser

	// =========================================================================
	// TẦNG 1: RECALL (TRUY XUẤT THÔ SONG SONG)
	// =========================================================================
	eg.Go(func() error { g1, _ = r.GetSocialFeed(egCtx, userID, 40); return nil })
	eg.Go(func() error { g2, _ = r.GetInterestFeed(egCtx, userID, 30); return nil })
	eg.Go(func() error { g3, _ = r.GetSemanticDiscoveryFeed(egCtx, userID, 15); return nil })
	eg.Go(func() error { g4, _ = r.GetGlobalTrendingFeed(egCtx, userID, 10); return nil })
	eg.Go(func() error { g5, _ = r.GetLocationFeed(egCtx, userID, 10); return nil })
	eg.Go(func() error { g6, _ = r.GetContactSyncFeed(egCtx, userID, 5); return nil })
	eg.Go(func() error { g7_pymk, _ = r.GetPeopleYouMayKnow(egCtx, userID, 5); return nil })

	// Lấy song song danh sách bài ĐÃ XEM từ Redis (để phục vụ Tầng 2)
	var viewedMap map[string]bool
	eg.Go(func() error {
		viewedMap, _ = r.redisRepo.GetViewedPosts(egCtx, userID)
		if viewedMap == nil {
			viewedMap = make(map[string]bool)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("feed recall failed: %w", err)
	}

	// GOM RỔ VÀ KHỬ TRÙNG LẶP (Deduplication)
	// Trộn tất cả G1 -> G6 vào 1 Map. Nếu trùng ID, CỘNG DỒN ĐIỂM.
	rawPool := make(map[string]float64)
	mergeFn := func(posts []graphEvent.ScoredPost) {
		for _, p := range posts {
			rawPool[p.PostID] += p.Score // Điểm gộp (W_a*A + W_i*I + W_t*T)
		}
	}
	mergeFn(g1)
	mergeFn(g2)
	mergeFn(g3)
	mergeFn(g4)
	mergeFn(g5)
	mergeFn(g6)

	// Thu thập mảng ID để đi lấy Metadata
	var postIDs []string
	for id := range rawPool {
		postIDs = append(postIDs, id)
	}

	// =========================================================================
	// TẦNG 2: FILTERING (LỌC RÁC & CHỐNG SPAM)
	// =========================================================================
	// Lấy Metadata (Author, Topic, Risk) của 110 bài trong 1 hit (Bulk Query)
	// metaMap, err := u.metaRepo.GetPostsMetadata(ctx, postIDs)
	// if err != nil {
	// 	return nil, err
	// }

	type CandidatePost struct {
		PostID string
		Score  float64
		Meta   graphEvent.PostMeta
	}
	var filteredCandidates []CandidatePost

	for id, baseScore := range rawPool {
		// 1. Lọc bài đã xem trong Redis
		if viewedMap[id] {
			continue
		}

		// meta, exists := metaMap[id]
		// if !exists {
		// 	continue
		// }

		// 2. Lọc Tín nhiệm (Chống Bot/Spam) - VD: Ngưỡng rủi ro = 0.8
		// if meta.RiskScore > 0.8 {
		// 	continue
		// }

		// Đưa vào danh sách ứng viên sạch
		filteredCandidates = append(filteredCandidates, CandidatePost{
			PostID: id,
			Score:  baseScore,
			Meta:   graphEvent.PostMeta{},
		})
	}

	// =========================================================================
	// TẦNG 3: RE-RANKING (ĐA DẠNG HÓA VÀ CHẤM ĐIỂM CUỐI)
	// =========================================================================

	// 1. Áp dụng công thức Hacker News: Score = AggregatedBase / (TimeDelta^1.5)
	now := time.Now().Unix()
	for i := range filteredCandidates {
		hoursOld := float64(now-filteredCandidates[i].Meta.CreatedAt) / 3600.0
		if hoursOld < 0 {
			hoursOld = 0
		}

		// Trọng lực (Gravity) = 1.5
		timePenalty := math.Pow(hoursOld+2.0, 1.5)
		filteredCandidates[i].Score = filteredCandidates[i].Score / timePenalty
	}

	// 2. Thuật toán Greedy - Intra-List Diversity (Phạt điểm trùng lặp)
	var finalFeed []CandidatePost

	// Lặp cho đến khi lấy đủ 20 bài HOẶC hết ứng viên
	for len(finalFeed) < 20 && len(filteredCandidates) > 0 {
		// Sắp xếp giảm dần theo điểm ở mỗi vòng lặp (vì điểm có thể bị phạt tụt xuống)
		sort.Slice(filteredCandidates, func(i, j int) bool {
			return filteredCandidates[i].Score > filteredCandidates[j].Score
		})

		topCandidate := filteredCandidates[0]

		// Đếm số lượng Tác giả/Topic liên tiếp ở cuối danh sách FinalFeed
		consecutiveAuthor := 0
		consecutiveTopic := 0

		if len(finalFeed) > 0 {
			lastPost := finalFeed[len(finalFeed)-1]
			if lastPost.Meta.AuthorID == topCandidate.Meta.AuthorID {
				consecutiveAuthor = 1
			}
			if lastPost.Meta.Topic == topCandidate.Meta.Topic {
				consecutiveTopic = 1
			}

			if len(finalFeed) > 1 {
				secondLastPost := finalFeed[len(finalFeed)-2]
				if consecutiveAuthor == 1 && secondLastPost.Meta.AuthorID == topCandidate.Meta.AuthorID {
					consecutiveAuthor = 2
				}
				if consecutiveTopic == 1 && secondLastPost.Meta.Topic == topCandidate.Meta.Topic {
					if len(finalFeed) > 2 && finalFeed[len(finalFeed)-3].Meta.Topic == topCandidate.Meta.Topic {
						consecutiveTopic = 3
					} else {
						consecutiveTopic = 2
					}
				}
			}
		}

		// ÁP DỤNG HÌNH PHẠT (PENALTY)
		isPenalized := false

		// Phạt 30% nếu trùng 2 tác giả liên tiếp
		if consecutiveAuthor >= 1 {
			filteredCandidates[0].Score *= 0.70
			isPenalized = true
		}

		// Phạt 50% nếu trùng 3 Topic liên tiếp
		if consecutiveTopic >= 2 {
			filteredCandidates[0].Score *= 0.50
			isPenalized = true
		}

		// Nếu bị phạt, không đẩy vào FinalFeed ngay, mà để nó ở lại Mảng Candidates
		// Vòng lặp tiếp theo sẽ sort lại để xem nó có còn đứng Top 1 nổi không.
		if isPenalized {
			// Chỉ phạt 1 lần cho vị trí này để tránh lặp vô hạn (Infinite loop trick)
			// Trong thực tế, bạn có thể pop nó ra, tìm bài top 2 đưa lên, rồi push nó vào lại.
			// Ở đây ta dùng chiến thuật: Nếu bị phạt, hoán đổi với bài thứ 2 (Swap)
			if len(filteredCandidates) > 1 {
				filteredCandidates[0], filteredCandidates[1] = filteredCandidates[1], filteredCandidates[0]
			}
			continue
		}

		// Nếu không vi phạm (hoặc đã chấp nhận hoán đổi), đẩy vào Feed cuối cùng
		finalFeed = append(finalFeed, topCandidate)
		filteredCandidates = filteredCandidates[1:] // Xóa bài đã chọn khỏi mảng ứng viên
	}

	// =========================================================================
	// RÁP GIAO DIỆN (WIDGET INJECTION)
	// =========================================================================
	var responseItems []graphEvent.FeedItem

	for i, p := range finalFeed {
		responseItems = append(responseItems, graphEvent.FeedItem{
			Type: "post",
			Data: p.PostID,
		})

		// Nhét Widget PYMK (G7) vào ngay sau bài viết thứ 3
		if i == 2 && len(g7_pymk) > 0 {
			responseItems = append(responseItems, graphEvent.FeedItem{
				Type: "pymk_widget",
				Data: g7_pymk,
			})
		}
	}

	return &graphEvent.NewsFeedResponse{
		Items:      responseItems,
		NextCursor: "encoded_timestamp_for_pagination",
	}, nil
}
