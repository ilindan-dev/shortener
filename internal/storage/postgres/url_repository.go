package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/ilindan-dev/shortener/internal/storage/postgres/db"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Ensures that URLRepository correctly implements the repo.URLRepository interface at compile time.
var _ repo.URLRepository = (*URLRepository)(nil)

// URLRepository implements the domain.repository.URLRepository interface
// using PostgreSQL as a backend.
type URLRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	logger  zerolog.Logger
}

// NewURLRepository creates a new instance of URLRepository.
func NewURLRepository(pool *pgxpool.Pool, logger *zerolog.Logger) *URLRepository {
	return &URLRepository{
		pool:    pool,
		queries: db.New(pool),
		logger:  logger.With().Str("layer", "postgres_repository").Logger(),
	}
}

// Create persists a new URL record in the database.
func (r *URLRepository) Create(ctx context.Context, originalURL string) (*model.URL, error) {
	createdDB, err := r.queries.CreateURL(ctx, originalURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			r.logger.Warn().Err(err).Str("url", originalURL).Msg("Failed to create URL due to duplicate")
			return nil, repo.ErrDuplicateRecord
		}
		r.logger.Error().Err(err).Str("url", originalURL).Msg("Failed to create URL")
		return nil, fmt.Errorf("postgres: CreateURL failed: %w", err)
	}

	return toDomainURL(createdDB), nil
}

// UpdateShortCode updates an existing URL record with its generated short code.
func (r *URLRepository) UpdateShortCode(ctx context.Context, id int64, shortCode string) error {
	params := db.UpdateURLShortCodeParams{
		ID:        id,
		ShortCode: pgtype.Text{String: shortCode, Valid: true},
	}

	err := r.queries.UpdateURLShortCode(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int64("id", id).Msg("Failed to update URL short code")
		return fmt.Errorf("postgres: UpdateURLShortCode failed: %w", err)
	}

	return nil
}

// GetByShortCode retrieves a single URL from the database by its unique short code.
func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	dbURL, err := r.queries.GetURLByShortCode(ctx, pgtype.Text{String: shortCode, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn().Str("short_code", shortCode).Msg("URL not found by short code")
			return nil, repo.ErrNotFound
		}
		r.logger.Error().Err(err).Str("short_code", shortCode).Msg("Failed to get URL by short code")
		return nil, fmt.Errorf("postgres: GetURLByShortCode failed: %w", err)
	}

	return toDomainURL(dbURL), nil
}

// toDomainURL converts a database model (from sqlc) to a domain model.
func toDomainURL(dbURL db.Url) *model.URL {
	domainModel := &model.URL{
		ID:          dbURL.ID,
		OriginalURL: dbURL.OriginalUrl,
		CreatedAt:   dbURL.CreatedAt.Time,
	}

	if dbURL.ShortCode.Valid {
		domainModel.ShortCode = dbURL.ShortCode.String
	}

	return domainModel
}
