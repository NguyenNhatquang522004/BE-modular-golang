package cassandra

import (
	"context"
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/infrastructure/concurrency"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/interaction/domain/entity"
	"github.com/gocql/gocql"
)

type ReactionsRepository struct {
	session   *gocql.Session
	pool      *concurrency.WorkerPool
	redisRepo IRepositoryShare.IRedis
}

func NewReactionsRepository(session *gocql.Session, pool *concurrency.WorkerPool, redisRepo IRepositoryShare.IRedis) *ReactionsRepository {
	return &ReactionsRepository{
		session:   session,
		pool:      pool,
		redisRepo: redisRepo,
	}
}
func (r *ReactionsRepository) CreateReaction(ctx context.Context, reaction *entity.EntityReaction) error {
	if reaction == nil {
		return fmt.Errorf("reaction payload is nil")
	}

	// 1. Tách context để đảm bảo lệnh Ghi (Write) không bị ngắt nếu HTTP Request bị Cancel/Timeout
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.EntityReaction{}.CassandratableEntityReaction()

	// 2. Query chuẩn bị (Sắp xếp đúng thứ tự cột)
	query := fmt.Sprintf(`
        INSERT INTO %s (target_id, user_id, target_type, reaction_code, created_at)
        VALUES (?, ?, ?, ?, ?)
    `, tableName)

	// 3. Thực thi query
	err := r.session.Query(query,
		reaction.TargetID,
		reaction.UserID,
		reaction.TargetType,
		reaction.ReactionCode,
		reaction.CreatedAt,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to create reaction for target %s by user %s: %w", reaction.TargetID, reaction.UserID.String(), err)
	}

	return nil
}
func (r *ReactionsRepository) CreateBulkReactions(ctx context.Context, reactions []*entity.EntityReaction) (int64, []*dto.ReactionBulkError, error) {
	if len(reactions) == 0 {
		return 0, nil, nil
	}

	// 1. Struct nội bộ dùng để hứng kết quả từ Goroutine một cách an toàn
	type taskResult struct {
		reaction *entity.EntityReaction
		err      error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để ngăn chặn chặn deadlock
	resultCh := make(chan taskResult, len(reactions))

	tableName := entity.EntityReaction{}.CassandratableEntityReaction()
	query := fmt.Sprintf(`
        INSERT INTO %s (target_id, user_id, target_type, reaction_code, created_at)
        VALUES (?, ?, ?, ?, ?)
    `, tableName)

	// 2. Đẩy task vào Worker Pool
	for _, item := range reactions {
		if item == nil {
			resultCh <- taskResult{reaction: nil, err: fmt.Errorf("reaction item is nil")}
			continue
		}

		payload := item // Copy biến để goroutine không trỏ sai vùng nhớ (Rất quan trọng ở Go < 1.22)

		err := r.pool.Run(ctx, func() {
			// Tách context bảo vệ ghi giống như single insert
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				payload.TargetID,
				payload.UserID,
				payload.TargetType,
				payload.ReactionCode,
				payload.CreatedAt,
			).WithContext(safeCtx).Exec()

			// Bắn kết quả về channel
			resultCh <- taskResult{reaction: payload, err: execErr}
		})

		// Xử lý khi worker pool bị đầy hoặc context gốc bị timeout không nhận thêm task
		if err != nil {
			resultCh <- taskResult{
				reaction: payload,
				err:      fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break // Dừng việc nhồi thêm task, những task đã vào pool vẫn sẽ chạy
		}
	}

	// 3. Chờ tất cả Worker hoàn thành và đóng Channel an toàn
	r.pool.Wait()
	close(resultCh)

	// 4. Tổng hợp dữ liệu (Lock-free vì chỉ 1 main goroutine đọc channel)
	var successCount int64
	var bulkErrors []*dto.ReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			// Xử lý extract thông tin để đưa vào BulkError
			targetID := "unknown_target"
			userID := "unknown_user"
			if res.reaction != nil {
				targetID = res.reaction.TargetID
				userID = res.reaction.UserID.String()
			}

			bulkErrors = append(bulkErrors, &dto.ReactionBulkError{
				TargetID: targetID,
				UserID:   userID,
				Error:    res.err.Error(),
			})
		} else {
			successCount++ // Tăng biến đếm thành công
		}
	}

	// 5. Đánh giá trạng thái trả về (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		// Trả về lỗi tổng quan nếu có bất kỳ record nào thất bại
		finalErr = fmt.Errorf("bulk insert completed with errors: %d/%d success, %d failed", successCount, len(reactions), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *ReactionsRepository) GetReactionsByTargetID(ctx context.Context, targetID string) ([]*entity.EntityReaction, error) {
	if targetID == "" {
		return nil, fmt.Errorf("targetID cannot be empty")
	}

	tableName := entity.EntityReaction{}.CassandratableEntityReaction()

	// Khai báo rõ các cột để tối ưu Serialize/Deserialize thay vì dùng SELECT *
	query := fmt.Sprintf(`
		SELECT target_id, user_id, target_type, reaction_code, created_at 
		FROM %s WHERE target_id = ?
	`, tableName)

	// Khởi tạo mảng kết quả
	// Không pre-allocate len() vì ta chưa biết trước có bao nhiêu dòng trả về
	var reactions []*entity.EntityReaction

	// Dùng Iter() và Scanner() để duyệt qua các dòng kết quả từ Cassandra
	// Bắt buộc truyền ctx gốc để ngắt kết nối DB sớm nếu HTTP request bị hủy
	iter := r.session.Query(query, targetID).WithContext(ctx).Iter()
	scanner := iter.Scanner()

	for scanner.Next() {
		var rct entity.EntityReaction

		// Map thẳng dữ liệu vào struct. Các kiểu Enum tự định nghĩa (TargetType, ReactionCode)
		// sẽ được gocql tự động map nếu underlying type tương thích (string/int).
		err := scanner.Scan(
			&rct.TargetID,
			&rct.UserID,
			&rct.TargetType,
			&rct.ReactionCode,
			&rct.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reaction for target %s: %w", targetID, err)
		}

		// Append pointer của biến cục bộ. (Từ Go 1.22 vòng lặp đã an toàn với pointer)
		reactions = append(reactions, &rct)
	}

	// Kiểm tra lỗi sau khi vòng lặp kết thúc (ví dụ: mất kết nối giữa chừng)
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("iteration error when fetching reactions for target %s: %w", targetID, err)
	}

	// Trả về slice rỗng (thay vì nil) nếu không có dữ liệu, giúp Frontend/Controller dễ map sang mảng JSON rỗng []
	if reactions == nil {
		return []*entity.EntityReaction{}, nil
	}

	return reactions, nil
}
func (r *ReactionsRepository) GetReactionsByUserID(ctx context.Context, userID gocql.UUID) ([]*entity.EntityReaction, error) {
	// Kiểm tra UUID rỗng (Nil UUID)
	var emptyUUID gocql.UUID
	if userID == emptyUUID {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	tableName := entity.EntityReaction{}.CassandratableEntityReaction()

	// Lưu ý: Nếu bạn dùng Materialized View, hãy đổi tableName thành tên của View đó.
	query := fmt.Sprintf(`
		SELECT target_id, user_id, target_type, reaction_code, created_at 
		FROM %s WHERE user_id = ?
	`, tableName)

	var reactions []*entity.EntityReaction

	iter := r.session.Query(query, userID).WithContext(ctx).Iter()
	scanner := iter.Scanner()

	for scanner.Next() {
		var rct entity.EntityReaction
		err := scanner.Scan(
			&rct.TargetID,
			&rct.UserID,
			&rct.TargetType,
			&rct.ReactionCode,
			&rct.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reaction for user %s: %w", userID.String(), err)
		}

		reactions = append(reactions, &rct)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("iteration error when fetching reactions for user %s: %w", userID.String(), err)
	}

	if reactions == nil {
		return []*entity.EntityReaction{}, nil
	}

	return reactions, nil
}
func (r *ReactionsRepository) DeleteReaction(ctx context.Context, targetID string, userID gocql.UUID) error {
	// 1. Fail-fast: Validate input
	if targetID == "" {
		return fmt.Errorf("targetID cannot be empty")
	}
	var emptyUUID gocql.UUID
	if userID == emptyUUID {
		return fmt.Errorf("userID cannot be empty")
	}

	// 2. Bảo vệ thao tác xóa.
	// Dù client ngắt kết nối, lệnh xóa vẫn phải được thực thi trọn vẹn để tránh rác dữ liệu.
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.EntityReaction{}.CassandratableEntityReaction()

	// Query xóa dựa vào trọn bộ Primary Key (TargetID + UserID)
	query := fmt.Sprintf("DELETE FROM %s WHERE target_id = ? AND user_id = ?", tableName)

	// 3. Thực thi
	err := r.session.Query(query, targetID, userID).WithContext(safeCtx).Exec()
	if err != nil {
		return fmt.Errorf("failed to delete reaction for target %s by user %s: %w", targetID, userID.String(), err)
	}

	return nil
}
func (r *ReactionsRepository) DeleteBulkReactions(ctx context.Context, targetIDs []*string, userIDs []*gocql.UUID) (int64, []*dto.ReactionBulkError, error) {
	// 1. FAIL-FAST: Kiểm tra tính hợp lệ của input
	if len(targetIDs) == 0 || len(userIDs) == 0 {
		return 0, nil, nil
	}

	// Bắt buộc 2 mảng phải có độ dài bằng nhau để ghép cặp chính xác
	if len(targetIDs) != len(userIDs) {
		return 0, nil, fmt.Errorf("mismatched input lengths: targetIDs length (%d) != userIDs length (%d)", len(targetIDs), len(userIDs))
	}

	totalTasks := len(targetIDs)
	tableName := entity.EntityReaction{}.CassandratableEntityReaction()

	// 2. Struct nội bộ chứa payload cho mỗi task và kết quả trả về
	type deletePayload struct {
		TargetID string
		UserID   gocql.UUID
	}

	type taskResult struct {
		payload deletePayload
		err     error
	}

	// Sử dụng Buffered Channel bằng đúng số lượng task để tránh deadlock
	resultCh := make(chan taskResult, totalTasks)

	// Câu lệnh CQL chuẩn để xóa 1 row dựa trên toàn bộ Primary Key
	query := fmt.Sprintf("DELETE FROM %s WHERE target_id = ? AND user_id = ?", tableName)

	// 3. Phân phối task vào Worker Pool
	for i := 0; i < totalTasks; i++ {
		// Xử lý an toàn con trỏ (Safe Pointer Dereference)
		if targetIDs[i] == nil || userIDs[i] == nil {
			tID := "nil_pointer"

			if targetIDs[i] != nil {
				tID = *targetIDs[i]
			}

			// Trả thẳng lỗi vào channel mà không cần gọi DB
			resultCh <- taskResult{
				payload: deletePayload{TargetID: tID, UserID: gocql.UUID{}}, // Gắn ID tạm để log
				err:     fmt.Errorf("nil pointer at index %d", i),
			}
			continue
		}

		// Khởi tạo payload hợp lệ
		payload := deletePayload{
			TargetID: *targetIDs[i],
			UserID:   *userIDs[i],
		}

		// Đưa task vào Worker Pool
		err := r.pool.Run(ctx, func() {
			// BEST PRACTICE: Trong Cassandra, thao tác Xóa (Delete) thực chất là thao tác Ghi (Tombstone).
			// Do đó, PHẢI bảo vệ Context giống như thao tác Insert/Update.
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query, payload.TargetID, payload.UserID).WithContext(safeCtx).Exec()

			// Bắn kết quả vào channel
			resultCh <- taskResult{
				payload: payload,
				err:     execErr,
			}
		})

		// Xử lý khi Worker Pool từ chối task (VD: Context gốc đã hết hạn trước khi kịp xếp hàng)
		if err != nil {
			resultCh <- taskResult{
				payload: payload,
				err:     fmt.Errorf("worker pool rejected delete task (context canceled/timeout): %w", err),
			}
			break // Cắt đứt vòng lặp để không nhồi thêm task vô ích
		}
	}

	// 4. Đồng bộ hóa: Chờ các goroutine chạy xong và đóng channel
	r.pool.Wait()
	close(resultCh)

	// 5. Gom kết quả (Lock-free Aggregation)
	var successCount int64
	var bulkErrors []*dto.ReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			bulkErrors = append(bulkErrors, &dto.ReactionBulkError{
				TargetID: res.payload.TargetID,
				UserID:   res.payload.UserID.String(), // Fallback về chuỗi rỗng nếu UUID chưa được gán chuẩn do lỗi nil pointer
				Error:    res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 6. Đánh giá trạng thái trả về
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk delete completed with errors: %d/%d success, %d failed", successCount, totalTasks, len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *ReactionsRepository) UpdateReaction(ctx context.Context, reaction *entity.EntityReaction) error {
	// 1. FAIL-FAST: Validate đầu vào
	if reaction == nil {
		return fmt.Errorf("reaction payload is nil")
	}
	if reaction.TargetID == "" || reaction.UserID.String() == "00000000-0000-0000-0000-000000000000" {
		return fmt.Errorf("missing primary key: target_id and user_id are required for update")
	}

	// 2. BẢO VỆ CONTEXT: Đảm bảo thao tác ghi không bị đứt gãy nếu HTTP Request bị ngắt
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.EntityReaction{}.CassandratableEntityReaction()

	// 3. Cú pháp UPDATE chuẩn của Cassandra (SET các data fields, WHERE bằng toàn bộ Primary Key)
	query := fmt.Sprintf(`
		UPDATE %s 
		SET target_type = ?, reaction_code = ?, created_at = ? 
		WHERE target_id = ? AND user_id = ?
	`, tableName)

	// 4. Thực thi query
	err := r.session.Query(query,
		reaction.TargetType,
		reaction.ReactionCode,
		reaction.CreatedAt, // Lưu ý: Nếu có trường UpdatedAt thì nên dùng UpdatedAt ở đây
		reaction.TargetID,
		reaction.UserID,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to update reaction for target %s by user %s: %w", reaction.TargetID, reaction.UserID.String(), err)
	}

	return nil
}
func (r *ReactionsRepository) UpdateBulkReactions(ctx context.Context, reactions []*entity.EntityReaction) (int64, []*dto.ReactionBulkError, error) {
	if len(reactions) == 0 {
		return 0, nil, nil
	}

	// 1. Struct nội bộ để vận chuyển kết quả qua Channel an toàn
	type taskResult struct {
		reaction *entity.EntityReaction
		err      error
	}

	// Khởi tạo Buffered Channel với dung lượng bằng số lượng task (Chống Deadlock)
	resultCh := make(chan taskResult, len(reactions))
	tableName := entity.EntityReaction{}.CassandratableEntityReaction()

	// Cú pháp UPDATE
	query := fmt.Sprintf(`
		UPDATE %s 
		SET target_type = ?, reaction_code = ?, created_at = ? 
		WHERE target_id = ? AND user_id = ?
	`, tableName)

	// 2. Phân phối công việc vào Worker Pool
	for i, item := range reactions {
		// Xử lý an toàn con trỏ nil
		if item == nil {
			resultCh <- taskResult{
				reaction: nil,
				err:      fmt.Errorf("reaction pointer at index %d is nil", i),
			}
			continue
		}

		// Copy biến cục bộ để Closure của Goroutine không tham chiếu sai vùng nhớ
		payload := item

		err := r.pool.Run(ctx, func() {
			// BEST PRACTICE: Tách context bảo vệ thao tác Ghi
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				payload.TargetType,
				payload.ReactionCode,
				payload.CreatedAt,
				payload.TargetID,
				payload.UserID,
			).WithContext(safeCtx).Exec()

			// Trả kết quả về channel
			resultCh <- taskResult{
				reaction: payload,
				err:      execErr,
			}
		})

		// Xử lý khi Pool bị quá tải hoặc Context gốc đã timeout/cancel không nhận thêm task
		if err != nil {
			resultCh <- taskResult{
				reaction: payload,
				err:      fmt.Errorf("worker pool rejected update task (context canceled): %w", err),
			}
			break // Cắt đứt vòng lặp để giữ an toàn cho hệ thống
		}
	}

	// 3. ĐỒNG BỘ: Đợi các worker đang chạy hoàn tất và đóng kênh
	r.pool.Wait()
	close(resultCh)

	// 4. TỔNG HỢP KẾT QUẢ (Lock-Free)
	var successCount int64
	var bulkErrors []*dto.ReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			// Trích xuất ID an toàn ngay cả khi payload bị nil
			tID := "unknown_target"
			uID := "unknown_user"
			if res.reaction != nil {
				tID = res.reaction.TargetID
				uID = res.reaction.UserID.String()
			}

			bulkErrors = append(bulkErrors, &dto.ReactionBulkError{
				TargetID: tID,
				UserID:   uID,
				Error:    res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 5. Kết luận trạng thái (Partial Success)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk update completed with errors: %d/%d success, %d failed", successCount, len(reactions), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *ReactionsRepository) PaginateReactionsByTargetID(ctx context.Context, targetID string, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Check Cache (Chỉ check khi gọi trang đầu tiên)
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			"reaction_cache_targetid_" + targetID,
			"reaction_cache_nextcursor_targetid_" + targetID,
			"reaction_cache_hasnext_targetid_" + targetID,
			"reaction_cache_limit_targetid_" + targetID,
		})

		if err == nil && datacache != nil {
			if reactionList, ok := datacache.([]*entity.EntityReaction); ok {
				return &dto.PaginationRes{
					Data:       reactionList,
					NextCursor: nextcursor,
					HasNext:    hasnext,
					Limit:      limitcache,
				}, nil
			}
		}
	}

	// 2. Decode Cursor (Lấy Page State)
	pageState, err := utils.DecodeCursorCassandra(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	// 3. Chuẩn bị Query cho Cassandra
	tableName := entity.EntityReaction{}.CassandratableEntityReaction()
	query := fmt.Sprintf(`
		SELECT target_id, user_id, target_type, reaction_code, created_at 
		FROM %s WHERE target_id = ?
	`, tableName)

	// Cấu hình query với PageSize và PageState
	q := r.session.Query(query, targetID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 4. Thực thi và duyệt kết quả
	iter := q.Iter()
	scanner := iter.Scanner()
	var reactions []*entity.EntityReaction

	for scanner.Next() {
		var rct entity.EntityReaction
		err := scanner.Scan(
			&rct.TargetID,
			&rct.UserID,
			&rct.TargetType,
			&rct.ReactionCode,
			&rct.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reaction for target %s: %w", targetID, err)
		}
		reactions = append(reactions, &rct)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination: %w", err)
	}

	// 5. Xử lý Next Cursor và HasNext
	// iter.PageState() trả về state cho trang kế tiếp. Nếu len > 0 nghĩa là còn dữ liệu.
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	// Nếu trả về slice nil thì gán bằng mảng rỗng để tránh lỗi null ở FE
	if reactions == nil {
		reactions = []*entity.EntityReaction{}
	}

	// 6. Cache kết quả cho trang đầu tiên
	if cursor == "" && len(reactions) > 0 {
		items := map[string]any{
			"reaction_cache_targetid_" + targetID:            reactions,
			"reaction_cache_nextcursor_targetid_" + targetID: nextCursorStr,
			"reaction_cache_hasnext_targetid_" + targetID:    hasNext,
			"reaction_cache_limit_targetid_" + targetID:      limit,
		}

		err = r.redisRepo.CustomizeSetCache(ctx, items)
		if err != nil {
			// Chỉ log lỗi chứ không chặn flow vì cache miss/error không nên làm hỏng main flow
			fmt.Printf("Failed to set cache for reaction pagination: %v\n", err)
		}
	}

	// 7. Trả về kết quả
	return &dto.PaginationRes{
		Data:       reactions,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
