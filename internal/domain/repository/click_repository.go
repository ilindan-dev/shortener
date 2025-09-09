package repository

import (
	"context"
	"github.com/ilindan-dev/shortener/internal/domain/model"
)

// ClickRepository defines the contract for storing click events.
type ClickRepository interface {
	// Create persists a new click event.
	Create(ctx context.Context, click *model.Click) error
}
