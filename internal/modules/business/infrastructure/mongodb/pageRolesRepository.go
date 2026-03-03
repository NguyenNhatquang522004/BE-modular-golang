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

type PagesRolesRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewPagesRolesRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *PagesRolesRepository {
	return &PagesRolesRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}

// Define the methods for the PageRolesRepository interface here
func (r *PagesRolesRepository) CreatePageRole(ctx context.Context, pageRole *entity.PageRole) error {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	if pageRole.ID.IsZero() {
		pageRole.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, pageRole)
	if err != nil {
		return err
	}
	return nil

}
func (r *PagesRolesRepository) CreateBulkPageRoles(ctx context.Context, pageRoles []*entity.PageRole) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	for _, pageRole := range pageRoles {
		if pageRole.ID.IsZero() {
			pageRole.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(pageRoles)
	result, err := collection.InsertMany(ctx, docs, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}

	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		// Kiểm tra xem có phải lỗi BulkWriteException không
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedDocs = append(failedDocs, &mongodbErrors.BulkError{
					ID:     pageRoles[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *PagesRolesRepository) GetPageRoleByID(ctx context.Context, pageRoleID string) (*entity.PageRole, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(pageRoleID)
	if err != nil {
		return nil, err
	}

	var pageRole entity.PageRole
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&pageRole)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err
	}
	return &pageRole, nil
}
func (r *PagesRolesRepository) GetPageRolesByPageID(ctx context.Context, pageID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	var pageRoles []*entity.PageRole
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"pagerole_cache_pageid_" + pageID,
			"pagerole_cache_nextcursor_pageid_" + pageID,
			"pagerole_cache_hasnext_pageid_" + pageID,
			"pagerole_cache_limit_pageid_" + pageID,
		})
		if err == nil && datacache != nil {
			if pageRolesList, ok := datacache.([]*entity.PageRole); ok {
				return &dto.PaginationRes{
					Data:       pageRolesList,
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
	if err = cursorDB.All(ctx, &pageRoles); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(pageRoles) > limit {
		hasNext = true
		pageRoles = pageRoles[:limit] // Lấy đúng số lượng cần thiết
		lastPageRole := pageRoles[len(pageRoles)-1]
		nextCursor = utils.EncodeCursorMongodb(lastPageRole.CreatedAt, lastPageRole.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"pagerole_cache_pageid_" + pageID:            pageRoles,
			"pagerole_cache_nextcursor_pageid_" + pageID: nextCursor,
			"pagerole_cache_hasnext_pageid_" + pageID:    hasNext,
			"pagerole_cache_limit_pageid_" + pageID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       pageRoles,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *PagesRolesRepository) GetPageRolesByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	var pageRoles []*entity.PageRole
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"pagerole_cache_userid_" + userID,
			"pagerole_cache_nextcursor_userid_" + userID,
			"pagerole_cache_hasnext_userid_" + userID,
			"pagerole_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if pageRolesList, ok := datacache.([]*entity.PageRole); ok {
				return &dto.PaginationRes{
					Data:       pageRolesList,
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
	if err = cursorDB.All(ctx, &pageRoles); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(pageRoles) > limit {
		hasNext = true
		pageRoles = pageRoles[:limit] // Lấy đúng số lượng cần thiết
		lastPageRole := pageRoles[len(pageRoles)-1]
		nextCursor = utils.EncodeCursorMongodb(lastPageRole.CreatedAt, lastPageRole.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"pagerole_cache_userid_" + userID:            pageRoles,
			"pagerole_cache_nextcursor_userid_" + userID: nextCursor,
			"pagerole_cache_hasnext_userid_" + userID:    hasNext,
			"pagerole_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       pageRoles,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *PagesRolesRepository) GetPageRolesByPageIDAndUserID(ctx context.Context, pageID string, userID string) (*entity.PageRole, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	finalpageID, _ := primitive.ObjectIDFromHex(pageID)
	query := bson.M{
		"page_id": finalpageID,
		"user_id": userID,
	}
	var pageRole entity.PageRole
	err := collection.FindOne(ctx, query).Decode(&pageRole)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err
	}
	return &pageRole, nil
}
func (r *PagesRolesRepository) UpdatePageRole(ctx context.Context, pageRole *entity.PageRole) error {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	filter := bson.M{"_id": pageRole.ID}
	update := bson.M{"$set": pageRole}
	var pageRoleUpdated *entity.PageRole
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, filter, update, options).Decode(&pageRoleUpdated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err
	}
	if pageRoleUpdated == nil {
		return mongo.ErrNoDocuments // Không tìm thấy document nào để cập nhật
	}
	return nil
}
func (r *PagesRolesRepository) UpdateBulkPageRoles(ctx context.Context, pageRoles []*entity.PageRole) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	var models []mongo.WriteModel
	for _, pageRole := range pageRoles {
		filter := bson.M{"_id": pageRole.ID}
		update := bson.M{"$set": pageRole}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}

	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedDocs = append(failedDocs, &mongodbErrors.BulkError{
					ID:     pageRoles[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return result.ModifiedCount, failedDocs, nil
		}
		return 0, nil, err
	}
	return result.ModifiedCount, failedDocs, nil
}
func (r *PagesRolesRepository) DeletePageRole(ctx context.Context, pageRoleID string) error {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(pageRoleID)
	if err != nil {
		return err
	}
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments // Không tìm thấy document nào để xóa
	}
	return nil
}
func (r *PagesRolesRepository) DeleteBulkPageRoles(ctx context.Context, pageRoleIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	var objIDs []primitive.ObjectID
	for _, id := range pageRoleIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, err
		}
		objIDs = append(objIDs, objID)
	}
	result, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	if result == nil && err != nil {
		return 0, nil, err
	}
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedDocs = append(failedDocs, &mongodbErrors.BulkError{
					ID:     pageRoleIDs[we.Index],
					Reason: we.Message,
				})
			}
			return result.DeletedCount, failedDocs, nil
		}
		return 0, nil, err
	}
	return result.DeletedCount, failedDocs, nil
}
func (r *PagesRolesRepository) DeletePageRoleByUserIDandPageID(ctx context.Context, userID string, pageID string) error {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	finalpageID, _ := primitive.ObjectIDFromHex(pageID)
	filter := bson.M{
		"page_id": finalpageID,
		"user_id": userID,
	}
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments // Không tìm thấy document nào để xóa
	}
	return nil
}
func (r *PagesRolesRepository) DeleteBulkPageRolesByUserIDsAndPageID(ctx context.Context, userIDs []string, pageID string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	finalpageID, _ := primitive.ObjectIDFromHex(pageID)
	filter := bson.M{
		"page_id": finalpageID,
		"user_id": bson.M{"$in": userIDs},
	}
	result, err := collection.DeleteMany(ctx, filter)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedDocs = append(failedDocs, &mongodbErrors.BulkError{
					ID:     userIDs[we.Index],
					Reason: we.Message,
				})
			}
			return result.DeletedCount, failedDocs, nil
		}
		return 0, nil, err
	}
	return result.DeletedCount, failedDocs, nil
}
func (r *PagesRolesRepository) DeleteBulkPageRolesByPageIDsAndUserID(ctx context.Context, pageIDs []string, userID string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	var objPageIDs []primitive.ObjectID
	for _, id := range pageIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, err
		}
		objPageIDs = append(objPageIDs, objID)
	}
	filter := bson.M{
		"page_id": bson.M{"$in": objPageIDs},
		"user_id": userID,
	}
	result, err := collection.DeleteMany(ctx, filter)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedDocs = append(failedDocs, &mongodbErrors.BulkError{
					ID:     pageIDs[we.Index],
					Reason: we.Message,
				})
			}
			return result.DeletedCount, failedDocs, nil
		}
		return 0, nil, err
	}
	return result.DeletedCount, failedDocs, nil
}
func (r *PagesRolesRepository) DeletePageRoleByPageID(ctx context.Context, pageID string) error {
	collection := r.client.Collection(entity.PageRole{}.CollectionName())
	parseID, ok := primitive.ObjectIDFromHex(pageID)
	if ok != nil {
		return ok
	}
	query := bson.M{
		"page_id": parseID,
	}
	_, err := collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
