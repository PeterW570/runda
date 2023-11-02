package database

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"peterweightman.com/runda/assets"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*sqlx.DB
}

type DbPoolConfig struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
	MaxLifetime  time.Duration
}

func New(dsn string, automigrate bool, cfg DbPoolConfig, logger *slog.Logger) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool.
	// Passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	// Set the maximum number of idle connections in the pool. Passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	// Set the maximum lifetime for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their lifetime.
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	if automigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, dsn)
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}

		logger.Info("database migrations applied")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
