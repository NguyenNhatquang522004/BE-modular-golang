package community

import (
	"context"
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleCommunity struct {
	// Dependency Injection (UseCases, Repositories...)
	client *mongo.Database
}

func NewModuleCommunity(client *mongo.Database) *ModuleCommunity {
	m := &ModuleCommunity{
		client: client,
	}
	err := m.InitMongo(client)
	if err != nil {
		panic("Failed to init MongoDB for Community Module: " + err.Error())
	}
	return m
}

// =============================================================================
// 1. MONGODB INITIALIZATION
// =============================================================================

func (m *ModuleCommunity) InitMongo(db *mongo.Database) error {
	ctx := context.Background()

	// 1. Collection: GROUPS
	if err := m.initGroupIndexes(ctx, db); err != nil {
		return err
	}

	// 2. Collection: GROUP MEMBERS
	if err := m.initGroupMemberIndexes(ctx, db); err != nil {
		return err
	}

	// 3. Collection: GROUP EVENTS
	if err := m.initGroupEventIndexes(ctx, db); err != nil {
		return err
	}

	// 4. Collection: GROUP FILES
	if err := m.initGroupFileIndexes(ctx, db); err != nil {
		return err
	}

	// 5. Collection: GROUP JOIN QUESTIONS
	if err := m.initGroupQuestionIndexes(ctx, db); err != nil {
		return err
	}

	log.Println(">>> Community Module: MongoDB Indexes Initialized Successfully")
	return nil
}

// --- Helper: Groups ---
func (m *ModuleCommunity) initGroupIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.Group{}.CollectionName())

	models := []mongo.IndexModel{
		// A. TEXT SEARCH: Tìm nhóm theo tên, mô tả, tags (Full-Text Search cơ bản trên Mongo)
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "description", Value: "text"},
				{Key: "tags", Value: "text"},
			},
			Options: options.Index().SetWeights(bson.M{
				"name": 10, // Ưu tiên tên nhóm
				"tags": 5,  // Sau đó đến tags
			}),
		},
		// B. UNIQUE SLUG: URL thân thiện (facebook.com/groups/hoi-yeu-meo)
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		// C. MY GROUPS: Lọc nhóm tôi đã tạo
		{
			Keys: bson.D{{Key: "creator_id", Value: 1}},
		},
		// D. DISCOVERY: Gợi ý nhóm phổ biến (Sort theo số lượng thành viên)
		{
			Keys: bson.D{{Key: "stats.member_count", Value: -1}},
		},
		// E. CATEGORY FILTER: Lọc nhóm theo danh mục và quyền riêng tư
		{
			Keys: bson.D{
				{Key: "category_id", Value: 1},
				{Key: "privacy", Value: 1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for Groups: %w", err)
	}
	return nil
}

// --- Helper: Group Members ---
func (m *ModuleCommunity) initGroupMemberIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.GroupMember{}.CollectionName())

	models := []mongo.IndexModel{
		// A. UNIQUE JOIN: Đảm bảo 1 user không join 1 nhóm 2 lần
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		// B. MY JOINED GROUPS: Danh sách nhóm tôi đã tham gia (Mới nhất lên đầu)
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "joined_at", Value: -1},
			},
		},
		// C. MEMBER LIST & FILTER: Lọc thành viên trong nhóm (VD: Tìm admin, tìm người bị ban)
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "role", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		// D. INACTIVE MEMBERS: Tìm thành viên tàu ngầm (ít tương tác) để lọc mem
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "last_active_at", Value: 1}, // Sort tăng dần (cũ nhất lên đầu)
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for GroupMembers: %w", err)
	}
	return nil
}

// --- Helper: Group Events ---
func (m *ModuleCommunity) initGroupEventIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.GroupEvent{}.CollectionName())

	models := []mongo.IndexModel{
		// A. GROUP CALENDAR: Lấy danh sách sự kiện của nhóm (Sắp diễn ra)
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "start_time", Value: 1},
			},
		},
		// B. GEO-SPATIAL SEARCH: Tìm sự kiện "Quanh đây" (Near Me)
		// Yêu cầu field location.coordinates phải là GeoJSON [lon, lat]
		{
			Keys: bson.D{{Key: "location.coordinates", Value: "2dsphere"}},
		},
		// C. UPCOMING EVENTS: Lọc sự kiện sắp tới trên toàn hệ thống (Global Discovery)
		{
			Keys: bson.D{{Key: "start_time", Value: 1}},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for GroupEvents: %w", err)
	}
	return nil
}

// --- Helper: Group Files ---
func (m *ModuleCommunity) initGroupFileIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.GroupFile{}.CollectionName())

	models := []mongo.IndexModel{
		// A. FILE LIST: Danh sách file trong nhóm (Mới nhất lên đầu)
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		// B. TOP DOWNLOADS: File tải nhiều nhất
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "download_count", Value: -1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for GroupFiles: %w", err)
	}
	return nil
}

// --- Helper: Group Join Questions ---
func (m *ModuleCommunity) initGroupQuestionIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.GroupJoinQuestion{}.CollectionName())

	models := []mongo.IndexModel{
		// A. QUESTION LIST: Lấy câu hỏi kiểm duyệt của nhóm (Theo thứ tự hiển thị)
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "order", Value: 1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for GroupJoinQuestions: %w", err)
	}
	return nil
}

// =============================================================================
// 2. ELASTICSEARCH INITIALIZATION
// =============================================================================

// func (m *ModuleCommunity) InitElastic(client *elasticsearch.Client) error {
// 	ctx := context.Background()
// 	indexName := entity.SearchGroup{}.IndexName() // "search_groups"

// 	// 1. Check Index Exists
// 	reqExists := esapi.IndicesExistsRequest{
// 		Index: []string{indexName},
// 	}
// 	resExists, err := reqExists.Do(ctx, client)
// 	if err != nil {
// 		return fmt.Errorf("check index exists error: %w", err)
// 	}
// 	defer resExists.Body.Close()

// 	if resExists.StatusCode == 200 {
// 		return nil // Index already exists
// 	}

// 	// 2. Define Mapping
// 	// Chú ý: Cấu hình Analyzer Tiếng Việt (vietnamese_folding)
// 	// Chú ý: Field location dùng type "geo_point" cho Geo-Search
// 	mapping := `{
// 		"settings": {
// 			"number_of_shards": 1,
// 			"number_of_replicas": 0,
// 			"analysis": {
// 				"analyzer": {
// 					"vietnamese_folding": {
// 						"tokenizer": "standard",
// 						"filter": ["lowercase", "asciifolding"]
// 					}
// 				}
// 			}
// 		},
// 		"mappings": {
// 			"properties": {
// 				"id": { "type": "keyword" },

// 				"name": {
// 					"type": "text",
// 					"analyzer": "vietnamese_folding",
// 					"search_analyzer": "vietnamese_folding",
// 					"boost": 2.0
// 				},
// 				"description": {
// 					"type": "text",
// 					"analyzer": "vietnamese_folding",
// 					"search_analyzer": "vietnamese_folding"
// 				},

// 				"tags": { "type": "keyword" },
// 				"privacy": { "type": "keyword" },

// 				"member_count": { "type": "integer" },

// 				"location": { "type": "geo_point" }
// 			}
// 		}
// 	}`

// 	// 3. Create Index
// 	reqCreate := esapi.IndicesCreateRequest{
// 		Index: indexName,
// 		Body:  strings.NewReader(mapping),
// 	}

// 	resCreate, err := reqCreate.Do(ctx, client)
// 	if err != nil {
// 		return fmt.Errorf("create index error: %w", err)
// 	}
// 	defer resCreate.Body.Close()

// 	if resCreate.IsError() {
// 		return fmt.Errorf("create index failed: %s", resCreate.String())
// 	}

// 	log.Printf(">>> Elastic Index [%s] initialized successfully\n", indexName)
// 	return nil
// }
