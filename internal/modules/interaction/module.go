package interaction

import (
	"context"
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleInteraction struct {
	// Các dependency khác nếu có (UseCases, Repos...)
	db      *mongo.Database
	session *gocql.Session
}

func NewModuleInteraction(db *mongo.Database, session *gocql.Session) *ModuleInteraction {
	m := &ModuleInteraction{
		db:      db,
		session: session,
	}
	err := m.InitMongo(db)
	if err != nil {
		log.Fatalf(">>> Interaction Module: Failed to initialize MongoDB indexes: %v", err)
	}

	err = m.InitCassandra(session)
	if err != nil {
		log.Fatalf(">>> Interaction Module: Failed to initialize Cassandra tables: %v", err)
	}
	log.Println(">>> Interaction Module: Initialization Completed Successfully")
	return m
}

// =============================================================================
// 1. MONGODB INITIALIZATION
// =============================================================================

func (m *ModuleInteraction) InitMongo(db *mongo.Database) error {
	ctx := context.Background()

	// 1. Collection: COMMENT
	if err := m.initCommentIndexes(ctx, db); err != nil {
		return err
	}

	// 2. Collection: COMMENT EDIT LOG (Audit)
	if err := m.initCommentEditLogIndexes(ctx, db); err != nil {
		return err
	}

	// 3. Collection: USER SAVED ITEM
	if err := m.initUserSavedItemIndexes(ctx, db); err != nil {
		return err
	}

	log.Println(">>> Interaction Module: MongoDB Indexes Initialized Successfully")
	return nil
}

// --- Helper: Comment ---
func (m *ModuleInteraction) initCommentIndexes(ctx context.Context, db *mongo.Database) error {
	// Lưu ý: Sửa receiver trong domain.Comment để hàm CollectionName public
	// Giả định tên collection là "Comment" dựa trên file comments.go
	coll := db.Collection(entity.Comment{}.CollectionnamComment())

	models := []mongo.IndexModel{
		// A. LOAD ROOT COMMENTS: Lấy comment của 1 bài viết (Sort cũ nhất -> mới nhất hoặc ngược lại)
		// Query: Find({ post_id: "...", parent_comment_id: null }).Sort({ created_at: 1 })
		{
			Keys: bson.D{
				{Key: "post_id", Value: 1},
				{Key: "parent_comment_id", Value: 1}, // Cần field này để phân biệt root comment và reply
				{Key: "created_at", Value: 1},
			},
		},
		// B. LOAD REPLIES: Lấy các câu trả lời của 1 comment cha
		// Query: Find({ parent_comment_id: "..." }).Sort({ created_at: 1 })
		{
			Keys: bson.D{
				{Key: "parent_comment_id", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
		// C. USER MENTIONS: Tìm các comment mà user X được tag
		// Query: Find({ mentions: "user_uuid" })
		{
			Keys: bson.D{{Key: "mentions", Value: 1}},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for Comment: %w", err)
	}
	return nil
}

// --- Helper: Comment Edit Log ---
func (m *ModuleInteraction) initCommentEditLogIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.EntityEditLog{}.CollectionName())

	models := []mongo.IndexModel{
		// A. VIEW HISTORY: Xem lịch sử sửa của 1 Comment/Post
		// Query: Find({ target_id: "...", target_collection: "..." }).Sort({ version: -1 })
		{
			Keys: bson.D{
				{Key: "target_id", Value: 1},
				{Key: "target_collection", Value: 1},
				{Key: "version", Value: -1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for CommentEntityEditLog: %w", err)
	}
	return nil
}

// --- Helper: User Saved Item ---
func (m *ModuleInteraction) initUserSavedItemIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())

	models := []mongo.IndexModel{
		// A. PREVENT DUPLICATES (Quan trọng nhất): 1 User không lưu 1 bài 2 lần
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "target_id", Value: 1},
				{Key: "target_type", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		// B. LIST SAVED ITEMS: Xem danh sách đã lưu của User X, filter theo Collection
		// Query: Find({ user_id: "...", collection_name: "..." }).Sort({ created_at: -1 })
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "collection_name", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for UserSavedItem: %w", err)
	}
	return nil
}

// =============================================================================
// 2. CASSANDRA INITIALIZATION
// =============================================================================

func (m *ModuleInteraction) InitCassandra(session *gocql.Session) error {
	// 1. Table: EntityReaction (Lưu lượt thả tim)
	if err := m.initEntityReactionTable(session); err != nil {
		return err
	}

	// 2. Table: UserReactionHistory (Nhật ký hoạt động)
	if err := m.initUserReactionHistoryTable(session); err != nil {
		return err
	}

	log.Println(">>> Interaction Module: Cassandra Tables Initialized Successfully")
	return nil
}

// --- Helper: Entity Reaction ---
func (m *ModuleInteraction) initEntityReactionTable(session *gocql.Session) error {
	// Lưu ý: TargetID là string (từ Mongo), UserID là UUID (từ Postgres)
	// PARTITION KEY: target_id (Để query tất cả like của 1 bài viết)
	// CLUSTERING KEY: user_id (Để check xem user X có like bài này chưa)
	query := `
		CREATE TABLE IF NOT EXISTS entity_reactions (
			target_id text,
			user_id uuid,
			target_type int,
			reaction_code int,
			created_at timestamp,
			PRIMARY KEY (target_id, user_id)
		);
	`
	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create table entity_reactions: %w", err)
	}
	return nil
}

// --- Helper: User Reaction History ---
func (m *ModuleInteraction) initUserReactionHistoryTable(session *gocql.Session) error {
	// PARTITION KEY: user_id (Gom lịch sử của 1 người vào 1 chỗ)
	// CLUSTERING KEY: created_at DESC (Mới nhất hiển thị trước)
	query := `
		CREATE TABLE IF NOT EXISTS user_reaction_history (
			user_id uuid,
			created_at timestamp,
			target_id text,
			target_type int,
			reaction_code int,
			PRIMARY KEY (user_id, created_at)
		) WITH CLUSTERING ORDER BY (created_at DESC);
	`
	// Lưu ý: Enum trong Go thường lưu là int trong DB.
	// target_type và reaction_code mình để int cho tối ưu.

	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create table user_reaction_history: %w", err)
	}
	return nil
}
