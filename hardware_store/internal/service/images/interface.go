package images

import (
	"context"
	"hardware_store/internal/model/images"

	"github.com/google/uuid"
)

type ImageService interface {
	CreateImage(ctx context.Context, image []byte, product uuid.UUID) (uuid.UUID, error)
	UpdateImage(ctx context.Context, id uuid.UUID, image []byte) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
	GetImage(ctx context.Context, id uuid.UUID) (images.Images, error)
	GetImageByProduct(ctx context.Context, product uuid.UUID) (images.Images, error)
}
