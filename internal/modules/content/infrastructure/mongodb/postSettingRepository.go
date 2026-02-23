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

type PostSettingRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewPostSettingRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *PostSettingRepository {
	return &PostSettingRepository{client: client, redisRepo: redisRepo}
}
func (r *PostSettingRepository) CreatePostSetting(ctx context.Context, postSetting *entity.PostSetting) (*entity.PostSetting, error) {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	if postSetting.ID.IsZero() {
		postSetting.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, postSetting)
	if err != nil {
		return nil, err
	}
	return postSetting, nil
}
func (r *PostSettingRepository) CreateBulkPostSetting(ctx context.Context, postSettings []*entity.PostSetting) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	for _, postlist := range postSettings {
		if postlist.ID.IsZero() {
			postlist.ID = primitive.NewObjectID()
		}
	}

	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(postSettings)
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
					ID:     postSettings[we.Index].PostID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return 0, nil, nil

}

func (r *PostSettingRepository) GetPostSettingByPostID(ctx context.Context, Postid string) (*entity.PostSetting, error) {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	var postSetting *entity.PostSetting
	filter := bson.M{"post_id": Postid}
	err := collection.FindOne(ctx, filter).Decode(&postSetting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, err // Lỗi khác xảy ra
	}
	return postSetting, nil
}
func (r *PostSettingRepository) GetPostSettingBulkByPostID(ctx context.Context, Postids []string) ([]*entity.PostSetting, error) {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	var postSettings []*entity.PostSetting
	objectIDs := make([]primitive.ObjectID, 0, len(Postids))
	for _, id := range Postids {
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
	if err = cursor.All(ctx, &postSettings); err != nil {
		return nil, err
	}
	return postSettings, nil
}
func (r *PostSettingRepository) UpdatePostSetting(ctx context.Context, Postid string, postSetting *entity.PostSetting) (*entity.PostSetting, error) {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	filter := bson.M{"post_id": Postid}
	update := bson.M{"$set": postSetting}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedPostSetting *entity.PostSetting
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedPostSetting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào để cập nhật
		}
		return nil, err // Lỗi khác xảy ra
	}
	return updatedPostSetting, nil
}
func (r *PostSettingRepository) UpdateBulkPostSettings(ctx context.Context, postSettings []*entity.PostSetting) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	models := make([]mongo.WriteModel, 0, len(postSettings))
	for _, ps := range postSettings {
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
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		// Dùng errors.As để ép kiểu err về mongo.BulkWriteException
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			// Duyệt qua danh sách các lỗi
			for _, we := range bulkErr.WriteErrors {
				// we.Index: Là chỉ số (index) trong slice 'posts' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(postSettings) {
					failedPost := postSettings[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
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
func (r *PostSettingRepository) DeletePostSetting(ctx context.Context, Postid string) error {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	filter := bson.M{"post_id": Postid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostSettingRepository) DeleteBulkPostSettings(ctx context.Context, postids []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range postids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid post ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	objectIDs := make([]primitive.ObjectID, 0, len(postids))
	for _, id := range postids {
		objID, _ := primitive.ObjectIDFromHex(id)
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"post_id": bson.M{"$in": objectIDs}}
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
				if failedIndex < len(postids) {
					failedPostID := postids[failedIndex]
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
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

func (r *PostSettingRepository) PanigationPostSettings(ctx context.Context, postid string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.PostSetting{}.CollectionNamePostsetting())
	var postSettings []*entity.PostSetting
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"postSetting_cache_postid_" + postid,
			"postSetting_cache_nextcursor_postid_" + postid,
			"postSetting_cache_hasnext_postid_" + postid,
			"postSetting_cache_limit_postid_" + postid,
		})
		if err == nil && datacache != nil {
			if psList, ok := datacache.([]*entity.PostSetting); ok {
				return &dto.PaginationRes{
					Data:       psList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	query := bson.M{"post_id": postid}

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
	if err = cursorDB.All(ctx, &postSettings); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(postSettings) > limit {
		hasNext = true
		postSettings = postSettings[:limit] // Lấy đúng số lượng cần thiết
		lastPost := postSettings[len(postSettings)-1]
		nextCursor = utils.EncodeCursorMongodb(lastPost.CreatedAt, lastPost.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"postSetting_cache_postid_" + postid:            postSettings,
			"postSetting_cache_nextcursor_postid_" + postid: nextCursor,
			"postSetting_cache_hasnext_postid_" + postid:    hasNext,
			"postSetting_cache_limit_postid_" + postid:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}

	return &dto.PaginationRes{
		Data:       postSettings,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
