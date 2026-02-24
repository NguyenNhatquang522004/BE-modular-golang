package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MediaAssetsRepository struct {
	// Add necessary fields for MongoDB connection and collection handling
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

// Implement methods for MediaAssetsRepository here, ensuring they satisfy the IMediaAssetsRepository interface defined in the domain layer.
func NewMediaAssetsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *MediaAssetsRepository {
	return &MediaAssetsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *MediaAssetsRepository) CreateMediaAsset(ctx context.Context, asset *entity.MediaAsset) error {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	asset.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, asset)
	if err != nil {
		return err
	}
	return nil
}

func (r *MediaAssetsRepository) CreateBulkMediaAssets(ctx context.Context, assets []*entity.MediaAsset) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	for _, asset := range assets {
		if asset.ID.IsZero() {
			asset.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(assets)
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
					ID:     assets[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *MediaAssetsRepository) GetMediaAssetByID(ctx context.Context, id string) (*entity.MediaAsset, error) {
	// Implementation for retrieving a media asset by ID from MongoDB
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	var asset *entity.MediaAsset
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&asset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy tài liệu, trả về nil
		}
		return nil, err // Lỗi khác xảy ra
	}
	return asset, err
}
func (r *MediaAssetsRepository) GetMediaAssetsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	var mediaAssets []*entity.MediaAsset
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"media_assets_cache_userid_" + userID,
			"media_assets_cache_nextcursor_userid_" + userID,
			"media_assets_cache_hasnext_userid_" + userID,
			"media_assets_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if mediaAssetsList, ok := datacache.([]*entity.MediaAsset); ok {
				return &dto.PaginationRes{
					Data:       mediaAssetsList,
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
	if err = cursorDB.All(ctx, &mediaAssets); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(mediaAssets) > limit {
		hasNext = true
		mediaAssets = mediaAssets[:limit] // Lấy đúng số lượng cần thiết
		lastMediaAsset := mediaAssets[len(mediaAssets)-1]
		nextCursor = utils.EncodeCursorMongodb(lastMediaAsset.CreatedAt, lastMediaAsset.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"media_assets_cache_userid_" + userID:            mediaAssets,
			"media_assets_cache_nextcursor_userid_" + userID: nextCursor,
			"media_assets_cache_hasnext_userid_" + userID:    hasNext,
			"media_assets_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       mediaAssets,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *MediaAssetsRepository) GetMediaAssetsByAlbumID(ctx context.Context, albumID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	var mediaAssets []*entity.MediaAsset
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"media_assets_cache_albumid_" + albumID,
			"media_assets_cache_nextcursor_albumid_" + albumID,
			"media_assets_cache_hasnext_albumid_" + albumID,
			"media_assets_cache_limit_albumid_" + albumID,
		})
		if err == nil && datacache != nil {
			if mediaAssetsList, ok := datacache.([]*entity.MediaAsset); ok {
				return &dto.PaginationRes{
					Data:       mediaAssetsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(albumID)
	query := bson.M{"album_id": finalid}

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
	if err = cursorDB.All(ctx, &mediaAssets); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(mediaAssets) > limit {
		hasNext = true
		mediaAssets = mediaAssets[:limit] // Lấy đúng số lượng cần thiết
		lastMediaAsset := mediaAssets[len(mediaAssets)-1]
		nextCursor = utils.EncodeCursorMongodb(lastMediaAsset.CreatedAt, lastMediaAsset.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"media_assets_cache_albumid_" + albumID:            mediaAssets,
			"media_assets_cache_nextcursor_albumid_" + albumID: nextCursor,
			"media_assets_cache_hasnext_albumid_" + albumID:    hasNext,
			"media_assets_cache_limit_albumid_" + albumID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       mediaAssets,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *MediaAssetsRepository) GetMediaAssetsByPostID(ctx context.Context, postID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	var mediaAssets []*entity.MediaAsset
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"comment_cache_postid_" + postID,
			"comment_cache_nextcursor_postid_" + postID,
			"comment_cache_hasnext_postid_" + postID,
			"comment_cache_limit_postid_" + postID,
		})
		if err == nil && datacache != nil {
			if mediaAssetsList, ok := datacache.([]*entity.MediaAsset); ok {
				return &dto.PaginationRes{
					Data:       mediaAssetsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(postID)
	query := bson.M{"post_id": finalid}

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
	if err = cursorDB.All(ctx, &mediaAssets); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(mediaAssets) > limit {
		hasNext = true
		mediaAssets = mediaAssets[:limit] // Lấy đúng số lượng cần thiết
		lastMediaAsset := mediaAssets[len(mediaAssets)-1]
		nextCursor = utils.EncodeCursorMongodb(lastMediaAsset.CreatedAt, lastMediaAsset.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"media_assets_cache_postid_" + postID:            mediaAssets,
			"media_assets_cache_nextcursor_postid_" + postID: nextCursor,
			"media_assets_cache_hasnext_postid_" + postID:    hasNext,
			"media_assets_cache_limit_postid_" + postID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       mediaAssets,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *MediaAssetsRepository) GetMediaAssetsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	var mediaAssets []*entity.MediaAsset
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"media_assets_cache_groupid_" + groupID,
			"media_assets_cache_nextcursor_groupid_" + groupID,
			"media_assets_cache_hasnext_groupid_" + groupID,
			"media_assets_cache_limit_groupid_" + groupID,
		})
		if err == nil && datacache != nil {
			if mediaAssetsList, ok := datacache.([]*entity.MediaAsset); ok {
				return &dto.PaginationRes{
					Data:       mediaAssetsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(groupID)
	query := bson.M{"group_id": finalid}

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
	if err = cursorDB.All(ctx, &mediaAssets); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(mediaAssets) > limit {
		hasNext = true
		mediaAssets = mediaAssets[:limit] // Lấy đúng số lượng cần thiết
		lastMediaAsset := mediaAssets[len(mediaAssets)-1]
		nextCursor = utils.EncodeCursorMongodb(lastMediaAsset.CreatedAt, lastMediaAsset.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"media_assets_cache_groupid_" + groupID:            mediaAssets,
			"media_assets_cache_nextcursor_groupid_" + groupID: nextCursor,
			"media_assets_cache_hasnext_groupid_" + groupID:    hasNext,
			"media_assets_cache_limit_groupid_" + groupID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       mediaAssets,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *MediaAssetsRepository) UpdateMediaAsset(ctx context.Context, asset *entity.MediaAsset) error {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	filter := bson.M{"_id": asset.ID}
	update := bson.M{"$set": asset}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedAsset *entity.MediaAsset
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedAsset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *MediaAssetsRepository) UpdateBulkMediaAssets(ctx context.Context, assets []*entity.MediaAsset) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(assets))
	for _, ps := range assets {
		filter := bson.M{"_id": ps.ID}
		update := bson.M{"$set": ps, "$currentDate": bson.M{"created_at": true}} // Cập nhật trường created_at
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(false))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	// Trường hợp 1: Lỗi hệ thống (Mất mạng, DB sập...) -> err khác nil, result nil
	if result == nil && err != nil {
		return 0, nil, fmt.Errorf("bulk write system error: %w", err)
	}
	// Trường hợp 2: Có lỗi xảy ra với một vài document (Partial Failure)
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		// Dùng errors.As để ép kiểu err về mongo.BulkWriteException
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			// Duyệt qua danh sách các lỗi
			for _, we := range bulkErr.WriteErrors {
				// we.Index: Là chỉ số (index) trong slice 'comments' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(assets) {
					failedAsset := assets[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedAsset.ID.Hex(), // Hoặc failedAsset.ID.String()
						Reason: we.Message,
					})
				}
			}
			// Mặc dù có err, nhưng BulkWrite vẫn có thể thành công một phần
			// Bạn có thể return err hoặc nil tùy logic nghiệp vụ
			// Ở đây mình return err để tầng trên biết là "không thành công 100%"
			return result.ModifiedCount, failedDocs, err
		}

		// Nếu lỗi không phải BulkWriteException (ví dụ context deadline exceeded)
		return result.ModifiedCount, nil, err
	}
	return result.ModifiedCount, nil, nil
}
func (r *MediaAssetsRepository) DeleteMediaAsset(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *MediaAssetsRepository) DeleteBulkMediaAssets(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range ids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid media asset ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, _ := primitive.ObjectIDFromHex(id)
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	result, err := collection.DeleteMany(ctx, filter)
	if result == nil && err != nil {
		return 0, nil, fmt.Errorf("bulk delete system error: %w", err)
	}
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedIndex := we.Index
				if failedIndex < len(ids) {
					failedID := ids[failedIndex]
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedID,
						Reason: we.Message,
					})
				}
			}
			return result.DeletedCount, failedDocs, nil

		}
		return 0, nil, err
	}
	return result.DeletedCount, nil, nil
}
func (r *MediaAssetsRepository) DeleteMediaAssetsByPostIDAndUserID(ctx context.Context, postID string, userID string) error {
	collection := r.client.Collection(entity.MediaAsset{}.CollectionName())
	finalPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return fmt.Errorf("invalid post ID format: %w", err)
	}
	filter := bson.M{"post_id": finalPostID, "user_id": userID}
	_, err = collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
