package category

import (
	"context"

	"github.com/google/uuid"
)

type Category struct {
	CategoryID uuid.UUID
	Category   string
}

type CategoryRepository interface {
	Insert(ctx context.Context, category Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, category Category) (Category, error)
}
