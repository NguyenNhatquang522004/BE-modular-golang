package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewConversationsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *ConversationsRepository {
	return &ConversationsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *ConversationsRepository) CreateConversation(ctx context.Context, conversation *entity.Conversation) error {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	if conversation.ID.IsZero() {
		conversation.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, conversation)
	if err != nil {
		return err
	}
	return nil
}
func (r *ConversationsRepository) CreateBulkConversations(ctx context.Context, conversations []entity.Conversation) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(conversations))
	for i := range conversations {
		if conversations[i].ID.IsZero() {
			conversations[i].ID = primitive.NewObjectID()
		}
		models = append(models, mongo.NewInsertOneModel().SetDocument(conversations[i]))
	}
	bulkOption := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, bulkOption)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErrors []*mongo.BulkWriteError
		if errors.As(err, &bulkErrors) {
			for _, be := range bulkErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     conversations[be.Index].ID.Hex(),
					Reason: be.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.InsertedCount, faildocs, nil
}
func (r *ConversationsRepository) GetConversationByID(ctx context.Context, id string) (*entity.Conversation, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var conversation entity.Conversation
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationsRepository) GetConversationByCreatorID(ctx context.Context, creatorID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	var conversations []*entity.Conversation
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"conversation_cache_creatorid_" + creatorID,
			"conversation_cache_nextcursor_creatorid_" + creatorID,
			"conversation_cache_hasnext_creatorid_" + creatorID,
			"conversation_cache_limit_creatorid_" + creatorID,
		})
		if err == nil && datacache != nil {
			if conversationsList, ok := datacache.([]*entity.Conversation); ok {
				return &dto.PaginationRes{
					Data:       conversationsList,
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
	if err = cursorDB.All(ctx, &conversations); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(conversations) > limit {
		hasNext = true
		conversations = conversations[:limit] // Lấy đúng số lượng cần thiết
		lastConversation := conversations[len(conversations)-1]
		nextCursor = utils.EncodeCursorMongodb(lastConversation.CreatedAt, lastConversation.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"conversation_cache_creatorid_" + creatorID:            conversations,
			"conversation_cache_nextcursor_creatorid_" + creatorID: nextCursor,
			"conversation_cache_hasnext_creatorid_" + creatorID:    hasNext,
			"conversation_cache_limit_creatorid_" + creatorID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       conversations,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *ConversationsRepository) GetConversationByOwnerID(ctx context.Context, ownerID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	var conversations []*entity.Conversation
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"conversation_cache_ownerid_" + ownerID,
			"conversation_cache_nextcursor_ownerid_" + ownerID,
			"conversation_cache_hasnext_ownerid_" + ownerID,
			"conversation_cache_limit_ownerid_" + ownerID,
		})
		if err == nil && datacache != nil {
			if conversationsList, ok := datacache.([]*entity.Conversation); ok {
				return &dto.PaginationRes{
					Data:       conversationsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	query := bson.M{"owner_id": ownerID}

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
	if err = cursorDB.All(ctx, &conversations); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(conversations) > limit {
		hasNext = true
		conversations = conversations[:limit] // Lấy đúng số lượng cần thiết
		lastConversation := conversations[len(conversations)-1]
		nextCursor = utils.EncodeCursorMongodb(lastConversation.CreatedAt, lastConversation.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"conversation_cache_ownerid_" + ownerID:            conversations,
			"conversation_cache_nextcursor_ownerid_" + ownerID: nextCursor,
			"conversation_cache_hasnext_ownerid_" + ownerID:    hasNext,
			"conversation_cache_limit_ownerid_" + ownerID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       conversations,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *ConversationsRepository) UpdateConversation(ctx context.Context, conversation *entity.Conversation) error {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	if conversation.ID.IsZero() {
		return fmt.Errorf("conversation ID is required for update")
	}
	filter := bson.M{"_id": conversation.ID}
	update := bson.M{"$set": conversation}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no conversation found with ID %s", conversation.ID.Hex())
	}
	return nil
}
func (r *ConversationsRepository) UpdateBulkConversations(ctx context.Context, conversations []entity.Conversation) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	models := make([]mongo.WriteModel, 0, len(conversations))
	for i := range conversations {
		if conversations[i].ID.IsZero() {
			return 0, nil, fmt.Errorf("conversation ID is required for update at index %d", i)
		}
		filter := bson.M{"_id": conversations[i].ID}
		update := bson.M{"$set": conversations[i]}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	bulkOption := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, bulkOption)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErrors []*mongo.BulkWriteError
		if errors.As(err, &bulkErrors) {
			for _, be := range bulkErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     conversations[be.Index].ID.Hex(),
					Reason: be.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.ModifiedCount, faildocs, nil
}
func (r *ConversationsRepository) DeleteConversation(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no conversation found with ID %s", id)
	}
	return nil
}
func (r *ConversationsRepository) DeleteBulkConversations(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid conversation ID %s: %w", id, err)
		}
		objectIDs = append(objectIDs, objectID)
	}
	result, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return 0, nil, err
	}
	return result.DeletedCount, nil, nil
}
func (r *ConversationsRepository) CheckConversationExists(ctx context.Context, userIDOne string, userIDTwo string) (*entity.Conversation, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	query := bson.M{
		"type": "private",
		"$or": []bson.M{
			{"creator_id": userIDOne, "owner_id": userIDTwo},
			{"creator_id": userIDTwo, "owner_id": userIDOne},
		},
	}
	var conversation entity.Conversation
	err := collection.FindOne(ctx, query).Decode(&conversation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy cuộc trò chuyện nào giữa hai người dùng
		}
		return nil, err // Lỗi khác xảy ra
	}
	return &conversation, nil // Trả về cuộc trò chuyện tìm thấy
}
func (r *ConversationsRepository) GetConversationByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// Implement the logic to get conversation by group ID
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	var conversations []*entity.Conversation
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"conversation_cache_groupID_" + groupID,
			"conversation_cache_nextcursor_groupID_" + groupID,
			"conversation_cache_hasnext_groupID_" + groupID,
			"conversation_cache_limit_groupID_" + groupID,
		})
		if err == nil && datacache != nil {
			if conversationsList, ok := datacache.([]*entity.Conversation); ok {
				return &dto.PaginationRes{
					Data:       conversationsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	parseID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %w", err)
	}
	query := bson.M{"related_group_id": parseID}
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
	if err = cursorDB.All(ctx, &conversations); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(conversations) > limit {
		hasNext = true
		conversations = conversations[:limit] // Lấy đúng số lượng cần thiết
		lastConversation := conversations[len(conversations)-1]
		nextCursor = utils.EncodeCursorMongodb(lastConversation.CreatedAt, lastConversation.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"conversation_cache_groupID_" + groupID:            conversations,
			"conversation_cache_nextcursor_groupID_" + groupID: nextCursor,
			"conversation_cache_hasnext_groupID_" + groupID:    hasNext,
			"conversation_cache_limit_groupID_" + groupID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       conversations,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *ConversationsRepository) CheckConversationExistsByPrivateChatKey(ctx context.Context, privateChatKey string) (*entity.Conversation, error) {
	collection := r.client.Collection(entity.Conversation{}.CollectionName())
	query := bson.M{
		"type":             "private",
		"private_chat_key": privateChatKey,
	}
	var conversation *entity.Conversation
	err := collection.FindOne(ctx, query).Decode(&conversation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy cuộc trò chuyện nào với privateChatKey này
		}
		return nil, err // Lỗi khác xảy ra
	}
	return conversation, nil // Trả về cuộc trò chuyện tìm thấy
}
