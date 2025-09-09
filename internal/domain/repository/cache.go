package repository

import (
	"context"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	"time"
)

// URLCache defines the contract for a caching layer.
type URLCache interface {
	// Get retrieves an item from the cache.
	Get(ctx context.Context, shortCode string) (*model.URL, error)

	// Set adds an item to the cache for a specified duration
	Set(ctx context.Context, url *model.URL, expiration time.Duration) error

	// Delete removes an item from the cache.
	Delete(ctx context.Context, shortCode string) error
}
