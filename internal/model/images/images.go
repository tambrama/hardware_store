package images

import (
	"context"

	"github.com/google/uuid"
)

type Images struct {
	ImageID uuid.UUID
	Image   []byte
}

type ImagesRepository interface {
	Insert(ctx context.Context, image Images) error
	Update(ctx context.Context, image Images) error
	Delete(ctx context.Context, imagesID uuid.UUID) error
	GetByProduct(ctx context.Context, productID uuid.UUID) (Images, error)
	GetById(ctx context.Context, imagesID uuid.UUID) (Images, error)
	Attach(ctx context.Context, imagesID, productID uuid.UUID) error
}
