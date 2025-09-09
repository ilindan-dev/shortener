package postgres

import (
	"context"
	"errors"
	"github.com/ilindan-dev/shortener/internal/domain/model"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/ilindan-dev/shortener/internal/storage/postgres/db"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"strconv"
	"time"
)

// Ensure NotificationRepository implements the interface
var _ repo.URLRepository = (*URLRepository)(nil)

// URLRepository implements the domain.repository.URLRepository interface using PostgreSQL as a backend.
type URLRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	logger  zerolog.Logger
}

func NewURLRepository(pool *pgxpool.Pool, logger *zerolog.Logger) *URLRepository {
	return &URLRepository{
		pool:    pool,
		queries: db.New(pool),
		logger:  logger.With().Str("layer", "postgres_repository").Logger(),
	}
}

func (r *URLRepository) Create(ctx context.Context, originalURL string) (*model.URL, error) {
	createdDB, err := r.queries.CreateURL(ctx, originalURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			r.logger.Warn().Err(err).Str("method", "URLRepository.Create").Str("url", originalURL).Msg("Failed to create URL")
			return nil, repo.ErrDuplicateRecord
		}
		r.logger.Error().Err(err).Str("method", "URLRepository.Create").Str("url", originalURL).Msg("Failed to create URL")
	}

	return toDomainModel(createdDB)
}

func (r *URLRepository) UpdateShortURL(ctx context.Context, id int64, shortCode string) error {
	params, err := toUpdateParams(id, shortCode)
	if err != nil {
		r.logger.Err(err).Str("method", "URLRepository.UpdateShortURL").Str("id", strconv.FormatInt(id, 10)).Msg("Failed to convert into update params")
		return err
	}

	err = r.queries.UpdateURLShortCode(ctx, params)
	if err != nil {
		r.logger.Err(err).Str("method", "URLRepository.UpdateShortURL").Str("id", strconv.FormatInt(id, 10)).Msg("Failed to update URL short code")
		return err
	}

	return nil
}

func (r *URLRepository) GetByShortURL(ctx context.Context, shortCode string) (*model.URL, error) {

}

func (r *URLRepository) GetAnalyticsByShortURL(ctx context.Context, shortCode string) (*model.AnalyticsView, error) {

}

func toDomainModel(dbURl db.Url) (*model.URL, error) {
	domainModel := &model.URL{
		ID:          dbURl.ID,
		OriginalURL: dbURl.OriginalUrl,
		CreatedAt:   dbURl.CreatedAt.Time,
	}

	if dbURl.ShortCode.Valid {
		domainModel.ShortCode = dbURl.ShortCode.String
	}

	return domainModel, nil
}

func toUpdateParams(id int64, shortCode string) (db.UpdateURLShortCodeParams, error) {
	return db.UpdateURLShortCodeParams{
		ID: id,
		ShortCode: pgtype.Text{
			String: shortCode,
			Valid:  true,
		},
	}, nil
}