package model

// AnalyticsView combines a URL with its associated click data.
// This is the model our service layer will work with.
type AnalyticsView struct {
	URL    URL
	Clicks []Click
}
