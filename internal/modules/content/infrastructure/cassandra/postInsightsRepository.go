package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/delivery/mapper"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/entity"
	"github.com/gocql/gocql"
)

type PostInsightsRepository struct {
	session *gocql.Session
	pool    IRepositoryShare.IWorkerPool
}

func NewPostInsightsRepository(session *gocql.Session, pool IRepositoryShare.IWorkerPool) *PostInsightsRepository {
	return &PostInsightsRepository{session: session, pool: pool}
}
func (r *PostInsightsRepository) CreatePostInsightInitPost(ctx context.Context, postinsight *entity.PostInsight) error {
	tableName := entity.PostInsight{}.Collectionnamepostinsight()
	if postinsight.PostID == (gocql.UUID{}) {
		postinsight.PostID = gocql.TimeUUID() // Tạo UUID mới nếu chưa có (Thường thì nên có sẵn từ bên ngoài, nhưng đây là biện pháp phòng ngừa)
	}
	query := fmt.Sprintf(`insert into %s ( post_id, reach, impressions, engagement_rate, reactions_total, comments_total, shares_total, clicks_total, video_views_3s, updated_at) values (?, 0, 0, 0.0, 0, 0, 0, 0, 0, ?)`, tableName)
	return r.session.Query(query, postinsight.PostID, time.Now()).WithContext(ctx).Exec()
}
func (r *PostInsightsRepository) CreatePostInsightInitPostBulk(ctx context.Context, postinsights []*entity.PostInsight) (int64, []*cassandraErrors.InsightBulkError, error) {
	if len(postinsights) == 0 {
		return 0, nil, nil
	}
	type taskResult struct {
		message *entity.PostInsight
		err     error
	}
	// 1. Parse toàn bộ UUID trước. Fail-fast: Nếu có 1 ID lỗi, dừng luôn không insert gì cả.
	uuids := make([]*entity.PostInsight, 0, len(postinsights))
	for _, postInsight := range postinsights {
		if postInsight.PostID == (gocql.UUID{}) {
			postInsight.PostID = gocql.TimeUUID() // Tạo UUID mới nếu chưa có (Thường thì nên có sẵn từ bên ngoài, nhưng đây là biện pháp phòng ngừa)
		}
		uuids = append(uuids, postInsight)
	}
	var faildocs []*cassandraErrors.InsightBulkError
	tableName := entity.PostInsight{}.Collectionnamepostinsight()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			post_id, reach, impressions, engagement_rate, 
			reactions_total, comments_total, shares_total, 
			clicks_total, video_views_3s, updated_at
		) VALUES (?, 0, 0, 0.0, 0, 0, 0, 0, 0, ?)
	`, tableName)

	now := time.Now()

	// 2. Tích hợp Worker Pool

	errCh := make(chan taskResult, len(uuids))

	for _, id := range uuids {
		item := id
		err := r.pool.Run(ctx, func() {
			// TÁCH CONTEXT Ở ĐÂY: Quyết tâm insert cho xong dù HTTP request đã bị ngắt!
			safeCtx := context.WithoutCancel(ctx)
			err := r.session.Query(
				query,
				item.PostID,
				item.Impressions, item.Reach, item.EngagementRate,
				item.ReactionsTotal, item.CommentsTotal, item.SharesTotal,
				item.ClicksTotal, item.VideoViews3s,
				now,
			).WithContext(safeCtx).Exec()
			if err != nil {
				errCh <- taskResult{message: nil, err: err}
			}
		})
		if err != nil {
			errCh <- taskResult{message: nil, err: err}
			// 2. BREAK NGAY LẬP TỨC! Không loop tiếp để tiết kiệm tài nguyên.
			break
		}
	}

	r.pool.Wait()
	close(errCh)
	// 3. Kiểm tra xem có lỗi nào xảy ra trong quá trình insert không
	// Trả về lỗi đầu tiên gặp phải
	for res := range errCh {
		if res.err != nil {
			PostID := ""
			if res.message != nil {
				PostID = res.message.PostID.String()
			}
			faildocs = append(faildocs, &cassandraErrors.InsightBulkError{
				PostID: PostID,
				Error:  res.err.Error(),
			})
		}

	}

	return int64(len(postinsights)), faildocs, nil
}
func (r *PostInsightsRepository) doUpdateInteraction(ctx context.Context, postID gocql.UUID, interactionType string, countDelta int) error {
	colName, err := mapper.MapInteractionTypeToColumn(interactionType)
	if err != nil {
		return err // Fail-fast nếu type sai
	}

	tableName := entity.PostInsight{}.Collectionnamepostinsight()

	// BƯỚC 1: Đọc giá trị hiện tại (Read)
	var currentCount int
	readQuery := fmt.Sprintf("SELECT %s FROM %s WHERE post_id = ?", colName, tableName)
	if err := r.session.Query(readQuery, postID).WithContext(ctx).Scan(&currentCount); err != nil {
		if err == gocql.ErrNotFound {
			currentCount = 0 // Coi như 0 nếu chưa có record
		} else {
			return fmt.Errorf("failed to read current %s for post %s: %w", colName, postID, err)
		}
	}

	// BƯỚC 2: Tính toán giá trị mới (Modify)
	newCount := currentCount + countDelta
	if newCount < 0 {
		newCount = 0 // Đảm bảo insight không bị số âm
	}

	// BƯỚC 3: Cập nhật xuống DB (Write)
	updateQuery := fmt.Sprintf("UPDATE %s SET %s = ?, updated_at = ? WHERE post_id = ?", tableName, colName)
	if err := r.session.Query(updateQuery, newCount, time.Now(), postID).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to update %s for post %s: %w", colName, postID, err)
	}

	return nil
}
func (r *PostInsightsRepository) UpdatePostInsightInteraction(ctx context.Context, PostID string, interactionType string, count int) error {
	// 1. Validate UUID ngay lập tức
	id, err := gocql.ParseUUID(PostID)
	if err != nil {
		return fmt.Errorf("invalid postID format '%s': %w", PostID, err)
	}

	// 2. Gọi logic dùng chung
	return r.doUpdateInteraction(ctx, id, interactionType, count)
}

// func (r *PostInsightsRepository) UpdatePostInsightInteractionBulk(ctx context.Context, reqs []*req.UpdatePostInsightsInteractionReq) (int64, []*cassandraErrors.InsightBulkError, error) {
// 	if len(reqs) == 0 {
// 		return 0, nil, nil
// 	}

// 	// 1. Giai đoạn tiền xử lý (Fail-fast): Validate toàn bộ dữ liệu trước khi đụng vào DB
// 	type taskResult struct {
// 		message *entity.PostInsight
// 		err     error
// 	}

// 	tasks := make([]taskResult, 0, len(reqs))
// 	for _, item := range reqs {
// 		if item == nil || item.PostInsightsReq == nil {
// 			continue // Bỏ qua nếu payload nil
// 		}
// 		id, err := gocql.ParseUUID(item.PostID) // Giả định item.PostID nằm trong PostInsightsReq
// 		if err != nil {
// 			return 0, nil, fmt.Errorf("bulk aborted - invalid postID '%s': %w", item.PostID, err)
// 		}
// 	}

// 	// 2. Chuẩn bị Worker Pool & Error Channel
// 	errCh := make(chan taskResult, len(tasks))

// 	// 3. Phân phối công việc vào Pool
// 	for _, t := range tasks {
// 		payload := t // Copy biến (An toàn cho Go < 1.22)

// 		err := r.pool.Run(ctx, func() {
// 			// BEST PRACTICE: Tách context để đảm bảo DB query chạy xong dù HTTP Request bị cancel
// 			safeCtx := context.WithoutCancel(ctx)

// 			// Thực thi logic Read-Modify-Write
// 			if err := r.doUpdateInteraction(safeCtx, payload.PostID, payload.InteractionType, payload.Count); err != nil {
// 				errCh <- taskResult{message: nil, err: err}
// 			}
// 		})

// 		// Nếu Pool từ chối task (vì Context gốc đã bị cancel hoặc timeout)
// 		if err != nil {
// 			errCh <- taskResult{message: nil, err: fmt.Errorf("worker pool rejected task for post %s: %w", payload.PostID, err)}
// 			break // Dừng việc nhồi thêm task, nhưng vẫn phải Wait() các task đang chạy
// 		}
// 	}

// 	// 4. Dọn dẹp & Đồng bộ
// 	r.pool.Wait()
// 	close(errCh)

// 	// 5. Gom lỗi (Chỉ trả về lỗi đầu tiên hoặc dùng errors.Join nếu dùng Go 1.20+)
// 	for err := range errCh {
// 		if err.err != nil {
// 			return 0, nil, err.err
// 		}
// 	}

//		return int64(len(reqs)), nil, nil
//	}
func (r *PostInsightsRepository) doUpdateLifeTime(ctx context.Context, postID gocql.UUID, metricType string, value float64) error {
	colName, err := mapper.MapMetricTypeToColumn(metricType)
	if err != nil {
		return err // Fail-fast nếu metricType bị sai
	}

	tableName := entity.PostInsight{}.Collectionnamepostinsight()

	// BƯỚC 1: Đọc giá trị hiện tại (Chỉ cần nếu bạn muốn + thêm vào giá trị cũ)
	var currentValue float64
	readQuery := fmt.Sprintf("SELECT %s FROM %s WHERE post_id = ?", colName, tableName)
	if err := r.session.Query(readQuery, postID).WithContext(ctx).Scan(&currentValue); err != nil {
		if err == gocql.ErrNotFound {
			currentValue = 0 // Bắt đầu từ 0 nếu chưa có dữ liệu
		} else {
			return fmt.Errorf("failed to read %s for post %s: %w", colName, postID, err)
		}
	}

	// BƯỚC 2: Tính toán (Cộng dồn giá trị mới)
	newValue := currentValue + value
	if newValue < 0 {
		newValue = 0 // Ngăn chặn chỉ số bị âm
	}

	// BƯỚC 3: Ghi xuống DB
	updateQuery := fmt.Sprintf("UPDATE %s SET %s = ?, updated_at = ? WHERE post_id = ?", tableName, colName)
	if err := r.session.Query(updateQuery, newValue, time.Now(), postID).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to update %s for post %s: %w", colName, postID, err)
	}

	return nil
}
func (r *PostInsightsRepository) UpdatePostInsightLifeTime(ctx context.Context, PostID string, metricType string, value float64) error {
	id, err := gocql.ParseUUID(PostID)
	if err != nil {
		return fmt.Errorf("invalid postID format '%s': %w", PostID, err)
	}

	return r.doUpdateLifeTime(ctx, id, metricType, value)
}

// func (r *PostInsightsRepository) UpdatePostInsightLifeTimeBulk(ctx context.Context, reqs []*req.UpdatePostInsightsLifeTimeReq) error {
// 	if len(reqs) == 0 {
// 		return nil
// 	}

// 	// 1. FAIL-FAST: Validate toàn bộ dữ liệu trước khi đẩy vào WorkerPool
// 	type validTask struct {
// 		PostID     gocql.UUID
// 		MetricType string
// 		Value      float64
// 	}

// 	tasks := make([]validTask, 0, len(reqs))
// 	for _, item := range reqs {
// 		if item == nil || item.PostInsightsReq == nil {
// 			continue // Bỏ qua payload rỗng để tránh panic nil pointer
// 		}

// 		id, err := gocql.ParseUUID(item.PostID)
// 		if err != nil {
// 			// Sai 1 cái là từ chối toàn bộ mảng (Fail-fast)
// 			return fmt.Errorf("bulk aborted - invalid postID '%s': %w", item.PostID, err)
// 		}

// 		tasks = append(tasks, validTask{
// 			PostID:     id,
// 			MetricType: item.MetricType,
// 			Value:      item.Value,
// 		})
// 	}

// 	// 2. Chạy Async với Worker Pool
// 	errCh := make(chan error, len(tasks))

// 	for _, t := range tasks {
// 		payload := t // Copy biến an toàn cho Goroutine (Bắt buộc với Go < 1.22)

// 		err := r.pool.Run(ctx, func() {
// 			// BẢO VỆ CONTEXT: Ngăn việc HTTP bị ngắt làm hỏng dữ liệu Bulk Insert
// 			safeCtx := context.WithoutCancel(ctx)

// 			// Gọi lại hàm doUpdateLifeTime dùng chung
// 			if err := r.doUpdateLifeTime(safeCtx, payload.PostID, payload.MetricType, payload.Value); err != nil {
// 				errCh <- err
// 			}
// 		})

// 		// 3. XỬ LÝ LỖI WORKER POOL
// 		if err != nil {
// 			errCh <- fmt.Errorf("worker pool rejected task for post %s (Context canceled/Timeout): %w", payload.PostID, err)
// 			break // Cắt đứt vòng lặp ngay lập tức
// 		}
// 	}

// 	// 4. CHỜ VÀ DỌN DẸP SẠCH SẼ
// 	r.pool.Wait() // Đợi các task đã xin được slot chạy xong
// 	close(errCh)

// 	// 5. Gom lỗi trả về (Trả về lỗi đầu tiên tìm thấy)
// 	for err := range errCh {
// 		if err != nil {
// 			return err
// 		}
// 	}

//		return nil
//	}
func (r *PostInsightsRepository) doGetInsight(ctx context.Context, postID gocql.UUID) (*entity.PostInsight, error) {
	tableName := entity.PostInsight{}.Collectionnamepostinsight()

	// Khai báo rõ các cột để tối ưu hoá tốc độ Scan thay vì dùng SELECT *
	query := fmt.Sprintf(`
		SELECT post_id, reach, impressions, engagement_rate, 
		       reactions_total, comments_total, shares_total, 
		       clicks_total, video_views_3s, updated_at 
		FROM %s WHERE post_id = ?`, tableName)

	var insight entity.PostInsight

	// Thực thi và map dữ liệu
	err := r.session.Query(query, postID).WithContext(ctx).Scan(
		&insight.PostID, &insight.Reach, &insight.Impressions, &insight.EngagementRate,
		&insight.ReactionsTotal, &insight.CommentsTotal, &insight.SharesTotal,
		&insight.ClicksTotal, &insight.VideoViews3s, &insight.UpdatedAt,
	)

	if err != nil {
		if err == gocql.ErrNotFound {
			// Best Practice: Trả về (nil, nil) nếu không tìm thấy để phân biệt với lỗi hệ thống
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query post insight for post %s: %w", postID, err)
	}

	return &insight, nil
}
func (r *PostInsightsRepository) GetPostInsightByPostID(ctx context.Context, PostID string) (*entity.PostInsight, error) {
	// 1. Validate định dạng UUID
	id, err := gocql.ParseUUID(PostID)
	if err != nil {
		return nil, fmt.Errorf("invalid postID format '%s': %w", PostID, err)
	}

	// 2. Gọi hàm core
	return r.doGetInsight(ctx, id)
}
func (r *PostInsightsRepository) GetPostInsightsByPostIDs(ctx context.Context, PostIDs []string) ([]*entity.PostInsight, error) {
	if len(PostIDs) == 0 {
		return []*entity.PostInsight{}, nil
	}

	// 1. FAIL-FAST: Validate toàn bộ UUID
	uuids := make([]gocql.UUID, 0, len(PostIDs))
	for _, idStr := range PostIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid postID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, id)
	}

	// 2. PRE-ALLOCATE MẢNG KẾT QUẢ (THỦ THUẬT KHÔNG DÙNG MUTEX)
	// Khởi tạo mảng có sẵn độ dài bằng len(uuids).
	// Mỗi goroutine sẽ tự điền vào đúng vị trí của nó (index) -> Không bị đụng độ biến (Race condition).
	results := make([]*entity.PostInsight, len(uuids))
	errCh := make(chan error, len(uuids))

	// 3. Sử dụng Worker Pool để đọc song song
	for i, id := range uuids {
		index := i // Copy biến cho goroutine
		postID := id

		err := r.pool.Run(ctx, func() {
			// LƯU Ý QUAN TRỌNG: Ở đây ta KHÔNG DÙNG context.WithoutCancel(ctx)
			// Giữ nguyên ctx gốc để nếu API bị hủy, DB sẽ lập tức ngừng query.

			insight, err := r.doGetInsight(ctx, postID)
			if err != nil {
				errCh <- err
				return
			}

			// Ghi thẳng vào index tương ứng, hoàn toàn Thread-Safe
			results[index] = insight
		})

		if err != nil {
			errCh <- fmt.Errorf("worker pool rejected read task for post %s: %w", postID, err)
			break
		}
	}

	// 4. Đồng bộ hóa và kiểm tra lỗi
	r.pool.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err // Có 1 query lỗi DB/Timeout thì báo lỗi luôn
		}
	}

	// 5. Lọc kết quả (Loại bỏ các giá trị nil sinh ra do bài viết không tồn tại trong DB)
	finalResults := make([]*entity.PostInsight, 0, len(results))
	for _, res := range results {
		if res != nil {
			finalResults = append(finalResults, res)
		}
	}

	return finalResults, nil
}
func (r *PostInsightsRepository) doDeleteInsight(ctx context.Context, postID gocql.UUID) error {
	tableName := entity.PostInsight{}.Collectionnamepostinsight()

	// Query xóa toàn bộ dòng (row) dựa trên Partition Key
	query := fmt.Sprintf("DELETE FROM %s WHERE post_id = ?", tableName)

	if err := r.session.Query(query, postID).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to delete post insight for post_id %s: %w", postID.String(), err)
	}

	return nil
}

func (r *PostInsightsRepository) DeletePostInsightByPostID(ctx context.Context, PostID string) error {
	// 1. Validate UUID ngay lập tức
	id, err := gocql.ParseUUID(PostID)
	if err != nil {
		return fmt.Errorf("invalid postID format '%s': %w", PostID, err)
	}

	// 2. Bảo vệ thao tác xóa (Vì Delete cũng là một dạng Write trong Cassandra)
	// Đảm bảo dù API bị timeout, lệnh xóa vẫn được đẩy xuống DB trọn vẹn
	safeCtx := context.WithoutCancel(ctx)

	// 3. Gọi hàm Core
	return r.doDeleteInsight(safeCtx, id)
}
func (r *PostInsightsRepository) DeletePostInsightsByPostIDBulk(ctx context.Context, PostIDs []string) (int64, []*cassandraErrors.InsightBulkError, error) {
	if len(PostIDs) == 0 {
		return 0, nil, nil
	}
	type taskResult struct {
		message *entity.PostInsight
		err     error
	}
	var faildocs []*cassandraErrors.InsightBulkError
	// 1. FAIL-FAST VALIDATION: Đảm bảo toàn bộ mảng ID hợp lệ trước khi chạm vào DB
	uuids := make([]*gocql.UUID, 0, len(PostIDs))
	for _, idStr := range PostIDs {
		id, err := gocql.ParseUUID(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("bulk delete aborted - invalid postID format '%s': %w", idStr, err)
		}
		uuids = append(uuids, &id)
	}

	// 2. Chuẩn bị kênh chứa lỗi
	errCh := make(chan taskResult, len(uuids))

	// 3. Phân phối task Xóa vào Worker Pool
	for _, id := range uuids {
		postID := id // Bắt buộc copy biến trong vòng lặp (Go < 1.22)

		err := r.pool.Run(ctx, func() {
			// BẢO VỆ CONTEXT: Không để thao tác bulk delete bị ngắt giữa chừng
			safeCtx := context.WithoutCancel(ctx)
			err := r.doDeleteInsight(safeCtx, *postID)
			if err != nil {
				errCh <- taskResult{message: nil, err: err}
				return
			}
			errCh <- taskResult{message: nil, err: nil}
		})

		// 4. XỬ LÝ LỖI POOL (Context gốc bị cancel trước khi xin được slot)
		if err != nil {
			errCh <- taskResult{message: nil, err: fmt.Errorf("worker pool rejected delete task for post %s: %w", postID.String(), err)}
			break // Cắt đứt vòng lặp để không cố đẩy thêm task
		}
	}

	// 5. ĐỒNG BỘ: Chờ các lệnh xóa đang chạy hoàn tất
	r.pool.Wait()
	close(errCh)

	// 6. Gom lỗi: Trả về lỗi đầu tiên (nếu có)
	for err := range errCh {
		if err.err != nil {
			postid := ""
			if err.message != nil {
				postid = err.message.PostID.String()
			}
			faildocs = append(faildocs, &cassandraErrors.InsightBulkError{
				PostID: postid,
				Error:  err.err.Error(),
			})
			// return 0, nil, err.err // Nếu muốn fail-fast ngay khi gặp lỗi, bỏ comment dòng này và bỏ qua việc gom lỗi vào faildocs
		}
	}
	return int64(len(PostIDs)), faildocs, nil
}
