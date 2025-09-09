package service

import (
	"context"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/ilindan-dev/shortener/pkg/base62"
	"github.com/rs/zerolog"
	"time"
)

// URLService encapsulates the business logic for URL shortening and analytics.
type URLService struct {
	urlRepo   repo.URLRepository
	clickRepo repo.ClickRepository
	cache     repo.URLCache
	logger    zerolog.Logger
}

// NewURLService creates a new instance of URLService.
func NewURLService(
	urlRepo repo.URLRepository,
	clickRepo repo.ClickRepository,
	cache repo.URLCache,
	logger *zerolog.Logger,
) *URLService {
	return &URLService{
		urlRepo:   urlRepo,
		clickRepo: clickRepo,
		cache:     cache,
		logger:    logger.With().Str("layer", "service").Logger(),
	}
}

// CreateShortURL orchestrates the entire process of creating a short URL.
func (s *URLService) CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error) {
	s.logger.Info().Str("original_url", originalURL).Msg("Creating new short URL")

	url, err := s.urlRepo.Create(ctx, originalURL)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create initial URL record")
		return nil, err
	}

	shortCode := base62.Encode(url.ID)
	url.ShortCode = shortCode

	if err := s.urlRepo.UpdateShortCode(ctx, url.ID, shortCode); err != nil {
		s.logger.Error().Err(err).Int64("url_id", url.ID).Msg("Failed to update URL with short code")
		// TODO: cleanup/retry mechanism here.
		return nil, err
	}

	if err := s.cache.Set(ctx, url, time.Hour*24*7); err != nil {
		s.logger.Error().Err(err).Int64("url_id", url.ID).Msg("Failed to warm up cache")
	}

	s.logger.Info().Str("short_code", shortCode).Int64("url_id", url.ID).Msg("Successfully created short URL")
	return url, nil
}

// ProcessRedirect finds the original URL for a given short code and records the click for analytics.
func (s *URLService) ProcessRedirect(ctx context.Context, shortCode, userAgent, ipAddress string) (*model.URL, error) {
	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	go func() {
		click := &model.Click{
			URLID:     url.ID,
			UserAgent: userAgent,
			IPAddress: ipAddress,
		}
		if err := s.clickRepo.Create(context.Background(), click); err != nil {
			s.logger.Error().Err(err).Int64("url_id", url.ID).Msg("Failed to record click")
		}
	}()

	return url, nil
}
