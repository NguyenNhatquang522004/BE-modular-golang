package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/errors/cassandraErrors"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/IRepositoryShare"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/dto"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/utils"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/gocql/gocql"
)

type PageDailyMetricsRepository struct {
	// Define the fields for the PageDailyMetricsRepository struct here
	session   *gocql.Session
	redisRepo IRepositoryShare.IRedis
	pool      IRepositoryShare.IWorkerPool
}

// Implement the methods for the PageDailyMetricsRepository struct here
func NewPageDailyMetricsRepository(session *gocql.Session, redisRepo IRepositoryShare.IRedis, pool IRepositoryShare.IWorkerPool) *PageDailyMetricsRepository {
	return &PageDailyMetricsRepository{
		session:   session,
		redisRepo: redisRepo,
		pool:      pool,
	}
}
func (r *PageDailyMetricsRepository) CreatePageDailyMetric(ctx context.Context, metric *entity.PageDailyMetric) error {
	tableName := entity.PageDailyMetric{}.TableName()
	query := fmt.Sprintf(`
		INSERT INTO %s (page_id, metric_date, Reach_Total, Reach_Paid, Reach_Organic,  ImpressionsTotal , new_followers, unfollows, profile_views, website_clicks, cta_clicks)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)
	err := r.session.Query(query,
		metric.ID,
		metric.MetricDate,
		metric.ReachTotal,
		metric.ReachPaid,
		metric.ReachOrganic,
		metric.ImpressionsTotal,
		metric.NewFollowers,
		metric.Unfollows,
		metric.ProfileViews,
		metric.WebsiteClicks,
		metric.CTAClicks,
	).WithContext(ctx).Exec()

	if err != nil {
		return fmt.Errorf("failed to insert page daily metric for page %s on date %s: %w", metric.ID.String(), metric.MetricDate.Format("2006-01-02"), err)
	}
	return nil
}
func (r *PageDailyMetricsRepository) CreateBulkPageDailyMetrics(ctx context.Context, metrics []*entity.PageDailyMetric) (int64, []*cassandraErrors.PageDailyMetricsBulkError, error) {
	tableName := entity.PageDailyMetric{}.TableName()
	if len(metrics) == 0 {
		return 0, nil, nil
	}
	type taskResult struct {
		message *entity.PageDailyMetric
		err     error
	}
	query := fmt.Sprintf(`
		INSERT INTO %s (page_id, metric_date, Reach_Total, Reach_Paid, Reach_Organic,  ImpressionsTotal , new_followers, unfollows, profile_views, website_clicks, cta_clicks)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName)
	resultsCh := make(chan taskResult, len(metrics))
	for i, item := range metrics {
		if item == nil {
			resultsCh <- taskResult{message: nil, err: fmt.Errorf("metric at index %d is nil", i)}
			continue
		}
		if item.ID == (gocql.UUID{}) || item.MetricDate.IsZero() {
			resultsCh <- taskResult{message: nil, err: fmt.Errorf("metric at index %d has invalid ID or MetricDate", i)}
			continue
		}
		payload := item
		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			err := r.session.Query(query,
				payload.ID,
				payload.MetricDate,
				payload.ReachTotal,
				payload.ReachPaid,
				payload.ReachOrganic,
				payload.ImpressionsTotal,
				payload.NewFollowers,
				payload.Unfollows,
				payload.ProfileViews,
				payload.WebsiteClicks,
				payload.CTAClicks,
			).WithContext(safeCtx).Exec()
			resultsCh <- taskResult{message: payload, err: err}
		})
		if err != nil {
			resultsCh <- taskResult{message: item, err: fmt.Errorf("failed to insert page daily metric for page %s on date %s: %w", item.ID.String(), item.MetricDate.Format("2006-01-02"), err)}
		}
	}
	r.pool.Wait()
	close(resultsCh)

	var successCount int64
	var bulkErrors []*cassandraErrors.PageDailyMetricsBulkError
	for result := range resultsCh {
		if result.err != nil {
			PageID := "unknow_page_id"
			var metricDate time.Time
			if result.message != nil {
				PageID = result.message.ID.String()
				metricDate = result.message.MetricDate
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.PageDailyMetricsBulkError{
				PageID:     PageID,
				MetricDate: metricDate,
				Error:      result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	return successCount, bulkErrors, nil
}
func (r *PageDailyMetricsRepository) GetPageDailyMetricsByPageID(ctx context.Context, pageID string, metricDate time.Time, cursor string, limit int) (*dto.PaginationRes, error) {
	// 1. Fail-fast validation
	if pageID == "" {
		return nil, fmt.Errorf("pageID cannot be empty")
	}

	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	// 2. Check Cache (Chỉ hit cache khi lấy trang đầu tiên)
	cacheKeyPrefix := fmt.Sprintf("page_daily_metrics_cache_pageid_%s_date_%s", pageID, metricDate.Format("2006-01-02"))
	if cursor == "" {
		datacache, nextcursor, hasnext, limitcache, err := r.redisRepo.CustomizeGetCache(ctx, []string{
			cacheKeyPrefix,
			cacheKeyPrefix + "_nextcursor",
			cacheKeyPrefix + "_hasnext",
			cacheKeyPrefix + "_limit",
		})

		if err == nil && datacache != nil {
			if msgList, ok := datacache.([]*entity.PageDailyMetric); ok {
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
	tableName := entity.PageDailyMetric{}.TableName()

	// Khai báo rõ các cột để tối ưu Serialize/Deserialize
	query := fmt.Sprintf(`
		SELECT page_id, metric_date, reach_total, reach_paid, reach_organic,
		       impressions_total, new_followers, unfollows, profile_views,
		       website_clicks, cta_clicks
		FROM %s WHERE page_id = ? AND metric_date = ?
	`, tableName)

	// LƯU Ý: Không dùng context.WithoutCancel ở đây.
	// Truyền ctx gốc để DB ngừng xử lý nếu client ngắt kết nối.
	finalid, err := gocql.ParseUUID(pageID)
	if err != nil {
		return nil, fmt.Errorf("invalid pageID: %w", err)
	}
	q := r.session.Query(query, finalid, metricDate).WithContext(ctx).PageSize(limit)
	if len(pageState) > 0 {
		q = q.PageState(pageState)
	}

	// 5. Thực thi và Map dữ liệu
	iter := q.Iter()
	scanner := iter.Scanner()
	var messages []*entity.PageDailyMetric

	for scanner.Next() {
		var m entity.PageDailyMetric
		// Gocql tự động map mảng LIST<TEXT> vào []string (attachments)
		// và xử lý nullable UUID (*gocql.UUID) thành nil khi NULL trong DB
		err := scanner.Scan(
			&m.ID,
			&m.MetricDate,
			&m.ReachTotal,
			&m.ReachPaid,
			&m.ReachOrganic,
			&m.ImpressionsTotal,
			&m.NewFollowers,
			&m.Unfollows,
			&m.ProfileViews,
			&m.WebsiteClicks,
			&m.CTAClicks,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page daily metric for page %s on date %s: %w", pageID, metricDate, err)
		}

		// Append con trỏ an toàn (Go >= 1.22)
		messages = append(messages, &m)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("database error during pagination for page %s on date %s: %w", pageID, metricDate, err)
	}

	// 6. Xử lý Next Cursor và HasNext
	nextPageState := iter.PageState()
	nextCursorStr := utils.EncodeCursorCassandra(nextPageState)
	hasNext := len(nextPageState) > 0

	// Trả về slice rỗng [] thay vì nil để Frontend không bị lỗi map null
	if messages == nil {
		messages = []*entity.PageDailyMetric{}
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
			fmt.Printf("failed to set cache for page daily metrics pagination for page %s on date %s: %v\n", pageID, metricDate, err)
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
func (r *PageDailyMetricsRepository) UpdatePageDailyMetric(ctx context.Context, metric *entity.PageDailyMetric) error {
	tablename := entity.PageDailyMetric{}.TableName()
	query := fmt.Sprintf(`
		UPDATE %s SET Reach_Total = ?, Reach_Paid = ?, Reach_Organic = ?, ImpressionsTotal = ?, new_followers = ?, unfollows = ?, profile_views = ?, website_clicks = ?, cta_clicks = ?
		WHERE page_id = ? AND metric_date = ?
	`, tablename)
	err := r.session.Query(query,
		metric.ReachTotal,
		metric.ReachPaid,
		metric.ReachOrganic,
		metric.ImpressionsTotal,
		metric.NewFollowers,
		metric.Unfollows,
		metric.ProfileViews,
		metric.WebsiteClicks,
		metric.CTAClicks,
		metric.ID,
		metric.MetricDate,
	).WithContext(ctx).Exec()

	if err != nil {
		return fmt.Errorf("failed to update page daily metric for page %s on date %s: %w", metric.ID.String(), metric.MetricDate.Format("2006-01-02"), err)
	}
	return nil
}
func (r *PageDailyMetricsRepository) UpdateBulkPageDailyMetrics(ctx context.Context, metrics []*entity.PageDailyMetric) (int64, []*cassandraErrors.PageDailyMetricsBulkError, error) {
	tableName := entity.PageDailyMetric{}.TableName()
	if len(metrics) == 0 {
		return 0, nil, nil
	}
	type taskResult struct {
		message *entity.PageDailyMetric
		err     error
	}
	query := fmt.Sprintf(`
		UPDATE %s SET Reach_Total = ?, Reach_Paid = ?, Reach_Organic = ?, ImpressionsTotal = ?, new_followers = ?, unfollows = ?, profile_views = ?, website_clicks = ?, cta_clicks = ?
		WHERE page_id = ? AND metric_date = ?
	`, tableName)
	resultsCh := make(chan taskResult, len(metrics))
	for i, item := range metrics {
		if item == nil {
			resultsCh <- taskResult{message: nil, err: fmt.Errorf("metric at index %d is nil", i)}
			continue
		}
		if item.ID == (gocql.UUID{}) || item.MetricDate.IsZero() {
			resultsCh <- taskResult{message: nil, err: fmt.Errorf("metric at index %d has invalid ID or MetricDate", i)}
			continue
		}
		payload := item
		err := r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			err := r.session.Query(query,
				payload.ReachTotal,
				payload.ReachPaid,
				payload.ReachOrganic,
				payload.ImpressionsTotal,
				payload.NewFollowers,
				payload.Unfollows,
				payload.ProfileViews,
				payload.WebsiteClicks,
				payload.CTAClicks,
				payload.ID,
				payload.MetricDate,
			).WithContext(safeCtx).Exec()
			resultsCh <- taskResult{message: payload, err: err}
		})
		if err != nil {
			resultsCh <- taskResult{message: item, err: fmt.Errorf("failed to update page daily metric for page %s on date %s: %w", item.ID.String(), item.MetricDate.Format("2006-01-02"), err)}
		}
	}
	r.pool.Wait()
	close(resultsCh)

	var successCount int64
	var bulkErrors []*cassandraErrors.PageDailyMetricsBulkError
	for result := range resultsCh {
		if result.err != nil {
			PageID := "unknow_page_id"
			var metricDate time.Time
			if result.message != nil {
				PageID = result.message.ID.String()
				metricDate = result.message.MetricDate
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.PageDailyMetricsBulkError{
				PageID:     PageID,
				MetricDate: metricDate,
				Error:      result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	return successCount, bulkErrors, nil
}
func (r *PageDailyMetricsRepository) DeletePageDailyMetric(ctx context.Context, pageID string, metricDate time.Time) error {
	tableName := entity.PageDailyMetric{}.TableName()
	query := fmt.Sprintf(`DELETE FROM %s WHERE page_id = ? AND metric_date = ?`, tableName)
	finalid, err := gocql.ParseUUID(pageID)
	if err != nil {
		return fmt.Errorf("invalid pageID: %w", err)
	}
	err = r.session.Query(query, finalid, metricDate).WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("failed to delete page daily metric for page %s on date %s: %w", pageID, metricDate, err)
	}
	return nil
}
func (r *PageDailyMetricsRepository) DeleteBulkPageDailyMetrics(ctx context.Context, pageID string, metricDates []time.Time) (int64, []*cassandraErrors.PageDailyMetricsBulkError, error) {
	tableName := entity.PageDailyMetric{}.TableName()
	if len(metricDates) == 0 {
		return 0, nil, nil
	}
	type taskResult struct {
		message *entity.PageDailyMetric
		err     error
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE page_id = ? AND metric_date = ?`, tableName)
	resultsCh := make(chan taskResult, len(metricDates))
	for i, metricDate := range metricDates {
		if metricDate.IsZero() {
			resultsCh <- taskResult{message: nil, err: fmt.Errorf("metricDate at index %d is invalid", i)}
			continue
		}
		finalid, err := gocql.ParseUUID(pageID)
		if err != nil {
			resultsCh <- taskResult{message: nil, err: fmt.Errorf("invalid pageID: %w", err)}
			continue
		}
		err = r.pool.Run(ctx, func() {
			safeCtx := context.WithoutCancel(ctx)
			err := r.session.Query(query, finalid, metricDate).WithContext(safeCtx).Exec()
			resultsCh <- taskResult{message: &entity.PageDailyMetric{ID: finalid, MetricDate: metricDate}, err: err}
		})
		if err != nil {
			resultsCh <- taskResult{message: &entity.PageDailyMetric{ID: finalid, MetricDate: metricDate}, err: fmt.Errorf("failed to delete page daily metric for page %s on date %s: %w", pageID, metricDate, err)}
		}
	}
	r.pool.Wait()
	close(resultsCh)

	var successCount int64
	var bulkErrors []*cassandraErrors.PageDailyMetricsBulkError
	for result := range resultsCh {
		if result.err != nil {
			PageID := "unknow_page_id"
			var metricDate time.Time
			if result.message != nil {
				PageID = result.message.ID.String()
				metricDate = result.message.MetricDate
			}
			bulkErrors = append(bulkErrors, &cassandraErrors.PageDailyMetricsBulkError{
				PageID:     PageID,
				MetricDate: metricDate,
				Error:      result.err.Error(),
			})
		} else {
			successCount++
		}
	}
	return successCount, bulkErrors, nil
}
func (r *PageDailyMetricsRepository) DeletePageDailyMetricsByPageID(ctx context.Context, pageID string) error {
	tableName := entity.PageDailyMetric{}.TableName()
	query := fmt.Sprintf(`DELETE FROM %s WHERE page_id = ?`, tableName)
	finalid, err := gocql.ParseUUID(pageID)
	if err != nil {
		return fmt.Errorf("invalid pageID: %w", err)
	}
	err = r.session.Query(query, finalid).WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("failed to delete page daily metrics for page %s: %w", pageID, err)
	}
	return nil
}
