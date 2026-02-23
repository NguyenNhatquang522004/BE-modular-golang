package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostEditLogsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewPostEditLogsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *PostEditLogsRepository {
	return &PostEditLogsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}

func (r *PostEditLogsRepository) CreatePostEditLog(ctx context.Context, postEditLog *entity.PostEntityEditLog) error {
	// Implementation of creating a single post edit log in MongoDB
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	postEditLog.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, postEditLog)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostEditLogsRepository) CreateBulkPostEditLog(ctx context.Context, postEditLogs []*entity.PostEntityEditLog) (int64, []*mongodbErrors.EditLogsBulkError, error) {
	// Implementation of creating multiple post edit logs in MongoDB with bulk operation
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	for _, postEditLog := range postEditLogs {
		if postEditLog.ID.IsZero() {
			postEditLog.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(postEditLogs)
	result, err := collection.InsertMany(ctx, docs, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}

	var failedDocs []*mongodbErrors.EditLogsBulkError
	if err != nil {
		// Kiểm tra xem có phải lỗi BulkWriteException không
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedDocs = append(failedDocs, &mongodbErrors.EditLogsBulkError{
					ID:       postEditLogs[we.Index].TargetID.Hex(),
					EditorID: postEditLogs[we.Index].EditorID,
					TargetID: postEditLogs[we.Index].TargetID.Hex(),
					Reason:   we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *PostEditLogsRepository) GetByTargetID(ctx context.Context, targetID string) (*entity.PostEntityEditLog, error) {
	// Implementation of retrieving a post edit log by target ID from MongoDB
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	var postEditLog = &entity.PostEntityEditLog{}
	filter := bson.M{"target_id": targetID}
	err := collection.FindOne(ctx, filter).Decode(&postEditLog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err // Lỗi khác xảy ra
	}
	return postEditLog, nil
}
func (r *PostEditLogsRepository) GetBulkByTargetID(ctx context.Context, targetIDs []string) ([]*entity.PostEntityEditLog, error) {
	// Implementation of retrieving multiple post edit logs by target IDs from MongoDB
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	var postEditLogs = []*entity.PostEntityEditLog{}
	objectIDs := make([]primitive.ObjectID, 0, len(targetIDs))
	for _, id := range targetIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"post_id": bson.M{"$in": objectIDs}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &postEditLogs); err != nil {
		return nil, err
	}
	return postEditLogs, nil
}
func (r *PostEditLogsRepository) GetByID(ctx context.Context, id string) (*entity.PostEntityEditLog, error) {
	// Implementation of retrieving a post edit log by its ID from MongoDB
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	var postEditLog *entity.PostEntityEditLog
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	err = collection.FindOne(ctx, filter).Decode(&postEditLog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err // Lỗi khác xảy ra
	}
	return postEditLog, nil
}
func (r *PostEditLogsRepository) GetBulkByID(ctx context.Context, ids []string) ([]*entity.PostEntityEditLog, error) {
	// Implementation of retrieving multiple post edit logs by their IDs from MongoDB
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	var postEditLogs []*entity.PostEntityEditLog
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &postEditLogs); err != nil {
		return nil, err
	}
	return postEditLogs, nil
}
func (r *PostEditLogsRepository) GetLatestVersion(ctx context.Context, targetID string) (int, error) {
	// Implementation of retrieving the latest version number of a post edit log for a given target ID from MongoDB
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	finalid, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return 0, err
	}
	filter := bson.M{"target_id": finalid}
	opts := options.FindOne().SetSort(bson.D{{"version", -1}})
	var postEditLog *entity.PostEntityEditLog
	err = collection.FindOne(ctx, filter, opts).Decode(&postEditLog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil // Không tìm thấy document nào
		}
		return 0, err // Lỗi khác xảy ra
	}
	return postEditLog.Version, nil
}
func (r *PostEditLogsRepository) GetAllVersions(ctx context.Context, targetID string) ([]*entity.PostEntityEditLog, error) {
	// Implementation of retrieving all versions of post edit logs for a given target ID from MongoDB
	// Implementation of retrieving the latest version number of a post edit log for a given target ID from MongoDB
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	finalid, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"target_id": finalid}
	opts := options.Find().SetSort(bson.D{{"version", -1}})
	var postEditLog = []*entity.PostEntityEditLog{}
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err // Lỗi khác xảy ra
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &postEditLog); err != nil {
		return nil, err
	}
	return postEditLog, nil
}
func (r *PostEditLogsRepository) UpdatePostEditLog(ctx context.Context, postEditLog *entity.PostEntityEditLog) error {
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	filter := bson.M{"post_id": postEditLog.TargetID, "version": postEditLog.Version}
	update := bson.M{"$set": postEditLog}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedPostEditLog *entity.PostEntityEditLog
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPostEditLog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *PostEditLogsRepository) UpdateBulkPostEditLog(ctx context.Context, postEditLogs []*entity.PostEntityEditLog) (int64, []*mongodbErrors.EditLogsBulkError, error) {
	// Implementation of updating multiple post edit logs in MongoDB with bulk operation
	collection := r.client.Collection(entity.PostEntityEditLog{}.Collectionnameposteditlog())
	models := make([]mongo.WriteModel, 0, len(postEditLogs))
	for _, ps := range postEditLogs {
		filter := bson.M{"post_id": ps.TargetID, "version": ps.Version}
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
	var failedDocs []*mongodbErrors.EditLogsBulkError
	if err != nil {
		// Dùng errors.As để ép kiểu err về mongo.BulkWriteException
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			// Duyệt qua danh sách các lỗi
			for _, we := range bulkErr.WriteErrors {
				// we.Index: Là chỉ số (index) trong slice 'posts' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(postEditLogs) {
					failedPost := postEditLogs[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.EditLogsBulkError{
						ID:       failedPost.ID.Hex(), // Hoặc failedPost.ID.String()
						TargetID: failedPost.TargetID.Hex(),
						EditorID: failedPost.EditorID,
						Reason:   we.Message,
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
func (r *PostEditLogsRepository) PaginationPostEditLog(ctx context.Context, targetID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// Implementation of paginating post edit logs for a given post ID from MongoDB

	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	var postExtension []*entity.PostExtension
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"postEditLogs_cache_postid_" + targetID,
			"postEditLogs_cache_nextcursor_postid_" + targetID,
			"postEditLogs_cache_hasnext_postid_" + targetID,
			"postEditLogs_cache_limit_postid_" + targetID,
		})
		if err == nil && datacache != nil {
			if psList, ok := datacache.([]*entity.PostExtension); ok {
				return &dto.PaginationRes{
					Data:       psList,
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
	if err = cursorDB.All(ctx, &postExtension); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(postExtension) > limit {
		hasNext = true
		postExtension = postExtension[:limit] // Lấy đúng số lượng cần thiết
		lastPost := postExtension[len(postExtension)-1]
		nextCursor = utils.EncodeCursorMongodb(lastPost.CreatedAt, lastPost.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"postEditLogs_cache_postid_" + targetID:            postExtension,
			"postEditLogs_cache_nextcursor_postid_" + targetID: nextCursor,
			"postEditLogs_cache_hasnext_postid_" + targetID:    hasNext,
			"postEditLogs_cache_limit_postid_" + targetID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       postExtension,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
