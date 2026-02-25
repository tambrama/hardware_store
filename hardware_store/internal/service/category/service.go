package category

import (
	"context"
	"fmt"
	"hardware_store/internal/model/category"
	"hardware_store/internal/model/tx"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Insert(ctx context.Context, category category.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (category.Category, error)
	GetAll(ctx context.Context) ([]category.Category, error)
	Update(ctx context.Context, category category.Category) (category.Category, error)
	UnsetCategory(ctx context.Context, category uuid.UUID) error
}

type categoryService struct {
	repo CategoryRepository
	tx   tx.Manager
}

func NewCategoryService(repo CategoryRepository, tx tx.Manager) *categoryService {
	return &categoryService{
		repo: repo,
		tx:   tx,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, category category.Category) error {
	return s.repo.Insert(ctx, category)
}
func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.repo.UnsetCategory(ctx, id); err != nil {
			return fmt.Errorf("failed to unset client address: %w", err)
		}
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete address: %w", err)
		}

		return nil
	})
}
func (s *categoryService) GetCategory(ctx context.Context, id uuid.UUID) (category.Category, error) {
	return s.repo.GetById(ctx, id)
}
func (s *categoryService) GetCategories(ctx context.Context) ([]category.Category, error) {
	return s.repo.GetAll(ctx)
}
func (s *categoryService) UpdateCategory(ctx context.Context, category category.Category) (category.Category, error) {
	return s.repo.Update(ctx, category)
}
