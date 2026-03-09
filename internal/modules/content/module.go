package content

import (
	"context"
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleContent struct {
	clientmongodb   *mongo.Database
	clientcassandra *gocql.Session
}

func NewModuleContent(mongoDB *mongo.Database, cassandraSession *gocql.Session) *ModuleContent {
	module := &ModuleContent{
		clientmongodb:   mongoDB,
		clientcassandra: cassandraSession,
	}

	// 1. Khởi tạo MongoDB Indexes (Nếu có lỗi sẽ panic để dev fix ngay)
	if err := module.InitMongo(mongoDB); err != nil {
		log.Panicf("Failed to initialize MongoDB indexes for Content Module: %v", err)
	}

	// 2. Khởi tạo Cassandra Schemas (Nếu có lỗi sẽ panic để dev fix ngay)
	if err := module.initCassandraSchemas(cassandraSession); err != nil {
		log.Panicf("Failed to initialize Cassandra schemas for Content Module: %v", err)
	}
	return module
}
func (m *ModuleContent) RegisterRoute(r *gin.RouterGroup) {

}
func (m *ModuleContent) initCassandraSchemas(session *gocql.Session) error {
	ctx := context.Background() // Context giữ chỗ (gocql thường dùng context trong query)
	_ = ctx

	// 1. INIT TABLE: POST INSIGHT (High Write Throughput)
	// Gọi hàm EnsureTableExists từ Entity PostInsight
	if err := (&entity.PostInsight{}).EnsureTableExists(session); err != nil {
		return fmt.Errorf("cassandra init failed for PostInsight: %w", err)
	}

	// log.Println(">>> Cassandra Tables Initialized")
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
