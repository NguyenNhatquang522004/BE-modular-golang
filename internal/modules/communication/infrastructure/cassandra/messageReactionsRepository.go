package cassandra

import (
	"context"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/communication/domain/entity"
	"github.com/gocql/gocql"
)

type MessageReactionsRepository struct {
	session   *gocql.Session
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
}

func NewMessageReactionsRepository(session *gocql.Session, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis) *MessageReactionsRepository {
	return &MessageReactionsRepository{
		session:   session,
		pool:      pool,
		redisRepo: redisRepo,
	}
}

// =========================================================================
// CREATE
// =========================================================================

func (r *MessageReactionsRepository) CreateReaction(ctx context.Context, reaction *entity.MessageReaction) error {
	// 1. Fail-fast validation
	if reaction == nil {
		return fmt.Errorf("reaction payload is nil")
	}

	var emptyUUID gocql.UUID
	if reaction.ConversationID == "" || reaction.MessageID == emptyUUID || reaction.UserID == emptyUUID {
		return fmt.Errorf("missing primary key components: conversation_id, message_id, and user_id are required")
	}

	// 2. Bảo vệ Context: Thao tác Ghi không được ngắt giữa chừng dù HTTP Request bị hủy
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.MessageReaction{}.TableName()

	// 3. Query với khai báo rõ từng cột (Explicit columns - best practice)
	// Cassandra INSERT là Upsert: nếu row đã tồn tại, reaction_code sẽ bị ghi đè (last-write-wins)
	query := fmt.Sprintf(`
		INSERT INTO %s (conversation_id, message_id, user_id, reaction_code, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	// 4. Thực thi query
	err := r.session.Query(query,
		reaction.ConversationID,
		reaction.MessageID,
		reaction.UserID,
		reaction.ReactionCode,
		reaction.CreatedAt,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to create reaction for user %s on message %s: %w",
			reaction.UserID.String(), reaction.MessageID.String(), err)
	}

	return nil
}

func (r *MessageReactionsRepository) CreateBulkReactions(ctx context.Context, reactions []*entity.MessageReaction) (int64, []*cassandraErrors.MessageReactionBulkError, error) {
	// 1. Fail-fast validation
	if len(reactions) == 0 {
		return 0, nil, nil
	}

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type taskResult struct {
		reaction *entity.MessageReaction
		err      error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để ngăn chặn deadlock
	resultCh := make(chan taskResult, len(reactions))

	tableName := entity.MessageReaction{}.TableName()
	query := fmt.Sprintf(`
		INSERT INTO %s (conversation_id, message_id, user_id, reaction_code, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	var emptyUUID gocql.UUID

	// 3. Phân phối task vào Worker Pool
	for i, item := range reactions {
		// Xử lý an toàn con trỏ nil
		if item == nil {
			resultCh <- taskResult{
				reaction: nil,
				err:      fmt.Errorf("reaction pointer at index %d is nil", i),
			}
			continue
		}

		// Kiểm tra khóa chính trước khi tốn tài nguyên gọi DB
		if item.ConversationID == "" || item.MessageID == emptyUUID || item.UserID == emptyUUID {
			resultCh <- taskResult{
				reaction: item,
				err:      fmt.Errorf("missing primary key components for insert at index %d", i),
			}
			continue
		}

		// Copy biến cục bộ an toàn cho Goroutine (Go < 1.22)
		payload := item

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT khi ghi
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				payload.ConversationID,
				payload.MessageID,
				payload.UserID,
				payload.ReactionCode,
				payload.CreatedAt,
			).WithContext(safeCtx).Exec()

			resultCh <- taskResult{reaction: payload, err: execErr}
		})

		// Xử lý khi worker pool bị đầy hoặc context gốc bị cancel/timeout
		if err != nil {
			resultCh <- taskResult{
				reaction: payload,
				err:      fmt.Errorf("worker pool rejected insert task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa: Chờ tất cả Worker hoàn thành và đóng Channel an toàn
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free vì chỉ 1 main goroutine đọc channel)
	var successCount int64
	var bulkErrors []*cassandraErrors.MessageReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			convID := "unknown_conversation"
			msgID := "unknown_message"
			userID := "unknown_user"

			if res.reaction != nil {
				convID = res.reaction.ConversationID
				msgID = res.reaction.MessageID.String()
				userID = res.reaction.UserID.String()
			}

			bulkErrors = append(bulkErrors, &cassandraErrors.MessageReactionBulkError{
				ConversationID: convID,
				MessageID:      msgID,
				UserID:         userID,
				Error:          res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái trả về (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk create reactions completed with errors: %d/%d success, %d failed",
			successCount, len(reactions), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}

// =========================================================================
// UPDATE
// =========================================================================

func (r *MessageReactionsRepository) UpdateReaction(ctx context.Context, reaction *entity.MessageReaction) error {
	// 1. Fail-fast validation
	if reaction == nil {
		return fmt.Errorf("reaction payload is nil")
	}

	var emptyUUID gocql.UUID
	// Kiểm tra BẮT BUỘC phải có đủ bộ Primary Key: ((conversation_id, message_id), user_id)
	if reaction.ConversationID == "" || reaction.MessageID == emptyUUID || reaction.UserID == emptyUUID {
		return fmt.Errorf("missing primary key components: conversation_id, message_id, and user_id are required for update")
	}

	// 2. Bảo vệ Context: Lệnh Ghi (Update) không được ngắt giữa chừng dù HTTP Request bị hủy
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.MessageReaction{}.TableName()

	// 3. Query chuẩn bị:
	// - SET các trường dữ liệu (Non-Primary Key columns)
	// - WHERE bắt buộc chứa đủ Composite Partition Key và Clustering Key
	query := fmt.Sprintf(`
		UPDATE %s
		SET reaction_code = ?, created_at = ?
		WHERE conversation_id = ? AND message_id = ? AND user_id = ?
	`, tableName)

	// 4. Thực thi truy vấn
	err := r.session.Query(query,
		// SET
		reaction.ReactionCode,
		reaction.CreatedAt,
		// WHERE (Primary Key)
		reaction.ConversationID,
		reaction.MessageID,
		reaction.UserID,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to update reaction for user %s on message %s: %w",
			reaction.UserID.String(), reaction.MessageID.String(), err)
	}

	return nil
}

func (r *MessageReactionsRepository) UpdateBulkReactions(ctx context.Context, reactions []*entity.MessageReaction) (int64, []*cassandraErrors.MessageReactionBulkError, error) {
	// 1. Fail-fast validation
	if len(reactions) == 0 {
		return 0, nil, nil
	}

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type taskResult struct {
		reaction *entity.MessageReaction
		err      error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để ngăn chặn deadlock
	resultCh := make(chan taskResult, len(reactions))

	tableName := entity.MessageReaction{}.TableName()
	query := fmt.Sprintf(`
		UPDATE %s
		SET reaction_code = ?, created_at = ?
		WHERE conversation_id = ? AND message_id = ? AND user_id = ?
	`, tableName)

	var emptyUUID gocql.UUID

	// 3. Phân phối task vào Worker Pool
	for i, item := range reactions {
		if item == nil {
			resultCh <- taskResult{
				reaction: nil,
				err:      fmt.Errorf("reaction pointer at index %d is nil", i),
			}
			continue
		}

		if item.ConversationID == "" || item.MessageID == emptyUUID || item.UserID == emptyUUID {
			resultCh <- taskResult{
				reaction: item,
				err:      fmt.Errorf("missing primary key components for update at index %d", i),
			}
			continue
		}

		// Copy biến cục bộ để Closure trong Goroutine trỏ đúng vùng nhớ (Go < 1.22)
		payload := item

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT KHI GHI
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				// SET
				payload.ReactionCode,
				payload.CreatedAt,
				// WHERE
				payload.ConversationID,
				payload.MessageID,
				payload.UserID,
			).WithContext(safeCtx).Exec()

			resultCh <- taskResult{reaction: payload, err: execErr}
		})

		// Xử lý khi Worker Pool từ chối (Context gốc đã Timeout/Cancel)
		if err != nil {
			resultCh <- taskResult{
				reaction: payload,
				err:      fmt.Errorf("worker pool rejected update task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa: Chờ tất cả Worker hoàn thành và đóng Channel an toàn
	r.pool.Wait()
	close(resultCh)

	// 5. Thu gom kết quả (Lock-free)
	var successCount int64
	var bulkErrors []*cassandraErrors.MessageReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			convID := "unknown_conversation"
			msgID := "unknown_message"
			userID := "unknown_user"

			if res.reaction != nil {
				convID = res.reaction.ConversationID
				msgID = res.reaction.MessageID.String()
				userID = res.reaction.UserID.String()
			}

			bulkErrors = append(bulkErrors, &cassandraErrors.MessageReactionBulkError{
				ConversationID: convID,
				MessageID:      msgID,
				UserID:         userID,
				Error:          res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk update reactions completed with errors: %d/%d success, %d failed",
			successCount, len(reactions), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}

// =========================================================================
// DELETE
// =========================================================================

func (r *MessageReactionsRepository) DeleteReaction(ctx context.Context, conversationID string, messageID string, userID string) error {
	// 1. Fail-fast validation
	if conversationID == "" || messageID == "" || userID == "" {
		return fmt.Errorf("conversationID, messageID, and userID cannot be empty")
	}

	// Parse UUID để đảm bảo tính hợp lệ trước khi gọi DB
	msgUUID, err := gocql.ParseUUID(messageID)
	if err != nil {
		return fmt.Errorf("invalid messageID format: %w", err)
	}

	userUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format: %w", err)
	}

	// 2. Bảo vệ Context: DELETE trong Cassandra là ghi Tombstone xuống đĩa
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.MessageReaction{}.TableName()

	// 3. Query: Yêu cầu chính xác 100% Primary Key để tránh xóa nhầm
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE conversation_id = ? AND message_id = ? AND user_id = ?
	`, tableName)

	// 4. Thực thi
	if err := r.session.Query(query, conversationID, msgUUID, userUUID).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete reaction for user %s on message %s: %w", userID, messageID, err)
	}

	return nil
}

func (r *MessageReactionsRepository) DeleteBulkReactions(ctx context.Context, conversationID string, messageID string, userIDs []string) (int64, []*cassandraErrors.MessageReactionBulkError, error) {
	// 1. Fail-fast validation
	if conversationID == "" || messageID == "" {
		return 0, nil, fmt.Errorf("conversationID and messageID cannot be empty")
	}
	if len(userIDs) == 0 {
		return 0, nil, nil
	}

	// Parse messageID một lần duy nhất (Fail-fast: tránh parse lặp lại trong goroutine)
	msgUUID, err := gocql.ParseUUID(messageID)
	if err != nil {
		return 0, nil, fmt.Errorf("bulk delete aborted - invalid messageID format: %w", err)
	}

	// Chuyển đổi toàn bộ mảng userID sang UUID trước để bắt lỗi sớm
	uuids := make([]gocql.UUID, 0, len(userIDs))
	for _, idStr := range userIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk delete aborted - invalid userID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. Chuẩn bị Struct và Channel an toàn cho Worker Pool
	type taskResult struct {
		userID gocql.UUID
		err    error
	}
	resultCh := make(chan taskResult, len(uuids))

	tableName := entity.MessageReaction{}.TableName()
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE conversation_id = ? AND message_id = ? AND user_id = ?
	`, tableName)

	// 3. Phân phối nhiệm vụ vào Worker Pool
	for _, uid := range uuids {
		// Copy biến cục bộ an toàn cho Goroutine (Go < 1.22)
		payloadUID := uid

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT vì DELETE là thao tác ghi Tombstone
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query, conversationID, msgUUID, payloadUID).WithContext(safeCtx).Exec()

			resultCh <- taskResult{userID: payloadUID, err: execErr}
		})

		// Xử lý khi Worker Pool từ chối task (Context gốc bị cancel/timeout)
		if err != nil {
			resultCh <- taskResult{
				userID: payloadUID,
				err:    fmt.Errorf("worker pool rejected delete task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa các Goroutine
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*cassandraErrors.MessageReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &cassandraErrors.MessageReactionBulkError{
				ConversationID: conversationID,
				MessageID:      messageID,
				UserID:         res.userID.String(),
				Error:          res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái cuối cùng
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete reactions completed with errors: %d/%d success, %d failed",
			successCount, len(userIDs), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}

func (r *MessageReactionsRepository) DeleteReactionsByMessageID(ctx context.Context, conversationID string, messageID string) error {
	// 1. Fail-fast validation
	if conversationID == "" || messageID == "" {
		return fmt.Errorf("conversationID and messageID cannot be empty")
	}

	msgUUID, err := gocql.ParseUUID(messageID)
	if err != nil {
		return fmt.Errorf("invalid messageID format: %w", err)
	}

	// 2. Bảo vệ Context: Partition Delete vẫn là thao tác ghi Tombstone
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.MessageReaction{}.TableName()

	// 3. Partition Delete: Chỉ cần Composite Partition Key để xóa TOÀN BỘ reaction của 1 tin nhắn.
	// Đây là thao tác O(1) trong Cassandra — cực kỳ nhanh.
	query := fmt.Sprintf("DELETE FROM %s WHERE conversation_id = ? AND message_id = ?", tableName)

	// 4. Thực thi
	if err := r.session.Query(query, conversationID, msgUUID).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete all reactions of message %s in conversation %s: %w",
			messageID, conversationID, err)
	}

	return nil
}

func (r *MessageReactionsRepository) DeleteBulkReactionsByMessageID(ctx context.Context, conversationID string, messageIDs []string) (int64, []*cassandraErrors.MessageReactionBulkError, error) {
	// 1. Fail-fast validation
	if conversationID == "" {
		return 0, nil, fmt.Errorf("conversationID cannot be empty")
	}
	if len(messageIDs) == 0 {
		return 0, nil, nil
	}

	// Chuyển đổi toàn bộ mảng UUID trước để bắt lỗi cấu trúc sớm (Fail-fast)
	uuids := make([]gocql.UUID, 0, len(messageIDs))
	for _, idStr := range messageIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk delete aborted - invalid messageID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. Chuẩn bị Struct và Channel an toàn
	type taskResult struct {
		messageID gocql.UUID
		err       error
	}
	resultCh := make(chan taskResult, len(uuids))

	tableName := entity.MessageReaction{}.TableName()
	// Mỗi lần xóa là 1 Partition Delete (xóa toàn bộ reaction của 1 message)
	query := fmt.Sprintf("DELETE FROM %s WHERE conversation_id = ? AND message_id = ?", tableName)

	// 3. Phân phối nhiệm vụ: Mỗi Worker chịu trách nhiệm 1 Partition Delete độc lập
	for _, msgID := range uuids {
		// Copy biến cục bộ an toàn cho Goroutine (Go < 1.22)
		payloadMsgID := msgID

		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query, conversationID, payloadMsgID).WithContext(safeCtx).Exec()

			resultCh <- taskResult{messageID: payloadMsgID, err: execErr}
		})

		if err != nil {
			resultCh <- taskResult{
				messageID: payloadMsgID,
				err:       fmt.Errorf("worker pool rejected partition delete task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*cassandraErrors.MessageReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &cassandraErrors.MessageReactionBulkError{
				ConversationID: conversationID,
				MessageID:      res.messageID.String(),
				UserID:         "all",
				Error:          res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái cuối cùng
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete reactions by messageID completed with errors: %d/%d success, %d failed",
			successCount, len(messageIDs), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}

// =========================================================================
// READ
// =========================================================================

func (r *MessageReactionsRepository) GetReactionsByMessageID(ctx context.Context, conversationID string, messageID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if conversationID == "" || messageID == "" {
		return nil, fmt.Errorf("conversationID and messageID cannot be empty")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	// Xác thực UUID sớm để tránh gọi DB vô ích
	if _, err := gocql.ParseUUID(messageID); err != nil {
		return nil, fmt.Errorf("invalid messageID format: %w", err)
	}

	// 2. Check Cache (Chỉ hit cache khi người dùng tải trang đầu tiên)
	cacheKeyPrefix := fmt.Sprintf("reactions_cache_convid_%s_msgid_%s", conversationID, messageID)
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			cacheKeyPrefix,
			cacheKeyPrefix + "_nextcursor",
			cacheKeyPrefix + "_hasnext",
			cacheKeyPrefix + "_limit",
		})

		if err == nil && datacache != nil {
			if reactionList, ok := datacache.([]*entity.MessageReaction); ok {
				return &dto.PaginationRes{
					Data:       reactionList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	// 3. Decode Cursor (Trích xuất Page State của Cassandra)
	pageState, err := utils.DecodeCursorCassandra(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	// 4. Chuẩn bị Query cho Cassandra
	// Query trực tiếp bằng Composite Partition Key → không cần ALLOW FILTERING
	tableName := entity.MessageReaction{}.TableName()

	// Khai báo rõ các cột để tối ưu Serialize/Deserialize
	query := fmt.Sprintf(`
		SELECT conversation_id, message_id, user_id, reaction_code, created_at
		FROM %s WHERE conversation_id = ? AND message_id = ?
	`, tableName)

	// LƯU Ý: Không dùng context.WithoutCancel ở đây.
	// Truyền ctx gốc để DB ngừng xử lý nếu client ngắt kết nối.
	q := r.session.Query(query, conversationID, messageID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 5. Thực thi và Map dữ liệu
	iter := q.Iter()
	scanner := iter.Scanner()
	var reactions []*entity.MessageReaction

	for scanner.Next() {
		var rx entity.MessageReaction
		err := scanner.Scan(
			&rx.ConversationID,
			&rx.MessageID,
			&rx.UserID,
			&rx.ReactionCode,
			&rx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reaction for message %s in conversation %s: %w",
				messageID, conversationID, err)
		}

		// Append con trỏ an toàn (Go >= 1.22)
		reactions = append(reactions, &rx)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for message %s in conversation %s: %w",
			messageID, conversationID, err)
	}

	// 6. Xử lý Next Cursor và HasNext
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	// Trả về slice rỗng [] thay vì nil để Frontend không bị lỗi map null
	if reactions == nil {
		reactions = []*entity.MessageReaction{}
	}

	// 7. Cache kết quả cho trang đầu tiên (Chỉ lưu nếu có data)
	if cursor == "" && len(reactions) > 0 {
		items := map[string]any{
			cacheKeyPrefix:                 reactions,
			cacheKeyPrefix + "_nextcursor": nextCursorStr,
			cacheKeyPrefix + "_hasnext":    hasNext,
			cacheKeyPrefix + "_limit":      limit,
		}

		if err := r.redisRepo.CustomizeSetCache(ctx, items); err != nil {
			// Chỉ log lỗi, không chặn main flow vì redis lỗi thì app vẫn phải chạy
			fmt.Printf("failed to set cache for reactions of message %s: %v\n", messageID, err)
		}
	}

	// 8. Trả về kết quả
	return &dto.PaginationRes{
		Data:       reactions,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}

func (r *MessageReactionsRepository) GetReactionByUser(ctx context.Context, conversationID string, messageID string, userID string) (*entity.MessageReaction, error) {
	// 1. Fail-fast validation
	if conversationID == "" || messageID == "" || userID == "" {
		return nil, fmt.Errorf("conversationID, messageID, and userID cannot be empty")
	}

	// Parse UUID để đảm bảo tính hợp lệ trước khi gọi DB
	msgUUID, err := gocql.ParseUUID(messageID)
	if err != nil {
		return nil, fmt.Errorf("invalid messageID format: %w", err)
	}

	userUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format: %w", err)
	}

	// 2. Truyền ctx gốc để DB ngừng xử lý nếu client ngắt kết nối (đây là Read)
	tableName := entity.MessageReaction{}.TableName()

	// 3. Single-row lookup bằng đủ Primary Key → truy vấn O(1), không cần ALLOW FILTERING
	query := fmt.Sprintf(`
		SELECT conversation_id, message_id, user_id, reaction_code, created_at
		FROM %s WHERE conversation_id = ? AND message_id = ? AND user_id = ?
	`, tableName)

	// 4. Thực thi và map kết quả
	var rx entity.MessageReaction
	err = r.session.Query(query, conversationID, msgUUID, userUUID).
		WithContext(ctx).
		Scan(
			&rx.ConversationID,
			&rx.MessageID,
			&rx.UserID,
			&rx.ReactionCode,
			&rx.CreatedAt,
		)

	// 5. Xử lý kết quả
	if err != nil {
		// ErrNotFound: User chưa thả reaction → trả về nil, nil (không phải lỗi hệ thống)
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get reaction for user %s on message %s: %w",
			userID, messageID, err)
	}

	return &rx, nil
}
func (r *MessageReactionsRepository) DeleteReactionByConversationID(ctx context.Context, conversationID string) error {
	// 1. Fail-fast validation
	if conversationID == "" {
		return fmt.Errorf("conversationID cannot be empty")
	}

	// 2. Bảo vệ Context: Partition Delete vẫn là thao tác ghi Tombstone
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.MessageReaction{}.TableName()

	// 3. Partition Delete: Xóa toàn bộ reaction của 1 conversation.
	// Đây là thao tác O(1) trong Cassandra — cực kỳ nhanh.
	query := fmt.Sprintf("DELETE FROM %s WHERE conversation_id = ?", tableName)

	// 4. Thực thi
	if err := r.session.Query(query, conversationID).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete all reactions of conversation %s: %w",
			conversationID, err)
	}

	return nil
}

func (r *MessageReactionsRepository) DeleteBulkReactionByConversationID(ctx context.Context, conversationIDs []string) error {
	tableName := entity.MessageReaction{}.TableName()
	query := fmt.Sprintf("DELETE FROM %s WHERE conversation_id = ?", tableName)

	for _, convID := range conversationIDs {
		if convID == "" {
			continue // Skip invalid conversationID
		}

		safeCtx := context.WithoutCancel(ctx)

		if err := r.session.Query(query, convID).WithContext(safeCtx).Exec(); err != nil {
			fmt.Printf("failed to delete reactions for conversation %s: %v\n", convID, err)
			// Log lỗi nhưng không dừng toàn bộ quá trình
		}
	}

	return nil
}

func (r *MessageReactionsRepository) UpdateOrInsertReaction(ctx context.Context, reaction *entity.MessageReaction) error {
	tableName := entity.MessageReaction{}.TableName()

	// 1. Fail-fast validation
	if reaction == nil {
		return fmt.Errorf("reaction payload is nil")
	}

	var emptyUUID gocql.UUID
	if reaction.ConversationID == "" || reaction.MessageID == emptyUUID || reaction.UserID == emptyUUID {
		return fmt.Errorf("missing primary key components: conversation_id, message_id, and user_id are required")
	}

	// 2. Bảo vệ Context: Thao tác Ghi không được ngắt giữa chừng dù HTTP Request bị hủy
	safeCtx := context.WithoutCancel(ctx)

	// 3. Query với khai báo rõ từng cột (Explicit columns - best practice)
	// Cassandra INSERT là Upsert: nếu row đã tồn tại, reaction_code sẽ bị ghi đè (last-write-wins)
	query := fmt.Sprintf(`
		INSERT INTO %s (conversation_id, message_id, user_id, reaction_code, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	// 4. Thực thi query
	err := r.session.Query(query,
		reaction.ConversationID,
		reaction.MessageID,
		reaction.UserID,
		reaction.ReactionCode,
		reaction.CreatedAt,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to upsert reaction for user %s on message %s: %w",
			reaction.UserID.String(), reaction.MessageID.String(), err)
	}

	return nil
}
