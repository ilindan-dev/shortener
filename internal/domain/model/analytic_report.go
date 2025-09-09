package model

type FullAnalyticsReport struct {
	URL               URL
	TotalClicks       int64
	ClicksByDay       []AggregatedStat
	ClicksByUserAgent []AggregatedStat
	RecentClicks      []Click
}
