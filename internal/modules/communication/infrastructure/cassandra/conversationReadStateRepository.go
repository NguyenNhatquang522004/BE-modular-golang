package cassandra

import (
	"context"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"github.com/gocql/gocql"
)

type ConversationReadStateRepository struct {
	session   *gocql.Session
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
}

func NewConversationReadStateRepository(session *gocql.Session, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis) *ConversationReadStateRepository {
	return &ConversationReadStateRepository{
		session:   session,
		pool:      pool,
		redisRepo: redisRepo,
	}
}
func (r *ConversationReadStateRepository) CreateConversationReadState(ctx context.Context, readState *entity.ConversationReadState) error {
	if readState == nil {
		return nil
	}
	if readState.ConversationID == "" || readState.UserID == (gocql.UUID{}) {
		return nil
	}
	safectx := context.WithoutCancel(ctx)
	err := r.session.Query(
		`INSERT INTO conversation_read_state (conversation_id, user_id, last_read_message_id, last_read_at) VALUES (?, ?, ?, ?)`,
		readState.ConversationID, readState.UserID, readState.LastReadMessageID, readState.LastReadAt,
	).WithContext(safectx).Exec()
	if err != nil {
		return err
	}
	return nil
}
func (r *ConversationReadStateRepository) CreateBulkConversationReadStates(ctx context.Context, readStates []*entity.ConversationReadState) (int64, []*cassandraErrors.ConversationReadStateBulkError, error) {
	if len(readStates) == 0 {
		return 0, nil, nil
	}
	type taskResult struct {
		readState *entity.ConversationReadState
		err       error
	}

	resutlsCh := make(chan *taskResult, len(readStates))
	tablename := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		INSERT INTO %s (conversation_id, user_id, last_read_message_id, last_read_at) VALUES (?, ?, ?, ?)
	`, tablename)
	var emptyUUID gocql.UUID
	for i, item := range readStates {
		if item == nil {
			resutlsCh <- &taskResult{readState: nil,
				err: fmt.Errorf("item at index %d is nil", i)}
		}
		if item.ConversationID == "" || item.UserID == emptyUUID {
			resutlsCh <- &taskResult{readState: item,
				err: fmt.Errorf("item at index %d has invalid ConversationID or UserID", i)}
		}
		payload := item
		err := r.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)

			err := r.session.Query(query,
				payload.ConversationID, payload.UserID, payload.LastReadMessageID, payload.LastReadAt,
			).WithContext(safectx).Exec()
			resutlsCh <- &taskResult{readState: payload, err: err}
		})

		if err != nil {
			resutlsCh <- &taskResult{readState: item, err: fmt.Errorf("failed to run task for item at index %d: %w", i, err)}
			break
		}
	}
	r.pool.Wait()
	close(resutlsCh)
	var successCount int64
	var bulkErrors []*cassandraErrors.ConversationReadStateBulkError
	for result := range resutlsCh {
		if result.err != nil {
			convID := "unknown_conversation"
			userID := "unknown_user"
			if result.readState != nil {
				convID = result.readState.ConversationID
				userID = result.readState.UserID.String()
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.ConversationReadStateBulkError{
				ConversationID: convID,
				UserID:         userID,
				Error:          result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("failed to create %d conversation read states", len(bulkErrors))
	}
	return successCount, bulkErrors, finalErr
}
func (r *ConversationReadStateRepository) UpdateConversationReadState(ctx context.Context, readState *entity.ConversationReadState) error {
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		UPDATE %s SET last_read_message_id = ?, last_read_at = ? WHERE conversation_id = ? AND user_id = ?
	`, tableName)
	safectx := context.WithoutCancel(ctx)
	err := r.session.Query(query,
		readState.LastReadMessageID, readState.LastReadAt, readState.ConversationID, readState.UserID,
	).WithContext(safectx).Exec()
	if err != nil {
		return err
	}
	return nil
}
func (r *ConversationReadStateRepository) UpdateBulkConversationReadStates(ctx context.Context, readStates []*entity.ConversationReadState) (int64, []*cassandraErrors.ConversationReadStateBulkError, error) {
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		UPDATE %s SET last_read_message_id = ?, last_read_at = ? WHERE conversation_id = ? AND user_id = ?
	`, tableName)
	type taskResult struct {
		readState *entity.ConversationReadState
		err       error
	}
	resutlsCh := make(chan *taskResult, len(readStates))
	var emptyUUID gocql.UUID
	for i, item := range readStates {
		if item == nil {
			resutlsCh <- &taskResult{readState: nil,
				err: fmt.Errorf("item at index %d is nil", i)}
		}
		if item.ConversationID == "" || item.UserID == emptyUUID {
			resutlsCh <- &taskResult{readState: item,
				err: fmt.Errorf("item at index %d has invalid ConversationID or UserID", i)}
		}
		payload := item
		err := r.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			err := r.session.Query(query,
				payload.LastReadMessageID, payload.LastReadAt, payload.ConversationID, payload.UserID,
			).WithContext(safectx).Exec()
			resutlsCh <- &taskResult{readState: payload, err: err}
		})
		if err != nil {
			resutlsCh <- &taskResult{readState: item, err: fmt.Errorf("failed to run task for item at index %d: %w", i, err)}
			break
		}

	}
	r.pool.Wait()
	close(resutlsCh)
	var successCount int64
	var bulkErrors []*cassandraErrors.ConversationReadStateBulkError
	for result := range resutlsCh {
		if result.err != nil {
			convID := "unknown_conversation"
			userID := "unknown_user"
			if result.readState != nil {
				convID = result.readState.ConversationID
				userID = result.readState.UserID.String()
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.ConversationReadStateBulkError{
				ConversationID: convID,
				UserID:         userID,
				Error:          result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("failed to update %d conversation read states", len(bulkErrors))
	}
	return successCount, bulkErrors, finalErr
}
func (r *ConversationReadStateRepository) DeleteConversationReadState(ctx context.Context, conversationID string, userID string) error {
	useridFinal, err := gocql.ParseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid userID: %s", userID)
	}
	if conversationID == "" {
		return fmt.Errorf("conversationID cannot be empty")
	}
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE conversation_id = ? AND user_id = ?
	`, tableName)
	safectx := context.WithoutCancel(ctx)
	err = r.session.Query(query, conversationID, useridFinal).WithContext(safectx).Exec()
	if err != nil {
		return err
	}
	return nil
}
func (r *ConversationReadStateRepository) DeleteBulkConversationReadStates(ctx context.Context, conversationIDs []string, userIDs []string) (int64, []*cassandraErrors.ConversationReadStateBulkError, error) {
	if len(conversationIDs) == 0 || len(userIDs) == 0 {
		return 0, nil, nil
	}
	if len(conversationIDs) != len(userIDs) {
		return 0, nil, fmt.Errorf("conversationIDs and userIDs must have the same length")
	}
	type taskResult struct {
		conversationID string
		userID         string
		err            error
	}
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE conversation_id = ? AND user_id = ?
	`, tableName)
	resultsCh := make(chan *taskResult, len(conversationIDs))
	for i := range conversationIDs {
		convID := conversationIDs[i]
		userID, err := gocql.ParseUUID(userIDs[i])
		if err != nil {
			resultsCh <- &taskResult{conversationID: convID, userID: userIDs[i], err: fmt.Errorf("invalid userID at index %d: %s", i, userIDs[i])}
			continue
		}
		if convID == "" {
			resultsCh <- &taskResult{conversationID: convID, userID: userIDs[i], err: fmt.Errorf("conversationID cannot be empty at index %d", i)}
			continue
		}
		err = r.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			err := r.session.Query(
				query,

				convID, userID,
			).WithContext(safectx).Exec()
			resultsCh <- &taskResult{conversationID: convID, userID: userIDs[i], err: err}
		})
		if err != nil {
			resultsCh <- &taskResult{conversationID: convID, userID: userIDs[i], err: fmt.Errorf("failed to run task for index %d: %w", i, err)}
			break
		}
	}
	r.pool.Wait()
	close(resultsCh)
	var bulkErrors []*cassandraErrors.ConversationReadStateBulkError
	var successCount int64
	for result := range resultsCh {
		if result.err != nil {
			convID := "unknown_conversation"
			userID := "unknown_user"
			if result.conversationID != "" {
				convID = result.conversationID
				userID = result.userID
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.ConversationReadStateBulkError{
				ConversationID: convID,
				UserID:         userID,
				Error:          result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("failed to delete %d conversation read states", len(bulkErrors))
	}
	return successCount, bulkErrors, finalErr
}

func (r *ConversationReadStateRepository) DeleteBulkConversationReadStatesByManyConversationID(ctx context.Context, conversationIDs []string, userID string) (int64, []*cassandraErrors.ConversationReadStateBulkError, error) {
	if len(conversationIDs) == 0 {
		return 0, nil, nil
	}
	if userID == "" {
		return 0, nil, fmt.Errorf("userID cannot be empty")
	}
	userIDFinal, err := gocql.ParseUUID(userID)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid userID: %s", userID)
	}
	type taskResult struct {
		conversationID string
		userID         string
		err            error
	}
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE conversation_id = ? AND user_id = ?
	`, tableName)
	resultsCh := make(chan *taskResult, len(conversationIDs))
	for i, convID := range conversationIDs {
		if convID == "" {
			resultsCh <- &taskResult{conversationID: convID, userID: userID, err: fmt.Errorf("conversationID cannot be empty at index %d", i)}
			continue
		}
		err = r.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			err := r.session.Query(
				query,
				convID, userIDFinal,
			).WithContext(safectx).Exec()
			resultsCh <- &taskResult{conversationID: convID, userID: userID, err: err}
		})
		if err != nil {
			resultsCh <- &taskResult{conversationID: convID, userID: userID, err: fmt.Errorf("failed to run task for index %d: %w", i, err)}
		}
	}
	r.pool.Wait()
	close(resultsCh)
	var bulkErrors []*cassandraErrors.ConversationReadStateBulkError
	var successCount int64
	for result := range resultsCh {
		if result.err != nil {
			convID := "unknown_conversation"
			userID := "unknown_user"
			if result.conversationID != "" {
				convID = result.conversationID
			}
			if result.userID != "" {
				userID = result.userID
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.ConversationReadStateBulkError{
				ConversationID: convID,
				UserID:         userID,
				Error:          result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("failed to delete %d conversation read states", len(bulkErrors))
	}
	return successCount, bulkErrors, finalErr
}
func (r *ConversationReadStateRepository) DeleteConversationReadStatesByConversationID(ctx context.Context, conversationID string) error {
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE conversation_id = ?
	`, tableName)
	safectx := context.WithoutCancel(ctx)
	err := r.session.Query(query, conversationID).WithContext(safectx).Exec()
	if err != nil {
		return err
	}
	return nil
}
func (r *ConversationReadStateRepository) DeleteConversationReadStatesByUserID(ctx context.Context, conversationID string, userID string) error {
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE conversation_id = ? AND user_id = ?
	`, tableName)
	safectx := context.WithoutCancel(ctx)
	err := r.session.Query(query, conversationID, userID).WithContext(safectx).Exec()
	if err != nil {
		return err
	}
	return nil
}
func (r *ConversationReadStateRepository) DeleteBulkConversationReadStatesByManyUserID(ctx context.Context, conversationID string, userIDs []string) (int64, []*cassandraErrors.ConversationReadStateBulkError, error) {
	tableName := (&entity.ConversationReadState{}).TableName()
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE conversation_id = ? AND user_id = ?
	`, tableName)
	type taskResult struct {
		conversationID string
		userID         string
		err            error
	}
	resultsCh := make(chan *taskResult, len(userIDs))
	for i, userID := range userIDs {
		if userID == "" {
			resultsCh <- &taskResult{conversationID: conversationID, userID: userID, err: fmt.Errorf("userID cannot be empty at index %d", i)}
			continue
		}
		userIDFinal, err := gocql.ParseUUID(userID)
		if err != nil {
			resultsCh <- &taskResult{conversationID: conversationID, userID: userID, err: fmt.Errorf("invalid userID at index %d: %s", i, userID)}
			continue
		}
		err = r.pool.Run(ctx, func() {
			safectx := context.WithoutCancel(ctx)
			err := r.session.Query(
				query,
				conversationID, userIDFinal,
			).WithContext(safectx).Exec()
			resultsCh <- &taskResult{conversationID: conversationID, userID: userID, err: err}
		})
		if err != nil {
			resultsCh <- &taskResult{conversationID: conversationID, userID: userID, err: fmt.Errorf("failed to run task for index %d: %w", i, err)}
		}
	}
	r.pool.Wait()
	close(resultsCh)
	var bulkErrors []*cassandraErrors.ConversationReadStateBulkError
	var successCount int64
	for result := range resultsCh {
		if result.err != nil {
			convID := "unknown_conversation"
			userID := "unknown_user"
			if result.conversationID != "" {
				convID = result.conversationID
			}
			if result.userID != "" {
				userID = result.userID
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.ConversationReadStateBulkError{
				ConversationID: convID,
				UserID:         userID,
				Error:          result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("failed to delete %d conversation read states", len(bulkErrors))
	}
	return successCount, bulkErrors, finalErr
}
