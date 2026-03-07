package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/community/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupMembersRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewGroupMembersRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *GroupMembersRepository {
	return &GroupMembersRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *GroupMembersRepository) CreateGroupMember(ctx context.Context, groupMember *entity.GroupMember) error {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	groupMember.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, groupMember)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupMembersRepository) CreateBulkGroupMembers(ctx context.Context, groupMembers []*entity.GroupMember) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	for _, groupMember := range groupMembers {
		if groupMember.ID.IsZero() {
			groupMember.ID = primitive.NewObjectID()
		}
	}
	opts := options.InsertMany().SetOrdered(false)
	docs := utils.ToInterfaceSlice(groupMembers)
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
					ID:     groupMembers[we.Index].ID.Hex(),
					Reason: we.Message,
				})
			}
			return int64(len(result.InsertedIDs)), failedDocs, nil
		}
		return 0, nil, err
	}
	return int64(len(result.InsertedIDs)), failedDocs, nil
}
func (r *GroupMembersRepository) GetGroupMemberByID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	var groupMembers []*entity.GroupMember
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"group_member_cache_groupid_" + groupID,
			"group_member_cache_nextcursor_groupid_" + groupID,
			"group_member_cache_hasnext_groupid_" + groupID,
			"group_member_cache_limit_groupid_" + groupID,
		})
		if err == nil && datacache != nil {
			if groupMembersList, ok := datacache.([]*entity.GroupMember); ok {
				return &dto.PaginationRes{
					Data:       groupMembersList,
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
	if err = cursorDB.All(ctx, &groupMembers); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(groupMembers) > limit {
		hasNext = true
		groupMembers = groupMembers[:limit] // Lấy đúng số lượng cần thiết
		lastGroupMember := groupMembers[len(groupMembers)-1]
		nextCursor = utils.EncodeCursorMongodb(lastGroupMember.CreatedAt, lastGroupMember.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"group_member_cache_groupid_" + groupID:            groupMembers,
			"group_member_cache_nextcursor_groupid_" + groupID: nextCursor,
			"group_member_cache_hasnext_groupid_" + groupID:    hasNext,
			"group_member_cache_limit_groupid_" + groupID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       groupMembers,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupMembersRepository) GetGroupMemberByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	var groupMembers []*entity.GroupMember
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"group_member_cache_userid_" + userID,
			"group_member_cache_nextcursor_userid_" + userID,
			"group_member_cache_hasnext_userid_" + userID,
			"group_member_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if groupMembersList, ok := datacache.([]*entity.GroupMember); ok {
				return &dto.PaginationRes{
					Data:       groupMembersList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	query := bson.M{"user_id": userID}

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
	if err = cursorDB.All(ctx, &groupMembers); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(groupMembers) > limit {
		hasNext = true
		groupMembers = groupMembers[:limit] // Lấy đúng số lượng cần thiết
		lastGroupMember := groupMembers[len(groupMembers)-1]
		nextCursor = utils.EncodeCursorMongodb(lastGroupMember.CreatedAt, lastGroupMember.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"group_member_cache_userid_" + userID:            groupMembers,
			"group_member_cache_nextcursor_userid_" + userID: nextCursor,
			"group_member_cache_hasnext_userid_" + userID:    hasNext,
			"group_member_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       groupMembers,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupMembersRepository) GetGroupMemberByUserIDAndGroupID(ctx context.Context, userID string, groupID string) (*entity.GroupMember, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	finalgroupID, _ := primitive.ObjectIDFromHex(groupID)
	query := bson.M{"user_id": userID, "group_id": finalgroupID}
	var groupMember entity.GroupMember
	err := collection.FindOne(ctx, query).Decode(&groupMember)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy kết quả
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &groupMember, nil
}
func (r *GroupMembersRepository) UpdateGroupMember(ctx context.Context, groupMember *entity.GroupMember) error {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	filter := bson.M{"_id": groupMember.ID}
	update := bson.M{"$set": groupMember}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedGroupMember *entity.GroupMember
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedGroupMember)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy document nào để cập nhật
		}
		return err // Lỗi khác xảy ra
	}
	return nil
}
func (r *GroupMembersRepository) UpdateBulkGroupMembers(ctx context.Context, groupMembers []*entity.GroupMember) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(groupMembers))
	for _, ps := range groupMembers {
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
				// we.Index: Là chỉ số (index) trong slice 'groupMembers' ban đầu bị lỗi
				failedIndex := we.Index

				if failedIndex < len(groupMembers) {
					failedGroupMember := groupMembers[failedIndex]

					// Ghi nhận lại ID và lý do lỗi
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedGroupMember.ID.Hex(), // Hoặc failedGroupMember.ID.String()
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
func (r *GroupMembersRepository) DeleteGroupMember(ctx context.Context, groupID string) error {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	finalid, _ := primitive.ObjectIDFromHex(groupID)
	filter := bson.M{"_id": finalid}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupMembersRepository) DeleteBulkGroupMembers(ctx context.Context, groupIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range groupIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid group ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	objectIDs := make([]primitive.ObjectID, 0, len(groupIDs))
	for _, id := range groupIDs {
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
				if failedIndex < len(groupIDs) {
					failedGroupID := groupIDs[failedIndex]
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     failedGroupID,
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
func (r *GroupMembersRepository) DeleteGroupMemberByUserIDAndGroupID(ctx context.Context, userID string, groupID string) error {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())

	// Validate và convert GroupID sang ObjectID
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	// Filter kết hợp cả 2 field để trúng chính xác record
	filter := bson.M{
		"user_id":  userID,
		"group_id": objGroupID,
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("database error during delete: %w", err)
	}

	return nil
}
func (r *GroupMembersRepository) DeleteBulkGroupMemberByUserIDAndGroupID(ctx context.Context, userID string, groupIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())

	// Chuẩn bị các operation cho BulkWrite
	models := make([]mongo.WriteModel, 0, len(groupIDs))

	for _, gID := range groupIDs {
		objGroupID, err := primitive.ObjectIDFromHex(gID)
		if err != nil {
			// Fail-fast nếu có ID không hợp lệ, hoặc bạn có thể log và skip tùy logic
			return 0, nil, fmt.Errorf("invalid group ID format: %s", gID)
		}

		filter := bson.M{
			"user_id":  userID,
			"group_id": objGroupID,
		}

		models = append(models, mongo.NewDeleteOneModel().SetFilter(filter))
	}

	// Tránh gọi DB nếu mảng rỗng
	if len(models) == 0 {
		return 0, nil, nil
	}

	// SetOrdered(false) để nếu xóa 1 group lỗi thì các group khác vẫn được tiến hành xóa
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)

	// Trường hợp 1: Lỗi hệ thống toàn cục
	if result == nil && err != nil {
		return 0, nil, fmt.Errorf("bulk delete system error: %w", err)
	}

	// Trường hợp 2: Lỗi cục bộ (Partial Failure)
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				if we.Index < len(groupIDs) {
					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     groupIDs[we.Index],
						Reason: we.Message,
					})
				}
			}
			// Trả về số lượng xóa thành công và danh sách lỗi chi tiết
			return result.DeletedCount, failedDocs, nil
		}
		// Lỗi hệ thống khác không thuộc BulkWriteException
		return result.DeletedCount, nil, err
	}

	return result.DeletedCount, nil, nil
}
func (r *GroupMembersRepository) DeleteBulkGroupMemberByManyUserIDAndGroupID(ctx context.Context, userID []string, groupIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())

	// 1. Validate và convert toàn bộ groupIDs sang ObjectID
	var validGroupIDs []primitive.ObjectID
	for _, gID := range groupIDs {
		objGroupID, err := primitive.ObjectIDFromHex(gID)
		if err != nil {
			// Fail-fast nếu có ID không hợp lệ định dạng
			return 0, nil, fmt.Errorf("invalid group ID format: %s", gID)
		}
		validGroupIDs = append(validGroupIDs, objGroupID)
	}

	// Cấu trúc struct phụ trợ để map index của lỗi với cặp UserID - GroupID
	type pair struct {
		userID  string
		groupID string
	}

	var pairs []pair
	var models []mongo.WriteModel

	// 2. Xây dựng danh sách operations (Tổ hợp chéo các Users và Groups)
	for _, uID := range userID {
		for i, objGroupID := range validGroupIDs {
			// Lưu lại thông tin cặp ID này để mapping nếu xảy ra lỗi
			pairs = append(pairs, pair{
				userID:  uID,
				groupID: groupIDs[i], // Giữ lại string ID nguyên bản để trả về client
			})

			filter := bson.M{
				"user_id":  uID,
				"group_id": objGroupID,
			}

			// Thêm operation DeleteOne cho cặp này
			models = append(models, mongo.NewDeleteOneModel().SetFilter(filter))
		}
	}

	// Tránh gọi DB nếu không có operation nào
	if len(models) == 0 {
		return 0, nil, nil
	}
	// 3. Thực thi BulkWrite không tuần tự (lỗi ở doc này không chặn doc khác)
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)

	// Trường hợp 1: Lỗi hệ thống toàn cục (Network, DB Down...)
	if result == nil && err != nil {
		return 0, nil, fmt.Errorf("bulk delete system error: %w", err)
	}

	// Trường hợp 2: Lỗi cục bộ (Partial Failure)
	var failedDocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, we := range bulkErr.WriteErrors {
				if we.Index < len(pairs) {
					p := pairs[we.Index]
					// Gom UserID và GroupID thành chuỗi ID để client dễ nhận diện lỗi
					identifier := fmt.Sprintf("userID:%s|groupID:%s", p.userID, p.groupID)

					failedDocs = append(failedDocs, &mongodbErrors.BulkError{
						ID:     identifier,
						Reason: we.Message,
					})
				}
			}
			// Vẫn return DeletedCount và mảng failedDocs vì BulkWrite thành công một phần
			return result.DeletedCount, failedDocs, nil
		}

		// Nếu lỗi không phải BulkWriteException (VD: timeout)
		return result.DeletedCount, nil, err
	}

	// Thành công toàn bộ
	return result.DeletedCount, nil, nil
}
func (r *GroupMembersRepository) GetGroupMembersByUserIDAndGroupID(ctx context.Context, userID string, groupID string) ([]*entity.GroupMember, error) {
	collection := r.client.Collection(entity.GroupMember{}.CollectionName())
	finalgroupID, _ := primitive.ObjectIDFromHex(groupID)
	query := bson.M{"user_id": userID, "group_id": finalgroupID}
	var groupMembers []*entity.GroupMember
	cursorDB, err := collection.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursorDB.Close(ctx)
	if err = cursorDB.All(ctx, &groupMembers); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}
	return groupMembers, nil
}
