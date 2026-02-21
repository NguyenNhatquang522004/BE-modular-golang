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

type GroupJoinQuestionsRepository struct {
	client    *mongo.Database
	redisRepo IRepositoryShare.IRedis
}

func NewGroupJoinQuestionsRepository(client *mongo.Database, redisRepo IRepositoryShare.IRedis) *GroupJoinQuestionsRepository {
	return &GroupJoinQuestionsRepository{
		client:    client,
		redisRepo: redisRepo,
	}
}
func (r *GroupJoinQuestionsRepository) CreateGroupJoinQuestion(ctx context.Context, question *entity.GroupJoinQuestion) error {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	_, err := collection.InsertOne(ctx, question)
	if err != nil {
		return err
	}
	return nil
}
func (r *GroupJoinQuestionsRepository) CreateBulkGroupJoinQuestions(ctx context.Context, questions []entity.GroupJoinQuestion) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())

	for i := range questions {
		if questions[i].ID.IsZero() {
			questions[i].ID = primitive.NewObjectID()

		}
	}
	docs := utils.ToInterfaceSlice(questions)
	result, err := collection.InsertMany(ctx, docs)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkerr mongo.BulkWriteException
		if errors.As(err, &bulkerr) {
			for _, writeError := range bulkerr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     questions[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return int64(len(result.InsertedIDs)), faildocs, nil
}
func (r *GroupJoinQuestionsRepository) GetGroupJoinQuestionByID(ctx context.Context, questionID string) (*entity.GroupJoinQuestion, error) {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	var question entity.GroupJoinQuestion
	objID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, primitive.M{"_id": objID}).Decode(&question)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *GroupJoinQuestionsRepository) GetGroupJoinQuestionsByGroupID(ctx context.Context, groupID string, cursor string, limit int) (*dto.PaginationRes, error) {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	var questions []*entity.GroupJoinQuestion
	querylimit := int64(limit + 1)
	// Check cache first
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"group_join_question_cache_groupid_" + groupID,
			"group_join_question_cache_nextcursor_groupid_" + groupID,
			"group_join_question_cache_hasnext_groupid_" + groupID,
			"group_join_question_cache_limit_groupid_" + groupID,
		})
		if err == nil && datacache != nil {
			if questionsList, ok := datacache.([]*entity.GroupJoinQuestion); ok {
				return &dto.PaginationRes{
					Data:       questionsList,
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
	if err = cursorDB.All(ctx, &questions); err != nil {
		return nil, fmt.Errorf("error decoding results: %w", err)
	}

	var nextCursor string
	hasNext := false
	if len(questions) > limit {
		hasNext = true
		questions = questions[:limit] // Lấy đúng số lượng cần thiết
		lastQuestion := questions[len(questions)-1]
		nextCursor = utils.EncodeCursorMongodb(lastQuestion.CreatedAt, lastQuestion.ID)
	}
	// Cache kết quả nếu là trang đầu tiên
	if cursor == "" {
		items := map[string]any{
			"group_join_question_cache_groupid_" + groupID:            questions,
			"group_join_question_cache_nextcursor_groupid_" + groupID: nextCursor,
			"group_join_question_cache_hasnext_groupid_" + groupID:    hasNext,
			"group_join_question_cache_limit_groupid_" + groupID:      limit,
		}
		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for pagination: %v\n", err)
		}
	}
	return &dto.PaginationRes{
		Data:       questions,
		NextCursor: nextCursor,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *GroupJoinQuestionsRepository) UpdateGroupJoinQuestion(ctx context.Context, question *entity.GroupJoinQuestion) error {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	filter := bson.M{"_id": question.ID}
	update := bson.M{"$set": question}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedQuestion entity.GroupJoinQuestion
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedQuestion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("question not found: %w", err)
		}
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}
func (r *GroupJoinQuestionsRepository) UpdateBulkGroupJoinQuestions(ctx context.Context, questions []entity.GroupJoinQuestion) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	var models []mongo.WriteModel
	for _, question := range questions {
		models = append(models, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": question.ID}).SetUpdate(bson.M{"$set": question}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkerr mongo.BulkWriteException
		if errors.As(err, &bulkerr) {
			for _, writeError := range bulkerr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     questions[writeError.Index].ID.Hex(),
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return int64(result.ModifiedCount), faildocs, nil
}
func (r *GroupJoinQuestionsRepository) DeleteGroupJoinQuestion(ctx context.Context, questionID string) error {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	objID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return err
	}
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("question not found")
	}
	return nil
}
func (r *GroupJoinQuestionsRepository) DeleteBulkGroupJoinQuestions(ctx context.Context, questionIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	for _, id := range questionIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return 0, nil, fmt.Errorf("invalid question ID: %s", id)
		}
	}
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	var models []mongo.WriteModel
	for _, id := range questionIDs {
		objID, _ := primitive.ObjectIDFromHex(id)
		models = append(models, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": objID}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkerr mongo.BulkWriteException
		if errors.As(err, &bulkerr) {
			for _, writeError := range bulkerr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     questionIDs[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return int64(result.DeletedCount), faildocs, nil
}
func (r *GroupJoinQuestionsRepository) DeleteGroupJoinQuestionsByGroupIDAndQuestionIDs(ctx context.Context, groupID string, questionIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	var objIDs []primitive.ObjectID
	for _, id := range questionIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid question ID: %s", id)
		}
		objIDs = append(objIDs, objID)
	}
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	groupObjID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid group ID: %s", groupID)
	}
	var models []mongo.WriteModel
	for _, id := range questionIDs {
		objID, _ := primitive.ObjectIDFromHex(id)
		models = append(models, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": objID, "group_id": groupObjID}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkerr mongo.BulkWriteException
		if errors.As(err, &bulkerr) {
			for _, writeError := range bulkerr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     questionIDs[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.DeletedCount, faildocs, nil
}

func (r *GroupJoinQuestionsRepository) DeleteGroupJoinQuestionsByGroupID(ctx context.Context, groupID string) (int64, error) {
	groupObjID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return 0, fmt.Errorf("invalid group ID: %s", groupID)
	}
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	result, err := collection.DeleteMany(ctx, bson.M{"group_id": groupObjID})
	if err != nil {
		return 0, err
	}
	// Invalidate cache
	cacheKey := "group_join_question_cache_groupid_" + groupID
	_, err = r.redisRepo.Del(ctx, cacheKey)
	if err != nil {
		fmt.Printf("Failed to invalidate cache for group %s: %v\n", groupID, err)
	}
	return result.DeletedCount, nil
}
func (r *GroupJoinQuestionsRepository) DeleteBulkGroupJoinQuestionsByGroupID(ctx context.Context, groupIDs []string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	var model []mongo.WriteModel
	for _, groupID := range groupIDs {
		id, err := primitive.ObjectIDFromHex(groupID)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid group ID: %s", groupID)
		}
		model = append(model, mongo.NewDeleteManyModel().SetFilter(bson.M{"group_id": id}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, model, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkerr mongo.BulkWriteException
		if errors.As(err, &bulkerr) {
			for _, writeError := range bulkerr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     groupIDs[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.DeletedCount, faildocs, nil
}
func (r *GroupJoinQuestionsRepository) DeleteGroupJoinQuestionByIDAndGroupID(ctx context.Context, questionID string, groupID string) error {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	questionObjID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return fmt.Errorf("invalid question ID: %s", questionID)
	}
	groupObjID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %s", groupID)
	}
	result, err := collection.DeleteOne(ctx, bson.M{"_id": questionObjID, "group_id": groupObjID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("question not found for the specified group")
	}
	return nil
}
func (r *GroupJoinQuestionsRepository) DeleteBulkGroupJoinQuestionByIDAndGroupID(ctx context.Context, questionIDs []string, groupID string) (int64, []*mongodbErrors.BulkError, error) {
	collection := r.client.Collection(entity.GroupJoinQuestion{}.CollectionName())
	groupObjID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid group ID: %s", groupID)
	}
	var models []mongo.WriteModel
	for _, questionID := range questionIDs {
		questionObjID, err := primitive.ObjectIDFromHex(questionID)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid question ID: %s", questionID)
		}
		models = append(models, mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": questionObjID, "group_id": groupObjID}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	result, err := collection.BulkWrite(ctx, models, opts)
	if result == nil && err != nil {
		return 0, nil, err
	}
	var faildocs []*mongodbErrors.BulkError
	if err != nil {
		var bulkerr mongo.BulkWriteException
		if errors.As(err, &bulkerr) {
			for _, writeError := range bulkerr.WriteErrors {
				faildocs = append(faildocs, &mongodbErrors.BulkError{
					ID:     questionIDs[writeError.Index],
					Reason: writeError.Message,
				})
			}
		} else {
			return 0, nil, err
		}
	}
	return result.DeletedCount, faildocs, nil
}
