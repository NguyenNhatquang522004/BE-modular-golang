package req

import "time"

// PageDailyMetricReq đại diện cho payload request từ Client.
// Best Practice: Sử dụng tag binding (như Gin/Validator) để validate dữ liệu ngay từ tầng transport.
type PageDailyMetricReq struct {
	ID               string    `json:"page_id" binding:"required,uuid"`
	MetricDate       time.Time `json:"metric_date" binding:"required"`
	
	ReachTotal       int64     `json:"reach_total" binding:"gte=0"`
	ReachPaid        int64     `json:"reach_paid" binding:"gte=0"`
	ReachOrganic     int64     `json:"reach_organic" binding:"gte=0"`
	
	ImpressionsTotal int64     `json:"impressions_total" binding:"gte=0"`
	NewFollowers     int       `json:"new_followers" binding:"gte=0"`
	Unfollows        int       `json:"unfollows" binding:"gte=0"`
	
	ProfileViews     int       `json:"profile_views" binding:"gte=0"`
	WebsiteClicks    int       `json:"website_clicks" binding:"gte=0"`
	CTAClicks        int       `json:"cta_clicks" binding:"gte=0"`
}