package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CallLogsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewCallLogsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *CallLogsRepository {
	return &CallLogsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *CallLogsRepository) CreateCallLog(ctx context.Context, callLog *entity.CallLog) error {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	if callLog.ID.IsZero() {
		callLog.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, callLog)
	if err != nil {
		return err
	}
	return nil
}
func (r *CallLogsRepository) CreateBulkCallLogs(ctx context.Context, callLogs []*entity.CallLog) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	for _, callLog := range callLogs {
		if callLog.ID.IsZero() {
			callLog.ID = primitive.NewObjectID()
		}
	}
	docs := utils.ToInterfaceSlice(callLogs)
	result, err := collection.InsertMany(ctx, docs)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     callLogs[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return int64(len(result.InsertedIDs)), faildocs, nil
}
func (r *CallLogsRepository) GetCallLogByID(ctx context.Context, id string) (*entity.CallLog, error) {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var callLog entity.CallLog
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&callLog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &callLog, nil
}

func (r *CallLogsRepository) GetCallLogsByConversationID(ctx context.Context, conversationID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	var callLogs []*entity.CallLog
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"calllogs_cache_conversationid_" + conversationID,
			"calllogs_cache_nextcursor_conversationid_" + conversationID,
			"calllogs_cache_hasnext_conversationid_" + conversationID,
			"calllogs_cache_limit_conversationid_" + conversationID,
		})
		if err == nil && datacache != nil {
			if callLogsList, ok := datacache.([]*entity.CallLog); ok {
				return &dto.PaginationRes{
					Data:       callLogsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(conversationID)
	query := bson.M{"conversation_id": finalid}

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
	if err = cursorDB.All(ctx, &callLogs); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(callLogs) > limit {
		hasNext = true
		callLogs = callLogs[:limit] // Lấy đúng số lượng cần thiết
		lastCallLog := callLogs[len(callLogs)-1]
		nextCursor = utils.EncodeCursorMongodb(lastCallLog.CreatedAt, lastCallLog.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"calllogs_cache_conversationid_" + conversationID:            callLogs,
			"calllogs_cache_nextcursor_conversationid_" + conversationID: nextCursor,
			"calllogs_cache_hasnext_conversationid_" + conversationID:    hasNext,
			"calllogs_cache_limit_conversationid_" + conversationID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       callLogs,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *CallLogsRepository) UpdateCallLog(ctx context.Context, callLog *entity.CallLog) error {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	result, err := collection.ReplaceOne(ctx, bson.M{"_id": callLog.ID}, callLog)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *CallLogsRepository) UpdateBulkCallLogs(ctx context.Context, callLogs []*entity.CallLog) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(callLogs))
	for _, callLog := range callLogs {
		models = append(models, mongo.NewReplaceOneModel().SetFilter(bson.M{"_id": callLog.ID}).SetReplacement(callLog))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     callLogs[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.ModifiedCount, faildocs, nil
}

func (r *CallLogsRepository) DeleteCallLog(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *CallLogsRepository) DeleteBulkCallLogs(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.CallLog{}.CollectionName())
	var models []mongo.WriteModel
	for i, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid ID at index %d: %w", i, err)
		}
		models = append(models, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": objID}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     ids[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.DeletedCount, faildocs, nil
}
