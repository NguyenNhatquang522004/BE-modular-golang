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

type ReelRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewReelRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *ReelRepository {
	return &ReelRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *ReelRepository) CreateReel(ctx context.Context, reel *entity.Reel) error {
	collection := r.client.Collection(entity.Reel{}.CollectionName())
	reel.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, reel)
	if err != nil {
		return err
	}
	return nil
}
func (r *ReelRepository) CreateBulkReels(ctx context.Context, reels []*entity.Reel) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.Reel{}.CollectionName())
	for _, reel := range reels {
		if reel.ID.IsZero() {
			reel.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(reels)
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
					ID:     reels[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *ReelRepository) GetReelByID(ctx context.Context, id string) (*entity.Reel, error) {
	collection := r.client.Collection(entity.Reel{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var reel entity.Reel
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&reel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err
	}
	return &reel, nil
}

func (r *ReelRepository) GetReelsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Reel{}.CollectionName())
	var reels []*entity.Reel
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"reel_cache_userid_" + userID,
			"reel_cache_nextcursor_userid_" + userID,
			"reel_cache_hasnext_userid_" + userID,
			"reel_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if reelsList, ok := datacache.([]*entity.Reel); ok {
				return &dto.PaginationRes{
					Data:       reelsList,
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
	if err = cursorDB.All(ctx, &reels); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(reels) > limit {
		hasNext = true
		reels = reels[:limit] // Lấy đúng số lượng cần thiết
		lastReel := reels[len(reels)-1]
		nextCursor = utils.EncodeCursorMongodb(lastReel.CreatedAt, lastReel.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"reel_cache_userid_" + userID:            reels,
			"reel_cache_nextcursor_userid_" + userID: nextCursor,
			"reel_cache_hasnext_userid_" + userID:    hasNext,
			"reel_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       reels,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *ReelRepository) UpdateReel(ctx context.Context, reel *entity.Reel) error {
	collection := r.client.Collection(entity.Reel{}.CollectionName())
	filter := bson.M{"_id": reel.ID}
	update := bson.M{"$set": reel}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedReel *entity.Reel
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedReel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *ReelRepository) UpdateBulkReels(ctx context.Context, reels []*entity.Reel) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.Reel{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(reels))
	for _, ps := range reels {
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
	var failedDocs []*dto.BulkError
	if err != nil {
		// Dùng errors.As để ép kiểu err về mongo.BulkWriteException
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			// Duyệt qua danh sách các lỗi
			for _, we := range bulkErr.WriteErrors {
				// we.Index: Là chỉ số (index) trong slice 'comments' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(reels) {
					failedReel := reels[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &dto.BulkError{
						ID:     failedReel.ID.Hex(), // Hoặc failedReel.ID.String()
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
func (r *ReelRepository) DeleteReel(ctx context.Context, id string) error {
		collection := r.client.Collection(entity.Reel{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *ReelRepository) DeleteBulkReels(ctx context.Context, ids []string) (int64, []*dto.BulkError, error) {
		for _, id := range ids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid reel ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.Reel{}.CollectionName())
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
	var failedDocs []*dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedIndex := we.Index
				if failedIndex < len(ids) {
					failedID := ids[failedIndex]
					failedDocs = append(failedDocs, &dto.BulkError{
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
