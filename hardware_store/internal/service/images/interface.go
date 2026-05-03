package images

import (
	"context"

	"github.com/google/uuid"
)

type ImageService interface {
	CreateImage(ctx context.Context, image []byte, product uuid.UUID) (uuid.UUID, error)
	UpdateImage(ctx context.Context, id uuid.UUID, image []byte) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
	GetImage(ctx context.Context, id uuid.UUID) ([]byte, error)
	GetImageByProduct(ctx context.Context, product uuid.UUID) ([]byte, error)
}
