package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/server/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/enum"
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

func (r *UserSettingRepository) GetUserSettings(userID string) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())
	filter := bson.M{"User_ID": userID}
	var userdata entity.UserSetting
	err := collection.FindOne(ctx, filter).Decode(&userdata)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()),
			response.WithStatus("404")), errors.New("user settings not found: " + err.Error())
	}

	return response.NewResponse(response.WithData(userdata),
		response.WithMessage("User settings retrieved successfully"),
		response.WithStatus("200")), nil
}

func (r *UserSettingRepository) UpdateUserSettings(userID string, settings *entity.UserSetting) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())
	filter := bson.M{"User_ID": userID}
	update := bson.M{"$set": settings}
	opts := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()),
			response.WithStatus("500")), errors.New("failed to update user settings: " + err.Error())
	}
	if result.MatchedCount == 0 {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("user settings not found"),
			response.WithStatus("404")), errors.New("user settings not found")
	}

	return response.NewResponse(response.WithData(""),
		response.WithMessage("User settings updated successfully"),
		response.WithStatus("200")), nil
}
func (r *UserSettingRepository) DeleteUserSettings(userID string) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())
	filter := bson.M{"User_ID": userID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()),
			response.WithStatus("500")), errors.New("failed to delete user settings: " + err.Error())
	}
	if result.DeletedCount == 0 {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("user settings not found"),
			response.WithStatus("404")), errors.New("user settings not found")
	}
	return response.NewResponse(response.WithData(""),
		response.WithMessage("User settings deleted successfully"),
		response.WithStatus("200")), nil
}
func (r *UserSettingRepository) CreateUserSettingDefault(userID string) (*response.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.UserSetting{}.CollectionName())

	defaultSettings := &entity.UserSetting{
		User_ID:                       userID,
		Theme_Mode:                    enum.ThemeMode(0),
		Font_Size:                     enum.FontSize(0),
		Compact_Mode:                  false,
		Lang_Code:                     enum.LangCode(0),
		Timezone:                      enum.Timezone(0),
		Auto_Translate:                false,
		Default_Post_Audience:         enum.PrivacyLevel(0),
		Default_Story_Audience:        enum.PrivacyLevel(0),
		Allow_Friend_Request_From:     enum.PrivacyLevel(0),
		Allow_Friend_List_View_From:   enum.PrivacyLevel(0),
		Allow_Email_Lookup_From:       enum.PrivacyLevel(0),
		Allow_Phone_Lookup_From:       enum.PrivacyLevel(0),
		Allow_Search_Engine_Indexing:  false,
		Allow_Timeline_Posting_From:   enum.PrivacyLevel(0),
		Review_Tags_Enabled:           false,
		Review_Timeline_Posts_Enabled: false,
		Notifications: &entity.NotificationSettings{
			EmailFrequency:   "daily",
			PushInteractions: true,
			PushFriends:      true,
			PushGroups:       true,
			PushEvents:       true,
			PushBirthdays:    true,
		},
		Updated_At: time.Now(),
	}

	_, err := collection.InsertOne(ctx, defaultSettings)
	if err != nil {
		return response.NewResponse(response.WithData(""),
			response.WithMessage(err.Error()),
			response.WithStatus("500")), errors.New("failed to create default user settings: " + err.Error())
	}

	return response.NewResponse(response.WithData(""),
		response.WithMessage("Default user settings created successfully"),
		response.WithStatus("201")), nil
}
