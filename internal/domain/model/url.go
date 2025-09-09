package model

import "time"

// URL is the domain model for a shortened link.
type URL struct {
	ID          int64
	OriginalURL string
	ShortCode   string
	CreatedAt   time.Time
}
