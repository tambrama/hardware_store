package postgres

import (
	"context"
	"hardware_store/internal/config"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewDB(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Не удалось подключиться к бд", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func AddDBLifecycle(lc fx.Lifecycle, pool *pgxpool.Pool) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})
}