package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoryRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewStoryRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *StoryRepository {
	return &StoryRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *StoryRepository) CreateStory(ctx context.Context, story *entity.Story) error {
	collection := r.client.Collection(entity.Story{}.CollectionName())
	story.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, story)
	if err != nil {
		return err
	}
	return nil
}
func (r *StoryRepository) CreateBulkStories(ctx context.Context, stories []*entity.Story) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.Story{}.CollectionName())
	for _, story := range stories {
		if story.ID.IsZero() {
			story.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(stories)
	result, err := collection.InsertMany(ctx, docs, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}

	var failedDocs []*dto.BulkError
	if err != nil {
		// Kiểm tra xem có phải lỗi BulkWriteException không
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedDocs = append(failedDocs, &dto.BulkError{
					ID:     stories[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *StoryRepository) GetStoryByID(ctx context.Context, id string) (*entity.Story, error) {
	// Implement the logic to retrieve a story by its ID from MongoDB
	collection := r.client.Collection(entity.Story{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var story entity.Story
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&story)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Story not found
		}
		return nil, err
	}
	return &story, nil
}

func (r *StoryRepository) GetBulkStoriesByID(ctx context.Context, ids []string) ([]*entity.Story, error) {
	// Implement the logic to retrieve multiple stories by their IDs from MongoDB
	var objIDs []*primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, &objID)
	}

	collection := r.client.Collection(entity.Story{}.CollectionName())
	cursor, err := collection.Find(ctx, primitive.M{"_id": primitive.M{"$in": objIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stories []*entity.Story
	if err = cursor.All(ctx, &stories); err != nil {
		return nil, err
	}

	return stories, nil
}
func (r *StoryRepository) GetStoriesByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Story{}.CollectionName())
	var stories []*entity.Story
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"story_cache_userid_" + userID,
			"story_cache_nextcursor_userid_" + userID,
			"story_cache_hasnext_userid_" + userID,
			"story_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if storiesList, ok := datacache.([]*entity.Story); ok {
				return &dto.PaginationRes{
					Data:       storiesList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	query := bson.M{"user_id": userID}

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
	if err = cursorDB.All(ctx, &stories); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(stories) > limit {
		hasNext = true
		stories = stories[:limit] // Lấy đúng số lượng cần thiết
		lastStory := stories[len(stories)-1]
		nextCursor = utils.EncodeCursorMongodb(lastStory.CreatedAt, lastStory.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"story_cache_userid_" + userID:            stories,
			"story_cache_nextcursor_userid_" + userID: nextCursor,
			"story_cache_hasnext_userid_" + userID:    hasNext,
			"story_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       stories,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *StoryRepository) UpdateStory(ctx context.Context, story *entity.Story) error {
	// Implement the logic to update a story in MongoDB
	return nil
}
func (r *StoryRepository) UpdateBulkStories(ctx context.Context, stories []*entity.Story) (int64, []*dto.BulkError, error) {
	// Implement the logic to update multiple stories in MongoDB
	return 0, nil, nil
}
func (r *StoryRepository) DeleteStory(ctx context.Context, id string) error {
	// Implement the logic to delete a story by its ID from MongoDB
	return nil
}
func (r *StoryRepository) DeleteBulkStories(ctx context.Context, ids []string) (int64, []*dto.BulkError, error) {
	// Implement the logic to delete multiple stories by their IDs from MongoDB
	return 0, nil, nil
}
