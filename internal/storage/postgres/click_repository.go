package postgres

import (
	"context"
	"fmt"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/ilindan-dev/shortener/internal/storage/postgres/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"net/netip"
)

// Ensures that ClickRepository correctly implements the repo.ClickRepository interface at compile time.
var _ repo.ClickRepository = (*ClickRepository)(nil)

// ClickRepository implements the domain.repository.ClickRepository interface
// using PostgreSQL as a backend.
type ClickRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	logger  zerolog.Logger
}

// NewClickRepository creates a new instance of ClickRepository.
func NewClickRepository(pool *pgxpool.Pool, logger *zerolog.Logger) *ClickRepository {
	return &ClickRepository{
		pool:    pool,
		queries: db.New(pool),
		logger:  logger.With().Str("layer", "postgres_click_repository").Logger(),
	}
}

// Create persists a new click event in the database.
func (r *ClickRepository) Create(ctx context.Context, click *model.Click) error {
	params, err := toDBCreateClickParams(click)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to map domain click to db params")
		return err
	}

	err = r.queries.CreateClick(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int64("url_id", click.URLID).Msg("Failed to create click")
		return fmt.Errorf("postgres: CreateClick failed: %w", err)
	}

	return nil
}

// toDBCreateClickParams converts a domain model.Click to the sqlc-generated parameters for creation.
func toDBCreateClickParams(click *model.Click) (db.CreateClickParams, error) {
	params := db.CreateClickParams{
		UrlID: click.URLID,
	}

	if click.UserAgent != "" {
		params.UserAgent = pgtype.Text{String: click.UserAgent, Valid: true}
	}

	if click.IPAddress != "" {
		addr, err := netip.ParseAddr(click.IPAddress)
		if err != nil {
			return db.CreateClickParams{}, err
		}
		params.IpAddress = &addr
	}

	return params, nil
}
