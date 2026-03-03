package sqlite

import (
	"auth-service/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/fx"
	_ "modernc.org/sqlite"
)

func NewSQLiteDB(cfg *config.Config) (*sql.DB, error) {
	const op = "storage.sqlite.New"
	// dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)", cfg.StoragePath)
	// db, err := sql.Open("sqlite", dsn)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }
	wd, _ := os.Getwd()
	absPath, _ := filepath.Abs(cfg.StoragePath)
	slog.Info("Opening SQLite DB",
		"op", op,
		"wd", wd,
		"storage_path", cfg.StoragePath,
		"abs_path", absPath,
	)

	// гарантируем, что директория под БД есть
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("%s: create dir %s: %w", op, dir, err)
	}

	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)", absPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
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
