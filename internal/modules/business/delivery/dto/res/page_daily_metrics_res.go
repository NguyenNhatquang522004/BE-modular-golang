package res

import "time"

// PageDailyMetricRes đại diện cho payload response trả về cho Client.
type PageDailyMetricRes struct {
	ID               string    `json:"page_id"`
	MetricDate       time.Time `json:"metric_date"`
	
	ReachTotal       int64     `json:"reach_total"`
	ReachPaid        int64     `json:"reach_paid"`
	ReachOrganic     int64     `json:"reach_organic"`
	
	ImpressionsTotal int64     `json:"impressions_total"`
	NewFollowers     int       `json:"new_followers"`
	Unfollows        int       `json:"unfollows"`
	
	ProfileViews     int       `json:"profile_views"`
	WebsiteClicks    int       `json:"website_clicks"`
	CTAClicks        int       `json:"cta_clicks"`
}