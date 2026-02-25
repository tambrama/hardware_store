package product

import (
	"context"
	"errors"
	"fmt"
	"hardware_store/internal/model/product"
	"hardware_store/internal/storage"
	"hardware_store/internal/storage/postgres/dto"
	"hardware_store/internal/storage/postgres/mapper"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewProductRepository(db *pgxpool.Pool, log *slog.Logger) *productRepository {
	return &productRepository{
		pool: db,
		log:  log,
	}
}

func (r *productRepository) Insert(ctx context.Context, product product.Product) error {
	dto := mapper.ProductToDTO(product)
	query := `INSERT INTO product 
	(product_id, name, category_id, price, available_stock, last_update_date, supplier_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, query, dto.ProductID, dto.Name, dto.CategoryID, dto.Price, dto.AvailableStock, dto.LastUpdateDate, dto.SupplierID)
	if err != nil {
		r.log.Error("failed to insert product",
			slog.Any("error", err),
			slog.String("product_id", dto.ProductID.String()),
		)
		return storage.ErrCreation
	}

	return nil
}

func (r *productRepository) UpdateBalance(ctx context.Context, id uuid.UUID, col int) (product.Product, error) {
	query := `UPDATE product 
	SET available_stock = available_stock - $2
	WHERE product_id = $1 AND available_stock >= $2
	RETURNING product_id, name, category_id, price, available_stock, last_update_date, supplier_id, image_id`

	var dto dto.ProductDTO

	err := r.pool.QueryRow(ctx, query, id, col).Scan(&dto.ProductID, &dto.Name, &dto.CategoryID, &dto.Price, &dto.AvailableStock, &dto.LastUpdateDate, &dto.SupplierID, &dto.ImageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return product.Product{}, storage.ErrProductNotFound
		}
		return product.Product{}, fmt.Errorf("ошибка обновления товара: %w", err)
	}
	return mapper.ProductFromDTO(dto), nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM product 
	WHERE product_id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return storage.ErrDelete
	}
	return nil
}

func (r *productRepository) GetById(ctx context.Context, id uuid.UUID) (product.Product, error) {
	query := `SELECT * FROM product 
	WHERE product_id = $1`

	var dto dto.ProductDTO

	err := r.pool.QueryRow(ctx, query, id).Scan(&dto.ProductID, &dto.Name, &dto.CategoryID, &dto.Price, &dto.AvailableStock, &dto.LastUpdateDate, &dto.SupplierID, &dto.ImageID)
	if err != nil {
		return product.Product{}, storage.ErrProductNotFound
	}
	return mapper.ProductFromDTO(dto), nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]product.Product, error) {
	query := `SELECT * FROM product
	WHERE available_stock > 0`

	row, err := r.pool.Query(ctx, query)
	if err != nil {
		return []product.Product{}, storage.ErrProductNotFound
	}
	defer row.Close()
	var products []product.Product
	for row.Next() {
		var dto dto.ProductDTO

		if err := row.Scan(&dto.ProductID, &dto.Name, &dto.CategoryID, &dto.Price, &dto.AvailableStock, &dto.LastUpdateDate, &dto.SupplierID, &dto.ImageID); err != nil {
			return []product.Product{}, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		products = append(products, mapper.ProductFromDTO(dto))
	}

	if err = row.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return products, nil
}
