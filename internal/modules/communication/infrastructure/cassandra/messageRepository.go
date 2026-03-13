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

type MessageRepository struct {
	session   *gocql.Session
	pool      IRepositoryShare.IWorkerPool
	redisRepo IRepositoryShare.IRedis
}

func NewMessageRepository(session *gocql.Session, pool IRepositoryShare.IWorkerPool, redisRepo IRepositoryShare.IRedis) *MessageRepository {
	return &MessageRepository{
		session:   session,
		pool:      pool,
		redisRepo: redisRepo,
	}
}

// =========================================================================
// CREATE
// =========================================================================

func (r *MessageRepository) CreateMessage(ctx context.Context, message *entity.Message) error {
	// 1. Fail-fast validation
	if message == nil {
		return fmt.Errorf("message payload is nil")
	}

	var emptyUUID gocql.UUID
	if message.ConversationID == "" || message.MessageID == emptyUUID {
		return fmt.Errorf("missing primary key components: conversation_id and message_id are required")
	}

	// 2. Bảo vệ Context: Thao tác Ghi không được ngắt giữa chừng dù HTTP Request bị hủy
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.Message{}.TableName()

	// 3. Query với khai báo rõ từng cột (Explicit columns - best practice)
	// Gocql tự động xử lý *gocql.UUID pointer thành null nếu nil
	query := fmt.Sprintf(`
		INSERT INTO %s (
			conversation_id, bucket, message_id,
			sender_id, type, content, attachments,
			is_edited, reply_to_message_id, story_ref_id,
			is_revoked, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	// 4. Thực thi query
	err := r.session.Query(query,
		message.ConversationID,
		message.Bucket,
		message.MessageID,
		message.SenderID,
		message.Type,
		message.Content,
		message.Attachments,
		message.IsEdited,
		message.ReplyToMessageID,
		message.StoryRefID,
		message.IsRevoked,
		message.CreatedAt,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to create message %s in conversation %s: %w", message.MessageID.String(), message.ConversationID, err)
	}

	return nil
}

func (r *MessageRepository) CreateBulkMessages(ctx context.Context, messages []*entity.Message) (int64, []*cassandraErrors.MessageBulkError, error) {
	// 1. Fail-fast validation
	if len(messages) == 0 {
		return 0, nil, nil
	}

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type taskResult struct {
		message *entity.Message
		err     error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để ngăn chặn deadlock
	resultCh := make(chan taskResult, len(messages))

	tableName := entity.Message{}.TableName()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			conversation_id, bucket, message_id,
			sender_id, type, content, attachments,
			is_edited, reply_to_message_id, story_ref_id,
			is_revoked, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	var emptyUUID gocql.UUID

	// 3. Phân phối task vào Worker Pool
	for i, item := range messages {
		// Xử lý an toàn con trỏ nil
		if item == nil {
			resultCh <- taskResult{
				message: nil,
				err:     fmt.Errorf("message pointer at index %d is nil", i),
			}
			continue
		}

		// Kiểm tra khóa chính trước khi tốn tài nguyên gọi DB
		if item.ConversationID == "" || item.MessageID == emptyUUID {
			resultCh <- taskResult{
				message: item,
				err:     fmt.Errorf("missing primary key components for insert at index %d", i),
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
				payload.Bucket,
				payload.MessageID,
				payload.SenderID,
				payload.Type,
				payload.Content,
				payload.Attachments,
				payload.IsEdited,
				payload.ReplyToMessageID,
				payload.StoryRefID,
				payload.IsRevoked,
				payload.CreatedAt,
			).WithContext(safeCtx).Exec()

			resultCh <- taskResult{
				message: payload,
				err:     execErr,
			}
		})

		// Xử lý khi worker pool bị đầy hoặc context gốc bị cancel/timeout
		if err != nil {
			resultCh <- taskResult{
				message: payload,
				err:     fmt.Errorf("worker pool rejected insert task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa: Chờ tất cả Worker hoàn thành và đóng Channel an toàn
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free vì chỉ 1 main goroutine đọc channel)
	var successCount int64
	var bulkErrors []*cassandraErrors.MessageBulkError

	for res := range resultCh {
		if res.err != nil {
			convID := "unknown_conversation"
			msgID := "unknown_message"

			if res.message != nil {
				convID = res.message.ConversationID
				msgID = res.message.MessageID.String()
			}

			bulkErrors = append(bulkErrors, &cassandraErrors.MessageBulkError{
				ConversationID: convID,
				MessageID:      msgID,
				Error:          res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái trả về (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk insert messages completed with errors: %d/%d success, %d failed", successCount, len(messages), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}

// =========================================================================
// UPDATE
// =========================================================================

func (r *MessageRepository) UpdateMessage(ctx context.Context, message *entity.Message) error {
	// 1. Fail-fast validation
	if message == nil {
		return fmt.Errorf("message payload is nil")
	}

	var emptyUUID gocql.UUID
	// Kiểm tra BẮT BUỘC phải có đủ bộ Primary Key: (conversation_id, bucket, message_id)
	if message.ConversationID == "" || message.MessageID == emptyUUID {
		return fmt.Errorf("missing primary key components: conversation_id and message_id are required for update")
	}

	// 2. Bảo vệ Context: Lệnh Ghi (Update) không được ngắt giữa chừng dù HTTP Request bị hủy
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.Message{}.TableName()

	// 3. Query chuẩn bị:
	// - SET các trường dữ liệu (Non-Primary Key columns)
	// - WHERE bắt buộc chứa đủ Partition Key và Clustering Key
	query := fmt.Sprintf(`
		UPDATE %s
		SET sender_id = ?, type = ?, content = ?, attachments = ?,
		    is_edited = ?, reply_to_message_id = ?, story_ref_id = ?,
		    is_revoked = ?, created_at = ?
		WHERE conversation_id = ? AND bucket = ? AND message_id = ?
	`, tableName)

	// 4. Thực thi truy vấn
	err := r.session.Query(query,
		// SET
		message.SenderID,
		message.Type,
		message.Content,
		message.Attachments,
		message.IsEdited,
		message.ReplyToMessageID,
		message.StoryRefID,
		message.IsRevoked,
		message.CreatedAt,
		// WHERE (Primary Key)
		message.ConversationID,
		message.Bucket,
		message.MessageID,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to update message %s in conversation %s: %w", message.MessageID.String(), message.ConversationID, err)
	}

	return nil
}

func (r *MessageRepository) UpdateBulkMessages(ctx context.Context, messages []*entity.Message) (int64, []*cassandraErrors.MessageBulkError, error) {
	// 1. Fail-fast validation
	if len(messages) == 0 {
		return 0, nil, nil
	}

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type taskResult struct {
		message *entity.Message
		err     error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để ngăn chặn deadlock
	resultCh := make(chan taskResult, len(messages))

	tableName := entity.Message{}.TableName()
	query := fmt.Sprintf(`
		UPDATE %s
		SET sender_id = ?, type = ?, content = ?, attachments = ?,
		    is_edited = ?, reply_to_message_id = ?, story_ref_id = ?,
		    is_revoked = ?, created_at = ?
		WHERE conversation_id = ? AND bucket = ? AND message_id = ?
	`, tableName)

	var emptyUUID gocql.UUID

	// 3. Phân phối task vào Worker Pool
	for i, item := range messages {
		if item == nil {
			resultCh <- taskResult{
				message: nil,
				err:     fmt.Errorf("message pointer at index %d is nil", i),
			}
			continue
		}

		if item.ConversationID == "" || item.MessageID == emptyUUID {
			resultCh <- taskResult{
				message: item,
				err:     fmt.Errorf("missing primary key components for update at index %d", i),
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
				payload.SenderID,
				payload.Type,
				payload.Content,
				payload.Attachments,
				payload.IsEdited,
				payload.ReplyToMessageID,
				payload.StoryRefID,
				payload.IsRevoked,
				payload.CreatedAt,
				// WHERE
				payload.ConversationID,
				payload.Bucket,
				payload.MessageID,
			).WithContext(safeCtx).Exec()

			resultCh <- taskResult{
				message: payload,
				err:     execErr,
			}
		})

		// Xử lý khi Worker Pool từ chối (Context gốc đã Timeout/Cancel)
		if err != nil {
			resultCh <- taskResult{
				message: payload,
				err:     fmt.Errorf("worker pool rejected update task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa: Chờ tất cả Worker hoàn thành và đóng Channel an toàn
	r.pool.Wait()
	close(resultCh)

	// 5. Thu gom kết quả (Lock-free)
	var successCount int64
	var bulkErrors []*cassandraErrors.MessageBulkError

	for res := range resultCh {
		if res.err != nil {
			convID := "unknown_conversation"
			msgID := "unknown_message"

			if res.message != nil {
				convID = res.message.ConversationID
				msgID = res.message.MessageID.String()
			}

			bulkErrors = append(bulkErrors, &cassandraErrors.MessageBulkError{
				ConversationID: convID,
				MessageID:      msgID,
				Error:          res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk update messages completed with errors: %d/%d success, %d failed", successCount, len(messages), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}

// =========================================================================
// DELETE
// =========================================================================

func (r *MessageRepository) DeleteMessage(ctx context.Context, conversationID string, bucket int, messageID string) error {
	// 1. Fail-fast validation
	if conversationID == "" || messageID == "" {
		return fmt.Errorf("conversationID and messageID cannot be empty")
	}

	// Parse messageID sang UUID để đảm bảo tính hợp lệ trước khi gọi DB
	msgUUID, err := gocql.ParseUUID(messageID)
	if err != nil {
		return fmt.Errorf("invalid messageID format: %w", err)
	}

	// 2. Bảo vệ Context: DELETE trong Cassandra là ghi Tombstone xuống đĩa
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.Message{}.TableName()

	// 3. Query: Yêu cầu chính xác 100% Primary Key (conversation_id, bucket, message_id)
	query := fmt.Sprintf("DELETE FROM %s WHERE conversation_id = ? AND bucket = ? AND message_id = ?", tableName)

	// 4. Thực thi
	if err := r.session.Query(query, conversationID, bucket, msgUUID).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete message %s in conversation %s: %w", messageID, conversationID, err)
	}

	return nil
}

func (r *MessageRepository) DeleteBulkMessages(ctx context.Context, conversationID string, bucket int, messageIDs []string) (int64, []*cassandraErrors.MessageBulkError, error) {
	// 1. Fail-fast validation
	if conversationID == "" || len(messageIDs) == 0 {
		return 0, nil, nil
	}

	// Chuyển đổi toàn bộ mảng messageID sang UUID trước để bắt lỗi sớm (Fail-fast)

	// 2. Chuẩn bị Struct và Channel an toàn cho Worker Pool
	type taskResult struct {
		message *entity.Message
		err     error
	}
	resultCh := make(chan taskResult, len(messageIDs))

	tableName := entity.Message{}.TableName()
	query := fmt.Sprintf("DELETE FROM %s WHERE conversation_id = ? AND bucket = ? AND message_id = ?", tableName)

	// 3. Phân phối nhiệm vụ vào Worker Pool
	for _, msgID := range messageIDs {
		// Copy biến cục bộ an toàn cho Goroutine (Go < 1.22)
		// Copy giá trị UUID vào biến cục bộ
		finalid, err := gocql.ParseUUID(msgID)
		payloadID := finalid
		if err != nil {
			resultCh <- taskResult{
				message: nil,
				err:     fmt.Errorf("invalid messageID format for %s: %w", msgID, err),
			}
			continue
		}
		err = r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT vì DELETE là thao tác ghi Tombstone
			safeCtx := context.WithoutCancel(ctx)
			execErr := r.session.Query(query, conversationID, bucket, payloadID).WithContext(safeCtx).Exec()
			resultCh <- taskResult{
				message: &entity.Message{MessageID: payloadID},
				err:     execErr,
			}
		})

		// Xử lý khi Worker Pool từ chối task (Context gốc bị cancel/timeout)
		if err != nil {
			resultCh <- taskResult{
				message: nil,
				err:     fmt.Errorf("worker pool rejected delete task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa các Goroutine
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*cassandraErrors.MessageBulkError

	for res := range resultCh {
		if res.err != nil {
			ConversationID := "unknown_conversation"
			MessageID := "unknown_message_id"

			if res.message != nil {
				MessageID = res.message.MessageID.String()
				ConversationID = conversationID
			}

			bulkErrors = append(bulkErrors, &cassandraErrors.MessageBulkError{
				ConversationID: ConversationID,
				MessageID:      MessageID,
				Error:          res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái cuối cùng
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete messages completed with errors: %d/%d success, %d failed", successCount, len(messageIDs), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}

// =========================================================================
// READ
// =========================================================================

func (r *MessageRepository) GetMessagesByConversationIDs(ctx context.Context, conversationID string, bucket int, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if conversationID == "" {
		return nil, fmt.Errorf("conversationID cannot be empty")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	// 2. Check Cache (Chỉ hit cache khi lấy trang đầu tiên)
	cacheKeyPrefix := fmt.Sprintf("messages_cache_convid_%s_bucket_%d", conversationID, bucket)
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			cacheKeyPrefix,
			cacheKeyPrefix + "_nextcursor",
			cacheKeyPrefix + "_hasnext",
			cacheKeyPrefix + "_limit",
		})

		if err == nil && datacache != nil {
			if msgList, ok := datacache.([]*entity.Message); ok {
				return &dto.PaginationRes{
					Data:       msgList,
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
	// Lọc theo CADT Partition Key (conversation_id, bucket) để tránh Full Cluster Scan
	tableName := entity.Message{}.TableName()

	// Khai báo rõ các cột để tối ưu Serialize/Deserialize
	query := fmt.Sprintf(`
		SELECT conversation_id, bucket, message_id,
		       sender_id, type, content, attachments,
		       is_edited, reply_to_message_id, story_ref_id,
		       is_revoked, created_at
		FROM %s WHERE conversation_id = ? AND bucket = ?
	`, tableName)

	// LƯU Ý: Không dùng context.WithoutCancel ở đây.
	// Truyền ctx gốc để DB ngừng xử lý nếu client ngắt kết nối.
	q := r.session.Query(query, conversationID, bucket).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 5. Thực thi và Map dữ liệu
	iter := q.Iter()
	scanner := iter.Scanner()
	var messages []*entity.Message

	for scanner.Next() {
		var m entity.Message
		// Gocql tự động map mảng LIST<TEXT> vào []string (attachments)
		// và xử lý nullable UUID (*gocql.UUID) thành nil khi NULL trong DB
		err := scanner.Scan(
			&m.ConversationID,
			&m.Bucket,
			&m.MessageID,
			&m.SenderID,
			&m.Type,
			&m.Content,
			&m.Attachments,
			&m.IsEdited,
			&m.ReplyToMessageID,
			&m.StoryRefID,
			&m.IsRevoked,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message in conversation %s bucket %d: %w", conversationID, bucket, err)
		}

		// Append con trỏ an toàn (Go >= 1.22)
		messages = append(messages, &m)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for conversation %s bucket %d: %w", conversationID, bucket, err)
	}

	// 6. Xử lý Next Cursor và HasNext
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	// Trả về slice rỗng [] thay vì nil để Frontend không bị lỗi map null
	if messages == nil {
		messages = []*entity.Message{}
	}

	// 7. Cache kết quả cho trang đầu tiên (Chỉ lưu nếu có data)
	if cursor == "" && len(messages) > 0 {
		items := map[string]any{
			cacheKeyPrefix:                 messages,
			cacheKeyPrefix + "_nextcursor": nextCursorStr,
			cacheKeyPrefix + "_hasnext":    hasNext,
			cacheKeyPrefix + "_limit":      limit,
		}

		if err := r.redisRepo.CustomizeSetCache(ctx, items); err != nil {
			// Chỉ log lỗi, không chặn main flow vì redis lỗi thì app vẫn phải chạy
			fmt.Printf("failed to set cache for messages pagination conversation %s bucket %d: %v\n", conversationID, bucket, err)
		}
	}

	// 8. Trả về kết quả
	return &dto.PaginationRes{
		Data:       messages,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *MessageRepository) DeleteMessagesByConversationID(ctx context.Context, conversationID string) error {
	safeCtx := context.WithoutCancel(ctx)
	tableName := entity.Message{}.TableName()

	// 3. Partition Delete: Xóa toàn bộ message của 1 conversation.
	// Đây là thao tác O(1) trong Cassandra — cực kỳ nhanh.
	query := fmt.Sprintf("DELETE FROM %s WHERE conversation_id = ?", tableName)

	// 4. Thực thi
	if err := r.session.Query(query, conversationID).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete all messages of conversation %s: %w",
			conversationID, err)
	}

	return nil
}
func (r *MessageRepository) GetMessagesByConversationID(ctx context.Context, conversationID string, bucket int) (*entity.Message, error) {
	// 1. Fail-fast validation
	if conversationID == "" {
		return nil, fmt.Errorf("conversationID cannot be empty")
	}

	// 2. Chuẩn bị Query cho Cassandra
	tableName := entity.Message{}.TableName()

	query := fmt.Sprintf(`
		SELECT conversation_id, bucket, message_id,
		       sender_id, type, content, attachments,
		       is_edited, reply_to_message_id, story_ref_id,
		       is_revoked, created_at
		FROM %s WHERE conversation_id = ? AND bucket = ?
	`, tableName)

	var message entity.Message

	err := r.session.Query(query, conversationID, bucket).WithContext(ctx).Scan(
		&message.ConversationID,
		&message.Bucket,
		&message.MessageID,
		&message.SenderID,
		&message.Type,
		&message.Content,
		&message.Attachments,
		&message.IsEdited,
		&message.ReplyToMessageID,
		&message.StoryRefID,
		&message.IsRevoked,
		&message.CreatedAt,
	)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil // Không tìm thấy message nào
		}
		return nil, fmt.Errorf("failed to get messages for conversation %s bucket %d: %w", conversationID, bucket, err)
	}

	return &message, nil
}
func (r *MessageRepository) DeleteBulkMessagesByConversationID(ctx context.Context, conversationIDs []string) error {
	tableName := entity.Message{}.TableName()
	querty := fmt.Sprintf(` DELETE FROM %s WHERE conversation_id = ?`, tableName)

	for _, conversationID := range conversationIDs {
		safeCtx := context.WithoutCancel(ctx)
		if err := r.session.Query(querty, conversationID).WithContext(safeCtx).Exec(); err != nil {
			return fmt.Errorf("failed to delete messages of conversation %s: %w", conversationID, err)
		}
	}

	return nil
}
func (r *MessageRepository) GetMessagesByMessageID(ctx context.Context, conversationID string, bucket int, messageID string) (*entity.Message, error) {
	// 1. Fail-fast validation
	if conversationID == "" || messageID == "" {
		return nil, fmt.Errorf("conversationID and messageID cannot be empty")
	}

	// 2. Parse messageID sang UUID để đảm bảo định dạng hợp lệ trước khi gọi DB
	msgUUID, err := gocql.ParseUUID(messageID)
	if err != nil {
		return nil, fmt.Errorf("invalid messageID format: %w", err)
	}

	// 3. Chuẩn bị Query cho Cassandra
	tableName := entity.Message{}.TableName()

	query := fmt.Sprintf(`
		SELECT conversation_id, bucket, message_id,
		       sender_id, type, content, attachments,
		       is_edited, reply_to_message_id, story_ref_id,
		       is_revoked, created_at
		FROM %s WHERE conversation_id = ? AND bucket = ? AND message_id = ?
	`, tableName)

	var message entity.Message

	err = r.session.Query(query, conversationID, bucket, msgUUID).WithContext(ctx).Scan(
		&message.ConversationID,
		&message.Bucket,
		&message.MessageID,
		&message.SenderID,
		&message.Type,
		&message.Content,
		&message.Attachments,
		&message.IsEdited,
		&message.ReplyToMessageID,
		&message.StoryRefID,
		&message.IsRevoked,
		&message.CreatedAt,
	)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil // Không tìm thấy message nào
		}
		return nil, fmt.Errorf("failed to get message %s for conversation %s bucket %d: %w", messageID, conversationID, bucket, err)
	}

	return &message, nil
}
func (r *MessageRepository) DeleteAllMessagesByMessageConversationID(ctx context.Context, conversationID string) error {
	tableName := entity.Message{}.TableName()
	query := fmt.Sprintf(`DELETE FROM %s WHERE conversation_id = ?`, tableName)

	safeCtx := context.WithoutCancel(ctx)
	if err := r.session.Query(query, conversationID).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete all messages of conversation %s: %w", conversationID, err)
	}

	return nil
}
