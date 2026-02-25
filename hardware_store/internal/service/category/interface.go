package category

import (
	"context"
	"hardware_store/internal/model/category"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, category category.Category) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategory(ctx context.Context, id uuid.UUID) (category.Category, error)
	GetCategories(ctx context.Context) ([]category.Category, error)
	UpdateCategory(ctx context.Context, categoty category.Category) (category.Category, error)
}
