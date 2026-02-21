package mongodb

import (
	"context"
	"errors"
	"fmt"

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

type CommentRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewCommentRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *CommentRepository {
	return &CommentRepository{client: client, redisRepo: redisRepo}
}
func (r *CommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) error {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	comment.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	return nil
}
func (r *CommentRepository) CreateBulkComments(ctx context.Context, comments []*entity.Comment) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	for _, comment := range comments {
		if comment.ID.IsZero() {
			comment.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(comments)
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
					ID:     comments[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *CommentRepository) GetCommentByID(ctx context.Context, commentID string) (*entity.Comment, error) {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	finalid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}
	var comment *entity.Comment
	err = collection.FindOne(ctx, bson.M{"_id": finalid}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return comment, nil
}
func (r *CommentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]*entity.Comment, error) {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	finalid, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}
	var comment []*entity.Comment
	cursor, err := collection.Find(ctx, bson.M{"PostID": finalid})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &comment); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}
	return comment, nil
}
func (r *CommentRepository) GetCommentsByUserID(ctx context.Context, userID string) ([]*entity.Comment, error) {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	var comment []*entity.Comment
	cursor, err := collection.Find(ctx, bson.M{"UserID": userID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy document nào
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &comment); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}
	return comment, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	filter := bson.M{"_id": comment.ID}
	update := bson.M{"$set": comment}
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
func (r *CommentRepository) UpdateBulkComments(ctx context.Context, comments []*entity.Comment) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	models := make([]mongo.WriteModel, 0, len(comments))
	for _, ps := range comments {
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

				if failedIndex < len(comments) {
					failedComment := comments[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedComment.ID.Hex(), // Hoặc failedComment.ID.String()
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
func (r *CommentRepository) DeleteComment(ctx context.Context, commentID string) error {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	finalid, _ := primitive.ObjectIDFromHex(commentID)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *CommentRepository) DeleteBulkComments(ctx context.Context, commentIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range commentIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid comment ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	objectIDs := make([]primitive.ObjectID, 0, len(commentIDs))
	for _, id := range commentIDs {
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
				if failedIndex < len(commentIDs) {
					failedCommentID := commentIDs[failedIndex]
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedCommentID,
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
func (r *CommentRepository) PaginationComments(ctx context.Context, postID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Comment{}.CollectionnamComment())
	var comments []*entity.Comment
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
			if commentsList, ok := datacache.([]*entity.Comment); ok {
				return &dto.PaginationRes{
					Data:       commentsList,
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
			"comment_cache_postid_" + postID:            comments,
			"comment_cache_nextcursor_postid_" + postID: nextCursor,
			"comment_cache_hasnext_postid_" + postID:    hasNext,
			"comment_cache_limit_postid_" + postID:      limit,
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
