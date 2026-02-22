package mongodb

import (
	"context"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/mongodbErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/notification/enum"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationTemplatesRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewNotificationTemplatesRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *NotificationTemplatesRepository {
	return &NotificationTemplatesRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *NotificationTemplatesRepository) CreateTemplate(ctx context.Context, template *entity.NotificationTemplate) error {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	if template.ID.IsZero() {
		template.ID = primitive.NewObjectID()
	}
	_, err := collection.InsertOne(ctx, template)
	if err != nil {
		return err
	}
	return nil
}
func (r *NotificationTemplatesRepository) CreateBulkTemplates(ctx context.Context, templates []*entity.NotificationTemplate) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	for _, template := range templates {
		if template.ID.IsZero() {
			template.ID = primitive.NewObjectID()
		}
	}
	docs := utils.ToInterfaceSlice(templates)
	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		var bulkErrors []*mongodbErrors.BulkError
		if we, ok := err.(mongo.BulkWriteException); ok {
			for _, writeError := range we.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     templates[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		}
		return 0, bulkErrors, err
	}
	return int64(len(result.InsertedIDs)), nil, nil
}
func (r *NotificationTemplatesRepository) GetTemplateByType(ctx context.Context, notificationType enum.NotificationType, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	var template []*entity.NotificationTemplate
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"NotificationTemplate_cache_type_" + notificationType.String(),
			"NotificationTemplate_cache_nextcursor_type_" + notificationType.String(),
			"NotificationTemplate_cache_hasnext_type_" + notificationType.String(),
			"NotificationTemplate_cache_limit_type_" + notificationType.String(),
		})
		if err == nil && datacache != nil {
			if pagesList, ok := datacache.([]*entity.NotificationTemplate); ok {
				return &dto.PaginationRes{
					Data:       pagesList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	query := bson.M{"type": notificationType.String()}

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
	if err = cursorDB.All(ctx, &template); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(template) > limit {
		hasNext = true
		template = template[:limit] // Lấy đúng số lượng cần thiết
		lastTemplate := template[len(template)-1]
		nextCursor = utils.EncodeCursorMongodb(lastTemplate.CreatedAt, lastTemplate.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"NotificationTemplate_cache_type_" + notificationType.String():            template,
			"NotificationTemplate_cache_nextcursor_type_" + notificationType.String(): nextCursor,
			"NotificationTemplate_cache_hasnext_type_" + notificationType.String():    hasNext,
			"NotificationTemplate_cache_limit_type_" + notificationType.String():      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       template,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *NotificationTemplatesRepository) GetTemplateByID(ctx context.Context, id string) (*entity.NotificationTemplate, error) {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	var template *entity.NotificationTemplate
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy template
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return template, nil
}
func (r *NotificationTemplatesRepository) GetTemplate(ctx context.Context, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	var template []*entity.NotificationTemplate
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"NotificationTemplate_cache",
			"NotificationTemplate_cache_nextcursor",
			"NotificationTemplate_cache_hasnext",
			"NotificationTemplate_cache_limit",
		})
		if err == nil && datacache != nil {
			if pagesList, ok := datacache.([]*entity.NotificationTemplate); ok {
				return &dto.PaginationRes{
					Data:       pagesList,
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
	if err = cursorDB.All(ctx, &template); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(template) > limit {
		hasNext = true
		template = template[:limit] // Lấy đúng số lượng cần thiết
		lastTemplate := template[len(template)-1]
		nextCursor = utils.EncodeCursorMongodb(lastTemplate.CreatedAt, lastTemplate.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"NotificationTemplate_cache":            template,
			"NotificationTemplate_cache_nextcursor": nextCursor,
			"NotificationTemplate_cache_hasnext":    hasNext,
			"NotificationTemplate_cache_limit":      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       template,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *NotificationTemplatesRepository) UpdateTemplate(ctx context.Context, template *entity.NotificationTemplate) error {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	if template.ID.IsZero() {
		return fmt.Errorf("template ID is required for update")
	}
	update := bson.M{
		"$set": template,
	}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": template.ID}, update)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}
func (r *NotificationTemplatesRepository) UpdateBulkTemplates(ctx context.Context, templates []*entity.NotificationTemplate) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	var models []mongo.WriteModel
	for _, template := range templates {
		if template.ID.IsZero() {
			return 0, nil, fmt.Errorf("template ID is required for update")
		}
		update := bson.M{
			"$set": template,
		}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": template.ID}).SetUpdate(update))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if err != nil {
		var bulkErrors []*mongodbErrors.BulkError
		if we, ok := err.(mongo.BulkWriteException); ok {
			for _, writeError := range we.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     templates[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		}
		return 0, bulkErrors, err
	}
	return result.ModifiedCount, nil, nil
}
func (r *NotificationTemplatesRepository) DeleteTemplate(ctx context.Context, id string) error {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}
func (r *NotificationTemplatesRepository) DeleteBulkTemplates(ctx context.Context, ids []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.NotificationTemplate{}.CollectionName())
	var models []mongo.WriteModel
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid ID format: %w", err)
		}
		models = append(models, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": oid}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if err != nil {
		var bulkErrors []*mongodbErrors.BulkError
		if we, ok := err.(mongo.BulkWriteException); ok {
			for _, writeError := range we.WriteErrors {
				bulkErrors = append(bulkErrors, &mongodbErrors.BulkError{
					ID:     ids[writeError.Index],
					Reason: writeError.Message,
				})
			}
		}
		return 0, bulkErrors, err
	}
	return result.DeletedCount, nil, nil
}
