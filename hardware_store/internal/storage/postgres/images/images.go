package images

import (
	"context"
	"errors"
	"fmt"
	"hardware_store/internal/model/images"
	"hardware_store/internal/storage"
	"hardware_store/internal/storage/postgres/dto"
	"hardware_store/internal/storage/postgres/mapper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type imagesRepository struct {
	pool *pgxpool.Pool
}

func NewImagesRepository(db *pgxpool.Pool) *imagesRepository {
	return &imagesRepository{
		pool: db,
	}
}

func (r *imagesRepository) Insert(ctx context.Context, image images.Images) error {
	imageDtO := mapper.ImageToDTO(image)
	query := `INSERT INTO images
	(image_id, image)
	VALUES ($1, $2)`

	_, err := r.pool.Exec(ctx, query, imageDtO.ImageID, imageDtO.Image)
	if err != nil {
		return storage.ErrCreation
	}
	return nil
}

func (r *imagesRepository) Update(ctx context.Context, image images.Images) error {
	query := `UPDATE images 
	SET image = $2
	WHERE image_id = $1`

	dto := mapper.ImageToDTO(image)

	_, err := r.pool.Exec(ctx, query, dto.ImageID, dto.Image)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrImageNotFound
		}
		return storage.ErrUpdate
	}

	return nil
}

func (r *imagesRepository) Delete(ctx context.Context, imagesID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryUPD := `UPDATE product
	SET image_id = NULL 
	WHERE image_id = $1`
	_, err = r.pool.Exec(ctx, queryUPD, imagesID)
	if err != nil {
		return storage.ErrUpdate
	}

	query := `DELETE FROM images 
	WHERE image_id = $1`
	_, err = r.pool.Exec(ctx, query, imagesID)
	if err != nil {
		return storage.ErrDelete
	}

	return tx.Commit(ctx)
}

func (r *imagesRepository) GetByProduct(ctx context.Context, productID uuid.UUID) (images.Images, error) {
	query := `SELECT i.image_id, i.image FROM images i
	JOIN product p ON p.image_id = i.image_id
	WHERE product_id = $1`

	var dto dto.ImagesDTO

	err := r.pool.QueryRow(ctx, query, productID).Scan(&dto.ImageID, &dto.Image)
	if err != nil {
		return images.Images{}, storage.ErrUpdate
	}

	return mapper.ImageFromDTO(dto), nil
}
func (r *imagesRepository) Attach(ctx context.Context, imagesID, productID uuid.UUID) error {
	query := `UPDATE product
        SET image_id = $1
        WHERE product_id = $2`

	res, err := r.pool.Exec(ctx, query, imagesID, productID)
	if err != nil {
		return fmt.Errorf("failed to attach image: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("product not found: %s", productID)
	}

	return nil
}

func (r *imagesRepository) GetById(ctx context.Context, imagesID uuid.UUID) (images.Images, error) {
	query := `SELECT image_id, image FROM images 
	WHERE image_id = $1`

	var dto dto.ImagesDTO

	err := r.pool.QueryRow(ctx, query, imagesID).Scan(&dto.ImageID, &dto.Image)
	if err != nil {
		return images.Images{}, storage.ErrImageNotFound
	}

	return mapper.ImageFromDTO(dto), nil
}
