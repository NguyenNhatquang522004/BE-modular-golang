package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewPostRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *PostRepository {
	return &PostRepository{client: client, redisRepo: redisRepo}
}

func (r *PostRepository) CreatePost(post *entity.Post) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) CreateBulkPosts(posts []*entity.Post) ([]*entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	for _, post := range posts {
		if post.ID.IsZero() {
			post.ID = primitive.NewObjectID()
		}
		// Đảm bảo các field bắt buộc khác (CreateAt...)
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(posts)
	_, err := collection.InsertMany(ctx, docs, opts)
	if err != nil {
		// Kiểm tra xem có phải lỗi BulkWriteException không
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			// Logic: Lọc ra những bài bị lỗi để trả về cho Usecase xử lý
			// (Ví dụ: báo cho user biết bài nào trùng slug)

			// Tạo map index bị lỗi để tra cứu nhanh
			failedIndexes := make(map[int]string)
			for _, we := range bulkErr.WriteErrors {
				failedIndexes[we.Index] = we.Message
			}

			// In log hoặc xử lý tùy nghiệp vụ
			// Ở đây mình ví dụ: vẫn return nil error (vì đã có cái thành công),
			// nhưng in log những cái thất bại.
			// Hoặc bạn có thể return custom error chứa danh sách failed.
			return posts, fmt.Errorf("inserted with failures: %d docs failed", len(failedIndexes))
		}

		// Lỗi hệ thống nghiêm trọng (Network, Auth...) -> Return lỗi luôn
		return nil, err
	}

	return posts, nil
}
func (r *PostRepository) GetPostByID(postID string) (*entity.Post, error) {
	// Implement the logic to retrieve a post by its ID from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	filter := bson.M{"_id": postID}
	var post entity.Post
	err := collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (r *PostRepository) GetPostsBulkByIDs(postIDs []string) ([]*entity.Post, error) {
	// Implement the logic to retrieve multiple posts by their IDs from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var objIDs []primitive.ObjectID
	for _, id := range postIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, oid)
		}
	}
	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*entity.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}
func (r *PostRepository) GetPostsByUserID(userID string) ([]*entity.Post, error) {
	// Implement the logic to retrieve posts by user ID from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	filter := bson.M{"user_id": userObjID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*entity.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) UpdatePost(post *entity.Post) (*entity.Post, error) {
	// Implement the logic to update a post in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	filter := bson.M{"_id": post.ID}
	update := bson.M{"$set": post,
		"$currentDate": bson.M{"updated_at": true},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return post, nil
}
func (r *PostRepository) UpdateBulkPosts(posts []*entity.Post) (int64, []dto.BulkError, error) {
	// Implement the logic to update multiple posts in MongoDB
	if len(posts) == 0 {
		return 0, nil, nil // Return early if there are no posts to update
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	defer cancel()
	models := make([]mongo.WriteModel, len(posts))
	for i, post := range posts {
		filter := bson.M{"_id": post.ID}
		post.UpdatedAt = time.Now()
		update := bson.M{
			"$set": post,
		}
		models[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update)
	}
	opts := options.BulkWrite().SetOrdered(false)

	result, err := collection.BulkWrite(ctx, models, opts)
	// Trường hợp 1: Lỗi hệ thống (Mất mạng, DB sập...) -> err khác nil, result nil
	if result == nil && err != nil {
		return 0, nil, fmt.Errorf("bulk write system error: %w", err)
	}
	// Trường hợp 2: Có lỗi xảy ra với một vài document (Partial Failure)
	var failedDocs []dto.BulkError
	if err != nil {
		// Dùng errors.As để ép kiểu err về mongo.BulkWriteException
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			// Duyệt qua danh sách các lỗi
			for _, we := range bulkErr.WriteErrors {
				// we.Index: Là chỉ số (index) trong slice 'posts' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(posts) {
					failedPost := posts[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, dto.BulkError{
						PostID: failedPost.ID.Hex(), // Hoặc failedPost.ID.String()
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
func (r *PostRepository) DeletePost(postID string) error {
	// Implement the logic to delete a post by its ID from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	filter := bson.M{"_id": postID}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeleteBulkPosts(postIDs []string) (int64, []dto.BulkError, error) {
	// Implement the logic to delete multiple posts by their IDs from MongoDB
	if len(postIDs) == 0 {
		return 0, nil, nil // Return early if there are no post IDs to delete
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	var objIDs []primitive.ObjectID
	for _, id := range postIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, oid)
		}
	}
	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	result, err := collection.DeleteMany(ctx, filter)
	if result == nil && err != nil {
		return 0, nil, fmt.Errorf("bulk delete system error: %w", err)
	}
	var failedDocs []dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedIndex := we.Index
				if failedIndex < len(postIDs) {
					failedPostID := postIDs[failedIndex]
					failedDocs = append(failedDocs, dto.BulkError{
						PostID: failedPostID,
						Reason: we.Message,
					})
				}
			}
			return result.DeletedCount, failedDocs, err
		}
		return result.DeletedCount, nil, err
	}
	return result.DeletedCount, nil, nil
}
func (r *PostRepository) PanigationPosts(userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// Implement the logic to paginate posts for a user from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.client.Collection(entity.Post{}.CollectionNamePost())
	queryliimit := limit + 1
	filter := bson.M{"user_id": userID}
	if cursor == "" {
		dataCache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"post_cache_user_" + userID,
			"post_cache_nextcursor_user_" + userID,
			"post_cache_hasnext_user_" + userID,
			"post_cache_limit_user_" + userID,
		})
		if err == nil && dataCache != nil {
			return &dto.PaginationRes{
				Data:       dataCache,
				NextCursor: nextcursor,
				HasNext:    hasnext,
				Limit:      limitcache,
			}, nil
		}
	}
	if cursor != "" {
		decodedCursor, err := utils.DecodeCursorMongodb(cursor)

		if err != nil {
			return nil, err
		}
		filter["$or"] = []bson.M{
			{
				"created_at": bson.M{"$lt": decodedCursor.CreatedAt},
			},
			{
				"created_at": decodedCursor.CreatedAt,
				"_id":        bson.M{"$lt": decodedCursor.PostID},
			},
		}
	}
	opts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1}, // Tie-breaker
		}).
		SetLimit(int64(queryliimit))
	data, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx)
	var posts []*entity.Post
	if err := data.All(ctx, &posts); err != nil {
		return nil, err
	}
	hasNext := false
	nextCursor := ""
	if len(posts) > limit {
		hasNext = true

		// Cắt bỏ phần tử dư thừa thứ (limit + 1)
		// Chỉ giữ lại đúng số lượng client yêu cầu
		posts = posts[:limit]

		// Lấy phần tử cuối cùng (của danh sách đã cắt) để tạo cursor
		lastPost := posts[len(posts)-1]
		nextCursor = utils.EncodeCursorMongodb(lastPost.CreatedAt, lastPost.ID)
	}
	items := map[string]any{
		"post_cache_user_" + userID:            posts,
		"post_cache_nextcursor_user_" + userID: nextCursor,
		"post_cache_hasnext_user_" + userID:    hasNext,
		"post_cache_limit_user_" + userID:      limit,
	}
	err = r.redisRepo.CustomizeSetCache(ctx, items)
	if err != nil {
		fmt.Printf("Failed to set cache for pagination: %v\n", err)
	}
	return &dto.PaginationRes{
		Data:       posts,      // Slice bài viết
		NextCursor: nextCursor, // Chuỗi mã hóa cho trang sau
		HasNext:    hasNext,    // True nếu còn trang sau
		Limit:      limit,
	}, nil
}
