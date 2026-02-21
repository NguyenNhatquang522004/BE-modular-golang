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

type GroupRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewGroupRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *GroupRepository {
	return &GroupRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *GroupRepository) CreateGroup(ctx context.Context, group *entity.Group) error {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	group.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, group)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupRepository) CreateBulkGroups(ctx context.Context, groups []*entity.Group) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	for _, group := range groups {
		if group.ID.IsZero() {
			group.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(groups)
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
					ID:     groups[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *GroupRepository) GetGroupByID(ctx context.Context, id string) (*entity.Group, error) {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var group entity.Group
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy nhóm
		}
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) GetBulkGroupsByIDs(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	var groups []*entity.Group
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"Group_cache_",
			"Group_cache_nextcursor_",
			"Group_cache_hasnext_",
			"Group_cache_limit_",
		})
		if err == nil && datacache != nil {
			if groupsList, ok := datacache.([]*entity.Group); ok {
				return &dto.PaginationRes{
					Data:       groupsList,
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
	if err = cursorDB.All(ctx, &groups); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(groups) > limit {
		hasNext = true
		groups = groups[:limit] // Lấy đúng số lượng cần thiết
		lastGroup := groups[len(groups)-1]
		nextCursor = utils.EncodeCursorMongodb(lastGroup.CreatedAt, lastGroup.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"group_cache_":            groups,
			"group_cache_nextcursor_": nextCursor,
			"group_cache_hasnext_":    hasNext,
			"group_cache_limit_":      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       groups,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupRepository) GetGroupsByCreatorID(ctx context.Context, creatorID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	var groups []*entity.Group
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"group_cache_creatorid_" + creatorID,
			"group_cache_nextcursor_creatorid_" + creatorID,
			"group_cache_hasnext_creatorid_" + creatorID,
			"group_cache_limit_creatorid_" + creatorID,
		})
		if err == nil && datacache != nil {
			if groupsList, ok := datacache.([]*entity.Group); ok {
				return &dto.PaginationRes{
					Data:       groupsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	query := bson.M{"creator_id": creatorID}

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
	if err = cursorDB.All(ctx, &groups); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(groups) > limit {
		hasNext = true
		groups = groups[:limit] // Lấy đúng số lượng cần thiết
		lastGroup := groups[len(groups)-1]
		nextCursor = utils.EncodeCursorMongodb(lastGroup.CreatedAt, lastGroup.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"group_cache_creatorid_" + creatorID:            groups,
			"group_cache_nextcursor_creatorid_" + creatorID: nextCursor,
			"group_cache_hasnext_creatorid_" + creatorID:    hasNext,
			"group_cache_limit_creatorid_" + creatorID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       groups,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupRepository) UpdateGroup(ctx context.Context, group *entity.Group) error {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	filter := bson.M{"_id": group.ID}
	update := bson.M{"$set": group}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedGroup *entity.Group
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedGroup)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *GroupRepository) UpdateBulkGroups(ctx context.Context, groups []*entity.Group) (int64, []*dto.BulkError, error) {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(groups))
	for _, ps := range groups {
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
				// we.Index: Là chỉ số (index) trong slice 'groups' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(groups) {
					failedGroup := groups[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &dto.BulkError{
						ID:     failedGroup.ID.Hex(), // Hoặc failedGroup.ID.String()
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
func (r *GroupRepository) DeleteGroup(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.Group{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupRepository) DeleteBulkGroups(ctx context.Context, ids []string) (int64, []*dto.BulkError, error) {
	for _, id := range ids {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid group ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.Group{}.CollectionName())
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
	var failedDocs []*dto.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				failedIndex := we.Index
				if failedIndex < len(ids) {
					failedID := ids[failedIndex]
					failedDocs = append(failedDocs, &dto.BulkError{
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
