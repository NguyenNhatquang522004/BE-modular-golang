package mapper

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/req"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/delivery/dto/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/business/domain/entity"
	"github.com/gocql/gocql"
)

// =========================================================================
// MAPPER REQ -> ENTITY
// =========================================================================

// ToEntityPageDailyMetric khởi tạo một Entity mới từ DTO Request.
func ToEntityPageDailyMetric(r *req.PageDailyMetricReq) (*entity.PageDailyMetric, error) {
	if r == nil {
		return nil, nil
	}

	// Parse chuỗi sang UUID của gocql. Bắt lỗi nếu client truyền chuỗi ID không hợp lệ.
	pageID, err := gocql.ParseUUID(r.ID)
	if err != nil {
		return nil, err
	}

	return &entity.PageDailyMetric{
		ID:               pageID,
		MetricDate:       r.MetricDate,
		ReachTotal:       r.ReachTotal,
		ReachPaid:        r.ReachPaid,
		ReachOrganic:     r.ReachOrganic,
		ImpressionsTotal: r.ImpressionsTotal,
		NewFollowers:     r.NewFollowers,
		Unfollows:        r.Unfollows,
		ProfileViews:     r.ProfileViews,
		WebsiteClicks:    r.WebsiteClicks,
		CTAClicks:        r.CTAClicks,
	}, nil
}

// UpdateToEntityPageDailyMetric ánh xạ toàn bộ dữ liệu từ DTO Request ghi đè vào Entity có sẵn.
// Lưu ý Best Practice cho Cassandra: Do ID và MetricDate là Primary Key,
// việc thay đổi 2 trường này bản chất là hành vi tạo bản ghi mới (Upsert) chứ không phải Update in-place.
func UpdateToEntityPageDailyMetric(r *req.PageDailyMetricReq, e *entity.PageDailyMetric) error {
	if r == nil || e == nil {
		return nil
	}

	pageID, err := gocql.ParseUUID(r.ID)
	if err != nil {
		return err
	}

	e.ID = pageID
	e.MetricDate = r.MetricDate
	e.ReachTotal = r.ReachTotal
	e.ReachPaid = r.ReachPaid
	e.ReachOrganic = r.ReachOrganic
	e.ImpressionsTotal = r.ImpressionsTotal
	e.NewFollowers = r.NewFollowers
	e.Unfollows = r.Unfollows
	e.ProfileViews = r.ProfileViews
	e.WebsiteClicks = r.WebsiteClicks
	e.CTAClicks = r.CTAClicks

	return nil
}

// =========================================================================
// MAPPER -> RES
// =========================================================================
// MAPPER -> RES
// =========================================================================

// ReqToResPageDailyMetric ánh xạ trực tiếp từ DTO Request sang DTO Response theo yêu cầu của bạn.
// Thường dùng khi muốn trả về nguyên vẹn dữ liệu client vừa gửi lên (ví dụ sau khi validate thành công).
func ReqToResPageDailyMetric(r *req.PageDailyMetricReq) *res.PageDailyMetricRes {
	if r == nil {
		return nil
	}

	return &res.PageDailyMetricRes{
		ID:               r.ID,
		MetricDate:       r.MetricDate,
		ReachTotal:       r.ReachTotal,
		ReachPaid:        r.ReachPaid,
		ReachOrganic:     r.ReachOrganic,
		ImpressionsTotal: r.ImpressionsTotal,
		NewFollowers:     r.NewFollowers,
		Unfollows:        r.Unfollows,
		ProfileViews:     r.ProfileViews,
		WebsiteClicks:    r.WebsiteClicks,
		CTAClicks:        r.CTAClicks,
	}
}

// EntityToResPageDailyMetric ánh xạ từ Entity sang DTO Response (Standard Best Practice).
func EntityToResPageDailyMetric(e *entity.PageDailyMetric) *res.PageDailyMetricRes {
	if e == nil {
		return nil
	}

	return &res.PageDailyMetricRes{
		ID:               e.ID.String(),
		MetricDate:       e.MetricDate,
		ReachTotal:       e.ReachTotal,
		ReachPaid:        e.ReachPaid,
		ReachOrganic:     e.ReachOrganic,
		ImpressionsTotal: e.ImpressionsTotal,
		NewFollowers:     e.NewFollowers,
		Unfollows:        e.Unfollows,
		ProfileViews:     e.ProfileViews,
		WebsiteClicks:    e.WebsiteClicks,
		CTAClicks:        e.CTAClicks,
	}
}
