package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SavedItemsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewSavedItemsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *SavedItemsRepository {
	return &SavedItemsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *SavedItemsRepository) CreateSaveItem(ctx context.Context, saveItem *entity.UserSavedItem) error {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	saveItem.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, saveItem)
	if err != nil {
		return err
	}
	return nil
}
func (r *SavedItemsRepository) CreateBulkSaveItem(ctx context.Context, saveItems []*entity.UserSavedItem) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	for _, saveItem := range saveItems {
		if saveItem.ID.IsZero() {
			saveItem.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(saveItems)
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
					ID:     saveItems[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *SavedItemsRepository) GetSaveItemByID(ctx context.Context, saveItemID string) (*entity.UserSavedItem, error) {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	objID, err := primitive.ObjectIDFromHex(saveItemID)
	if err != nil {
		return nil, err
	}
	var saveItem entity.UserSavedItem
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&saveItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err
	}
	return &saveItem, nil
}
func (r *SavedItemsRepository) GetBulkSaveItemByID(ctx context.Context, saveItemIDs []string) ([]*entity.UserSavedItem, error) {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	var objIDs []primitive.ObjectID
	for _, id := range saveItemIDs {
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
	var saveItems []*entity.UserSavedItem
	if err = cursor.All(ctx, &saveItems); err != nil {
		return nil, err
	}
	return saveItems, nil
}
func (r *SavedItemsRepository) GetSaveItemsByUserID(ctx context.Context, userID string) ([]*entity.UserSavedItem, error) {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	filter := bson.M{"user_id": userID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var saveItems []*entity.UserSavedItem
	if err = cursor.All(ctx, &saveItems); err != nil {
		return nil, err
	}
	return saveItems, nil
}
func (r *SavedItemsRepository) UpdateSaveItem(ctx context.Context, saveItem *entity.UserSavedItem) error {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	filter := bson.M{"_id": saveItem.ID}
	update := bson.M{"$set": saveItem}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedSaveItem *entity.UserSavedItem
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedSaveItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *SavedItemsRepository) UpdateBulkSaveItem(ctx context.Context, saveItems []*entity.UserSavedItem) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	models := make([]mongo.WriteModel, 0, len(saveItems))
	for _, ps := range saveItems {
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
				// we.Index: Là chỉ số (index) trong slice 'saveItems' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(saveItems) {
					failedSaveItem := saveItems[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &dto.BulkError{
						ID:     failedSaveItem.ID.Hex(), // Hoặc failedSaveItem.ID.String()
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
func (r *SavedItemsRepository) DeleteSaveItem(ctx context.Context, saveItemID string) error {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	finalid, _ := primitive.ObjectIDFromHex(saveItemID)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *SavedItemsRepository) DeleteBulkSaveItem(ctx context.Context, saveItemIDs []string) (int64, []*dto.BulkError, error) {
	for _, id := range saveItemIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid save item ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	objectIDs := make([]primitive.ObjectID, 0, len(saveItemIDs))
	for _, id := range saveItemIDs {
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
				if failedIndex < len(saveItemIDs) {
					failedSaveItemID := saveItemIDs[failedIndex]
					failedDocs = append(failedDocs, &dto.BulkError{
						ID:     failedSaveItemID,
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
func (r *SavedItemsRepository) PaginationSaveItem(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.UserSavedItem{}.CollectionnamUserSavedItem())
	var comments []*entity.UserSavedItem
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"saveItem_cache_postid_" + userID,
			"saveItem_cache_nextcursor_postid_" + userID,
			"saveItem_cache_hasnext_postid_" + userID,
			"saveItem_cache_limit_postid_" + userID,
		})
		if err == nil && datacache != nil {
			if commentsList, ok := datacache.([]*entity.UserSavedItem); ok {
				return &dto.PaginationRes{
					Data:       commentsList,
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
	if err = cursorDB.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(comments) > limit {
		hasNext = true
		comments = comments[:limit] // Lấy đúng số lượng cần thiết
		lastComment := comments[len(comments)-1]
		nextCursor = utils.EncodeCursorMongodb(lastComment.CreatedAt, lastComment.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"saveItem_cache_postid_" + userID:            comments,
			"saveItem_cache_nextcursor_postid_" + userID: nextCursor,
			"saveItem_cache_hasnext_postid_" + userID:    hasNext,
			"saveItem_cache_limit_postid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       comments,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
