package mongodb

import (
	"context"
	"errors"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserSettingRepository struct {
	client *mongo.Database
}

func NewUserSettingRepository(client *mongo.Database) *UserSettingRepository {
	return &UserSettingRepository{
		client: client,
	}
}

func (r *UserSettingRepository) GetUserSettings(ctx context.Context, userID string) (*entity.UserSetting, error) {
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())
	filter := bson.M{"User_ID": userID}
	var userdata entity.UserSetting
	err := collection.FindOne(ctx, filter).Decode(&userdata)
	if err != nil {
		return nil, errors.New("user settings not found: " + err.Error())
	}

	return &userdata, nil
}

func (r *UserSettingRepository) UpdateUserSettings(ctx context.Context, userID string, settings *entity.UserSetting) error {
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())
	filter := bson.M{"User_ID": userID}
	update := bson.M{"$set": settings}
	opts := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errors.New("failed to update user settings: " + err.Error())
	}
	if result.MatchedCount == 0 {
		return errors.New("user settings not found")
	}

	return nil
}
func (r *UserSettingRepository) DeleteUserSettings(ctx context.Context, userID string) error {
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())
	filter := bson.M{"User_ID": userID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("failed to delete user settings: " + err.Error())
	}
	if result.DeletedCount == 0 {
		return errors.New("user settings not found")
	}
	return nil
}
func (r *UserSettingRepository) CreateUserSettingDefault(ctx context.Context, userSetting *entity.UserSetting) error {
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())
	_, err := collection.InsertOne(ctx, userSetting)
	if err != nil {
		return errors.New("failed to create default user settings: " + err.Error())
	}
	return nil

}
