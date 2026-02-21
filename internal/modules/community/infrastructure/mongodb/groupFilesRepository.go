package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupFilesRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewGroupFilesRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *GroupFilesRepository {
	return &GroupFilesRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *GroupFilesRepository) CreateGroupFile(ctx context.Context, file *entity.GroupFile) error {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	_, err := collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupFilesRepository) CreateBulkGroupFiles(ctx context.Context, files []entity.GroupFile) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	docs := utils.ToInterfaceSlice(files)
	var model []mongo.WriteModel
	for _, doc := range docs {
		model = append(model, mongo.NewInsertOneModel().SetDocument(doc))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &dto.BulkError{
					ID:     files[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.InsertedCount, faildocs, nil
}
func (r *GroupFilesRepository) GetGroupFileByID(ctx context.Context, id string) (*entity.GroupFile, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var file *entity.GroupFile
	err = collection.FindOne(ctx, primitive.M{"_id": idObj}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return file, nil
}
func (r *GroupFilesRepository) GetGroupFilesByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	var files []*entity.GroupFile
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"file_cache_groupid_" + groupID,
			"file_cache_nextcursor_groupid_" + groupID,
			"file_cache_hasnext_groupid_" + groupID,
			"file_cache_limit_groupid_" + groupID,
		})
		if err == nil && datacache != nil {
			if filesList, ok := datacache.([]*entity.GroupFile); ok {
				return &dto.PaginationRes{
					Data:       filesList,
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
	if err = cursorDB.All(ctx, &files); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(files) > limit {
		hasNext = true
		files = files[:limit] // Lấy đúng số lượng cần thiết
		lastFile := files[len(files)-1]
		nextCursor = utils.EncodeCursorMongodb(lastFile.CreatedAt, lastFile.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"file_cache_groupid_" + groupID:            files,
			"file_cache_nextcursor_groupid_" + groupID: nextCursor,
			"file_cache_hasnext_groupid_" + groupID:    hasNext,
			"file_cache_limit_groupid_" + groupID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       files,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupFilesRepository) GetGroupFilesByGroupIDAndUploaderID(ctx context.Context, groupID string, UploaderID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	var files []*entity.GroupFile
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"file_cache_groupid_" + groupID + "_uploaderid_" + UploaderID,
			"file_cache_nextcursor_groupid_" + groupID + "_uploaderid_" + UploaderID,
			"file_cache_hasnext_groupid_" + groupID + "_uploaderid_" + UploaderID,
			"file_cache_limit_groupid_" + groupID + "_uploaderid_" + UploaderID,
		})
		if err == nil && datacache != nil {
			if filesList, ok := datacache.([]*entity.GroupFile); ok {
				return &dto.PaginationRes{
					Data:       filesList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(groupID)
	query := bson.M{"group_id": finalid, "uploader_id": UploaderID}

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
	if err = cursorDB.All(ctx, &files); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(files) > limit {
		hasNext = true
		files = files[:limit] // Lấy đúng số lượng cần thiết
		lastFile := files[len(files)-1]
		nextCursor = utils.EncodeCursorMongodb(lastFile.CreatedAt, lastFile.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"file_cache_groupid_" + groupID + "_uploaderid_" + UploaderID:            files,
			"file_cache_nextcursor_groupid_" + groupID + "_uploaderid_" + UploaderID: nextCursor,
			"file_cache_hasnext_groupid_" + groupID + "_uploaderid_" + UploaderID:    hasNext,
			"file_cache_limit_groupid_" + groupID + "_uploaderid_" + UploaderID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       files,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupFilesRepository) UpdateGroupFileDownloadCount(ctx context.Context, id string, newCount int) error {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"download_count": newCount}}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": idObj}, update)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupFilesRepository) DeleteGroupFile(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": idObj})
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupFilesRepository) deleteBulkGroupFiles(ctx context.Context, ids []string) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	var model []mongo.WriteModel
	for _, id := range ids {
		idObj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid ID format: %s", id)
		}
		model = append(model, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": idObj}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &dto.BulkError{
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

func (r *GroupFilesRepository) DeleteGroupFilesByGroupID(ctx context.Context, groupID string) error {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	finalid, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = collection.DeleteMany(ctx, bson.M{"group_id": finalid})
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupFilesRepository) DeleteBulkGroupFilesByGroupID(ctx context.Context, groupIDs []string) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	var model []mongo.WriteModel
	for _, groupID := range groupIDs {
		finalid, err := primitive.ObjectIDFromHex(groupID)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid group ID format: %s", groupID)
		}
		model = append(model, mongo.NewDeleteManyModel().SetFilter(bson.M{"group_id": finalid}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &dto.BulkError{
					ID:     groupIDs[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return 0, nil, nil
}
func (r *GroupFilesRepository) DeleteGroupFileByIDAndGroupID(ctx context.Context, id string, groupID string) error {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	groupIDObj, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": idObj, "group_id": groupIDObj})
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupFilesRepository) DeleteBulkGroupFileByIDAndGroupID(ctx context.Context, ids []string, groupID string) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	groupIDObj, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid group ID format: %s", groupID)
	}
	var model []mongo.WriteModel
	for _, id := range ids {
		idObj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid ID format: %s", id)
		}
		model = append(model, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": idObj, "group_id": groupIDObj}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &dto.BulkError{
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

func (r *GroupFilesRepository) DeleteGroupFilesByGroupIDAndUploaderID(ctx context.Context, groupID string, uploaderID string) error {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	groupIDObj, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = collection.DeleteMany(ctx, bson.M{"group_id": groupIDObj, "uploader_id": uploaderID})
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupFilesRepository) DeleteBulkGroupFilesByGroupIDAndUploaderID(ctx context.Context, groupID string, uploaderIDs []string) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.GroupFile{}.CollectionName())
	groupIDObj, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid group ID format: %s", groupID)
	}
	var model []mongo.WriteModel
	for _, uploaderID := range uploaderIDs {
		model = append(model, mongo.NewDeleteManyModel().SetFilter(bson.M{"group_id": groupIDObj, "uploader_id": uploaderID}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &dto.BulkError{
					ID:     uploaderIDs[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.DeletedCount, faildocs, nil
}
