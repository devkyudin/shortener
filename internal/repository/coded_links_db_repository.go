package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/devkyudin/shortener/internal/config"
	"github.com/devkyudin/shortener/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

type CodedLinksDbRepository struct {
	db  *sql.DB
	log *logger.Container
}

func NewCodedLinksDbRepository(cfg *config.Config, log *logger.Container) (*CodedLinksDbRepository, error) {
	db, err := sql.Open("sqlite3", *cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	repository := &CodedLinksDbRepository{
		db:  db,
		log: log,
	}

	return repository, nil
}

func (r *CodedLinksDbRepository) Close() error {
	return r.db.Close()
}

func (r *CodedLinksDbRepository) Ping(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		r.log.Logger.ErrorContext(ctx, "Ошибка во время пинга БД", slog.Any("error", err))
		return err
	}

	return nil
}
