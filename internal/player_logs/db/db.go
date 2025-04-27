package db

import (
	"context"
	"database/sql"
	"go.uber.org/fx"
	"skyhawk/internal/player_logs/config"
)

func NewDB(lc fx.Lifecycle, cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})
	return db, nil
}
