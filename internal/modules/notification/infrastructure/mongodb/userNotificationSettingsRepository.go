package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserNotificationSettingsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewUserNotificationSettingsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *UserNotificationSettingsRepository {
	return &UserNotificationSettingsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *UserNotificationSettingsRepository) CreateUserNotificationSettings(ctx context.Context, settings *entity.UserNotificationSetting) error {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())
	if settings.ID.IsZero() {
		settings.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, settings)
	if err != nil {
		return err
	}

	return nil
}
func (r *UserNotificationSettingsRepository) CreateBulkUserNotificationSettings(ctx context.Context, settings []*entity.UserNotificationSetting) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())
	docs := utils.ToInterfaceSlice(settings)
	result, err := collection.InsertMany(ctx, docs)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkWriteException mongo.BulkWriteException
		if errors.As(err, &bulkWriteException) {
			for _, writeError := range bulkWriteException.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     settings[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return int64(len(settings)), faildocs, nil
}

func (r *UserNotificationSettingsRepository) GetUserNotificationSettingsByUserID(ctx context.Context, userID string) (*entity.UserNotificationSetting, error) {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())
	var result *entity.UserNotificationSetting
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
func (r *UserNotificationSettingsRepository) UpdateUserNotificationSettings(ctx context.Context, settings *entity.UserNotificationSetting) error {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())
	if settings.ID.IsZero() {
		return errors.New("ID is required for update")
	}
	filter := bson.M{"user_id": settings.UserID}
	update := bson.M{"$set": bson.M{
		"settings":   settings.Settings,
		"fcm_tokens": settings.FCMTokens,
	}}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, options)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return mongo.ErrNoDocuments
		}
		return err
	}
	return nil
}
func (r *UserNotificationSettingsRepository) UpdateBulkUserNotificationSettings(ctx context.Context, settings []*entity.UserNotificationSetting) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())
	var models []mongo.WriteModel
	for _, setting := range settings {
		if setting.ID.IsZero() {
			return 0, nil, errors.New("ID is required for update")
		}
		filter := bson.M{"user_id": setting.UserID}
		update := bson.M{"$set": bson.M{
			"settings":   setting.Settings,
			"fcm_tokens": setting.FCMTokens,
		}}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	result, err := collection.BulkWrite(ctx, models)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkWriteException mongo.BulkWriteException
		if errors.As(err, &bulkWriteException) {
			for _, writeError := range bulkWriteException.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     settings[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.ModifiedCount, faildocs, nil
}
func (r *UserNotificationSettingsRepository) DeleteUserNotificationSettings(ctx context.Context, userID string) error {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())
	result, err := collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *UserNotificationSettingsRepository) DeleteBulkUserNotificationSettings(ctx context.Context, userIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())

	result, err := collection.DeleteMany(ctx, bson.M{"user_id": bson.M{"$in": userIDs}})
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkWriteException mongo.BulkWriteException
		if errors.As(err, &bulkWriteException) {
			for _, writeError := range bulkWriteException.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     userIDs[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.DeletedCount, faildocs, nil
}
func (r *UserNotificationSettingsRepository) GetUserNotificationSettingsByDateOfBirth(ctx context.Context, month int, day int, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.UserNotificationSetting{}.CollectionName())
	var settings []*entity.UserNotificationSetting
	querylimit := int64(limit + 1)
	// Check cache first
	query := bson.M{"$expr": bson.M{
		"$and": []bson.M{
			{"$eq": []interface{}{bson.M{"$month": "$birthday"}, month}},
			{"$eq": []interface{}{bson.M{"$dayOfMonth": "$birthday"}, day}},
		},
	}}

	if cursor != "" {
		decodedCursor, err := utils.DecodeCursorMongodb(cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		query["$or"] = []bson.M{
			{"created_at": bson.M{"$lt": decodedCursor.CreatedAt}},
			{"created_at": decodedCursor.CreatedAt, "_id": bson.M{"$lt": decodedCursor.PostID}},
		}
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1}, // Tie-breaker
		}).
		SetLimit(int64(querylimit))
	cursorDB, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursorDB.Close(ctx)
	if err = cursorDB.All(ctx, &settings); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(settings) > limit {
		hasNext = true
		settings = settings[:limit] // Lấy đúng số lượng cần thiết
		lastSetting := settings[len(settings)-1]
		nextCursor = utils.EncodeCursorMongodb(lastSetting.CreatedAt, lastSetting.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	return &dto.PaginationRes{
		Data:       settings,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil

}
