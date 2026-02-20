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

type ReactionsHistoryRepository struct {
	session   *gocql.Session
	pool      *concurrency.WorkerPool
	redisRepo IRepositoryShare.IRedis
}

func NewReactionsHistoryRepository(session *gocql.Session, pool *concurrency.WorkerPool, redisRepo IRepositoryShare.IRedis) *ReactionsHistoryRepository {
	return &ReactionsHistoryRepository{
		session:   session,
		pool:      pool,
		redisRepo: redisRepo,
	}
}
func (r *ReactionsHistoryRepository) CreateReactionHistory(ctx context.Context, reaction *entity.UserReactionHistory) error {
	// 1. FAIL-FAST: Validate đầu vào để tránh panic hoặc query lỗi
	if reaction == nil {
		return fmt.Errorf("reaction history payload is nil")
	}

	// 2. BẢO VỆ CONTEXT: Tách context để đảm bảo thao tác Ghi hoàn tất dù API bị ngắt
	safeCtx := context.WithoutCancel(ctx)

	tableName := entity.UserReactionHistory{}.CassandratableUserReactionHistory()

	// 3. Chuẩn bị query theo đúng thứ tự các trường trong struct
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, created_at, target_id, target_type, reaction_code)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	// 4. Thực thi truy vấn
	err := r.session.Query(query,
		reaction.UserID,
		reaction.CreatedAt,
		reaction.TargetID,
		reaction.TargetType,
		reaction.ReactionCode,
	).WithContext(safeCtx).Exec()

	if err != nil {
		return fmt.Errorf("failed to insert reaction history for user %s on target %s: %w", reaction.UserID.String(), reaction.TargetID, err)
	}

	return nil
}
func (r *ReactionsHistoryRepository) CreateBulkReactionHistory(ctx context.Context, reactions []*entity.UserReactionHistory) (int64, []*dto.ReactionBulkError, error) {
	if len(reactions) == 0 {
		return 0, nil, nil
	}

	// 1. Struct nội bộ vận chuyển kết quả an toàn qua Channel
	type taskResult struct {
		reaction *entity.UserReactionHistory
		err      error
	}

	// Channel đệm (Buffered Channel) bằng đúng số task để ngăn ngừa Deadlock
	resultCh := make(chan taskResult, len(reactions))
	tableName := entity.UserReactionHistory{}.CassandratableUserReactionHistory()

	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, created_at, target_id, target_type, reaction_code)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	// 2. Phân phối công việc vào Worker Pool
	for i, item := range reactions {
		// Xử lý an toàn với phần tử nil trong mảng
		if item == nil {
			resultCh <- taskResult{
				reaction: nil,
				err:      fmt.Errorf("reaction history item at index %d is nil", i),
			}
			continue
		}

		// Copy biến cục bộ để bảo vệ vùng nhớ trong closure (goroutine)
		payload := item

		err := r.pool.Run(ctx, func() {
			// Tách Context bảo vệ Write-operation
			safeCtx := context.WithoutCancel(ctx)

			execErr := r.session.Query(query,
				payload.UserID,
				payload.CreatedAt,
				payload.TargetID,
				payload.TargetType,
				payload.ReactionCode,
			).WithContext(safeCtx).Exec()

			// Bắn kết quả về Channel trung tâm
			resultCh <- taskResult{
				reaction: payload,
				err:      execErr,
			}
		})

		// Xử lý khi Pool từ chối task (VD: Context gốc timeout trước khi kịp xếp hàng)
		if err != nil {
			resultCh <- taskResult{
				reaction: payload,
				err:      fmt.Errorf("worker pool rejected task (context canceled/timeout): %w", err),
			}
			break // Dừng nhồi task mới
		}
	}

	// 3. ĐỒNG BỘ: Chờ tất cả Worker xong việc và đóng kênh
	r.pool.Wait()
	close(resultCh)

	// 4. GOM KẾT QUẢ (Không cần Mutex vì chỉ có 1 Goroutine chính đọc Channel)
	var successCount int64
	var bulkErrors []*dto.ReactionBulkError

	for res := range resultCh {
		if res.err != nil {
			// Trích xuất ID an toàn ngay cả khi pointer bị nil
			uID := "unknown_user"
			tID := "unknown_target"
			if res.reaction != nil {
				uID = res.reaction.UserID.String()
				tID = res.reaction.TargetID
			}

			bulkErrors = append(bulkErrors, &dto.ReactionBulkError{
				UserID:   uID,
				TargetID: tID,
				Error:    res.err.Error(),
			})
		} else {
			successCount++
		}
	}

	// 5. Kết luận trạng thái (Thành công một phần hoặc toàn bộ)
	var finalErr error
	if len(bulkErrors) > 0 {
		finalErr = fmt.Errorf("bulk history insert completed with errors: %d/%d success, %d failed", successCount, len(reactions), len(bulkErrors))
	}

	return successCount, bulkErrors, finalErr
}
func (r *ReactionsHistoryRepository) GetReactionHistoryByUserID(ctx context.Context, userID gocql.UUID) ([]*entity.UserReactionHistory, error) {
	// Kiểm tra UUID rỗng
	var emptyUUID gocql.UUID
	if userID == emptyUUID {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	tableName := entity.UserReactionHistory{}.CassandratableUserReactionHistory()

	// Khai báo rõ các cột. Vì CreatedAt là Clustering Key, Cassandra sẽ tự động trả về
	// theo thứ tự sắp xếp bạn đã định nghĩa trong DB (thường là DESC - mới nhất lên đầu).
	query := fmt.Sprintf(`
		SELECT user_id, created_at, target_id, target_type, reaction_code 
		FROM %s WHERE user_id = ? LIMIT 100
	`, tableName)

	var history []*entity.UserReactionHistory

	// Sử dụng ctx gốc: Nếu User tắt app/reload, DB sẽ ngừng quét ngay lập tức.
	iter := r.session.Query(query, userID).WithContext(ctx).Iter()
	scanner := iter.Scanner()

	for scanner.Next() {
		var h entity.UserReactionHistory
		err := scanner.Scan(
			&h.UserID,
			&h.CreatedAt,
			&h.TargetID,
			&h.TargetType,
			&h.ReactionCode,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history for user %s: %w", userID.String(), err)
		}
		history = append(history, &h)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error during history iteration for user %s: %w", userID.String(), err)
	}

	// Trả về mảng rỗng [] thay vì null nếu không có dữ liệu
	if history == nil {
		return []*entity.UserReactionHistory{}, nil
	}

	return history, nil
}
func (r *ReactionsHistoryRepository) PanigationReactionHistoryByUserID(ctx context.Context, userID gocql.UUID, cursor string, limit int) (*dto.PaginationRes, error) {
	var emptyUUID gocql.UUID
	if userID == emptyUUID {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	// Giải mã PageState từ Cursor (Base64)
	pageState, err := utils.DecodeCursorCassandra(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	tableName := entity.UserReactionHistory{}.CassandratableUserReactionHistory()
	query := fmt.Sprintf(`
		SELECT user_id, created_at, target_id, target_type, reaction_code 
		FROM %s WHERE user_id = ?
	`, tableName)

	// Thiết lập PageSize và PageState cho query
	q := r.session.Query(query, userID).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	iter := q.Iter()
	scanner := iter.Scanner()
	var history []*entity.UserReactionHistory

	for scanner.Next() {
		var h entity.UserReactionHistory
		err := scanner.Scan(
			&h.UserID,
			&h.CreatedAt,
			&h.TargetID,
			&h.TargetType,
			&h.ReactionCode,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error in pagination for user %s: %w", userID.String(), err)
		}
		history = append(history, &h)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error in history pagination: %w", err)
	}

	// Lấy PageState của trang kế tiếp
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	if history == nil {
		history = []*entity.UserReactionHistory{}
	}

	return &dto.PaginationRes{
		Data:       history,
		NextCursor: nextCursorStr,
		HasNext:    hasNext,
		Limit:      limit,
	}, nil
}
