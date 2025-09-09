package redis

import (
	"context"
	"errors"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/rs/zerolog"
	"time"
)

// Ensures that CachedURLRepository correctly implements the repo.URLRepository interface at compile time.
var _ repo.URLRepository = (*CachedURLRepository)(nil)

// CachedURLRepository is a decorator for a URLRepository that adds a caching layer.
type CachedURLRepository struct {
	primaryRepo repo.URLRepository
	cache       repo.URLCache
	logger      zerolog.Logger
	ttl         time.Duration
}

// NewCachedURLRepository creates a new instance of the cached repository decorator.
func NewCachedURLRepository(
	primaryRepo repo.URLRepository,
	cache repo.URLCache,
	logger *zerolog.Logger,
) *CachedURLRepository {
	return &CachedURLRepository{
		primaryRepo: primaryRepo,
		cache:       cache,
		logger:      logger.With().Str("layer", "cached_repository").Logger(),
		ttl:         time.Hour * 24 * 7,
	}
}

// Create first persists the URL in the primary repository, then warms up the cache.
func (r *CachedURLRepository) Create(ctx context.Context, originalURL string) (*model.URL, error) {
	return r.primaryRepo.Create(ctx, originalURL)
}

// UpdateShortCode updates the primary repository and then warms up the cache.
func (r *CachedURLRepository) UpdateShortCode(ctx context.Context, id int64, shortCode string) error {
	if err := r.primaryRepo.UpdateShortCode(ctx, id, shortCode); err != nil {
		return err
	}
	return nil
}

// GetByShortCode implements the cache-aside pattern.
func (r *CachedURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	cachedURL, err := r.cache.Get(ctx, shortCode)
	if err == nil {
		r.logger.Info().Str("short_code", shortCode).Msg("Cache hit")
		return cachedURL, nil
	}

	if !errors.Is(err, repo.ErrNotFound) {
		r.logger.Error().Err(err).Str("short_code", shortCode).Msg("Cache get error, falling back to primary repository")
	} else {
		r.logger.Info().Str("short_code", shortCode).Msg("Cache miss")
	}

	dbURL, err := r.primaryRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Set(ctx, dbURL, r.ttl); err != nil {
		r.logger.Error().Err(err).Str("short_code", dbURL.ShortCode).Msg("Failed to set cache after DB fetch")
	}

	return dbURL, nil
}
