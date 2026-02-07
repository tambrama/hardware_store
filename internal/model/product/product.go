package product

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ProductID      uuid.UUID
	Name           string
	CategoryID     uuid.UUID
	Price          float64
	AvailableStock int
	LastUpdateDate time.Time
	SupplierID     uuid.UUID
	ImageID        *uuid.UUID
}

type ProductRepository interface {
	Insert(ctx context.Context, product Product) error
	UpdateBalance(ctx context.Context, id uuid.UUID, col int) (Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	UnsetCategory(ctx context.Context, category uuid.UUID) error
}
