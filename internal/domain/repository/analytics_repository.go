package repository

import (
	"context"
	"github.com/ilindan-dev/shortener/internal/domain/model"
)

// AnalyticsRepository defines the contract for retrieving aggregated analytics data.
type AnalyticsRepository interface {
	// GetRawClicks
	GetRawClicks(ctx context.Context, urlID int64) ([]model.Click, error)

	// GetClicksByPeriod
	GetClicksByPeriod(ctx context.Context, urlID int64, period string) ([]model.AggregatedStat, error)

	// GetClicksByUserAgent
	GetClicksByUserAgent(ctx context.Context, urlID int64) ([]model.AggregatedStat, error)

	// GetClicksByPeriodAndUserAgent
	GetClicksByPeriodAndUserAgent(ctx context.Context, urlID int64, period string) ([]model.AggregatedStatDetailed, error)
}
