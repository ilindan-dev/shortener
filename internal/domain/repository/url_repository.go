package repository

import (
	"context"
	"github.com/ilindan-dev/shortener/internal/domain/model"
)

// URLRepository defines the contract for URL persistence.
type URLRepository interface {
	// Create persists a new URL with just the original URL and returns the created record.
	Create(ctx context.Context, originalURL string) (*model.URL, error)

	// UpdateShortCode updates an existing URL record with its generated short URL.
	UpdateShortCode(ctx context.Context, id int64, shortCode string) error

	// GetByShortCode retrieves a URL by its unique shortened URL string.
	GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
}
