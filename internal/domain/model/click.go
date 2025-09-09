package model

import "time"

// Click is the domain model for a single redirect event.
type Click struct {
	ID        int64
	URLID     int64
	UserAgent string
	IPAddress string
	CreatedAt time.Time
}
