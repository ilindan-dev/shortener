package postgres

import (
	"context"
	"fmt"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/ilindan-dev/shortener/internal/storage/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Ensures AnalyticsRepository implements the interface.
var _ repo.AnalyticsRepository = (*AnalyticsRepository)(nil)

// AnalyticsRepository implements the domain.repository.AnalyticsRepository interface.
type AnalyticsRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	logger  zerolog.Logger
}

// NewAnalyticsRepository creates a new instance of AnalyticsRepository.
func NewAnalyticsRepository(pool *pgxpool.Pool, logger *zerolog.Logger) *AnalyticsRepository {
	return &AnalyticsRepository{
		pool:    pool,
		queries: db.New(pool),
		logger:  logger.With().Str("layer", "postgres_analytics_repository").Logger(),
	}
}

// GetRawClicks fetches all raw click events for a given URL ID.
func (r *AnalyticsRepository) GetRawClicks(ctx context.Context, urlID int64) ([]model.Click, error) {
	dbClicks, err := r.queries.GetClicksByURLID(ctx, urlID)
	if err != nil {
		r.logger.Error().Err(err).Int64("url_id", urlID).Msg("Failed to get raw clicks")
		return nil, fmt.Errorf("postgres: GetClicksByURLID failed: %w", err)
	}
	return toDomainClicks(dbClicks), nil
}

// GetClicksByPeriod fetches click counts aggregated by a time period.
func (r *AnalyticsRepository) GetClicksByPeriod(ctx context.Context, urlID int64, period string) ([]model.AggregatedStat, error) {
	params := db.GetClicksByPeriodParams{
		UrlID:  urlID,
		Period: period,
	}
	rows, err := r.queries.GetClicksByPeriod(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int64("url_id", urlID).Str("period", period).Msg("Failed to get clicks by period")
		return nil, fmt.Errorf("postgres: GetClicksByPeriod failed: %w", err)
	}

	return toAggregatedStatsFromTime(rows), nil
}

// GetClicksByUserAgent fetches click counts aggregated by user agent.
func (r *AnalyticsRepository) GetClicksByUserAgent(ctx context.Context, urlID int64) ([]model.AggregatedStat, error) {
	rows, err := r.queries.GetClicksByUserAgent(ctx, urlID)
	if err != nil {
		r.logger.Error().Err(err).Int64("url_id", urlID).Msg("Failed to get clicks by user agent")
		return nil, fmt.Errorf("postgres: GetClicksByUserAgent failed: %w", err)
	}

	return toAggregatedStatsFromString(rows), nil
}

// GetClicksByPeriodAndUserAgent fetches click counts aggregated by both time period and user agent.
func (r *AnalyticsRepository) GetClicksByPeriodAndUserAgent(ctx context.Context, urlID int64, period string) ([]model.AggregatedStatDetailed, error) {
	params := db.GetClicksByPeriodAndUserAgentParams{
		Period: period,
		UrlID:  urlID,
	}
	rows, err := r.queries.GetClicksByPeriodAndUserAgent(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int64("url_id", urlID).Str("period", period).Msg("Failed to get detailed analytics")
		return nil, fmt.Errorf("postgres: GetClicksByPeriodAndUserAgent failed: %w", err)
	}

	return toAggregatedStatsDetailed(rows), nil
}

// --- Mapper Functions ---

func toAggregatedStatsFromTime(rows []db.GetClicksByPeriodRow) []model.AggregatedStat {
	stats := make([]model.AggregatedStat, len(rows))
	for i, row := range rows {
		stats[i] = model.AggregatedStat{
			Key:   row.Key.Time.Format("2006-01-02"),
			Value: row.Value,
		}
	}
	return stats
}

func toAggregatedStatsFromString(rows []db.GetClicksByUserAgentRow) []model.AggregatedStat {
	stats := make([]model.AggregatedStat, len(rows))
	for i, row := range rows {
		stats[i] = model.AggregatedStat{
			Key:   row.Key,
			Value: row.Value,
		}
	}
	return stats
}

func toAggregatedStatsDetailed(rows []db.GetClicksByPeriodAndUserAgentRow) []model.AggregatedStatDetailed {
	stats := make([]model.AggregatedStatDetailed, len(rows))
	for i, row := range rows {
		stats[i] = model.AggregatedStatDetailed{
			TimeKey: row.TimeKey.Time.Format("2006-01-02"),
			UAKey:   row.UaKey,
			Value:   row.Value,
		}
	}
	return stats
}

func toDomainClicks(dbClicks []db.Click) []model.Click {
	if len(dbClicks) == 0 {
		return []model.Click{}
	}

	clicks := make([]model.Click, 0, len(dbClicks))
	for _, dbClick := range dbClicks {
		clicks = append(clicks, toDomainClick(dbClick))
	}
	return clicks
}

func toDomainClick(dbClick db.Click) model.Click {
	click := model.Click{
		ID:        dbClick.ID,
		URLID:     dbClick.UrlID,
		CreatedAt: dbClick.CreatedAt.Time,
	}

	if dbClick.UserAgent.Valid {
		click.UserAgent = dbClick.UserAgent.String
	}

	if dbClick.IpAddress.IsValid() {
		click.IPAddress = dbClick.IpAddress.String()
	}

	return click
}
