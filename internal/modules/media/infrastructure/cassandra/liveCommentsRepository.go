package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/concurrency"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/media/domain/entity"
	"github.com/gocql/gocql"
)

type LiveCommentsRepository struct {
	// Add necessary fields for Cassandra connection and table handling
	session   *gocql.Session
	pool      *concurrency.WorkerPool
	redisRepo IRepositoryShare.IRedis
}

// Implement methods for LiveCommentsRepository here, ensuring they satisfy the ILiveCommentsRepository interface defined in the domain layer.
func NewLiveCommentsRepository(session *gocql.Session, pool *concurrency.WorkerPool, redisRepo IRepositoryShare.IRedis) *LiveCommentsRepository {
	return &LiveCommentsRepository{
		session:   session,
		pool:      pool,
		redisRepo: redisRepo,
	}
}
func (r *LiveCommentsRepository) CreateLiveComment(ctx context.Context, comment *entity.LiveComment) error {
	// 1. Fail-fast validation
	if comment == nil {
		return fmt.Errorf("live comment payload is nil")
	}

	var emptyUUID gocql.UUID
	if comment.StreamID == emptyUUID || comment.CommentID == emptyUUID || comment.CreatedAt.IsZero() {
		return fmt.Errorf("missing primary key components: stream_id, created_at, and comment_id are required")
	}

	// 2. Bảo vệ Context: Chat nhảy rất nhanh, nhưng 1 khi đã gọi DB thì phải ghi cho xong
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.LiveComment{}.TableName()

	// 3. Query chuẩn bị: Gocql tự động map mảng []string của Go sang SET<TEXT> hoặc LIST<TEXT> của Cassandra
	query := fmt.Sprintf(`
		INSERT INTO %s (
			stream_id, created_at, comment_id, 
			user_id, user_nickname, user_avatar_url, user_badges, 
			content, is_pinned
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	// 4. Thực thi query
	err := r.session.Query(query,
		comment.StreamID,
		comment.CreatedAt,
		comment.CommentID,
		comment.UserID,
		comment.UserNickname,
		comment.UserAvatarURL,
		comment.UserBadges,
		comment.Content,
		comment.IsPinned,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to create live comment %s for stream %s: %w", comment.CommentID.String(), comment.StreamID.String(), err)
	}

	return nil
}

func (r *LiveCommentsRepository) CreateBulkLiveComments(ctx context.Context, comments []*entity.LiveComment) (int64, []*dto.LiveCommentBulkError, error) {
	// 1. Fail-fast validation
	if len(comments) == 0 {
		return 0, nil, nil
	}

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type taskResult struct {
		comment *entity.LiveComment
		err     error
	}

	// Sử dụng Buffered Channel chống deadlock
	resultCh := make(chan taskResult, len(comments))

	tableName := entity.LiveComment{}.TableName()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			stream_id, created_at, comment_id, 
			user_id, user_nickname, user_avatar_url, user_badges, 
			content, is_pinned
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	var emptyUUID gocql.UUID

	// 3. Phân phối task vào Worker Pool
	for i, item := range comments {
		// Xử lý an toàn con trỏ nil
		if item == nil {
			resultCh <- taskResult{
				comment: nil,
				err:     fmt.Errorf("comment pointer at index %d is nil", i),
			}
			continue
		}

		// Kiểm tra khóa chính trước khi tốn tài nguyên gọi DB
		if item.StreamID == emptyUUID || item.CommentID == emptyUUID || item.CreatedAt.IsZero() {
			resultCh <- taskResult{
				comment: item,
				err:     fmt.Errorf("missing primary key components for insert"),
			}
			continue
		}

		// Copy biến cục bộ an toàn cho Goroutine (Go < 1.22)
		payload := item

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				payload.StreamID,
				payload.CreatedAt,
				payload.CommentID,
				payload.UserID,
				payload.UserNickname,
				payload.UserAvatarURL,
				payload.UserBadges,
				payload.Content,
				payload.IsPinned,
			).WithContext(safeCtx).Exec()

			// Bắn kết quả về channel
			resultCh <- taskResult{
				comment: payload,
				err:     execErr,
			}
		})

		// Xử lý khi worker pool bị đầy hoặc context gốc bị cancel/timeout
		if err != nil {
			resultCh <- taskResult{
				comment: payload,
				err:     fmt.Errorf("worker pool rejected task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*dto.LiveCommentBulkError

	for res := range resultCh {
		if res.err != nil {
			// Trích xuất ID an toàn
			sID := "unknown_stream"
			cID := "unknown_comment"
			uID := "unknown_user"

			if res.comment != nil {
				sID = res.comment.StreamID.String()
				cID = res.comment.CommentID.String()
				uID = res.comment.UserID.String()
			}

			bulkErrors = append(bulkErrors, &dto.LiveCommentBulkError{
				StreamID:  sID,
				CommentID: cID,
				UserID:    uID,
				Error:     res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk insert live comments completed with errors: %d/%d success, %d failed", successCount, len(comments), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *LiveCommentsRepository) GetLiveCommentsByStreamID(ctx context.Context, streamID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if streamID == "" {
		return nil, fmt.Errorf("streamID cannot be empty")
	}

	// Xác thực UUID sớm để tránh gọi DB vô ích
	if _, err := gocql.ParseUUID(streamID); err != nil {
		return nil, fmt.Errorf("invalid streamID format: %w", err)
	}

	// 2. Check Cache (Chỉ hit cache khi lấy trang đầu tiên - lúc người dùng mới vào phòng live)
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"livecomment_cache_streamid_" + streamID,
			"livecomment_cache_nextcursor_streamid_" + streamID,
			"livecomment_cache_hasnext_streamid_" + streamID,
			"livecomment_cache_limit_streamid_" + streamID,
		})

		if err == nil && datacache != nil {
			if commentList, ok := datacache.([]*entity.LiveComment); ok {
				return &dto.PaginationRes{
					Data:       commentList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	// 3. Decode Cursor (Lấy Page State của Cassandra)
	pageState, err := utils.DecodeCursorCassandra(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	// 4. Chuẩn bị Query cho Cassandra
	tableName := entity.LiveComment{}.TableName()

	// Khai báo rõ các cột để tối ưu Serialize/Deserialize
	query := fmt.Sprintf(`
		SELECT stream_id, created_at, comment_id, user_id, user_nickname, 
		       user_avatar_url, user_badges, content, is_pinned 
		FROM %s WHERE stream_id = ?
	`, tableName)

	// LƯU Ý: Không dùng context.WithoutCancel ở đây. Truyền ctx gốc để DB ngừng xử lý nếu client ngắt kết nối.
	q := r.session.Query(query, streamID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 5. Thực thi và Map dữ liệu
	iter := q.Iter()
	scanner := iter.Scanner()
	var comments []*entity.LiveComment

	for scanner.Next() {
		var c entity.LiveComment
		// gocql sẽ tự động map mảng SET<TEXT> trong DB vào []string (user_badges)
		err := scanner.Scan(
			&c.StreamID,
			&c.CreatedAt,
			&c.CommentID,
			&c.UserID,
			&c.UserNickname,
			&c.UserAvatarURL,
			&c.UserBadges,
			&c.Content,
			&c.IsPinned,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan live comment for stream %s: %w", streamID, err)
		}

		// Append con trỏ an toàn (Go >= 1.22)
		comments = append(comments, &c)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for stream %s: %w", streamID, err)
	}

	// 6. Xử lý Next Cursor và HasNext
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	// Trả về slice rỗng thay vì nil để Frontend không bị lỗi map null
	if comments == nil {
		comments = []*entity.LiveComment{}
	}

	// 7. Cache kết quả cho trang đầu tiên
	if cursor == "" && len(comments) > 0 {
		items := map[string]any{
			"livecomment_cache_streamid_" + streamID:            comments,
			"livecomment_cache_nextcursor_streamid_" + streamID: nextCursorStr,
			"livecomment_cache_hasnext_streamid_" + streamID:    hasNext,
			"livecomment_cache_limit_streamid_" + streamID:      limit,
		}

		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			// Chỉ log lỗi chứ không chặn main flow vì redis lỗi thì app vẫn phải chạy
			fmt.Printf("Failed to set cache for live comment pagination: %v\n", err)
		}
	}

	// 8. Kết quả
	return &dto.PaginationRes{
		Data:       comments,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}

func (r *LiveCommentsRepository) GetLiveCommentsByUserID(ctx context.Context, streamID string, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if streamID == "" || userID == "" {
		return nil, fmt.Errorf("streamID and userID cannot be empty")
	}

	if _, err := gocql.ParseUUID(streamID); err != nil {
		return nil, fmt.Errorf("invalid streamID format: %w", err)
	}
	if _, err := gocql.ParseUUID(userID); err != nil {
		return nil, fmt.Errorf("invalid userID format: %w", err)
	}

	// 2. Check Cache
	cacheKeyPrefix := "livecomment_cache_streamid_" + streamID + "_userid_" + userID
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			cacheKeyPrefix,
			cacheKeyPrefix + "_nextcursor",
			cacheKeyPrefix + "_hasnext",
			cacheKeyPrefix + "_limit",
		})

		if err == nil && datacache != nil {
			if commentList, ok := datacache.([]*entity.LiveComment); ok {
				return &dto.PaginationRes{
					Data:       commentList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	// 3. Decode Cursor
	pageState, err := utils.DecodeCursorCassandra(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	tableName := entity.LiveComment{}.TableName()

	// LƯU Ý QUAN TRỌNG VỀ CASSANDRA:
	// Bảng hiện tại có PK là ((stream_id), created_at, comment_id).
	// user_id KHÔNG NẰM TRONG KHÓA CHÍNH.
	// Tuy nhiên, do ta đã lọc trước bằng Partition Key (stream_id = ?),
	// ta CÓ THỂ dùng ALLOW FILTERING để quét user_id một cách an toàn bên trong 1 partition (1 buổi live).
	// Nếu tính năng lấy comment theo user này được gọi RẤT NHIỀU,
	// bạn NÊN TẠO THÊM 1 Materialized View (VD: live_comments_by_user).
	query := fmt.Sprintf(`
		SELECT stream_id, created_at, comment_id, user_id, user_nickname, 
		       user_avatar_url, user_badges, content, is_pinned 
		FROM %s WHERE stream_id = ? AND user_id = ? ALLOW FILTERING
	`, tableName)

	q := r.session.Query(query, streamID, userID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 4. Thực thi và Map dữ liệu
	iter := q.Iter()
	scanner := iter.Scanner()
	var comments []*entity.LiveComment

	for scanner.Next() {
		var c entity.LiveComment
		err := scanner.Scan(
			&c.StreamID,
			&c.CreatedAt,
			&c.CommentID,
			&c.UserID,
			&c.UserNickname,
			&c.UserAvatarURL,
			&c.UserBadges,
			&c.Content,
			&c.IsPinned,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan live comment for stream %s, user %s: %w", streamID, userID, err)
		}
		comments = append(comments, &c)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for stream %s, user %s: %w", streamID, userID, err)
	}

	// 5. Xử lý Cursor & HasNext
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	if comments == nil {
		comments = []*entity.LiveComment{}
	}

	// 6. Set Cache (Trang đầu tiên)
	if cursor == "" && len(comments) > 0 {
		items := map[string]any{
			cacheKeyPrefix:                 comments,
			cacheKeyPrefix + "_nextcursor": nextCursorStr,
			cacheKeyPrefix + "_hasnext":    hasNext,
			cacheKeyPrefix + "_limit":      limit,
		}

		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for user's live comment pagination: %v\n", err)
		}
	}

	return &dto.PaginationRes{
		Data:       comments,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *LiveCommentsRepository) GetALLliveCommentsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format: %w", err)
	}

	// 2. Check Cache (Chỉ hit cache khi user yêu cầu trang đầu tiên)
	cacheKeyPrefix := "all_livecomments_cache_userid_" + userID
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			cacheKeyPrefix,
			cacheKeyPrefix + "_nextcursor",
			cacheKeyPrefix + "_hasnext",
			cacheKeyPrefix + "_limit",
		})

		if err == nil && datacache != nil {
			if commentList, ok := datacache.([]*entity.LiveComment); ok {
				return &dto.PaginationRes{
					Data:       commentList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	// 3. Decode Cursor
	pageState, err := utils.DecodeCursorCassandra(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	// 4. Chuẩn bị Query cho Cassandra
	// LƯU Ý TỐI QUAN TRỌNG: Không query trực tiếp vào bảng entity.LiveComment{}.TableName()
	// Bắt buộc phải dùng Materialized View để tránh Full Cluster Scan.
	mvName := "live_comments_by_user"

	query := fmt.Sprintf(`
		SELECT stream_id, created_at, comment_id, user_id, user_nickname, 
		       user_avatar_url, user_badges, content, is_pinned 
		FROM %s WHERE user_id = ?
	`, mvName)

	// Giữ nguyên ctx để DB ngắt truy vấn nếu Client hủy Request
	q := r.session.Query(query, uID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 5. Thực thi và duyệt dữ liệu
	iter := q.Iter()
	scanner := iter.Scanner()
	var comments []*entity.LiveComment

	for scanner.Next() {
		var c entity.LiveComment
		err := scanner.Scan(
			&c.StreamID,
			&c.CreatedAt,
			&c.CommentID,
			&c.UserID,
			&c.UserNickname,
			&c.UserAvatarURL,
			&c.UserBadges,
			&c.Content,
			&c.IsPinned,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan live comment from MV for user %s: %w", userID, err)
		}

		// Map con trỏ an toàn
		comments = append(comments, &c)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for user %s: %w", userID, err)
	}

	// 6. Xử lý phân trang (Next Cursor & HasNext)
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	// Trả về mảng rỗng [] thay vì null cho frontend
	if comments == nil {
		comments = []*entity.LiveComment{}
	}

	// 7. Lưu Cache kết quả trang đầu tiên
	if cursor == "" && len(comments) > 0 {
		items := map[string]any{
			cacheKeyPrefix:                 comments,
			cacheKeyPrefix + "_nextcursor": nextCursorStr,
			cacheKeyPrefix + "_hasnext":    hasNext,
			cacheKeyPrefix + "_limit":      limit,
		}

		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			// Bỏ qua lỗi Cache, không làm gián đoạn luồng chính
			fmt.Printf("Failed to set cache for all user live comments pagination: %v\n", err)
		}
	}

	// 8. Trả kết quả
	return &dto.PaginationRes{
		Data:       comments,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *LiveCommentsRepository) UpdateLiveComment(ctx context.Context, comment *entity.LiveComment) error {
	// 1. Fail-fast validation
	if comment == nil {
		return fmt.Errorf("live comment payload is nil")
	}

	var emptyUUID gocql.UUID
	// Kiểm tra BẮT BUỘC phải có đủ 3 thành phần của Primary Key
	if comment.StreamID == emptyUUID || comment.CommentID == emptyUUID || comment.CreatedAt.IsZero() {
		return fmt.Errorf("missing primary key components: stream_id, created_at, and comment_id are required for update")
	}

	// 2. Bảo vệ Context: Lệnh Ghi (Update) không được ngắt giữa chừng dù HTTP Request bị hủy
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.LiveComment{}.TableName()

	// 3. Query chuẩn bị:
	// - SET các trường dữ liệu (Non-Primary Key columns)
	// - WHERE bắt buộc chứa đủ Partition Key và Clustering Keys
	query := fmt.Sprintf(`
		UPDATE %s 
		SET user_id = ?, user_nickname = ?, user_avatar_url = ?, user_badges = ?, content = ?, is_pinned = ?
		WHERE stream_id = ? AND created_at = ? AND comment_id = ?
	`, tableName)

	// 4. Thực thi truy vấn
	err := r.session.Query(query,
		// Các trường SET
		comment.UserID,
		comment.UserNickname,
		comment.UserAvatarURL,
		comment.UserBadges,
		comment.Content,
		comment.IsPinned,
		// Các trường WHERE (Primary Key)
		comment.StreamID,
		comment.CreatedAt,
		comment.CommentID,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to update live comment %s for stream %s: %w", comment.CommentID.String(), comment.StreamID.String(), err)
	}

	return nil
}

func (r *LiveCommentsRepository) UpdateBulkLiveComments(ctx context.Context, comments []*entity.LiveComment) (int64, []*dto.LiveCommentBulkError, error) {
	// 1. Fail-fast validation
	if len(comments) == 0 {
		return 0, nil, nil
	}

	// 2. Khởi tạo Channel và Struct nội bộ chống Deadlock
	type taskResult struct {
		comment *entity.LiveComment
		err     error
	}
	resultCh := make(chan taskResult, len(comments))

	tableName := entity.LiveComment{}.TableName()
	query := fmt.Sprintf(`
		UPDATE %s 
		SET user_id = ?, user_nickname = ?, user_avatar_url = ?, user_badges = ?, content = ?, is_pinned = ?
		WHERE stream_id = ? AND created_at = ? AND comment_id = ?
	`, tableName)

	var emptyUUID gocql.UUID

	// 3. Phân bổ Task vào Worker Pool
	for i, item := range comments {
		// Bỏ qua an toàn nếu pointer nil
		if item == nil {
			resultCh <- taskResult{
				comment: nil,
				err:     fmt.Errorf("comment pointer at index %d is nil", i),
			}
			continue
		}

		// Kiểm tra khóa chính trước khi gọi DB để tối ưu tài nguyên
		if item.StreamID == emptyUUID || item.CommentID == emptyUUID || item.CreatedAt.IsZero() {
			resultCh <- taskResult{
				comment: item,
				err:     fmt.Errorf("missing primary key components for update"),
			}
			continue
		}

		// Copy biến cục bộ để Closure trong Goroutine trỏ đúng vùng nhớ
		payload := item

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT KHI GHI
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				// SET
				payload.UserID,
				payload.UserNickname,
				payload.UserAvatarURL,
				payload.UserBadges,
				payload.Content,
				payload.IsPinned,
				// WHERE
				payload.StreamID,
				payload.CreatedAt,
				payload.CommentID,
			).WithContext(safeCtx).Exec()

			// Trả kết quả về channel
			resultCh <- taskResult{
				comment: payload,
				err:     execErr,
			}
		})

		// Xử lý khi Worker Pool từ chối (Context gốc đã Timeout/Cancel)
		if err != nil {
			resultCh <- taskResult{
				comment: payload,
				err:     fmt.Errorf("worker pool rejected update task: %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa
	r.pool.Wait()
	close(resultCh)

	// 5. Thu gom kết quả (Lock-free)
	var successCount int64
	var bulkErrors []*dto.LiveCommentBulkError

	for res := range resultCh {
		if res.err != nil {
			// Fallback ID an toàn
			sID := "unknown_stream"
			cID := "unknown_comment"
			uID := "unknown_user"

			if res.comment != nil {
				sID = res.comment.StreamID.String()
				cID = res.comment.CommentID.String()
				uID = res.comment.UserID.String()
			}

			bulkErrors = append(bulkErrors, &dto.LiveCommentBulkError{
				StreamID:  sID,
				CommentID: cID,
				UserID:    uID,
				Error:     res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk update live comments completed with errors: %d/%d success, %d failed", successCount, len(comments), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *LiveCommentsRepository) DeleteLiveComment(ctx context.Context, streamID string, commentID string, createdAt time.Time) error {
	// 1. Fail-fast validation
	if streamID == "" || commentID == "" {
		return fmt.Errorf("streamID and commentID cannot be empty")
	}
	if createdAt.IsZero() {
		return fmt.Errorf("createdAt timestamp is required for deletion")
	}

	sID, err := gocql.ParseUUID(streamID)
	if err != nil {
		return fmt.Errorf("invalid streamID format: %w", err)
	}

	cID, err := gocql.ParseUUID(commentID)
	if err != nil {
		return fmt.Errorf("invalid commentID format: %w", err)
	}

	// 2. Bảo vệ Context (Thao tác DELETE thực chất là ghi Tombstone xuống đĩa)
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.LiveComment{}.TableName()

	// 3. Query: Yêu cầu chính xác 100% Primary Key
	query := fmt.Sprintf("DELETE FROM %s WHERE stream_id = ? AND created_at = ? AND comment_id = ?", tableName)

	// 4. Thực thi
	if err := r.session.Query(query, sID, createdAt, cID).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete live comment %s in stream %s: %w", commentID, streamID, err)
	}

	return nil
}

func (r *LiveCommentsRepository) DeleteBulkLiveComments(ctx context.Context, streamID string, commentIDs []string, createdAt time.Time) (int64, []*dto.LiveCommentBulkError, error) {
	// 1. Fail-fast validation
	if streamID == "" || len(commentIDs) == 0 {
		return 0, nil, nil
	}
	if createdAt.IsZero() {
		return 0, nil, fmt.Errorf("createdAt timestamp is required for bulk deletion")
	}

	sID, err := gocql.ParseUUID(streamID)
	if err != nil {
		return 0, nil, fmt.Errorf("bulk delete aborted - invalid streamID format: %w", err)
	}

	// Chuyển đổi toàn bộ mảng CommentID sang UUID trước để bắt lỗi sớm
	uuids := make([]gocql.UUID, 0, len(commentIDs))
	for _, idStr := range commentIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk delete aborted - invalid commentID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. Chuẩn bị Struct và Channel an toàn
	type taskResult struct {
		commentID gocql.UUID
		err       error
	}
	resultCh := make(chan taskResult, len(uuids))

	tableName := entity.LiveComment{}.TableName()
	query := fmt.Sprintf("DELETE FROM %s WHERE stream_id = ? AND created_at = ? AND comment_id = ?", tableName)

	// 3. Phân phối nhiệm vụ vào Worker Pool
	for _, cID := range uuids {
		payloadCID := cID // Copy biến an toàn cho Goroutine (Go < 1.22)

		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			
			// Bắn lệnh xóa xuống DB
			execErr := r.session.Query(query, sID, createdAt, payloadCID).WithContext(safeCtx).Exec()
			
			resultCh <- taskResult{
				commentID: payloadCID,
				err:       execErr,
			}
		})

		// Xử lý khi Pool từ chối do quá tải hoặc Context HTTP bị Timeout/Cancel
		if err != nil {
			resultCh <- taskResult{
				commentID: payloadCID,
				err:       fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa các Goroutine
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*dto.LiveCommentBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &dto.LiveCommentBulkError{
				StreamID:  streamID,
				CommentID: res.commentID.String(),
				UserID:    "unknown_user", // Do ta không cần UserID để xóa, ta set fallback label
				Error:     res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete live comments completed with errors: %d/%d success", successCount, len(uuids))
	}

	return successCount, bulkErrors, finalErr
}

func (r *LiveCommentsRepository) DeleteLiveCommentsByTimeRange(ctx context.Context, streamID string, startTime time.Time, endTime time.Time) error {
	// 1. Fail-fast validation
	if streamID == "" {
		return fmt.Errorf("streamID cannot be empty")
	}
	if startTime.IsZero() || endTime.IsZero() {
		return fmt.Errorf("startTime and endTime cannot be zero")
	}
	if startTime.After(endTime) {
		return fmt.Errorf("startTime must be before or equal to endTime")
	}

	sID, err := gocql.ParseUUID(streamID)
	if err != nil {
		return fmt.Errorf("invalid streamID format '%s': %w", streamID, err)
	}

	// 2. Bảo vệ Context vì đây là thao tác Ghi (Tombstone)
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.LiveComment{}.TableName()

	// 3. Tận dụng sức mạnh Range Deletion của Cassandra
	// Chỉ cần chỉ định Partition Key (stream_id) và khoảng của Clustering Key 1 (created_at)
	query := fmt.Sprintf("DELETE FROM %s WHERE stream_id = ? AND created_at >= ? AND created_at <= ?", tableName)

	// 4. Thực thi truy vấn
	if err := r.session.Query(query, sID, startTime, endTime).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete live comments for stream %s between %v and %v: %w", streamID, startTime, endTime, err)
	}

	return nil
}

func (r *LiveCommentsRepository) DeleteBulkLiveCommentsByTimeRange(ctx context.Context, streamIDs []string, startTime time.Time, endTime time.Time) (int64, []*dto.LiveCommentBulkError, error) {
	// 1. Fail-fast validation
	if len(streamIDs) == 0 {
		return 0, nil, nil
	}
	if startTime.IsZero() || endTime.IsZero() {
		return 0, nil, fmt.Errorf("startTime and endTime cannot be zero")
	}
	if startTime.After(endTime) {
		return 0, nil, fmt.Errorf("startTime must be before or equal to endTime")
	}

	// Ép kiểu toàn bộ UUID trước để bắt lỗi cấu trúc sớm
	uuids := make([]gocql.UUID, 0, len(streamIDs))
	for _, idStr := range streamIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk range delete aborted - invalid streamID '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. Chuẩn bị Struct và Channel an toàn cho Worker Pool
	type taskResult struct {
		streamID gocql.UUID
		err      error
	}
	resultCh := make(chan taskResult, len(uuids))

	tableName := entity.LiveComment{}.TableName()
	query := fmt.Sprintf("DELETE FROM %s WHERE stream_id = ? AND created_at >= ? AND created_at <= ?", tableName)

	// 3. Phân phối nhiệm vụ (Mỗi Worker chịu trách nhiệm dọn dẹp 1 Stream)
	for _, id := range uuids {
		payloadID := id // Copy biến an toàn

		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			
			// Bắn lệnh Range Deletion xuống DB
			execErr := r.session.Query(query, payloadID, startTime, endTime).WithContext(safeCtx).Exec()
			
			resultCh <- taskResult{
				streamID: payloadID,
				err:      execErr,
			}
		})

		// Xử lý khi Pool bị nghẽn hoặc Context Timeout
		if err != nil {
			resultCh <- taskResult{
				streamID: payloadID,
				err:      fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ hóa
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*dto.LiveCommentBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &dto.LiveCommentBulkError{
				StreamID:  res.streamID.String(),
				// Đánh dấu rõ đây là tác vụ dọn dẹp hàng loạt, không nhắm tới 1 comment hay 1 user cụ thể
				CommentID: "range_deletion", 
				UserID:    "all_users_in_range",
				Error:     res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái cuối cùng
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk range delete completed with errors: %d/%d success", successCount, len(uuids), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}