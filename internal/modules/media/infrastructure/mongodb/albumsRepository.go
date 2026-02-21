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

type AlbumsRepository struct {
	// Add necessary fields for MongoDB connection and collection handling
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

// Implement methods for AlbumsRepository here, ensuring they satisfy the IAlbumsRepository interface defined in the domain layer.
func NewAlbumsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *AlbumsRepository {
	return &AlbumsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}

func (r *AlbumsRepository) CreateAlbum(ctx context.Context, album *entity.Album) error {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	album.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, album)
	if err != nil {
		return err
	}
	return nil
}
func (r *AlbumsRepository) CreateBulkAlbums(ctx context.Context, albums []*entity.Album) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	for _, album := range albums {
		if album.ID.IsZero() {
			album.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(albums)
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
					ID:     albums[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *AlbumsRepository) GetAlbumByID(ctx context.Context, id string) (*entity.Album, error) {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var album *entity.Album
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&album)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (r *AlbumsRepository) GetAlbumsByIDs(ctx context.Context, ids []string) ([]*entity.Album, error) {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, objID)
	}

	cursor, err := collection.Find(ctx, primitive.M{"_id": primitive.M{"$in": objIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var albums []*entity.Album
	if err = cursor.All(ctx, &albums); err != nil {
		return nil, err
	}

	return albums, nil
}
func (r *AlbumsRepository) GetAlbumsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	var albums []*entity.Album
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"album_cache_userid_" + userID,
			"album_cache_nextcursor_userid_" + userID,
			"album_cache_hasnext_userid_" + userID,
			"album_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if albumsList, ok := datacache.([]*entity.Album); ok {
				return &dto.PaginationRes{
					Data:       albumsList,
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
	if err = cursorDB.All(ctx, &albums); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(albums) > limit {
		hasNext = true
		albums = albums[:limit] // Lấy đúng số lượng cần thiết
		lastAlbum := albums[len(albums)-1]
		nextCursor = utils.EncodeCursorMongodb(lastAlbum.CreatedAt, lastAlbum.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"album_cache_userid_" + userID:            albums,
			"album_cache_nextcursor_userid_" + userID: nextCursor,
			"album_cache_hasnext_userid_" + userID:    hasNext,
			"album_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       albums,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *AlbumsRepository) GetAlbumsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	var albums []*entity.Album
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"album_cache_groupid_" + groupID,
			"album_cache_nextcursor_groupid_" + groupID,
			"album_cache_hasnext_groupid_" + groupID,
			"album_cache_limit_groupid_" + groupID,
		})
		if err == nil && datacache != nil {
			if albumsList, ok := datacache.([]*entity.Album); ok {
				return &dto.PaginationRes{
					Data:       albumsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(groupID)
	query := bson.M{"group_id": finalid}

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
	if err = cursorDB.All(ctx, &albums); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(albums) > limit {
		hasNext = true
		albums = albums[:limit] // Lấy đúng số lượng cần thiết
		lastAlbum := albums[len(albums)-1]
		nextCursor = utils.EncodeCursorMongodb(lastAlbum.CreatedAt, lastAlbum.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"album_cache_groupid_" + groupID:            albums,
			"album_cache_nextcursor_groupid_" + groupID: nextCursor,
			"album_cache_hasnext_groupid_" + groupID:    hasNext,
			"album_cache_limit_groupid_" + groupID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       albums,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *AlbumsRepository) UpdateAlbum(ctx context.Context, album *entity.Album) error {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	filter := bson.M{"_id": album.ID}
	update := bson.M{"$set": album}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedAlbum *entity.Album
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedAlbum)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *AlbumsRepository) UpdateBulkAlbums(ctx context.Context, albums []*entity.Album) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(albums))
	for _, ps := range albums {
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
				// we.Index: Là chỉ số (index) trong slice 'albums' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(albums) {
					failedAlbum := albums[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedAlbum.ID.Hex(), // Hoặc failedAlbum.ID.String()
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
func (r *AlbumsRepository) DeleteAlbum(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.Album{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *AlbumsRepository) DeleteBulkAlbums(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range ids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid album ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.Album{}.CollectionName())
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
