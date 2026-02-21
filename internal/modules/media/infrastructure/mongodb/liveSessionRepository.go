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

type LiveSessionRepository struct {
	// Add necessary fields for MongoDB connection and collection handling
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

// Implement methods for LiveSessionRepository here, ensuring they satisfy the ILiveSessionRepository interface defined in the domain layer.
func NewLiveSessionRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *LiveSessionRepository {
	return &LiveSessionRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *LiveSessionRepository) CreateLiveSession(ctx context.Context, liveSession *entity.LiveSession) error {
	// Implementation for creating a live session in MongoDB
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	liveSession.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, liveSession)
	if err != nil {
		return err
	}
	return nil

}

func (r *LiveSessionRepository) CreateBulkLiveSessions(ctx context.Context, liveSessions []*entity.LiveSession) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	for _, liveSession := range liveSessions {
		if liveSession.ID.IsZero() {
			liveSession.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(liveSessions)
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
					ID:     liveSessions[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *LiveSessionRepository) GetLiveSessionByID(ctx context.Context, id string) (*entity.LiveSession, error) {
	// Implementation for retrieving a live session by ID from MongoDB
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var liveSession entity.LiveSession
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&liveSession)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err
	}
	return &liveSession, nil
}
func (r *LiveSessionRepository) GetLiveSessionsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	var liveSessions []*entity.LiveSession
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"live_session_cache_userid_" + userID,
			"live_session_cache_nextcursor_userid_" + userID,
			"live_session_cache_hasnext_userid_" + userID,
			"live_session_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if liveSessionsList, ok := datacache.([]*entity.LiveSession); ok {
				return &dto.PaginationRes{
					Data:       liveSessionsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(userID)
	query := bson.M{"user_id": finalid}

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
	if err = cursorDB.All(ctx, &liveSessions); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(liveSessions) > limit {
		hasNext = true
		liveSessions = liveSessions[:limit] // Lấy đúng số lượng cần thiết
		lastLiveSession := liveSessions[len(liveSessions)-1]
		nextCursor = utils.EncodeCursorMongodb(lastLiveSession.CreatedAt, lastLiveSession.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"live_session_cache_userid_" + userID:            liveSessions,
			"live_session_cache_nextcursor_userid_" + userID: nextCursor,
			"live_session_cache_hasnext_userid_" + userID:    hasNext,
			"live_session_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       liveSessions,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *LiveSessionRepository) GetLiveSessionsByCategoryID(ctx context.Context, categoryID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	var liveSessions []*entity.LiveSession
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"live_session_cache_categoryid_" + categoryID,
			"live_session_cache_nextcursor_categoryid_" + categoryID,
			"live_session_cache_hasnext_categoryid_" + categoryID,
			"live_session_cache_limit_categoryid_" + categoryID,
		})
		if err == nil && datacache != nil {
			if liveSessionsList, ok := datacache.([]*entity.LiveSession); ok {
				return &dto.PaginationRes{
					Data:       liveSessionsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(categoryID)
	query := bson.M{"category_id": finalid}

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
	if err = cursorDB.All(ctx, &liveSessions); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(liveSessions) > limit {
		hasNext = true
		liveSessions = liveSessions[:limit] // Lấy đúng số lượng cần thiết
		lastLiveSession := liveSessions[len(liveSessions)-1]
		nextCursor = utils.EncodeCursorMongodb(lastLiveSession.CreatedAt, lastLiveSession.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"live_session_cache_categoryid_" + categoryID:            liveSessions,
			"live_session_cache_nextcursor_categoryid_" + categoryID: nextCursor,
			"live_session_cache_hasnext_categoryid_" + categoryID:    hasNext,
			"live_session_cache_limit_categoryid_" + categoryID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       liveSessions,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *LiveSessionRepository) UpdateLiveSession(ctx context.Context, liveSession *entity.LiveSession) error {
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	filter := bson.M{"_id": liveSession.ID}
	update := bson.M{"$set": liveSession}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedLiveSession *entity.LiveSession
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedLiveSession)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *LiveSessionRepository) UpdateBulkLiveSessions(ctx context.Context, liveSessions []*entity.LiveSession) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(liveSessions))
	for _, ps := range liveSessions {
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
				// we.Index: Là chỉ số (index) trong slice 'liveSessions' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(liveSessions) {
					failedLiveSession := liveSessions[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedLiveSession.ID.Hex(), // Hoặc failedLiveSession.ID.String()
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
func (r *LiveSessionRepository) DeleteLiveSession(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *LiveSessionRepository) DeleteBulkLiveSessions(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range ids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid live session ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
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
func (r *LiveSessionRepository) CheckLiveSessionExists(ctx context.Context, id string) (bool, error) {
	// Implementation for checking if a live session exists in MongoDB
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	count, err := collection.CountDocuments(ctx, bson.M{"_id": objID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *LiveSessionRepository) CheckLiveSessionExistsByStreamKey(ctx context.Context, HostUserID string, streamKey string) (bool, error) {
	// Implementation for checking if a live session exists by stream key in MongoDB
	collection := r.client.Collection(entity.LiveSession{}.CollectionName())
	count, err := collection.CountDocuments(ctx, bson.M{"host_user_id": HostUserID, "stream_key": streamKey})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
