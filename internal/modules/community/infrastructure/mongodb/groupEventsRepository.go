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

type GroupEventsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewGroupEventsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *GroupEventsRepository {
	return &GroupEventsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}

func (r *GroupEventsRepository) CreateGroupEvent(ctx context.Context, event *entity.GroupEvent) error {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupEventsRepository) CreateBulkGroupEvents(ctx context.Context, events []*entity.GroupEvent) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	docs := utils.ToInterfaceSlice(events)
	var model []mongo.WriteModel
	for _, doc := range docs {
		model = append(model, mongo.NewInsertOneModel().SetDocument(doc))
	}
	bulkOption := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, bulkOption)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bunkErr mongo.BulkWriteException
		if errors.As(err, &bunkErr) {
			for _, writeError := range bunkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     string(events[writeError.Index].ID.Hex()),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.InsertedCount, faildocs, nil
}
func (r *GroupEventsRepository) GetGroupEventByID(ctx context.Context, id string) (*entity.GroupEvent, error) {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var event entity.GroupEvent
	err = collection.FindOne(ctx, bson.M{"_id": idObj}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *GroupEventsRepository) GetGroupEventsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	var events []*entity.GroupEvent
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"group_event_cache_groupid_" + groupID,
			"group_event_cache_nextcursor_groupid_" + groupID,
			"group_event_cache_hasnext_groupid_" + groupID,
			"group_event_cache_limit_groupid_" + groupID,
		})
		if err == nil && datacache != nil {
			if eventsList, ok := datacache.([]*entity.GroupEvent); ok {
				return &dto.PaginationRes{
					Data:       eventsList,
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
	if err = cursorDB.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(events) > limit {
		hasNext = true
		events = events[:limit] // Lấy đúng số lượng cần thiết
		lastEvent := events[len(events)-1]
		nextCursor = utils.EncodeCursorMongodb(lastEvent.CreatedAt, lastEvent.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"group_event_cache_groupid_" + groupID:            events,
			"group_event_cache_nextcursor_groupid_" + groupID: nextCursor,
			"group_event_cache_hasnext_groupid_" + groupID:    hasNext,
			"group_event_cache_limit_groupid_" + groupID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       events,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupEventsRepository) UpdateGroupEvent(ctx context.Context, event *entity.GroupEvent) error {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	var updatedDoc *entity.GroupEvent
	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": event}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, filter, update, options).Decode(&updatedDoc)
	if err != nil {
		return err

	}
	return nil
}
func (r *GroupEventsRepository) UpdateBulkGroupEvents(ctx context.Context, events []*entity.GroupEvent) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	var model []mongo.WriteModel
	for _, event := range events {
		filter := bson.M{"_id": event.ID}
		update := bson.M{"$set": event}
		model = append(model, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	bulkOption := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, bulkOption)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bunkErr mongo.BulkWriteException
		if errors.As(err, &bunkErr) {
			for _, writeError := range bunkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     string(events[writeError.Index].ID.Hex()),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.ModifiedCount, faildocs, nil
}
func (r *GroupEventsRepository) UpdateGroupEventDownloadCount(ctx context.Context, id string, newCount int) error {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	var updatedDoc *entity.GroupEvent
	filter := bson.M{"_id": idObj}
	update := bson.M{"$set": bson.M{"download_count": newCount}}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, filter, update, options).Decode(&updatedDoc)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupEventsRepository) DeleteGroupEvent(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
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
func (r *GroupEventsRepository) DeleteBulkGroupEvents(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	var model []mongo.WriteModel
	for _, id := range ids {
		idObj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, err
		}
		model = append(model, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": idObj}))
	}
	bulkOption := options.BulkWrite().SetOrdered(false)
	result, err := r.client.Collection(entity.GroupEvent{}.CollectionName()).BulkWrite(ctx, model, bulkOption)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bunkErr mongo.BulkWriteException
		if errors.As(err, &bunkErr) {
			for _, writeError := range bunkErr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
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
func (r *GroupEventsRepository) DeleteGroupEventsByGroupID(ctx context.Context, groupID string) error {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
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
func (r *GroupEventsRepository) DeleteGroupEventsByCreatorIDAndGroupID(ctx context.Context, creatorID string, groupID string) error {
	collection := r.client.Collection(entity.GroupEvent{}.CollectionName())
	finalid, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	_, err = collection.DeleteMany(ctx, bson.M{"group_id": finalid, "creator_id": creatorID})
	if err != nil {
		return err
	}
	return nil
}
