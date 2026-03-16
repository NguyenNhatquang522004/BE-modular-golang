package neo4j

import (
	"context"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/sharedEnums"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/graph/domain/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepository struct {
	// Có thể thêm các trường như Neo4j Driver, Logger, Config nếu cần
	driver neo4j.DriverWithContext
	aiRepo IRepositoryShare.IAI
	cfg    *configs.Config
}

func NewGraphRepository(driver neo4j.DriverWithContext, aiRepo IRepositoryShare.IAI, cfg *configs.Config) *GraphRepository {
	return &GraphRepository{
		driver: driver,
		aiRepo: aiRepo,
		cfg:    cfg,
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
		SET rel.short_term_score = totalTopicScore,
		    rel.last_calculated_at = $now
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
