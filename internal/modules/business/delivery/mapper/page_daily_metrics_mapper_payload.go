package mapper

import (
	"fmt"
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/shared/events/businessEvent"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/gocql/gocql"
)

func ToEntityPageDailyMetricPayload(p *businessEvent.PageDailyMetricPayload) (*entity.PageDailyMetric, error) {
	// ==========================================
	// 0. Bắt lỗi nil pointer (Best Practice bảo vệ hàm)
	// ==========================================
	if p == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}

	// ==========================================
	// 1. Xử lý UUID an toàn
	// ==========================================
	pageID, err := gocql.ParseUUID(p.PageID)
	if err != nil {
		return nil, fmt.Errorf("invalid page_id format: %w", err)
	}

	// ==========================================
	// 2. Chuẩn hóa Thời gian (Time Normalization)
	// Đảm bảo đưa về đúng 00:00:00 UTC như đã thảo luận
	// ==========================================
	normalizedDate := time.Date(
		p.MetricDate.Year(),
		p.MetricDate.Month(),
		p.MetricDate.Day(),
		0, 0, 0, 0, time.UTC,
	)

	// ==========================================
	// 3. Gán dữ liệu an toàn từ Pointer sang Value
	// ==========================================
	ent := &entity.PageDailyMetric{
		ID:         pageID,
		MetricDate: normalizedDate,

		// Dùng hàm getVal để tự động lấy giá trị (nếu có) hoặc gán bằng 0 (nếu nil)
		ReachTotal:   getVal(p.ReachTotal),
		ReachPaid:    getVal(p.ReachPaid),
		ReachOrganic: getVal(p.ReachOrganic),

		ImpressionsTotal: getVal(p.ImpressionsTotal),
		NewFollowers:     getVal(p.NewFollowers),
		Unfollows:        getVal(p.Unfollows),

		ProfileViews:  getVal(p.ProfileViews),
		WebsiteClicks: getVal(p.WebsiteClicks),
		CTAClicks:     getVal(p.CTAClicks),
	}

	return ent, nil
}

// ==========================================
// HELPER FUNCTION (Best Practice cho Go 1.18+)
// ==========================================

// getVal là hàm generic giúp đọc giá trị từ con trỏ một cách an toàn nhất.
func getVal[T any](ptr *T) T {
	if ptr == nil {
		var zero T // Sẽ tự động là 0 với int/int64, "" với string,...
		return zero
	}
	return *ptr
}
