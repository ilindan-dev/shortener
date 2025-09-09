package service

import (
	"context"
	"fmt"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// AnalyticsService provides business logic for URL analytics.
type AnalyticsService struct {
	urlRepo       repo.URLRepository
	analyticsRepo repo.AnalyticsRepository
	logger        zerolog.Logger
}

// NewAnalyticsService creates a new instance of AnalyticsService.
func NewAnalyticsService(
	urlRepo repo.URLRepository,
	analyticsRepo repo.AnalyticsRepository,
	logger *zerolog.Logger,
) *AnalyticsService {
	return &AnalyticsService{
		urlRepo:       urlRepo,
		analyticsRepo: analyticsRepo,
		logger:        logger.With().Str("layer", "analytics_service").Logger(),
	}
}

// GetFullAnalyticsReport fetches and aggregates all analytics data for a given short code.
func (s *AnalyticsService) GetFullAnalyticsReport(ctx context.Context, shortCode string) (*model.FullAnalyticsReport, error) {
	s.logger.Info().Str("short_code", shortCode).Msg("Fetching full analytics report")

	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	report := &model.FullAnalyticsReport{
		URL: *url,
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		stats, err := s.analyticsRepo.GetClicksByPeriod(gCtx, url.ID, "day")
		if err != nil {
			s.logger.Error().Err(err).Int64("url_id", url.ID).Msg("Failed to fetch clicks by day")
			return fmt.Errorf("could not fetch daily stats: %w", err)
		}
		report.ClicksByDay = stats
		return nil
	})

	g.Go(func() error {
		stats, err := s.analyticsRepo.GetClicksByUserAgent(gCtx, url.ID)
		if err != nil {
			s.logger.Error().Err(err).Int64("url_id", url.ID).Msg("Failed to fetch clicks by user agent")
			return fmt.Errorf("could not fetch user agent stats: %w", err)
		}
		report.ClicksByUserAgent = stats
		return nil
	})

	g.Go(func() error {
		clicks, err := s.analyticsRepo.GetRawClicks(gCtx, url.ID)
		if err != nil {
			s.logger.Error().Err(err).Int64("url_id", url.ID).Msg("Failed to fetch raw clicks")
			return fmt.Errorf("could not fetch raw clicks: %w", err)
		}
		report.RecentClicks = clicks
		report.TotalClicks = int64(len(clicks))
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	s.logger.Info().Str("short_code", shortCode).Msg("Successfully fetched analytics report")
	return report, nil
}
