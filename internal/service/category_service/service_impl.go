package categoryservice

import (
	"context"
	"fmt"
	"hardware_store/internal/model/category"
	"hardware_store/internal/model/product"
	"hardware_store/internal/model/tx"

	"github.com/google/uuid"
)

type categoryService struct {
	repo        category.CategoryRepository
	repoProduct product.ProductRepository
	tx          tx.Manager
}

func NewCategoryService(repo category.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, category category.Category) error {
	return s.repo.Insert(ctx, category)
}
func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.repoProduct.UnsetCategory(ctx, id); err != nil {
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
