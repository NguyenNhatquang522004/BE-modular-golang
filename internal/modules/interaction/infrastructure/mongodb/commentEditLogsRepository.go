package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentEditLogsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewCommentEditLogsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *CommentEditLogsRepository {
	return &CommentEditLogsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *CommentEditLogsRepository) CreateEditLog(ctx context.Context, editlog *entity.EntityEditLog) error {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	editlog.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, editlog)
	if err != nil {
		return err
	}
	return nil
}
func (r *CommentEditLogsRepository) CreateBulkEditLogs(ctx context.Context, editLogs []*entity.EntityEditLog) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	for _, editLog := range editLogs {
		if editLog.ID.IsZero() {
			editLog.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(editLogs)
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
					ID:     editLogs[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *CommentEditLogsRepository) GetEditLogByID(ctx context.Context, editLogID string) (*entity.EntityEditLog, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	var editLog entity.EntityEditLog
	objID, err := primitive.ObjectIDFromHex(editLogID)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&editLog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy document, trả về nil
		}
		return nil, err // Lỗi khác xảy ra
	}
	return nil, nil
}
func (r *CommentEditLogsRepository) GetBulkEditLogsByID(ctx context.Context, editLogIDs []string) ([]*entity.EntityEditLog, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	var editLogs []*entity.EntityEditLog
	var objIDs []primitive.ObjectID
	for _, id := range editLogIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, objID)
	}
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &editLogs); err != nil {
		return nil, err
	}
	return editLogs, nil
}
func (r *CommentEditLogsRepository) GetEditLogsByTargetID(ctx context.Context, targetID string) ([]*entity.EntityEditLog, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	var editLogs []*entity.EntityEditLog
	cursor, err := collection.Find(ctx, bson.M{"target_id": targetID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &editLogs); err != nil {
		return nil, err
	}
	return editLogs, nil
}

func (r *CommentEditLogsRepository) GetLatestEditLogByTargetID(ctx context.Context, targetID string) (*entity.EntityEditLog, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	var editLog entity.EntityEditLog
	err := collection.FindOne(ctx, bson.M{"target_id": targetID}, options.FindOne().SetSort(bson.M{"version": -1})).Decode(&editLog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy document, trả về nil
		}
		return nil, err // Lỗi khác xảy ra
	}
	return &editLog, nil
}
func (r *CommentEditLogsRepository) GetVersionBulkEditLogsByTargetIDs(ctx context.Context, targetIDs []string) ([]*entity.EntityEditLog, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	var editLogs []*entity.EntityEditLog
	cursor, err := collection.Find(ctx, bson.M{"target_id": bson.M{"$in": targetIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &editLogs); err != nil {
		return nil, err
	}
	return editLogs, nil
}
func (r *CommentEditLogsRepository) GetVersionBulkEditLogsByTargetID(ctx context.Context, targetID string) ([]*entity.EntityEditLog, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	var editLogs []*entity.EntityEditLog
	cursor, err := collection.Find(ctx, bson.M{"target_id": targetID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &editLogs); err != nil {
		return nil, err
	}
	return editLogs, nil
}
func (r *CommentEditLogsRepository) UpdateEditLog(ctx context.Context, editlog *entity.EntityEditLog) error {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	filter := bson.M{"comment_id": editlog.ID}
	update := bson.M{"$set": editlog}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedComment *entity.Comment
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedComment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *CommentEditLogsRepository) UpdateBulkEditLogs(ctx context.Context, editLogs []*entity.EntityEditLog) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(editLogs))
	for _, ps := range editLogs {
		filter := bson.M{"comment_id": ps.ID}
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

				if failedIndex < len(editLogs) {
					failedEditLog := editLogs[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedEditLog.ID.Hex(), // Hoặc failedEditLog.ID.String()
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
func (r *CommentEditLogsRepository) DeleteEditLog(ctx context.Context, editLogID string) error {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(editLogID)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *CommentEditLogsRepository) DeleteBulkEditLogs(ctx context.Context, editLogsIds []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range editLogsIds {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid edit log ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	objectIDs := make([]primitive.ObjectID, 0, len(editLogsIds))
	for _, id := range editLogsIds {
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
				if failedIndex < len(editLogsIds) {
					failedEditLogID := editLogsIds[failedIndex]
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedEditLogID,
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
func (r *CommentEditLogsRepository) PaginationEditLogs(ctx context.Context, targetID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	var editLogs []*entity.EntityEditLog
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"comment_edit_log_cache_targetid_" + targetID,
			"comment_edit_log_cache_nextcursor_targetid_" + targetID,
			"comment_edit_log_cache_hasnext_targetid_" + targetID,
			"comment_edit_log_cache_limit_targetid_" + targetID,
		})
		if err == nil && datacache != nil {
			if editLogsList, ok := datacache.([]*entity.EntityEditLog); ok {
				return &dto.PaginationRes{
					Data:       editLogsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(targetID)
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
	if err = cursorDB.All(ctx, &editLogs); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(editLogs) > limit {
		hasNext = true
		editLogs = editLogs[:limit] // Lấy đúng số lượng cần thiết
		lastEditLog := editLogs[len(editLogs)-1]
		nextCursor = utils.EncodeCursorMongodb(lastEditLog.CreatedAt, lastEditLog.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"comment_edit_log_cache_targetid_" + targetID:            editLogs,
			"comment_edit_log_cache_nextcursor_targetid_" + targetID: nextCursor,
			"comment_edit_log_cache_hasnext_targetid_" + targetID:    hasNext,
			"comment_edit_log_cache_limit_targetid_" + targetID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       editLogs,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *CommentEditLogsRepository) DeleteEditLogsByTargetID(ctx context.Context, targetID string) error {
	collection := r.client.Collection(entity.EntityEditLog{}.CollectionName())
	filter := bson.M{"target_id": targetID}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
