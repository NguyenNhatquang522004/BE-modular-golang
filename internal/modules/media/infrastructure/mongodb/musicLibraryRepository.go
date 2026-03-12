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

type MusicLibraryRepository struct {
	// Add necessary fields for MongoDB connection and collection handling
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

// Implement methods for MusicLibraryRepository here, ensuring they satisfy the IMusicLibraryRepository interface defined in the domain layer.
func NewMusicLibraryRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *MusicLibraryRepository {
	return &MusicLibraryRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *MusicLibraryRepository) CreateMusicLibrary(ctx context.Context, musicLibrary *entity.MusicLibrary) error {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	musicLibrary.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, musicLibrary)
	if err != nil {
		return err
	}
	return nil
}
func (r *MusicLibraryRepository) CreateBulkMusicLibraries(ctx context.Context, musicLibraries []*entity.MusicLibrary) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	for _, musicLibrary := range musicLibraries {
		if musicLibrary.ID.IsZero() {
			musicLibrary.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(musicLibraries)
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
					ID:     musicLibraries[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *MusicLibraryRepository) GetMusicLibraryByID(ctx context.Context, id string) (*entity.MusicLibrary, error) {
	// Implement the logic to retrieve a MusicLibrary document by its ID from MongoDB
	// Return the found MusicLibrary and any error encountered
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var musicLibrary entity.MusicLibrary
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&musicLibrary)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &musicLibrary, nil
}

func (r *MusicLibraryRepository) GetMusicLibraries(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	var musicLibraries []*entity.MusicLibrary
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"MusicLibrary_cache_postid_",
			"MusicLibrary_cache_nextcursor_postid_",
			"MusicLibrary_cache_hasnext_postid_",
			"MusicLibrary_cache_limit_postid_",
		})
		if err == nil && datacache != nil {
			if musicLibrariesList, ok := datacache.([]*entity.MusicLibrary); ok {
				return &dto.PaginationRes{
					Data:       musicLibrariesList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	query := bson.M{}
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
	if err = cursorDB.All(ctx, &musicLibraries); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(musicLibraries) > limit {
		hasNext = true
		musicLibraries = musicLibraries[:limit] // Lấy đúng số lượng cần thiết
		lastMusicLibrary := musicLibraries[len(musicLibraries)-1]
		nextCursor = utils.EncodeCursorMongodb(lastMusicLibrary.CreatedAt, lastMusicLibrary.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {

		items := map[string]any{
			"MusicLibrary_cache_postid_":            musicLibraries,
			"MusicLibrary_cache_nextcursor_postid_": nextCursor,
			"MusicLibrary_cache_hasnext_postid_":    hasNext,
			"MusicLibrary_cache_limit_postid_":      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       musicLibraries,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *MusicLibraryRepository) GetMusicLibrariesByArtist(ctx context.Context, artist string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	var musicLibraries []*entity.MusicLibrary
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"musicLibrary_cache_artist_" + artist,
			"musicLibrary_cache_nextcursor_artist_" + artist,
			"musicLibrary_cache_hasnext_artist_" + artist,
			"musicLibrary_cache_limit_artist_" + artist,
		})
		if err == nil && datacache != nil {
			if musicLibrariesList, ok := datacache.([]*entity.MusicLibrary); ok {
				return &dto.PaginationRes{
					Data:       musicLibrariesList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	query := bson.M{"artist": artist}

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
	if err = cursorDB.All(ctx, &musicLibraries); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(musicLibraries) > limit {
		hasNext = true
		musicLibraries = musicLibraries[:limit] // Lấy đúng số lượng cần thiết
		lastMusicLibrary := musicLibraries[len(musicLibraries)-1]
		nextCursor = utils.EncodeCursorMongodb(lastMusicLibrary.CreatedAt, lastMusicLibrary.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"musicLibrary_cache_artist_" + artist:            musicLibraries,
			"musicLibrary_cache_nextcursor_artist_" + artist: nextCursor,
			"musicLibrary_cache_hasnext_artist_" + artist:    hasNext,
			"musicLibrary_cache_limit_artist_" + artist:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       musicLibraries,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *MusicLibraryRepository) UpdateMusicLibrary(ctx context.Context, musicLibrary *entity.MusicLibrary) error {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	filter := bson.M{"_id": musicLibrary.ID}
	update := bson.M{"$set": musicLibrary}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedMusicLibrary *entity.MusicLibrary
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedMusicLibrary)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *MusicLibraryRepository) UpdateBulkMusicLibraries(ctx context.Context, musicLibraries []*entity.MusicLibrary) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(musicLibraries))
	for _, ps := range musicLibraries {
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
				// we.Index: Là chỉ số (index) trong slice 'musicLibraries' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(musicLibraries) {
					failedMusicLibrary := musicLibraries[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedMusicLibrary.ID.Hex(), // Hoặc failedMusicLibrary.ID.String()
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
func (r *MusicLibraryRepository) DeleteMusicLibrary(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *MusicLibraryRepository) DeleteBulkMusicLibraries(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range ids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid music library ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
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

func (r *MusicLibraryRepository) DeleteMusicLibraryByArtist(ctx context.Context, artistID string) error {
	collection := r.client.Collection(entity.MusicLibrary{}.CollectionName())
	filter := bson.M{"artist_id": artistID}
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
