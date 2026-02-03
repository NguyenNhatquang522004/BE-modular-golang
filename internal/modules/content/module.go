package content

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleContent struct {
}

func (m *ModuleContent) RegisterRoute(r *gin.RouterGroup) {

}

// InitElastic: Khởi tạo Index và Mapping cho Search Module
func (m *ModuleContent) InitElastic(client *elasticsearch.Client) error {
	ctx := context.Background()
	indexName := entity.SearchPost{}.IndexName()

	// 1. Kiểm tra xem Index đã tồn tại chưa
	// Dùng esapi để gọi API check exists
	reqExists := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}
	resExists, err := reqExists.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("check index exists error: %w", err)
	}
	defer resExists.Body.Close()

	// Nếu Index đã tồn tại (StatusCode 200) -> Bỏ qua (hoặc xử lý migration nếu cần)
	if resExists.StatusCode == 200 {
		return nil
	}

	// 2. Định nghĩa Mapping & Settings (JSON)
	// Đây là phần QUAN TRỌNG NHẤT
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"vietnamese_folding": {
						"tokenizer": "standard",
						"filter": ["lowercase", "asciifolding"] 
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				
				"content": { 
					"type": "text",
					"analyzer": "vietnamese_folding",
					"search_analyzer": "vietnamese_folding"
				},
				
				"hashtags": { "type": "keyword" },
				"author_id": { "type": "keyword" },
				"group_id": { "type": "keyword" },
				"page_id": { "type": "keyword" },
				
				"media_types": { "type": "keyword" },
				"privacy": { "type": "keyword" },
				
				"created_at": { "type": "date" },
				
				"likes_count": { "type": "integer" },
				"comments_count": { "type": "integer" }
			}
		}
	}`

	// 3. Tạo Index
	reqCreate := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	resCreate, err := reqCreate.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("create index error: %w", err)
	}
	defer resCreate.Body.Close()

	if resCreate.IsError() {
		return fmt.Errorf("create index failed: %s", resCreate.String())
	}

	fmt.Printf(">>> Elastic Index [%s] initialized successfully with Vietnamese Analyzer\n", indexName)
	return nil
}

// InitMongo: Khởi tạo toàn bộ Collection và Index cho Content Module
func (m *ModuleContent) InitMongo(db *mongo.Database) error {
	ctx := context.Background()

	// =========================================================================
	// 1. INIT COLLECTION: POST
	// =========================================================================
	if err := m.initPostIndexes(ctx, db); err != nil {
		return err
	}

	// =========================================================================
	// 2. INIT COLLECTION: POST MEDIA (1-1 with Post)
	// =========================================================================
	if err := m.initPostMediaIndexes(ctx, db); err != nil {
		return err
	}

	// =========================================================================
	// 3. INIT COLLECTION: POST EXTENSION (1-1 with Post)
	// =========================================================================
	if err := m.initPostExtensionIndexes(ctx, db); err != nil {
		return err
	}

	// =========================================================================
	// 4. INIT COLLECTION: POST SETTING (1-1 with Post + Scheduler)
	// =========================================================================
	if err := m.initPostSettingIndexes(ctx, db); err != nil {
		return err
	}

	// =========================================================================
	// 5. INIT COLLECTION: EDIT LOG (Audit)
	// =========================================================================
	if err := m.initEditLogIndexes(ctx, db); err != nil {
		return err
	}

	log.Println(">>> Content Module: MongoDB Indexes Initialized Successfully")
	return nil
}

// -----------------------------------------------------------------------------
// HELPER FUNCTIONS (Clean Code)
// -----------------------------------------------------------------------------

func (m *ModuleContent) initPostIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.Post{}.CollectionNamePost())

	models := []mongo.IndexModel{
		// 1. UNIQUE SLUG: Bắt buộc để làm SEO link (VD: /post/chao-ngay-moi)
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		// 2. USER FEED: Lấy danh sách bài viết của 1 User (Sắp xếp mới nhất)
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		// 3. NEWS FEED FILTER: Lọc bài viết Public + Active để hiển thị ra Feed chung
		{
			Keys: bson.D{
				{Key: "status", Value: 1},        // Active
				{Key: "privacy.scope", Value: 1}, // Public
				{Key: "created_at", Value: -1},   // Mới nhất
			},
		},
		// 4. HASHTAG SEARCH: Tìm bài viết theo hashtag (Multikey Index)
		{
			Keys: bson.D{{Key: "hashtags", Value: 1}},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for Post: %w", err)
	}
	return nil
}

func (m *ModuleContent) initPostMediaIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.PostMedia{}.CollectionNamePostMedia())

	models := []mongo.IndexModel{
		// 1. ENFORCE 1-1 RELATIONSHIP: 1 Post chỉ có 1 Media Doc
		{
			Keys:    bson.D{{Key: "post_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for PostMedia: %w", err)
	}
	return nil
}

func (m *ModuleContent) initPostExtensionIndexes(ctx context.Context, db *mongo.Database) error {
	// Lưu ý: Sửa lại domain PostExtension để hàm CollectionName là Public
	coll := db.Collection(entity.PostExtension{}.Collectionnamepostextension())

	models := []mongo.IndexModel{
		// 1. ENFORCE 1-1 RELATIONSHIP
		{
			Keys:    bson.D{{Key: "post_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		// 2. SHARE TRACKING: Tìm tất cả bài viết đã share từ bài gốc X
		// (Hỗ trợ tính năng: "Xem ai đã share bài viết này")
		{
			Keys: bson.D{
				{Key: "share_data.original_post_id", Value: 1},
				{Key: "id", Value: -1}, // Sort mới nhất
			},
			// Partial Index: Chỉ index những dòng có share_data (tiết kiệm RAM)
			Options: options.Index().SetPartialFilterExpression(bson.M{
				"share_data": bson.M{"$exists": true},
			}),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for PostExtension: %w", err)
	}
	return nil
}

func (m *ModuleContent) initPostSettingIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.PostSetting{}.CollectionNamePostsetting())

	models := []mongo.IndexModel{
		// 1. ENFORCE 1-1 RELATIONSHIP
		{
			Keys:    bson.D{{Key: "post_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		// 2. SCHEDULER WORKER: Quét bài viết cần đăng
		// Query: Find({ "schedule.is_scheduled": true, "schedule.publish_time": { $lte: now } })
		{
			Keys: bson.D{
				{Key: "schedule.is_scheduled", Value: 1},
				{Key: "schedule.publish_time", Value: 1},
			},
			// Chỉ index những bài có đặt lịch
			Options: options.Index().SetPartialFilterExpression(bson.M{
				"schedule.is_scheduled": true,
			}),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for PostSetting: %w", err)
	}
	return nil
}

func (m *ModuleContent) initEditLogIndexes(ctx context.Context, db *mongo.Database) error {
	// Lưu ý: Sửa lại domain PostEntityEditLog để hàm CollectionName là Public
	coll := db.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())

	models := []mongo.IndexModel{
		// 1. VIEW HISTORY: Lấy lịch sử sửa đổi của 1 Post/Comment
		// Query: Find({ target_id: "...", target_collection: "posts" }).Sort({ version: -1 })
		{
			Keys: bson.D{
				{Key: "target_id", Value: 1},
				{Key: "target_collection", Value: 1},
				{Key: "version", Value: -1},
			},
		},
		// 2. EDITOR TRACKING: Xem user X đã sửa những bài nào (Audit)
		{
			Keys: bson.D{
				{Key: "editor_id", Value: 1},
				{Key: "edited_at", Value: -1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for PostEntityEditLog: %w", err)
	}
	return nil
}
