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

type ConversationParticipantsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewConversationParticipantsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *ConversationParticipantsRepository {
	return &ConversationParticipantsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *ConversationParticipantsRepository) CreateConversationParticipant(ctx context.Context, participant *entity.ConversationParticipant) error {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	if participant.ID.IsZero() {
		participant.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, participant)
	if err != nil {
		return err
	}
	return nil
}
func (r *ConversationParticipantsRepository) CreateBulkConversationParticipants(ctx context.Context, participants []entity.ConversationParticipant) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	docs := utils.ToInterfaceSlice(participants)
	result, err := collection.InsertMany(ctx, docs)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     participants[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return int64(len(result.InsertedIDs)), faildocs, nil
}
func (r *ConversationParticipantsRepository) GetConversationParticipantsByConversationID(ctx context.Context, conversationID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	var participants []*entity.ConversationParticipant
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"conversation_participants_cache_conversationid_" + conversationID,
			"conversation_participants_cache_nextcursor_conversationid_" + conversationID,
			"conversation_participants_cache_hasnext_conversationid_" + conversationID,
			"conversation_participants_cache_limit_conversationid_" + conversationID,
		})
		if err == nil && datacache != nil {
			if participantsList, ok := datacache.([]*entity.ConversationParticipant); ok {
				return &dto.PaginationRes{
					Data:       participantsList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}
	finalid, _ := primitive.ObjectIDFromHex(conversationID)
	query := bson.M{"conversation_id": finalid}

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
	if err = cursorDB.All(ctx, &participants); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(participants) > limit {
		hasNext = true
		participants = participants[:limit] // Lấy đúng số lượng cần thiết
		lastParticipant := participants[len(participants)-1]
		nextCursor = utils.EncodeCursorMongodb(lastParticipant.CreatedAt, lastParticipant.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"conversation_participants_cache_conversationid_" + conversationID:            participants,
			"conversation_participants_cache_nextcursor_conversationid_" + conversationID: nextCursor,
			"conversation_participants_cache_hasnext_conversationid_" + conversationID:    hasNext,
			"conversation_participants_cache_limit_conversationid_" + conversationID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       participants,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *ConversationParticipantsRepository) GetConversationParticipantsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	var participants []*entity.ConversationParticipant
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"conversation_participants_cache_userid_" + userID,
			"conversation_participants_cache_nextcursor_userid_" + userID,
			"conversation_participants_cache_hasnext_userid_" + userID,
			"conversation_participants_cache_limit_userid_" + userID,
		})
		if err == nil && datacache != nil {
			if participantsList, ok := datacache.([]*entity.ConversationParticipant); ok {
				return &dto.PaginationRes{
					Data:       participantsList,
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
	if err = cursorDB.All(ctx, &participants); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(participants) > limit {
		hasNext = true
		participants = participants[:limit] // Lấy đúng số lượng cần thiết
		lastParticipant := participants[len(participants)-1]
		nextCursor = utils.EncodeCursorMongodb(lastParticipant.CreatedAt, lastParticipant.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"conversation_participants_cache_userid_" + userID:            participants,
			"conversation_participants_cache_nextcursor_userid_" + userID: nextCursor,
			"conversation_participants_cache_hasnext_userid_" + userID:    hasNext,
			"conversation_participants_cache_limit_userid_" + userID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       participants,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *ConversationParticipantsRepository) UpdateConversationParticipant(ctx context.Context, participant *entity.ConversationParticipant) error {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	filter := bson.M{"_id": participant.ID}
	update := bson.M{"$set": participant}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ConversationParticipantsRepository) UpdateBulkConversationParticipants(ctx context.Context, participants []entity.ConversationParticipant) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	var models []mongo.WriteModel
	for _, participant := range participants {
		filter := bson.M{"_id": participant.ID}
		update := bson.M{"$set": participant}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}

	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     participants[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}

	return result.ModifiedCount, faildocs, nil
}

func (r *ConversationParticipantsRepository) DeleteConversationParticipant(ctx context.Context, conversationID string, userID string) error {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	query := bson.M{
		"conversation_id": conversationID,
		"user_id":         userID,
	}
	result, err := collection.DeleteOne(ctx, query)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ConversationParticipantsRepository) DeleteBulkConversationParticipants(ctx context.Context, conversationIDs []string, userIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.ConversationParticipant{}.CollectionName())
	var models []mongo.WriteModel
	for _, conversationID := range conversationIDs {
		for _, userID := range userIDs {
			query := bson.M{
				"conversation_id": conversationID,
				"user_id":         userID,
			}
			models = append(models, mongo.NewDeleteOneModel().SetFilter(query))
		}
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}

	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkErr mongo.BulkWriteException
		if errors.As(err, &bulkErr) {
			for _, writeError := range bulkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     fmt.Sprintf("conversation_id: %s, user_id: %s", conversationIDs[writeError.Index/len(userIDs)], userIDs[writeError.Index%len(userIDs)]),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}

	return result.DeletedCount, faildocs, nil
}
