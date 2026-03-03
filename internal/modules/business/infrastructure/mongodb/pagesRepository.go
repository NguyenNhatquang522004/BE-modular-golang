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

type PagesRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewPagesRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *PagesRepository {
	return &PagesRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *PagesRepository) CreatePage(ctx context.Context, page *entity.Page) error {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	if page.ID.IsZero() {
		page.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, page)
	return err
}
func (r *PagesRepository) CreateBulkPages(ctx context.Context, pages []*entity.Page) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	if len(pages) == 0 {
		return 0, []*mongodbErrors.BulkError{}, nil
	}
	for _, page := range pages {
		if page.ID.IsZero() {
			page.ID = primitive.NewObjectID()
		}
	}
	docs := utils.ToInterfaceSlice(pages)
	result, err := collection.InsertMany(ctx, docs)
	if result == nil && err != nil {
		return 0, []*mongodbErrors.BulkError{}, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(
					faildocs, &mongodbErrors.BulkError{
						ID:     pages[writeError.Index].ID.Hex(),
						Reason: writeError.Message,
					},
				)
			}
			return int64(len(result.InsertedIDs)), faildocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), faildocs, nil
}
func (r *PagesRepository) GetPageByID(ctx context.Context, pageID string) (*entity.Page, error) {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	finalid, err := primitive.ObjectIDFromHex(pageID)
	if err != nil {
		return nil, errors.New("invalid page ID format")
	}
	var page entity.Page
	err = collection.FindOne(ctx, bson.M{"_id": finalid}).Decode(&page)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("page not found: %w", err)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &page, nil
}
func (r *PagesRepository) GetPageCursorByID(ctx context.Context, pageID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	var pages []*entity.Page
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"page_cache_pageid_" + pageID,
			"page_cache_nextcursor_pageid_" + pageID,
			"page_cache_hasnext_pageid_" + pageID,
			"page_cache_limit_pageid_" + pageID,
		})
		if err == nil && datacache != nil {
			if pagesList, ok := datacache.([]*entity.Page); ok {
				return &dto.PaginationRes{
					Data:       pagesList,
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
	if err = cursorDB.All(ctx, &pages); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(pages) > limit {
		hasNext = true
		pages = pages[:limit] // Lấy đúng số lượng cần thiết
		lastPage := pages[len(pages)-1]
		nextCursor = utils.EncodeCursorMongodb(lastPage.CreatedAt, lastPage.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"page_cache_pageid_" + pageID:            pages,
			"page_cache_nextcursor_pageid_" + pageID: nextCursor,
			"page_cache_hasnext_pageid_" + pageID:    hasNext,
			"page_cache_limit_pageid_" + pageID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       pages,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *PagesRepository) UpdatePage(ctx context.Context, page *entity.Page) error {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	if page.ID.IsZero() {
		return errors.New("page ID is required for update")
	}
	filter := bson.M{"_id": page.ID}
	update := bson.M{"$set": page}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *PagesRepository) UpdateBulkPages(ctx context.Context, pages []*entity.Page) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	if len(pages) == 0 {
		return 0, []*mongodbErrors.BulkError{}, nil
	}
	var models []mongo.WriteModel
	for _, page := range pages {
		if page.ID.IsZero() {
			return 0, nil, errors.New("page ID is required for update")
		}
		filter := bson.M{"_id": page.ID}
		update := bson.M{"$set": page}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, []*mongodbErrors.BulkError{}, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(
					faildocs, &mongodbErrors.BulkError{
						ID:     pages[writeError.Index].ID.Hex(),
						Reason: writeError.Message,
					},
				)
			}
			return result.ModifiedCount, faildocs, nil
		}
		return 0, nil, err
	}
	return result.ModifiedCount, faildocs, nil
}
func (r *PagesRepository) DeletePage(ctx context.Context, pageID string) error {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	finalid, err := primitive.ObjectIDFromHex(pageID)
	if err != nil {
		return errors.New("invalid page ID format")
	}
	result, err := collection.DeleteOne(ctx, bson.M{"_id": finalid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *PagesRepository) DeleteBulkPages(ctx context.Context, pageIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Page{}.CollectionName())
	if len(pageIDs) == 0 {
		return 0, []*mongodbErrors.BulkError{}, nil
	}
	var objectIDs []primitive.ObjectID
	for _, id := range pageIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid page ID format: %s", id)
		}
		objectIDs = append(objectIDs, oid)
	}
	result, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(
					faildocs, &mongodbErrors.BulkError{
						ID:     pageIDs[writeError.Index],
						Reason: writeError.Message,
					},
				)
			}
			return result.DeletedCount, faildocs, nil
		}
		return 0, nil, err
	}
	return result.DeletedCount, nil, nil
}
