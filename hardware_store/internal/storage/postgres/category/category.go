package category

import (
	"context"
	"errors"
	"fmt"
	"hardware_store/internal/model/category"
	"hardware_store/internal/storage"
	"hardware_store/internal/storage/postgres/dto"
	"hardware_store/internal/storage/postgres/mapper"
	"hardware_store/internal/storage/postgres/tx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *categoryRepository {
	return &categoryRepository{
		pool: db,
	}
}

func (r *categoryRepository) Insert(ctx context.Context, category category.Category) error {
	query := `INSERT INTO category
	(category_id, category)
	VALUES ($1, $2)`

	dto := mapper.CategoryToDTO(category)
	_, err := r.pool.Exec(ctx, query, dto.CategoryID, dto.Category)
	if err != nil {
		return storage.ErrCreation
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM category 
	WHERE category_id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return storage.ErrDelete
	}

	return nil
}

func (r *categoryRepository) GetById(ctx context.Context, id uuid.UUID) (category.Category, error) {
	query := `SELECT * FROM category 
	WHERE category_id = $1`

	var dto dto.CategoryDTO

	err := r.pool.QueryRow(ctx, query, id).Scan(&dto.CategoryID, &dto.Category)
	if err != nil {
		return category.Category{}, storage.ErrCategoryNotFound
	}

	return mapper.CategoryFromDTO(dto), nil
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]category.Category, error) {
	query := `SELECT * FROM category`

	row, err := r.pool.Query(ctx, query)
	if err != nil {
		return []category.Category{}, storage.ErrCategoryNotFound
	}
	defer row.Close()
	var categories []category.Category
	for row.Next() {
		var dto dto.CategoryDTO

		if err := row.Scan(&dto.CategoryID, &dto.Category); err != nil {
			return []category.Category{}, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		categories = append(categories, mapper.CategoryFromDTO(dto))
	}

	if err = row.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, cat category.Category) (category.Category, error) {
	query := `UPDATE category 
	SET category = $2
	WHERE category_id = $1
	RETURNING category_id, category`

	categoryDTO := mapper.CategoryToDTO(cat)
	var dtoUPD dto.CategoryDTO

	err := r.pool.QueryRow(ctx, query, categoryDTO.CategoryID, categoryDTO.Category).Scan(&dtoUPD.CategoryID, &dtoUPD.Category)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return category.Category{}, storage.ErrCategoryNotFound
		}
		return category.Category{}, storage.ErrUpdate
	}

	return mapper.CategoryFromDTO(dtoUPD), nil
}

func (r *categoryRepository) UnsetCategory(ctx context.Context, category uuid.UUID) error {
	exec := tx.FromContext(ctx, r.pool)

	query := `UPDATE product 
	SET category_id = NULL
	WHERE category_id = $1`

	_, err := exec.Exec(ctx, query, category)
	if err != nil {
		return fmt.Errorf("failed to unset category %s from products: %w", category, err)
	}

	return nil
}
