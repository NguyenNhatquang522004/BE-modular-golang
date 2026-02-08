package identity

import (
	"context"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ModuleIdentity struct {
	client *mongo.Database
	db     *gorm.DB
}

func NewModuleIdentity(client *mongo.Database, db *gorm.DB) *ModuleIdentity {
	return &ModuleIdentity{
		client: client,
		db:     db,
	}
}

func (m *ModuleIdentity) RegisterRoute(r *gin.RouterGroup) {

}

func (m *ModuleIdentity) InitPostgres() error {
	return m.db.AutoMigrate(&entity.User{},
		&entity.UserSession{},
		&entity.UserRole{},
	)
}

func (m *ModuleIdentity) InitMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Lấy collection từ domain thông qua helper method
	collection := m.client.Collection(entity.UserSetting{}.CollectionName())
	// Đánh Index cho UserSetting
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		// 1. Unique Index: { user_id: 1 }
		// Đảm bảo tính toàn vẹn: 1 User chỉ có 1 bản ghi Settings duy nhất.
		// Tối ưu tốc độ: Truy vấn cấu hình khi user đăng nhập cực nhanh.
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return err
	}

	return nil
}
