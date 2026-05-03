package storage

import (
	"context"
	"generator-service/internal/models"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Storage struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewDB(dbName string, lc fx.Lifecycle, logger *slog.Logger) *pgxpool.Pool {
	logger.Info("connecting to database", slog.String("url", dbName))
	pool, err := pgxpool.New(context.Background(), dbName)
	if err != nil {
		log.Fatal("Не удалось подключиться к бд", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing database pool")
			pool.Close()
			return nil
		},
	})

	return pool
}

func NewStorage(pool *pgxpool.Pool, log *slog.Logger) *Storage {
	return &Storage{
		pool: pool,
		log:  log,
	}
}
func (s *Storage) Get(ctx context.Context) (*models.Product, error) {
	query := `
		SELECT 
		product_id,
		name, price, available_stock
		FROM product
		ORDER BY RANDOM()
		LIMIT 1
	`

	var product models.Product
	err := s.pool.QueryRow(ctx, query).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
	)

	if err != nil {
		s.log.Error("failed to fetch random product",
			slog.String("query", "get_random_product"),
			slog.Any("error", err),
		)
		return nil, err
	}
	s.log.Debug("random product fetched for update",
		slog.String("product_id", product.ID),
		slog.String("name", product.Name),
		slog.Float64("price", product.Price),
		slog.Int("stock", product.Stock),
	)

	return &product, nil
}
