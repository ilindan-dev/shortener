package http

import "time"

// CreateURLRequest defines the structure for a new URL shortening request.
type CreateURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// URLResponse defines the structure for a successful URL creation response.
type URLResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// ClickDTO defines a simplified view of a click for the analytics response.
type ClickDTO struct {
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"user_agent"`
}

// AnalyticsResponse defines the structure for the full analytics report.
type AnalyticsResponse struct {
	OriginalURL       string     `json:"original_url"`
	ShortURL          string     `json:"short_url"`
	TotalClicks       int64      `json:"total_clicks"`
	ClicksByDay       []StatItem `json:"clicks_by_day"`
	ClicksByUserAgent []StatItem `json:"clicks_by_user_agent"`
	RecentClicks      []ClickDTO `json:"recent_clicks"`
}

// StatItem is a generic structure for aggregated data.
type StatItem struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

// ErrorResponse defines a standard structure for API error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}
