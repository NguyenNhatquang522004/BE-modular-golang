package feed

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleFeed struct {
	client *mongo.Database
}

func NewModuleFeed(db *mongo.Database) *ModuleFeed {
	module := &ModuleFeed{
		client: db,
	}
	if err := module.InitMongo(db); err != nil {
		panic("Failed to initialize MongoDB indexes for ModuleFeed: " + err.Error())
	}
	return module
}
func (m *ModuleFeed) InitMongo(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(entity.CollectionSearchHistories)

	// Định nghĩa các Indexes theo Best Practice
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			// 1. Compound Index: { user_id: 1, created_at: -1 }
			// Tối ưu cho truy vấn: "Lấy 10 từ khóa tìm kiếm gần nhất của User A"
			// Sắp xếp created_at: -1 giúp việc lấy dữ liệu mới nhất cực nhanh mà không cần Sort lại
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			// 2. TTL Index (Time-To-Live): Tự động xóa dữ liệu sau 30 ngày (2,592,000 giây)
			// Giúp duy trì dung lượng DB gọn nhẹ, chỉ giữ lại lịch sử tìm kiếm gần đây
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(2592000),
		},
	})

	return err
}
