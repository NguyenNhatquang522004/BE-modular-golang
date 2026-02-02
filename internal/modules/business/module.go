package business

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ModuleBusiness struct {
	// Dependency Injection (UseCases, Repositories...)
}

func (m *ModuleBusiness) InitMongo(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- 1. Collection: Pages ---
	pagesCol := db.Collection(entity.CollectionPages)
	_, err := pagesCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			// Unique Slug: Quan trọng cho SEO và định danh URL
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			// Text Index: Cho phép tìm kiếm theo tên Page
			Keys: bson.D{{Key: "name", Value: "text"}},
		},
		{
			// 2dsphere Index: Hỗ trợ tìm cửa hàng/địa điểm gần đây
			Keys: bson.D{{Key: "address.coordinates", Value: "2dsphere"}},
		},
		{
			// Lọc danh sách page của 1 user nhanh hơn
			Keys: bson.D{{Key: "creator_user_id", Value: 1}},
		},
	})
	if err != nil {
		return err
	}

	// --- 2. Collection: PageRoles ---
	rolesCol := db.Collection(entity.CollectionPageRoles)
	_, err = rolesCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		// Compound Unique Index: Đảm bảo 1 user chỉ có 1 role duy nhất trên 1 page
		Keys:    bson.D{{Key: "page_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return err
}

func (m *ModuleBusiness) InitPostgres(db *gorm.DB) error {
	// 1. Auto Migrate các bảng theo đúng thứ tự phụ thuộc
	err := db.AutoMigrate(
		&entity.AdAccount{},
		&entity.AdCampaign{},
		&entity.Ad{},
	)
	if err != nil {
		return err
	}

	// 2. Thiết lập Additional Constraints & Indexes nâng cao (nếu GORM tag chưa bao phủ hết)

	// AdAccount: Đảm bảo Currency luôn in hoa (Ví dụ về Business Rule)
	// GORM đã tạo index cho OwnerUserID và Status qua tags.

	// AdCampaign: Index phức hợp cho việc lọc Campaign theo Account + Status
	db.Exec("CREATE INDEX IF NOT EXISTS idx_campaigns_account_status ON ad_campaigns (account_id, status)")

	// Ad: Index cho TargetPostID vì đây là điểm nối giữa Postgres và Mongo
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ads_target_post ON ads (target_post_id)")

	return nil
}
func (m *ModuleBusiness) InitCassandra(session *gocql.Session) error {
	// Khởi tạo bảng PageDailyMetrics
	// Partition Key: page_id (Dữ liệu 1 page nằm cùng 1 node)
	// Clustering Key: metric_date (Sắp xếp dữ liệu theo thời gian giảm dần)
	query := `
	CREATE TABLE IF NOT EXISTS page_daily_metrics (
		page_id uuid,
		metric_date date,
		reach_total bigint,
		reach_paid bigint,
		reach_organic bigint,
		impressions_total bigint,
		new_followers int,
		unfollows int,
		profile_views int,
		website_clicks int,
		cta_clicks int,
		PRIMARY KEY ((page_id), metric_date)
	) WITH CLUSTERING ORDER BY (metric_date DESC);`

	return session.Query(query).Exec()
}
