package repository

import (
	"context"
	"log/slog"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CodedLinksDbRepository struct {
	dbPool *pgxpool.Pool
	log    *logger.Container
}

func NewCodedLinksDbRepository(cfg *config.Config, log *logger.Container) (*CodedLinksDbRepository, error) {
	pool, err := pgxpool.New(context.Background(), *cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	repository := &CodedLinksDbRepository{
		dbPool: pool,
		log:    log,
	}

	return repository, nil
}

func (r *CodedLinksDbRepository) Ping(ctx context.Context) error {
	if err := r.dbPool.Ping(ctx); err != nil {
		r.log.Logger.ErrorContext(ctx, "Ошибка во время пинга БД", slog.Any("error", err))
		return err
	}

	return nil
}
