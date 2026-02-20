package media

import (
	"context"
	"fmt"
	"log"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleMedia struct {
	// Dependency injections (Repo, Usecase...) sẽ nằm ở đây
	session *gocql.Session
	client  *mongo.Database
}

func NewModuleMedia(session *gocql.Session, client *mongo.Database) *ModuleMedia {
	m := &ModuleMedia{
		session: session,
		client:  client,
	}
	if err := m.InitMongo(client); err != nil {
		log.Fatalf(">>> Media Module: Failed to initialize MongoDB - %v", err)
	}
	if err := m.InitCassandra(session); err != nil {
		log.Fatalf(">>> Media Module: Failed to initialize Cassandra - %v", err)
	}
	log.Println(">>> Media Module: Initialization Completed")
	return m
}

// =============================================================================
// 1. MONGODB INITIALIZATION
// =============================================================================

func (m *ModuleMedia) InitMongo(db *mongo.Database) error {
	ctx := context.Background()

	// 1. Collection: ALBUMS
	if err := m.initAlbumIndexes(ctx, db); err != nil {
		return err
	}

	// 2. Collection: LIVE SESSIONS
	if err := m.initLiveSessionIndexes(ctx, db); err != nil {
		return err
	}

	// 3. Collection: MUSIC LIBRARY
	if err := m.initMusicLibraryIndexes(ctx, db); err != nil {
		return err
	}

	// 4. Collection: REELS
	if err := m.initReelIndexes(ctx, db); err != nil {
		return err
	}

	// 5. Collection: STORIES
	if err := m.initStoryIndexes(ctx, db); err != nil {
		return err
	}

	// 6. Collection: MEDIA ASSETS
	if err := m.initMediaAssetIndexes(ctx, db); err != nil {
		return err
	}

	log.Println(">>> Media Module: MongoDB Indexes Initialized Successfully")
	return nil
}

// --- Helper: Albums ---
func (m *ModuleMedia) initAlbumIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.Album{}.CollectionName())

	models := []mongo.IndexModel{
		// A. LIST USER ALBUMS: Lấy danh sách album của User theo loại (VD: chỉ lấy album ảnh Profile)
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "type", Value: 1},
			},
		},
		// B. LIST GROUP ALBUMS: Lấy danh sách album của Nhóm
		{
			Keys: bson.D{{Key: "group_id", Value: 1}},
			// Partial Index: Chỉ index các album thuộc về Group (group_id tồn tại) -> Tiết kiệm RAM
			Options: options.Index().SetPartialFilterExpression(bson.M{
				"group_id": bson.M{"$exists": true},
			}),
		},
		// C. SORT RECENT: Sắp xếp album theo thời gian cập nhật ảnh cuối cùng
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "last_asset_added_at", Value: -1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for Albums: %w", err)
	}
	return nil
}

// --- Helper: Live Sessions ---
func (m *ModuleMedia) initLiveSessionIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.LiveSession{}.CollectionName())

	models := []mongo.IndexModel{
		// A. CHECK LIVE STATUS: Kiểm tra User X có đang Live không?
		{
			Keys: bson.D{
				{Key: "host_user_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		// B. BROWSE CATEGORY: Xem danh sách livestream đang active trong category Game/Music...
		{
			Keys: bson.D{
				{Key: "category_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		// C. LIST ALL ACTIVE: API lấy danh sách tất cả stream đang diễn ra (cho trang chủ)
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for LiveSessions: %w", err)
	}
	return nil
}

// --- Helper: Music Library ---
func (m *ModuleMedia) initMusicLibraryIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.MusicLibrary{}.CollectionName())

	models := []mongo.IndexModel{
		// A. FULL-TEXT SEARCH: Tìm bài hát theo tên, ca sĩ, hoặc lời bài hát
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "artist", Value: "text"},
				{Key: "lyrics_snippet", Value: "text"},
			},
			Options: options.Index().SetWeights(bson.M{
				"title":          10, // Ưu tiên tìm theo tên bài hát nhất
				"artist":         5,
				"lyrics_snippet": 1,
			}),
		},
		// B. TRENDING MUSIC: Top bài hát được sử dụng nhiều nhất
		{
			Keys: bson.D{{Key: "usage_count", Value: -1}},
		},
		// C. GENRE FILTER: Lọc theo thể loại (Multikey Index vì genres là array)
		{
			Keys: bson.D{{Key: "genre", Value: 1}}, // Lưu ý field json là "genre" nhưng bson struct là "genre"
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for MusicLibrary: %w", err)
	}
	return nil
}

// --- Helper: Reels ---
func (m *ModuleMedia) initReelIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.Reel{}.CollectionName())

	models := []mongo.IndexModel{
		// A. USER PROFILE: Tab Reels trên trang cá nhân
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		// B. DISCOVERY: Tìm kiếm theo Hashtag
		{
			Keys: bson.D{{Key: "hashtags", Value: 1}},
		},
		// C. AUDIO PAGE: Danh sách các video sử dụng cùng 1 bài nhạc
		{
			Keys: bson.D{{Key: "audio_meta.track_id", Value: 1}},
			Options: options.Index().SetPartialFilterExpression(bson.M{
				"audio_meta.track_id": bson.M{"$exists": true},
			}),
		},
		// D. PROCESSING QUEUE: Worker tìm các video đang xử lý để transcode
		{
			Keys: bson.D{{Key: "processing_status", Value: 1}},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for Reels: %w", err)
	}
	return nil
}

// --- Helper: Stories ---
func (m *ModuleMedia) initStoryIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.Story{}.CollectionName())

	models := []mongo.IndexModel{
		// A. FEED BAR: Load story của bạn bè (Logic: Query list user_id IN [...])
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		// B. CLEANUP / FILTER: Lọc các story còn hạn (chưa hết 24h) hoặc tìm story cũ để Archive
		{
			Keys: bson.D{
				{Key: "expires_at", Value: 1},
				{Key: "is_archived", Value: 1},
			},
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for Stories: %w", err)
	}
	return nil
}

// --- Helper: Media Assets ---
func (m *ModuleMedia) initMediaAssetIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(entity.MediaAsset{}.CollectionName())

	models := []mongo.IndexModel{
		// A. USER LIBRARY: Tab Ảnh/Video trên trang cá nhân
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		// B. ALBUM VIEW: Load ảnh trong 1 album theo thứ tự sắp xếp
		{
			Keys: bson.D{
				{Key: "album_id", Value: 1},
				{Key: "order", Value: 1},
			},
			Options: options.Index().SetPartialFilterExpression(bson.M{
				"album_id": bson.M{"$exists": true},
			}),
		},
		// C. GROUP MEDIA: Tab Media trong Group
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetPartialFilterExpression(bson.M{
				"group_id": bson.M{"$exists": true},
			}),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("failed to create indexes for MediaAssets: %w", err)
	}
	return nil
}

// =============================================================================
// 2. CASSANDRA INITIALIZATION
// =============================================================================

func (m *ModuleMedia) InitCassandra(session *gocql.Session) error {
	// 1. Table: LIVE COMMENTS (High Throughput Write)
	if err := m.initLiveCommentsTable(session); err != nil {
		return err
	}
	if err := m.InitViewsCassandraLiveCommentsTable(session); err != nil {
		return err
	}

	// 2. Table: STORY VIEWS (Write Heavy + Read List)
	if err := m.initStoryViewsTable(session); err != nil {
		return err
	}

	log.Println(">>> Media Module: Cassandra Tables Initialized Successfully")
	return nil
}

// --- Helper: Live Comments ---
func (m *ModuleMedia) initLiveCommentsTable(session *gocql.Session) error {
	// PARTITION KEY: stream_id (Gom comment của 1 buổi live vào 1 node)
	// CLUSTERING KEY: created_at DESC (Mới nhất lên đầu), comment_id (Uniqueness)
	query := `
		CREATE TABLE IF NOT EXISTS live_comments (
			stream_id uuid,
			created_at timestamp,
			comment_id uuid,
			user_id uuid,
			user_nickname text,
			user_avatar_url text,
			user_badges list<text>, 
			content text,
			is_pinned boolean,
			PRIMARY KEY (stream_id, created_at, comment_id)
		) WITH CLUSTERING ORDER BY (created_at DESC);
	`
	// Note: user_badges dùng list<text> để tương thích []string
	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create table live_comments: %w", err)
	}
	return nil
}
func (m *ModuleMedia) InitViewsCassandraLiveCommentsTable(session *gocql.Session) error {
	// Lệnh 1: Tạo Keyspace (nếu chưa có) - Thường làm thủ công hoặc ở bước riêng,
	// nhưng có thể để ở đây nếu dùng cho test cục bộ.

	// Lệnh 2: Tạo Materialized View
	createMVQuery := `
		CREATE MATERIALIZED VIEW IF NOT EXISTS live_comments_by_user AS
			SELECT *
			FROM live_comments
			WHERE user_id IS NOT NULL 
			  AND created_at IS NOT NULL 
			  AND stream_id IS NOT NULL 
			  AND comment_id IS NOT NULL
			PRIMARY KEY (user_id, created_at, stream_id, comment_id)
			WITH CLUSTERING ORDER BY (created_at DESC);
	`

	log.Println("Checking and initializing Cassandra schemas...")

	// Thực thi lệnh. Việc có IF NOT EXISTS giúp lệnh này an toàn dù chạy nhiều lần.
	if err := session.Query(createMVQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create materialized view live_comments_by_user: %w", err)
	}

	log.Println("Cassandra schemas initialized successfully.")
	return nil
}

// --- Helper: Story Views ---
func (m *ModuleMedia) initStoryViewsTable(session *gocql.Session) error {
	// PARTITION KEY: story_id (Gom tất cả view của 1 story)
	// CLUSTERING KEY: viewed_at DESC (Mới nhất lên đầu), viewer_id (1 user chỉ tính 1 lần tại 1 thời điểm)
	query := `
		CREATE TABLE IF NOT EXISTS story_views (
			story_id uuid,
			viewed_at timestamp,
			viewer_id uuid,
			viewer_name text,
			viewer_avatar_url text,
			interaction_type int,
			reaction_code text,
			poll_option_index int,
			PRIMARY KEY (story_id, viewed_at, viewer_id)
		) WITH CLUSTERING ORDER BY (viewed_at DESC);
	`
	// Note: interaction_type là Enum, lưu int trong DB cho gọn
	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create table story_views: %w", err)
	}
	return nil
}
