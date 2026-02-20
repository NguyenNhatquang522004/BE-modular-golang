package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostExtensionRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewPostExtensionRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *PostExtensionRepository {
	return &PostExtensionRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}

func (r *PostExtensionRepository) CreatePostExtension(ctx context.Context, postExtension *entity.PostExtension) error {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	postExtension.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, postExtension)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostExtensionRepository) CreateBulkPostExtension(ctx context.Context, postExtensions []*entity.PostExtension) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	for _, postExtension := range postExtensions {
		if postExtension.ID.IsZero() {
			postExtension.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(postExtensions)
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
					ID:     postExtensions[we.Index].PostID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *PostExtensionRepository) GetByPostID(ctx context.Context, postID string) (*entity.PostExtension, error) {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	var postExtension *entity.PostExtension
	filter := bson.M{"post_id": postID}
	err := collection.FindOne(ctx, filter).Decode(&postExtension)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err // Lỗi khác xảy ra
	}
	return postExtension, nil

}
func (r *PostExtensionRepository) GetBulkByPostIDs(ctx context.Context, postIDs []string) ([]*entity.PostExtension, error) {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	var postExtensions []*entity.PostExtension
	objectIDs := make([]primitive.ObjectID, 0, len(postIDs))
	for _, id := range postIDs {
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
	if err = cursor.All(ctx, &postExtensions); err != nil {
		return nil, err
	}
	return postExtensions, nil
}
func (r *PostExtensionRepository) GetShareDataByPostID(ctx context.Context, postID string) (*entity.ShareData, error) {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	var postExtension *entity.PostExtension
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"post_id": objectID,
		"share_data": bson.M{"$exists": true}} // Chỉ lấy những document có share_data khác null
	opts := options.FindOne().SetProjection(bson.M{
		"share_data": 1,
		"_id":        0,
	})
	err = collection.FindOne(ctx, filter, opts).Decode(&postExtension)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err // Lỗi khác xảy ra
	}
	return postExtension.ShareData, nil
}
func (r *PostExtensionRepository) GetBulkShareDataByPostIDs(ctx context.Context, postIDs []string) ([]*entity.ShareData, error) {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	var postExtensions []*entity.PostExtension
	objectIDs := make([]primitive.ObjectID, 0, len(postIDs))
	for _, id := range postIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"post_id": bson.M{"$in": objectIDs},
		"share_data": bson.M{"$exists": true}} // Chỉ lấy những document có share_data khác null
	opts := options.Find().SetProjection(bson.M{
		"share_data": 1,
		"_id":        0,
	})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &postExtensions); err != nil {
		return nil, err
	}
	shareDataList := make([]*entity.ShareData, 0, len(postExtensions))
	for _, pe := range postExtensions {
		shareDataList = append(shareDataList, pe.ShareData)
	}
	return shareDataList, nil
}
func (r *PostExtensionRepository) UpdatePostExtension(ctx context.Context, postExtension *entity.PostExtension) error {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	filter := bson.M{"post_id": postExtension.PostID}
	update := bson.M{"$set": postExtension}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedPostExtension *entity.PostExtension
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPostExtension)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *PostExtensionRepository) UpdateBulkPostExtension(ctx context.Context, postExtensions []*entity.PostExtension) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	models := make([]mongo.WriteModel, 0, len(postExtensions))
	for _, ps := range postExtensions {
		filter := bson.M{"post_id": ps.PostID}
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
				// we.Index: Là chỉ số (index) trong slice 'posts' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(postExtensions) {
					failedPost := postExtensions[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &dto.BulkError{
						ID:     failedPost.ID.Hex(), // Hoặc failedPost.ID.String()
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
func (r *PostExtensionRepository) DeleteByPostID(ctx context.Context, postID string) error {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	finalid, _ := primitive.ObjectIDFromHex(postID)
	filter := bson.M{"post_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostExtensionRepository) DeleteBulkByPostIDs(ctx context.Context, postIDs []string) (int64, []*dto.BulkError, error) {
	for _, id := range postIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid post ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	objectIDs := make([]primitive.ObjectID, 0, len(postIDs))
	for _, id := range postIDs {
		objID, _ := primitive.ObjectIDFromHex(id)
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"post_id": bson.M{"$in": objectIDs}}
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
				if failedIndex < len(postIDs) {
					failedPostID := postIDs[failedIndex]
					failedDocs = append(failedDocs, &dto.BulkError{
						ID:     failedPostID,
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
func (r *PostExtensionRepository) PaginationPostExtension(ctx context.Context, postID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.PostExtension{}.Collectionnamepostextension())
	var postExtension []*entity.PostExtension
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"postExtension_cache_postid_" + postID,
			"postExtension_cache_nextcursor_postid_" + postID,
			"postExtension_cache_hasnext_postid_" + postID,
			"postExtension_cache_limit_postid_" + postID,
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
			"postExtension_cache_postid_" + postID:            postExtension,
			"postExtension_cache_nextcursor_postid_" + postID: nextCursor,
			"postExtension_cache_hasnext_postid_" + postID:    hasNext,
			"postExtension_cache_limit_postid_" + postID:      limit,
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
