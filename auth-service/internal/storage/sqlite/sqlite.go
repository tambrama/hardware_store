package sqlite

import (
	"auth-service/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"go.uber.org/fx"
	_ "modernc.org/sqlite"
)

func NewSQLiteDB(cfg *config.Config) (*sql.DB, error) {
	const op = "storage.sqlite.New"
	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)", cfg.StoragePath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return db, nil
}

func AddLifecycle(lc fx.Lifecycle, storage *sql.DB, log *slog.Logger) {
	lc.Append(fx.Hook{
		// OnStart: func(ctx context.Context) error {
		// 	log.Info("Running database migrations...")
		// 	if err := goose.SetDialect("sqlite3"); err != nil {
		// 		return err
		// 	}
		// 	if err := goose.Up(storage, "migrations"); err != nil {
		// 		log.Error("migration failed:", slog.Any("error", err))
		// 		return err
		// 	}
		// 	log.Info("Migrations completed successfully")
		// 	return nil
		// },
		OnStop: func(ctx context.Context) error {
			log.Info("Closing SQLite database connection")
			return storage.Close()
		},
	})
}
