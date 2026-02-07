package productservice

import (
	"context"
	"hardware_store/internal/model/product"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product product.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateProduct(ctx context.Context, id uuid.UUID, col int) (product.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (product.Product, error)
	GetProducts(ctx context.Context) ([]product.Product, error)
}