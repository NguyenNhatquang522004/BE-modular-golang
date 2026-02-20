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

type StoryViewRepository struct {
	// Add necessary fields for Cassandra connection and table handling
	session   *gocql.Session
	pool      *concurrency.WorkerPool
	redisRepo IRepositoryShare.IRedis
}

// Implement methods for StoryViewRepository here, ensuring they satisfy the IStoryViewRepository interface defined in the domain layer.
func NewStoryViewRepository(session *gocql.Session, pool *concurrency.WorkerPool, redisRepo IRepositoryShare.IRedis) *StoryViewRepository {
	return &StoryViewRepository{
		session:   session,
		pool:      pool,
		redisRepo: redisRepo,
	}
}
func (r *StoryViewRepository) CreateStoryView(ctx context.Context, storyView *entity.StoryView) error {
	// 1. Fail-fast: Validate input
	if storyView == nil {
		return fmt.Errorf("storyView payload is nil")
	}

	var emptyUUID gocql.UUID
	if storyView.StoryID == emptyUUID || storyView.ViewerID == emptyUUID {
		return fmt.Errorf("missing primary key components: story_id and viewer_id are required")
	}

	// 2. Bảo vệ Context: Đảm bảo lệnh Ghi (Write) không bị đứt gãy nếu HTTP Request bị Cancel/Timeout
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.StoryView{}.TableName()

	// 3. Query chuẩn bị: Khai báo rõ thứ tự cột (Explicit columns)
	query := fmt.Sprintf(`
		INSERT INTO %s (
			story_id, viewed_at, viewer_id, 
			viewer_name, viewer_avatar_url, 
			interaction_type, reaction_code, poll_option_index
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	// 4. Thực thi query
	// Gocql tự động xử lý pointer *int (poll_option_index) thành null nếu nil
	err := r.session.Query(query,
		storyView.StoryID,
		storyView.ViewedAt,
		storyView.ViewerID,
		storyView.ViewerName,
		storyView.ViewerAvatarURL,
		storyView.InteractionType,
		storyView.ReactionCode,
		storyView.PollOptionIndex,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to create story view for story %s by viewer %s: %w", storyView.StoryID.String(), storyView.ViewerID.String(), err)
	}

	return nil
}
func (r *StoryViewRepository) CreateBulkStoryViews(ctx context.Context, storyViews []*entity.StoryView) (int64, []*dto.StoryViewBulkError, error) {
	// 1. Fail-fast validation
	if len(storyViews) == 0 {
		return 0, nil, nil
	}

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type taskResult struct {
		storyView *entity.StoryView
		err       error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để ngăn chặn deadlock
	resultCh := make(chan taskResult, len(storyViews))

	tableName := entity.StoryView{}.TableName()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			story_id, viewed_at, viewer_id, 
			viewer_name, viewer_avatar_url, 
			interaction_type, reaction_code, poll_option_index
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)

	// 3. Phân phối task vào Worker Pool
	for i, item := range storyViews {
		// Xử lý an toàn con trỏ nil
		if item == nil {
			resultCh <- taskResult{
				storyView: nil,
				err:       fmt.Errorf("storyView pointer at index %d is nil", i),
			}
			continue
		}

		// Bắt buộc copy biến cục bộ để Closure của Goroutine không tham chiếu sai vùng nhớ (Go < 1.22)
		payload := item

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT: Tách context bảo vệ lệnh ghi giống như single insert
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				payload.StoryID,
				payload.ViewedAt,
				payload.ViewerID,
				payload.ViewerName,
				payload.ViewerAvatarURL,
				payload.InteractionType,
				payload.ReactionCode,
				payload.PollOptionIndex,
			).WithContext(safeCtx).Exec()

			// Bắn kết quả về channel
			resultCh <- taskResult{
				storyView: payload,
				err:       execErr,
			}
		})

		// Xử lý khi worker pool bị đầy hoặc context gốc bị timeout không nhận thêm task
		if err != nil {
			resultCh <- taskResult{
				storyView: payload,
				err:       fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break // Dừng việc nhồi thêm task, những task đã vào pool vẫn sẽ chạy trọn vẹn
		}
	}

	// 4. Đồng bộ hóa: Chờ tất cả Worker hoàn thành và đóng Channel an toàn
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free vì chỉ 1 main goroutine đọc channel)
	var successCount int64
	var bulkErrors []*dto.StoryViewBulkError

	for res := range resultCh {
		if res.err != nil {
			// Trích xuất ID an toàn ngay cả khi payload bị nil
			sID := "unknown_story"
			uID := "unknown_viewer"

			if res.storyView != nil {
				sID = res.storyView.StoryID.String()
				uID = res.storyView.ViewerID.String()
			}

			bulkErrors = append(bulkErrors, &dto.StoryViewBulkError{
				StoryID: sID,
				UserID:  uID,
				Error:   res.err.Error(),
			})
		} else {
			successCount++ // Tăng biến đếm thành công
		}
	}

	// 6. Đánh giá trạng thái trả về (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk insert story views completed with errors: %d/%d success, %d failed", successCount, len(storyViews), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *StoryViewRepository) GetStoryViewsByStoryID(ctx context.Context, storyID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if storyID == "" {
		return nil, fmt.Errorf("storyID cannot be empty")
	}

	// 2. Check Cache (Chỉ hit cache khi người dùng tải trang đầu tiên)
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"storyview_cache_storyid_" + storyID,
			"storyview_cache_nextcursor_storyid_" + storyID,
			"storyview_cache_hasnext_storyid_" + storyID,
			"storyview_cache_limit_storyid_" + storyID,
		})

		if err == nil && datacache != nil {
			if viewList, ok := datacache.([]*entity.StoryView); ok {
				return &dto.PaginationRes{
					Data:       viewList,
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
	tableName := entity.StoryView{}.TableName()

	// Khai báo rõ các cột để tối ưu hoá tốc độ Scan thay vì dùng SELECT *
	query := fmt.Sprintf(`
		SELECT story_id, viewed_at, viewer_id, viewer_name, viewer_avatar_url, 
		       interaction_type, reaction_code, poll_option_index 
		FROM %s WHERE story_id = ?
	`, tableName)

	// Truyền ctx gốc vào để nếu HTTP request hủy, DB sẽ ngừng query ngay
	q := r.session.Query(query, storyID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 5. Thực thi và duyệt kết quả
	iter := q.Iter()
	scanner := iter.Scanner()

	// Không pre-allocate len() vì ta chưa biết trước có bao nhiêu dòng thực tế
	var views []*entity.StoryView

	for scanner.Next() {
		var sv entity.StoryView

		// Map thẳng dữ liệu vào struct. Field PollOptionIndex là con trỏ (*int)
		// gocql sẽ tự động gán nil nếu giá trị trong db là null.
		err := scanner.Scan(
			&sv.StoryID,
			&sv.ViewedAt,
			&sv.ViewerID,
			&sv.ViewerName,
			&sv.ViewerAvatarURL,
			&sv.InteractionType,
			&sv.ReactionCode,
			&sv.PollOptionIndex,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan story view for story %s: %w", storyID, err)
		}

		// Append con trỏ của biến cục bộ (An toàn trên Go >= 1.22)
		views = append(views, &sv)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for story %s: %w", storyID, err)
	}

	// 6. Xử lý Next Cursor và HasNext
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	// Nếu slice nil, gán bằng mảng rỗng [] để tránh lỗi JSON null ở Frontend
	if views == nil {
		views = []*entity.StoryView{}
	}

	// 7. Lưu Cache kết quả cho trang đầu tiên (Chỉ lưu nếu có data)
	if cursor == "" && len(views) > 0 {
		items := map[string]any{
			"storyview_cache_storyid_" + storyID:            views,
			"storyview_cache_nextcursor_storyid_" + storyID: nextCursorStr,
			"storyview_cache_hasnext_storyid_" + storyID:    hasNext,
			"storyview_cache_limit_storyid_" + storyID:      limit,
		}

		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			// Chỉ log lỗi thay vì chặn flow chính vì cache miss/error không nên làm hỏng luồng
			fmt.Printf("Failed to set cache for story view pagination: %v\n", err)
		}
	}

	// 8. Trả về kết quả
	return &dto.PaginationRes{
		Data:       views,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}

func (r *StoryViewRepository) GetStoryViewsByUserID(ctx context.Context, userID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	// 2. Check Cache
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"storyview_cache_userid_" + userID,
			"storyview_cache_nextcursor_userid_" + userID,
			"storyview_cache_hasnext_userid_" + userID,
			"storyview_cache_limit_userid_" + userID,
		})

		if err == nil && datacache != nil {
			if viewList, ok := datacache.([]*entity.StoryView); ok {
				return &dto.PaginationRes{
					Data:       viewList,
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

	tableName := entity.StoryView{}.TableName()

	// LƯU Ý QUAN TRỌNG:
	// Trong thiết kế Primary Key hiện tại, viewer_id là Clustering Key thứ 2.
	// Để query trực tiếp bằng viewer_id mà không gây ful-scan (ALLOW FILTERING),
	// bạn CẦN SỬ DỤNG MATERIALIZED VIEW. Hãy thay tableName thành tên view đó nếu có.
	// Ví dụ: tableName = "story_views_by_user"

	query := fmt.Sprintf(`
		SELECT story_id, viewed_at, viewer_id, viewer_name, viewer_avatar_url, 
		       interaction_type, reaction_code, poll_option_index 
		FROM %s WHERE viewer_id = ?
	`, tableName)

	q := r.session.Query(query, userID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 4. Thực thi và map dữ liệu
	iter := q.Iter()
	scanner := iter.Scanner()
	var views []*entity.StoryView

	for scanner.Next() {
		var sv entity.StoryView
		err := scanner.Scan(
			&sv.StoryID,
			&sv.ViewedAt,
			&sv.ViewerID,
			&sv.ViewerName,
			&sv.ViewerAvatarURL,
			&sv.InteractionType,
			&sv.ReactionCode,
			&sv.PollOptionIndex,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan story view for user %s: %w", userID, err)
		}
		views = append(views, &sv)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for user %s: %w", userID, err)
	}

	// 5. Build Cursor & HasNext
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	if views == nil {
		views = []*entity.StoryView{}
	}

	// 6. Set Cache (Trang đầu)
	if cursor == "" && len(views) > 0 {
		items := map[string]any{
			"storyview_cache_userid_" + userID:            views,
			"storyview_cache_nextcursor_userid_" + userID: nextCursorStr,
			"storyview_cache_hasnext_userid_" + userID:    hasNext,
			"storyview_cache_limit_userid_" + userID:      limit,
		}

		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			fmt.Printf("Failed to set cache for user story view pagination: %v\n", err)
		}
	}

	return &dto.PaginationRes{
		Data:       views,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
func (r *StoryViewRepository) UpdateStoryView(ctx context.Context, storyView *entity.StoryView) error {
	// 1. Fail-fast validation: Bắt buộc phải có đủ bộ Primary Key (Partition Key + toàn bộ Clustering Keys)
	if storyView == nil {
		return fmt.Errorf("storyView payload is nil")
	}

	var emptyUUID gocql.UUID
	if storyView.StoryID == emptyUUID || storyView.ViewerID == emptyUUID || storyView.ViewedAt.IsZero() {
		return fmt.Errorf("missing primary key components: story_id, viewed_at, and viewer_id are all required for update")
	}

	// 2. Bảo vệ Context: Đảm bảo thao tác cập nhật không bị đứt gãy nếu HTTP Request bị ngắt
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.StoryView{}.TableName()

	// 3. Query chuẩn bị: SET các data column, WHERE bắt buộc phải chứa đủ 3 thành phần của PK
	query := fmt.Sprintf(`
		UPDATE %s 
		SET viewer_name = ?, viewer_avatar_url = ?, interaction_type = ?, reaction_code = ?, poll_option_index = ?
		WHERE story_id = ? AND viewed_at = ? AND viewer_id = ?
	`, tableName)

	// 4. Thực thi query
	err := r.session.Query(query,
		storyView.ViewerName,
		storyView.ViewerAvatarURL,
		storyView.InteractionType,
		storyView.ReactionCode,
		storyView.PollOptionIndex,
		// Primary Key Components cho WHERE clause
		storyView.StoryID,
		storyView.ViewedAt,
		storyView.ViewerID,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to update story view for story %s by viewer %s at %v: %w",
			storyView.StoryID.String(), storyView.ViewerID.String(), storyView.ViewedAt, err)
	}

	return nil
}

func (r *StoryViewRepository) UpdateBulkStoryViews(ctx context.Context, storyViews []*entity.StoryView) (int64, []*dto.StoryViewBulkError, error) {
	// 1. Fail-fast validation
	if len(storyViews) == 0 {
		return 0, nil, nil
	}

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type taskResult struct {
		storyView *entity.StoryView
		err       error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để ngăn chặn deadlock
	resultCh := make(chan taskResult, len(storyViews))

	tableName := entity.StoryView{}.TableName()
	query := fmt.Sprintf(`
		UPDATE %s 
		SET viewer_name = ?, viewer_avatar_url = ?, interaction_type = ?, reaction_code = ?, poll_option_index = ?
		WHERE story_id = ? AND viewed_at = ? AND viewer_id = ?
	`, tableName)

	var emptyUUID gocql.UUID

	// 3. Phân phối task vào Worker Pool
	for i, item := range storyViews {
		// Xử lý an toàn con trỏ nil
		if item == nil {
			resultCh <- taskResult{
				storyView: nil,
				err:       fmt.Errorf("storyView pointer at index %d is nil", i),
			}
			continue
		}

		// Kiểm tra tính hợp lệ của Primary Key trước khi đưa vào Pool để đỡ tốn tài nguyên DB
		if item.StoryID == emptyUUID || item.ViewerID == emptyUUID || item.ViewedAt.IsZero() {
			resultCh <- taskResult{
				storyView: item,
				err:       fmt.Errorf("missing primary key components for update"),
			}
			continue
		}

		// Bắt buộc copy biến cục bộ để Closure của Goroutine không tham chiếu sai vùng nhớ
		payload := item

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				payload.ViewerName,
				payload.ViewerAvatarURL,
				payload.InteractionType,
				payload.ReactionCode,
				payload.PollOptionIndex,
				// WHERE clause
				payload.StoryID,
				payload.ViewedAt,
				payload.ViewerID,
			).WithContext(safeCtx).Exec()

			// Bắn kết quả về channel
			resultCh <- taskResult{
				storyView: payload,
				err:       execErr,
			}
		})

		// Xử lý khi worker pool từ chối task (context cancel/timeout)
		if err != nil {
			resultCh <- taskResult{
				storyView: payload,
				err:       fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break // Dừng việc nhồi thêm task
		}
	}

	// 4. Đồng bộ hóa: Chờ tất cả Worker hoàn thành và đóng Channel an toàn
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*dto.StoryViewBulkError

	for res := range resultCh {
		if res.err != nil {
			// Trích xuất ID an toàn ngay cả khi payload bị nil
			sID := "unknown_story"
			uID := "unknown_viewer"

			if res.storyView != nil {
				sID = res.storyView.StoryID.String()
				uID = res.storyView.ViewerID.String()
			}

			bulkErrors = append(bulkErrors, &dto.StoryViewBulkError{
				StoryID: sID,
				UserID:  uID,
				Error:   res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái trả về (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk update story views completed with errors: %d/%d success, %d failed", successCount, len(storyViews), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *StoryViewRepository) DeleteStoryView(ctx context.Context, storyID string) error {
	// 1. Fail-fast validation
	if storyID == "" {
		return fmt.Errorf("storyID cannot be empty")
	}

	// Cần parse string sang UUID để đảm bảo dữ liệu hợp lệ trước khi gọi DB
	id, err := gocql.ParseUUID(storyID)
	if err != nil {
		return fmt.Errorf("invalid storyID format '%s': %w", storyID, err)
	}

	// 2. Bảo vệ Context: Delete trong Cassandra là ghi Tombstone, nên vẫn cần bảo vệ
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.StoryView{}.TableName()

	// 3. Query: Xóa toàn bộ Partition (Toàn bộ lượt xem của 1 Story)
	query := fmt.Sprintf("DELETE FROM %s WHERE story_id = ?", tableName)

	// 4. Thực thi
	if err := r.session.Query(query, id).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete all views for story %s: %w", storyID, err)
	}

	return nil
}

func (r *StoryViewRepository) DeleteBulkStoryViews(ctx context.Context, storyIDs []string) (int64, []*dto.StoryViewBulkError, error) {
	// 1. Fail-fast validation
	if len(storyIDs) == 0 {
		return 0, nil, nil
	}

	// Parse toàn bộ UUID trước (Fail-fast: Lỗi 1 ID thì từ chối luôn cả mảng để đảm bảo an toàn)
	uuids := make([]gocql.UUID, 0, len(storyIDs))
	for _, idStr := range storyIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk delete aborted - invalid storyID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. Struct vận chuyển kết quả
	type taskResult struct {
		storyID gocql.UUID
		err     error
	}

	// Khởi tạo channel chống Deadlock
	resultCh := make(chan taskResult, len(uuids))

	tableName := entity.StoryView{}.TableName()
	query := fmt.Sprintf("DELETE FROM %s WHERE story_id = ?", tableName)

	// 3. Phân phối vào Worker Pool
	for _, id := range uuids {
		// Copy biến an toàn cho Goroutine
		payloadID := id

		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			execErr := r.session.Query(query, payloadID).WithContext(safeCtx).Exec()

			resultCh <- taskResult{
				storyID: payloadID,
				err:     execErr,
			}
		})

		// Bắt lỗi Pool từ chối (Do context bị timeout/cancel)
		if err != nil {
			resultCh <- taskResult{
				storyID: payloadID,
				err:     fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break // Không cố nhồi thêm
		}
	}

	// 4. Đồng bộ
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu
	var successCount int64
	var bulkErrors []*dto.StoryViewBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &dto.StoryViewBulkError{
				StoryID: res.storyID.String(),
				// Điền "all_viewers" hoặc rỗng vì thao tác này xóa nguyên partition, không nhắm tới 1 user cụ thể
				UserID: "all_viewers",
				Error:  res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete story views completed with errors: %d/%d success, %d failed", successCount, len(uuids), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *StoryViewRepository) DeleteStoryViewsByUserID(ctx context.Context, userID string) error {
	// 1. Fail-fast validation
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format '%s': %w", userID, err)
	}

	// LƯU Ý CASSANDRA: viewer_id không phải là Partition Key.
	// Bắt buộc phải dùng Materialized View để tìm danh sách (story_id, viewed_at) của user này.
	// Giả định bạn có một Materialized View tên là 'story_views_by_user'
	mvName := "story_views_by_user"
	baseTableName := entity.StoryView{}.TableName()

	// 2. READ-PHASE: Truy vấn MV để lấy Primary Key gốc
	readQuery := fmt.Sprintf("SELECT story_id, viewed_at FROM %s WHERE viewer_id = ?", mvName)

	type pkInfo struct {
		storyID  gocql.UUID
		viewedAt time.Time
	}
	var pks []pkInfo

	// Không dùng context.WithoutCancel ở đây để HTTP cancel thì Read dừng ngay lập tức
	iter := r.session.Query(readQuery, uID).WithContext(ctx).Iter()
	scanner := iter.Scanner()

	for scanner.Next() {
		var sID gocql.UUID
		var vAt time.Time
		if err := scanner.Scan(&sID, &vAt); err != nil {
			return fmt.Errorf("failed to scan pk info from MV for user %s: %w", userID, err)
		}
		pks = append(pks, pkInfo{storyID: sID, viewedAt: vAt})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("database error when reading user %s views: %w", userID, err)
	}

	// Không có gì để xóa
	if len(pks) == 0 {
		return nil
	}

	// 3. DELETE-PHASE: Dùng Worker Pool để bắn nhiều lệnh xóa độc lập xuống Base Table
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE story_id = ? AND viewed_at = ? AND viewer_id = ?", baseTableName)
	errCh := make(chan error, len(pks))

	for _, pk := range pks {
		payload := pk // Copy biến cục bộ an toàn

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT GHI
			safeCtx := context.WithoutCancel(ctx)
			execErr := r.session.Query(deleteQuery, payload.storyID, payload.viewedAt, uID).WithContext(safeCtx).Exec()
			errCh <- execErr
		})

		if err != nil {
			errCh <- fmt.Errorf("worker pool rejected delete task for user %s: %w", userID, err)
			break
		}
	}

	// 4. Đồng bộ và gom lỗi
	r.pool.Wait()
	close(errCh)

	for e := range errCh {
		if e != nil {
			return e // Trả về lỗi đầu tiên gặp phải
		}
	}

	return nil
}

func (r *StoryViewRepository) DeleteStoryViewsByUserIAndStoryID(ctx context.Context, userID string, storyID string) error {
	// 1. Fail-fast validation
	if userID == "" || storyID == "" {
		return fmt.Errorf("userID and storyID cannot be empty")
	}

	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format: %w", err)
	}

	sID, err := gocql.ParseUUID(storyID)
	if err != nil {
		return fmt.Errorf("invalid storyID format: %w", err)
	}

	tableName := entity.StoryView{}.TableName()

	// LƯU Ý CASSANDRA: Ta đã biết story_id (Partition) và viewer_id (Clustering thứ 2),
	// nhưng bị thiếu viewed_at (Clustering thứ 1).
	// Giải pháp: Thực hiện Read-Before-Write dùng ALLOW FILTERING trong phạm vi 1 Partition.

	// 2. READ-PHASE: Lấy danh sách viewed_at
	readQuery := fmt.Sprintf("SELECT viewed_at FROM %s WHERE story_id = ? AND viewer_id = ? ALLOW FILTERING", tableName)

	iter := r.session.Query(readQuery, sID, uID).WithContext(ctx).Iter()
	scanner := iter.Scanner()
	var timestamps []time.Time

	for scanner.Next() {
		var vAt time.Time
		if err := scanner.Scan(&vAt); err != nil {
			return fmt.Errorf("failed to scan viewed_at for story %s user %s: %w", storyID, userID, err)
		}
		timestamps = append(timestamps, vAt)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("database error when reading timestamps for story %s user %s: %w", storyID, userID, err)
	}

	// Không có lượt xem nào của user này trong story này
	if len(timestamps) == 0 {
		return nil
	}

	// 3. DELETE-PHASE: Dùng Worker Pool xóa chính xác từng bản ghi
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE story_id = ? AND viewed_at = ? AND viewer_id = ?", tableName)
	errCh := make(chan error, len(timestamps))

	for _, t := range timestamps {
		payloadTime := t // Copy biến

		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			execErr := r.session.Query(deleteQuery, sID, payloadTime, uID).WithContext(safeCtx).Exec()
			errCh <- execErr
		})

		if err != nil {
			errCh <- fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err)
			break
		}
	}

	r.pool.Wait()
	close(errCh)

	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return nil
}
func (r *StoryViewRepository) DeleteStoryViewsByStoryIDAndUserID(ctx context.Context, storyID string, userID string) error {
	// 1. Fail-fast validation
	if storyID == "" || userID == "" {
		return fmt.Errorf("storyID and userID cannot be empty")
	}

	sID, err := gocql.ParseUUID(storyID)
	if err != nil {
		return fmt.Errorf("invalid storyID format: %w", err)
	}

	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid userID format: %w", err)
	}

	tableName := entity.StoryView{}.TableName()

	// 2. READ-PHASE: Lấy danh sách các lần xem (viewed_at) của user này trong story này
	// An toàn khi dùng ALLOW FILTERING vì ta đã chỉ định rõ Partition Key (story_id)
	readQuery := fmt.Sprintf("SELECT viewed_at FROM %s WHERE story_id = ? AND viewer_id = ? ALLOW FILTERING", tableName)

	iter := r.session.Query(readQuery, sID, uID).WithContext(ctx).Iter()
	scanner := iter.Scanner()
	var timestamps []time.Time

	for scanner.Next() {
		var vAt time.Time
		if err := scanner.Scan(&vAt); err != nil {
			return fmt.Errorf("failed to scan viewed_at for story %s, user %s: %w", storyID, userID, err)
		}
		timestamps = append(timestamps, vAt)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("database error when reading timestamps for story %s, user %s: %w", storyID, userID, err)
	}

	// Nếu user chưa từng xem story này -> Không cần làm gì thêm
	if len(timestamps) == 0 {
		return nil
	}

	// 3. DELETE-PHASE: Đẩy lệnh xóa xuống Worker Pool
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE story_id = ? AND viewed_at = ? AND viewer_id = ?", tableName)
	errCh := make(chan error, len(timestamps))

	for _, t := range timestamps {
		payloadTime := t // Copy biến an toàn cho Goroutine (Go < 1.22)

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT: Không để thao tác Xóa bị ngắt giữa chừng
			safeCtx := context.WithoutCancel(ctx)
			execErr := r.session.Query(deleteQuery, sID, payloadTime, uID).WithContext(safeCtx).Exec()
			errCh <- execErr
		})

		if err != nil {
			errCh <- fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err)
			break
		}
	}

	// 4. Đồng bộ và thu gom lỗi
	r.pool.Wait()
	close(errCh)

	for e := range errCh {
		if e != nil {
			return e // Trả về lỗi đầu tiên gặp phải
		}
	}

	return nil
}

func (r *StoryViewRepository) DeleteBulkStoryViewsByUserIAndStoryID(ctx context.Context, userID string, storyID []string) (int64, []*dto.StoryViewBulkError, error) {
	// 1. Fail-fast validation
	if userID == "" || len(storyID) == 0 {
		return 0, nil, nil
	}

	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return 0, nil, fmt.Errorf("bulk delete aborted - invalid userID format '%s': %w", userID, err)
	}

	uuids := make([]gocql.UUID, 0, len(storyID))
	for _, idStr := range storyID {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk delete aborted - invalid storyID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. Khởi tạo Channel và các Query
	type taskResult struct {
		storyID gocql.UUID
		err     error
	}
	resultCh := make(chan taskResult, len(uuids))

	tableName := entity.StoryView{}.TableName()
	readQuery := fmt.Sprintf("SELECT viewed_at FROM %s WHERE story_id = ? AND viewer_id = ? ALLOW FILTERING", tableName)
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE story_id = ? AND viewed_at = ? AND viewer_id = ?", tableName)

	// 3. Phân phối task vào Worker Pool
	// Ở Bulk Delete này, mỗi Worker sẽ chịu trách nhiệm cho 1 StoryID (Bao gồm cả Đọc và Xóa)
	for _, sID := range uuids {
		payloadSID := sID

		err := r.pool.Run(ctx, func() {
			// BƯỚC A: READ-PHASE (Giữ nguyên Ctx để hủy sớm nếu client ngắt kết nối)
			iter := r.session.Query(readQuery, payloadSID, uID).WithContext(ctx).Iter()
			scanner := iter.Scanner()

			var timestamps []time.Time
			for scanner.Next() {
				var vAt time.Time
				if err := scanner.Scan(&vAt); err == nil {
					timestamps = append(timestamps, vAt)
				}
			}

			if err := scanner.Err(); err != nil {
				resultCh <- taskResult{storyID: payloadSID, err: fmt.Errorf("read phase failed: %w", err)}
				return
			}

			if len(timestamps) == 0 {
				// Không có lượt xem nào -> Xem như xóa thành công
				resultCh <- taskResult{storyID: payloadSID, err: nil}
				return
			}

			// BƯỚC B: DELETE-PHASE (Bảo vệ Ctx vì đây là thao tác Ghi)
			safeCtx := context.WithoutCancel(ctx)
			for _, t := range timestamps {
				if execErr := r.session.Query(deleteQuery, payloadSID, t, uID).WithContext(safeCtx).Exec(); execErr != nil {
					resultCh <- taskResult{storyID: payloadSID, err: fmt.Errorf("delete phase failed at timestamp %v: %w", t, execErr)}
					return // Dừng ngay nếu 1 record bị lỗi
				}
			}

			// Thành công toàn bộ
			resultCh <- taskResult{storyID: payloadSID, err: nil}
		})

		if err != nil {
			resultCh <- taskResult{storyID: payloadSID, err: fmt.Errorf("worker pool rejected task: %w", err)}
			break
		}
	}

	// 4. Đồng bộ
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*dto.StoryViewBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &dto.StoryViewBulkError{
				StoryID: res.storyID.String(),
				UserID:  userID,
				Error:   res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Kết luận trạng thái
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete completed with errors: %d/%d success", successCount, len(uuids))
	}

	return successCount, bulkErrors, finalErr
}

func (r *StoryViewRepository) DeleteBulkStoryViewsByStoryIDAndUserID(ctx context.Context, storyID []string, userID string) (int64, []*dto.StoryViewBulkError, error) {
	// BEST PRACTICE: Tránh lặp lại mã (DRY - Don't Repeat Yourself).
	// Vì logic của hàm này giống y hệt hàm trên, ta chỉ cần định tuyến (route)
	// lại tham số và gọi thẳng vào hàm đã được implement trọn vẹn.
	return r.DeleteBulkStoryViewsByUserIAndStoryID(ctx, userID, storyID)
}
func (r *StoryViewRepository) DeleteStoryViewsByViewedAt(ctx context.Context, storyID string, viewedAt time.Time) error {
	// 1. Fail-fast validation
	if storyID == "" {
		return fmt.Errorf("storyID cannot be empty")
	}
	if viewedAt.IsZero() {
		return fmt.Errorf("viewedAt threshold cannot be zero")
	}

	id, err := gocql.ParseUUID(storyID)
	if err != nil {
		return fmt.Errorf("invalid storyID format '%s': %w", storyID, err)
	}

	// 2. Bảo vệ Context vì đây là thao tác Ghi (Tombstone)
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.StoryView{}.TableName()

	// 3. Tận dụng sức mạnh Range Deletion của Cassandra
	// Vì viewed_at là Clustering Key số 1 ngay sau Partition Key (story_id),
	// ta hoàn toàn hợp lệ khi dùng toán tử <= để xóa hàng loạt các lượt xem cũ.
	query := fmt.Sprintf("DELETE FROM %s WHERE story_id = ? AND viewed_at <= ?", tableName)

	// 4. Thực thi
	if err := r.session.Query(query, id, viewedAt).WithContext(safeCtx).Exec(); err != nil {
		return fmt.Errorf("failed to delete story views older than %v for story %s: %w", viewedAt, storyID, err)
	}

	return nil
}

func (r *StoryViewRepository) DeleteBulkStoryViewsByViewedAt(ctx context.Context, storyIDs []string, viewedAt time.Time) (int64, []*dto.StoryViewBulkError, error) {
	// 1. Fail-fast validation
	if len(storyIDs) == 0 {
		return 0, nil, nil
	}
	if viewedAt.IsZero() {
		return 0, nil, fmt.Errorf("viewedAt threshold cannot be zero")
	}

	// Chuyển đổi toàn bộ mảng UUID để bắt lỗi sớm
	uuids := make([]gocql.UUID, 0, len(storyIDs))
	for _, idStr := range storyIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk delete aborted - invalid storyID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. Chuẩn bị Struct và Channel an toàn
	type taskResult struct {
		storyID gocql.UUID
		err     error
	}
	resultCh := make(chan taskResult, len(uuids))

	tableName := entity.StoryView{}.TableName()
	query := fmt.Sprintf("DELETE FROM %s WHERE story_id = ? AND viewed_at <= ?", tableName)

	// 3. Phân phối nhiệm vụ vào Worker Pool
	for _, id := range uuids {
		payloadID := id // Copy biến an toàn

		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			
			// Bắn lệnh xóa theo khoảng thời gian
			execErr := r.session.Query(query, payloadID, viewedAt).WithContext(safeCtx).Exec()
			
			resultCh <- taskResult{
				storyID: payloadID,
				err:     execErr,
			}
		})

		// Xử lý khi Pool từ chối do quá tải hoặc Timeout
		if err != nil {
			resultCh <- taskResult{
				storyID: payloadID,
				err:     fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break
		}
	}

	// 4. Đồng bộ
	r.pool.Wait()
	close(resultCh)

	// 5. Tổng hợp dữ liệu (Lock-free)
	var successCount int64
	var bulkErrors []*dto.StoryViewBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &dto.StoryViewBulkError{
				StoryID: res.storyID.String(),
				UserID:  "range_deletion", // Đánh dấu đây là lỗi của quá trình xóa diện rộng, không phải 1 user cụ thể
				Error:   res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete by viewed_at completed with errors: %d/%d success", successCount, len(uuids))
	}

	return successCount, bulkErrors, finalErr
}