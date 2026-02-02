package social

import (
	"context"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/social/domain/entity"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ModuleSocial struct {
}

func (m *ModuleSocial) InitPostgres(db *gorm.DB) error {
	return db.AutoMigrate(&entity.Followers{},
		&entity.Friendships{},
		&entity.User_Blocks{},
	)
}

func (m *ModuleSocial) InitMongo(db *mongo.Database) error {
	ctx := context.Background()
	collProfile := db.Collection(entity.Profiles{}.CollectionNameProfiles())
	indexModelUsesetting := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collProfile.Indexes().CreateOne(ctx, indexModelUsesetting)

	if err != nil {
		return fmt.Errorf("init index search_history failed: %v", err)
	}
	return nil
}

func (m *ModuleSocial) RegisterRoute(r *gin.RouterGroup) {

}
