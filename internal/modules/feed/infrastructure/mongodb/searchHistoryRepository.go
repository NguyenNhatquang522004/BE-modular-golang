package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/feed/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchHistoryRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewSearchHistoryRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *SearchHistoryRepository {
	return &SearchHistoryRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *SearchHistoryRepository) CreateSearchHistory(ctx context.Context, searchHistory *entity.SearchHistory) error {
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	_, err := collection.InsertOne(ctx, searchHistory)
	if err != nil {
		return err
	}
	return nil
}
func (r *SearchHistoryRepository) CreateBulkSearchHistory(ctx context.Context, searchHistories []entity.SearchHistory) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	for i := range searchHistories {
		if searchHistories[i].ID.IsZero() {
			searchHistories[i].ID = primitive.NewObjectID() // Tạo ID mới cho mỗi bản ghi mới
		}
	}
	docs := utils.ToInterfaceSlice(searchHistories)
	options := options.InsertMany().SetOrdered(false) // ordered=false để tiếp tục chèn các bản ghi còn lại nếu có lỗi
	result, err := collection.InsertMany(ctx, docs, options)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var bulkErrors []*mongodbErrors.BulkError
	if err != nil {
		var writeException mongo.BulkWriteException
		if errors.As(err, &writeException) {
			for _, writeError := range writeException.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     searchHistories[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		}
	}

	return int64(len(result.InsertedIDs)), bulkErrors, nil
}
func (r *SearchHistoryRepository) GetSearchHistoriesByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"search_history_cache_userid_" + userID,
			"search_history_cache_nextcursor_userid_" + userID,
			"search_history_cache_hasnext_userid_" + userID,
			"search_history_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if pagesList, ok := datacache.([]*entity.SearchHistory); ok {
				return &dto.PaginationRes{
					Data:       pagesList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	var searchHistories []*entity.SearchHistory
	querylimit := int64(limit + 1)
	query := bson.M{"user_id": userID}
	if cursor != "" {
		decodedCursor, err := utils.DecodeCursorMongodb(cursor)
		if err != nil {
			return nil, err
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
	cursorMongo, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursorMongo.Close(ctx)
	if err := cursorMongo.All(ctx, &searchHistories); err != nil {
		return nil, err
	}
	var nextCursor string
	hasNext := false
	if int64(len(searchHistories)) == querylimit {
		hasNext = true
		lastItem := searchHistories[len(searchHistories)-1]
		nextCursor = utils.EncodeCursorMongodb(lastItem.CreatedAt, lastItem.ID)
		searchHistories = searchHistories[:len(searchHistories)-1] // Loại bỏ phần tử cuối cùng để trả về đúng limit
	}
	if cursor == "" {
		items := map[string]any{
			"search_history_cache_userid_" + userID:            searchHistories,
			"search_history_cache_nextcursor_userid_" + userID: nextCursor,
			"search_history_cache_hasnext_userid_" + userID:    hasNext,
			"search_history_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       searchHistories,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *SearchHistoryRepository) GetSearchHistoryByUserIDTop(ctx context.Context, userID string) ([]*entity.SearchHistory, error) {
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	var searchHistories []*entity.SearchHistory
	query := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1}, // Tie-breaker
		}).
		SetLimit(10)
	cursorMongo, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursorMongo.Close(ctx)
	if err := cursorMongo.All(ctx, &searchHistories); err != nil {
		return nil, err
	}
	return searchHistories, nil
}
func (r *SearchHistoryRepository) UpdateSearchHistory(ctx context.Context, searchHistory *entity.SearchHistory) error {
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	if searchHistory.ID.IsZero() {
		return errors.New("ID is required for update")
	}
	filter := bson.M{"_id": searchHistory.ID}
	update := bson.M{"$set": searchHistory}
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
func (r *SearchHistoryRepository) UpdateBulkSearchHistory(ctx context.Context, searchHistories []entity.SearchHistory) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	var bulkErrors []*mongodbErrors.BulkError
	models := make([]mongo.WriteModel, len(searchHistories))
	for i, sh := range searchHistories {
		models[i] = mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": sh.ID}).SetUpdate(bson.M{"$set": sh}).SetUpsert(false)
	}
	bulkOption := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, bulkOption)
	if err != nil {
		var writeException mongo.BulkWriteException
		if errors.As(err, &writeException) {
			for _, writeError := range writeException.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     searchHistories[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.ModifiedCount, bulkErrors, nil
}
func (r *SearchHistoryRepository) DeleteSearchHistory(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *SearchHistoryRepository) DeleteBulkSearchHistory(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.SearchHistory{}.CollectionName())
	var model []mongo.WriteModel
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, err
		}
		model = append(model, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": objectID}))
	}
	bulkOption := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, bulkOption)
	if result == nil && err != nil {
		return 0, nil, err
	}

	var bulkErrors []*mongodbErrors.BulkError
	if err != nil {
		var writeException mongo.BulkWriteException
		if errors.As(err, &writeException) {
			for _, writeError := range writeException.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     ids[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.DeletedCount, bulkErrors, nil
}
