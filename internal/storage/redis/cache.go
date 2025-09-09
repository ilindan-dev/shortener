package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/ilindan-dev/shortener/pkg/keybuilder"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"time"
)

// Ensures that URLCache correctly implements the repo.URLCache interface at compile time.
var _ repo.URLCache = (*URLCache)(nil)

// URLCache implements the domain.repository.URLCache interface using Redis.
type URLCache struct {
	redis  *goredis.Client
	logger zerolog.Logger
}

// NewURLCache creates a new instance of URLCache.
func NewURLCache(logger *zerolog.Logger, redis *goredis.Client) *URLCache {
	return &URLCache{
		redis:  redis,
		logger: logger.With().Str("layer", "redis_cache").Logger(),
	}
}

// Get retrieves a URL from the cache by its short code.
func (c *URLCache) Get(ctx context.Context, shortCode string) (*model.URL, error) {
	key := keybuilder.URLCacheKey(shortCode)
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			c.logger.Info().Str("key", key).Str("cache", "miss").Msg("URL not found in cache")
			return nil, repo.ErrNotFound
		}
		c.logger.Error().Err(err).Str("key", key).Msg("Failed to get key from Redis")
		return nil, err
	}

	var url model.URL
	if err := json.Unmarshal([]byte(val), &url); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("Failed to unmarshal URL from cache")
		return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	c.logger.Info().Str("key", key).Str("cache", "hit").Msg("URL found in cache")
	return &url, nil
}

// Set adds a URL to the cache with a specified expiration time.
func (c *URLCache) Set(ctx context.Context, url *model.URL, expiration time.Duration) error {
	if url.ShortCode == "" {
		return errors.New("cannot cache URL with empty short code")
	}

	key := keybuilder.URLCacheKey(url.ShortCode)
	urlBytes, err := json.Marshal(url)
	if err != nil {
		c.logger.Error().Err(err).Str("short_code", url.ShortCode).Msg("Failed to marshal URL for cache")
		return fmt.Errorf("failed to marshal URL: %w", err)
	}

	if err := c.redis.Set(ctx, key, urlBytes, expiration).Err(); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("Failed to set key in Redis")
		return err
	}

	c.logger.Info().Str("key", key).Msg("URL successfully set in cache")
	return nil
}

// Delete removes a URL from the cache.
func (c *URLCache) Delete(ctx context.Context, shortCode string) error {
	key := keybuilder.URLCacheKey(shortCode)

	result, err := c.redis.Del(ctx, key).Result()
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("Failed to execute delete command on Redis")
		return err
	}

	if result == 0 {
		c.logger.Info().Str("key", key).Msg("Attempted to delete key from cache, but it was not found")
	} else {
		c.logger.Info().Str("key", key).Msg("Successfully deleted key from Redis")
	}

	return nil
}
