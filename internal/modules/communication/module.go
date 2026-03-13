package communication

import (
	"context"
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleCommunication struct {
	// Dependencies...
	client  *mongo.Database
	session *gocql.Session
}

func NewModuleCommunication(client *mongo.Database, session *gocql.Session) *ModuleCommunication {
	m := &ModuleCommunication{
		client:  client,
		session: session,
	}
	if err := m.InitMongo(client); err != nil {
		log.Fatalf("Failed to initialize MongoDB for Communication Module: %v", err)
	}
	if err := m.InitCassandra(session); err != nil {
		log.Fatalf("Failed to initialize Cassandra for Communication Module: %v", err)
	}
	return m
}

// =============================================================================
// 1. MONGODB INITIALIZATION
// =============================================================================

func (m *ModuleCommunication) InitMongo(db *mongo.Database) error {
	ctx := context.Background()

	// 1. Collection: CONVERSATIONS
	if err := m.initConversationIndexes(ctx, db); err != nil {
		return err
	}

	// 2. Collection: CONVERSATION PARTICIPANTS
	if err := m.initParticipantIndexes(ctx, db); err != nil {
		return err
	}

	// 3. Collection: CALL LOGS
	if err := m.initCallLogIndexes(ctx, db); err != nil {
		return err
	}

	log.Println(">>> Communication Module: MongoDB Indexes Initialized Successfully")
	return nil
}

// --- Helper: Conversations ---
func (m *ModuleCommunication) initConversationIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.Conversation{}.CollectionName())

	models := []mongo.IndexModel{
		// A. COMMUNITY CHANNELS: Tìm nhóm chat liên kết với 1 Group cộng đồng
		{
			Keys: bson.D{{Key: "related_group_id", Value: 1}},
			// Partial Index: Chỉ index những cuộc hội thoại có liên kết group
			Options: options.Index().SetPartialFilterExpression(bson.M{
				"related_group_id": bson.M{"$exists": true},
			}),
		},
		// B. CREATOR LOOKUP: Tìm các nhóm do user tạo (Optional)
		{
			Keys: bson.D{{Key: "creator_id", Value: 1}},
		},
		// C. PRIVATE CHAT 1-1: Đảm bảo không trùng lặp cuộc hội thoại giữa 2 người (Idempotency)
		{
			Keys: bson.D{{Key: "private_chat_key", Value: 1}},
			Options: options.Index().
				SetUnique(true). // Bắt buộc Unique để chống Race Condition khi tạo chat 1-1
				SetSparse(true), // Rất quan trọng: Bỏ qua các Group Chat (vì Group chat không có field private_chat_key)
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for Conversations: %w", err)
	}
	return nil
}

// --- Helper: Participants (QUAN TRỌNG NHẤT) ---
func (m *ModuleCommunication) initParticipantIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.ConversationParticipant{}.CollectionName())

	models := []mongo.IndexModel{
		// A. INBOX LIST: Load danh sách hội thoại của 1 User
		// Sắp xếp theo lần cuối online/hoạt động để hiển thị inbox active lên đầu
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "last_seen_at", Value: -1},
			},
		},
		// B. MEMBER LIST: Load thành viên của 1 nhóm
		{
			Keys: bson.D{{Key: "conversation_id", Value: 1}},
		},
		// C. UNIQUE CONSTRAINT: Đảm bảo 1 user không thể tham gia 1 nhóm 2 lần
		// Đây là index quan trọng để tránh bug duplicate member
		{
			Keys: bson.D{
				{Key: "conversation_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for ConversationParticipants: %w", err)
	}
	return nil
}

// --- Helper: Call Logs ---
func (m *ModuleCommunication) initCallLogIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.CallLog{}.CollectionName())

	models := []mongo.IndexModel{
		// A. GROUP CALL HISTORY: Lịch sử cuộc gọi trong 1 nhóm
		{
			Keys: bson.D{
				{Key: "conversation_id", Value: 1},
				{Key: "started_at", Value: -1},
			},
		},
		// B. USER CALL HISTORY: Lịch sử cuộc gọi của user (Multikey Index trên mảng participants)
		{
			Keys: bson.D{
				{Key: "participants", Value: 1},
				{Key: "started_at", Value: -1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for CallLogs: %w", err)
	}
	return nil
}

// =============================================================================
// 2. CASSANDRA INITIALIZATION
// =============================================================================

func (m *ModuleCommunication) InitCassandra(session *gocql.Session) error {
	// 1. Table: MESSAGES (Core Chat)
	if err := m.initMessagesTable(session); err != nil {
		return err
	}

	// 2. Table: MESSAGE REACTIONS
	if err := m.initMessageReactionsTable(session); err != nil {
		return err
	}

	// 3. Table: READ STATE (Seen status)
	if err := m.initReadStateTable(session); err != nil {
		return err
	}

	log.Println(">>> Communication Module: Cassandra Tables Initialized Successfully")
	return nil
}

// --- Helper: Messages ---
func (m *ModuleCommunication) initMessagesTable(session *gocql.Session) error {
	// PARTITION KEY: conversation_id, bucket (Chia nhỏ partition để tránh row quá lớn)
	// CLUSTERING KEY: message_id DESC (Load tin nhắn mới nhất trước -> Hành vi cuộn trang chat)
	query := `
		CREATE TABLE IF NOT EXISTS messages (
			conversation_id text,
			bucket int,
			message_id timeuuid,
			sender_id uuid,
			type int,
			content text,
			attachments list<text>,
			reply_to_message_id uuid,
			story_ref_id uuid,
			is_revoked boolean,
			created_at timestamp,
			PRIMARY KEY ((conversation_id, bucket), message_id)
		) WITH CLUSTERING ORDER BY (message_id DESC);
	`
	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create table messages: %w", err)
	}
	return nil
}

// --- Helper: Message Reactions ---
func (m *ModuleCommunication) initMessageReactionsTable(session *gocql.Session) error {
	// PARTITION KEY: conversation_id, message_id (Gom tất cả reaction của 1 tin nhắn vào 1 chỗ)
	// CLUSTERING KEY: user_id (Đảm bảo 1 user chỉ thả 1 reaction)
	query := `
		CREATE TABLE IF NOT EXISTS message_reactions (
			conversation_id text,
			message_id timeuuid,
			user_id uuid,
			reaction_code int,
			created_at timestamp,
			PRIMARY KEY ((conversation_id, message_id), user_id)
		);
	`
	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create table message_reactions: %w", err)
	}
	return nil
}

// --- Helper: Read State ---
func (m *ModuleCommunication) initReadStateTable(session *gocql.Session) error {
	// PARTITION KEY: conversation_id (Để query: "Trong nhóm này, ai đã đọc đến đâu?")
	// CLUSTERING KEY: user_id
	query := `
		CREATE TABLE IF NOT EXISTS conversation_read_state (
			conversation_id text,
			user_id uuid,
			last_read_message_id timeuuid,
			last_read_at timestamp,
			PRIMARY KEY ((conversation_id), user_id)
		);
	`
	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create table conversation_read_state: %w", err)
	}
	return nil
}
