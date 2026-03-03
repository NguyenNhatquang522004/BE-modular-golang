package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageFollowersRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewPageFollowersRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *PageFollowersRepository {
	return &PageFollowersRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *PageFollowersRepository) CreateFollower(ctx context.Context, follower *entity.PageFollower) error {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	if follower.ID.IsZero() {
		follower.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, follower)
	return err
}
func (r *PageFollowersRepository) CreateBulkFollowers(ctx context.Context, followers []*entity.PageFollower) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	docs := utils.ToInterfaceSlice(followers)
	result, err := collection.InsertMany(ctx, docs)
	if result == nil && err != nil {
		return 0, nil, err
	}
	insertedCount := int64(len(result.InsertedIDs))
	var bulkErrors []*mongodbErrors.BulkError
	if err != nil {
		var mongoBulkErr mongo.BulkWriteException
		if errors.As(err, &mongoBulkErr) {
			for _, writeErr := range mongoBulkErr.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     followers[writeErr.Index].ID.Hex(),
					Reason: writeErr.Message,
				})
			}
			return insertedCount, bulkErrors, nil
		}
		return 0, nil, err
	}
	return insertedCount, bulkErrors, nil
}
func (r *PageFollowersRepository) GetFollowersByPageID(ctx context.Context, pageID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	var followers []*entity.PageFollower
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"follower_cache_pageid_" + pageID,
			"follower_cache_nextcursor_pageid_" + pageID,
			"follower_cache_hasnext_pageid_" + pageID,
			"follower_cache_limit_pageid_" + pageID,
		})
		if err == nil && datacache != nil {
			if followersList, ok := datacache.([]*entity.PageFollower); ok {
				return &dto.PaginationRes{
					Data:       followersList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(pageID)
	query := bson.M{"page_id": finalid}

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
	if err = cursorDB.All(ctx, &followers); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(followers) > limit {
		hasNext = true
		followers = followers[:limit] // Lấy đúng số lượng cần thiết
		lastFollower := followers[len(followers)-1]
		nextCursor = utils.EncodeCursorMongodb(lastFollower.CreatedAt, lastFollower.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"follower_cache_pageid_" + pageID:            followers,
			"follower_cache_nextcursor_pageid_" + pageID: nextCursor,
			"follower_cache_hasnext_pageid_" + pageID:    hasNext,
			"follower_cache_limit_pageid_" + pageID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       followers,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *PageFollowersRepository) GetFollowerByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	var followers []*entity.PageFollower
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"follower_cache_userid_" + userID,
			"follower_cache_nextcursor_userid_" + userID,
			"follower_cache_hasnext_userid_" + userID,
			"follower_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if followersList, ok := datacache.([]*entity.PageFollower); ok {
				return &dto.PaginationRes{
					Data:       followersList,
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
	if err = cursorDB.All(ctx, &followers); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(followers) > limit {
		hasNext = true
		followers = followers[:limit] // Lấy đúng số lượng cần thiết
		lastFollower := followers[len(followers)-1]
		nextCursor = utils.EncodeCursorMongodb(lastFollower.CreatedAt, lastFollower.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"follower_cache_userid_" + userID:            followers,
			"follower_cache_nextcursor_userid_" + userID: nextCursor,
			"follower_cache_hasnext_userid_" + userID:    hasNext,
			"follower_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       followers,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *PageFollowersRepository) UpdateFollower(ctx context.Context, follower *entity.PageFollower) error {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	filter := bson.M{"_id": follower.ID}
	update := bson.M{"$set": follower}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After) // Trả về document sau khi cập nhật
	var updatedFollower *entity.PageFollower
	err := collection.FindOneAndUpdate(ctx, filter, update, options).Decode(&updatedFollower)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("follower not found: %w", err)
		}
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *PageFollowersRepository) UpdateBulkFollowers(ctx context.Context, followers []*entity.PageFollower) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	var models []mongo.WriteModel
	for _, follower := range followers {
		filter := bson.M{"_id": follower.ID}
		update := bson.M{"$set": follower}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	modifiedCount := result.ModifiedCount
	var bulkErrors []*mongodbErrors.BulkError
	if err != nil {
		var mongoBulkErr mongo.BulkWriteException
		if errors.As(err, &mongoBulkErr) {
			for _, writeErr := range mongoBulkErr.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     followers[writeErr.Index].ID.Hex(),
					Reason: writeErr.Message,
				})
			}
			return modifiedCount, bulkErrors, nil
		}
		return 0, nil, err
	}
	return modifiedCount, bulkErrors, nil
}
func (r *PageFollowersRepository) DeleteFollower(ctx context.Context, pageID string, userID string) error {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	finalPageID, err := primitive.ObjectIDFromHex(pageID)
	if err != nil {
		return fmt.Errorf("invalid page ID: %w", err)
	}
	filter := bson.M{"page_id": finalPageID, "user_id": userID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("follower not found")
	}
	return nil
}
func (r *PageFollowersRepository) DeleteBulkFollowers(ctx context.Context, pageID string, userIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	finalPageID, err := primitive.ObjectIDFromHex(pageID)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid page ID: %w", err)
	}
	filter := bson.M{"page_id": finalPageID, "user_id": bson.M{"$in": userIDs}}
	result, err := collection.DeleteMany(ctx, filter)
	if result == nil && err != nil {
		return 0, nil, err
	}
	deletedCount := result.DeletedCount
	var bulkErrors []*mongodbErrors.BulkError
	if err != nil {
		var mongoBulkErr mongo.BulkWriteException
		if errors.As(err, &mongoBulkErr) {
			for _, writeErr := range mongoBulkErr.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     userIDs[writeErr.Index],
					Reason: writeErr.Message,
				})
			}
			return deletedCount, bulkErrors, nil
		}
		return 0, nil, err
	}
	return deletedCount, bulkErrors, nil
}
func (r *PageFollowersRepository) DeleteBulkFollowersByPageID(ctx context.Context, pageID string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	finalPageID, err := primitive.ObjectIDFromHex(pageID)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid page ID: %w", err)
	}
	filter := bson.M{"page_id": finalPageID}
	result, err := collection.DeleteMany(ctx, filter)
	if result == nil && err != nil {
		return 0, nil, err
	}
	deletedCount := result.DeletedCount
	var bulkErrors []*mongodbErrors.BulkError
	if err != nil {
		var mongoBulkErr mongo.BulkWriteException
		if errors.As(err, &mongoBulkErr) {
			for _, writeErr := range mongoBulkErr.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     pageID,
					Reason: writeErr.Message,
				})
			}
			return deletedCount, bulkErrors, nil
		}
		return 0, nil, err
	}
	return deletedCount, bulkErrors, nil
}

func (r *PageFollowersRepository) DeleteFollowerByPageID(ctx context.Context, pageID string) error {
	collection := r.client.Collection(entity.PageFollower{}.CollectionName())
	finalPageID, err := primitive.ObjectIDFromHex(pageID)
	if err != nil {
		return fmt.Errorf("invalid page ID: %w", err)
	}
	filter := bson.M{"page_id": finalPageID}
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no followers found for this page")
	}
	return nil
}
