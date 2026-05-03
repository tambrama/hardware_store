package images

import (
	"context"
	"fmt"
	"hardware_store/internal/clients/photo"
	"hardware_store/internal/logger"
	"hardware_store/internal/model/images"
	"hardware_store/internal/model/tx"
	"hardware_store/internal/storage/cache"
	"log/slog"

	"github.com/google/uuid"
)

type ImagesRepository interface {
	Insert(ctx context.Context, image images.Images) error
	Update(ctx context.Context, image images.Images) error
	Delete(ctx context.Context, imagesID uuid.UUID) error
	GetByProduct(ctx context.Context, productID uuid.UUID) (images.Images, error)
	GetById(ctx context.Context, imagesID uuid.UUID) (images.Images, error)
	Attach(ctx context.Context, imagesID, productID uuid.UUID) error
}

type imageService struct {
	repo        ImagesRepository
	clientPhoto *photo.PhotoServiceClient
	cache       *cache.ImageCache
	tx          tx.Manager
	logger      *slog.Logger
}

func NewImageService(repo ImagesRepository, clientPhoto *photo.PhotoServiceClient, cache *cache.ImageCache, tx tx.Manager, logger *slog.Logger) *imageService {
	return &imageService{
		repo:        repo,
		clientPhoto: clientPhoto,
		cache:       cache,
		tx:          tx,
		logger:      logger,
	}
}

func (s *imageService) CreateImage(ctx context.Context, image []byte, product uuid.UUID) (uuid.UUID, error) {
	const op = "images.CreateImage"

	imgResp, err := s.clientPhoto.CreatePhoto(ctx, image)
	if err != nil {
		s.logger.Error("Failed to create photo via PhotoServiceClient",
			logger.Err(err),
			slog.String("product_id", product.String()),
		)
		return uuid.Nil, fmt.Errorf("%s: create in photo-service: %w", op, err)
	}

	err = s.repo.Attach(ctx, imgResp.ImageID, product)
	if err != nil {
		fmt.Printf("Attach error: %v\n", err)
		return uuid.Nil, fmt.Errorf("%s: attach to product: %w", op, err)
	}

	s.cache.Set(ctx, product, images.Images{ImageID: imgResp.ImageID})

	return imgResp.ImageID, nil
}

func (s *imageService) UpdateImage(ctx context.Context, id uuid.UUID, image []byte) error {
	const op = "images.UpdateImage"

	_, err := s.clientPhoto.UpdatePhoto(ctx, image, id.String())
	if err != nil {
		return fmt.Errorf("%s: update photo: %w", op, err)
	}
	s.cache.Delete(ctx, id)
	return nil
}

func (s *imageService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	const op = "images.DeleteImage"
	err := s.clientPhoto.DeletePhoto(ctx, id.String())
	if err != nil {
		return fmt.Errorf("%s: delete photo: %w", op, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: detach from product: %w", op, err)
	}

	s.cache.Delete(ctx, id)

	return nil
}

func (s *imageService) GetImage(ctx context.Context, id uuid.UUID) ([]byte, error) {
	const op = "images.GetImage"

	img, err := s.clientPhoto.GetPhoto(ctx, id.String())
	if err != nil {
		return []byte{}, fmt.Errorf("%s: update in photo-service: %w", op, err)
	}
	return img, nil
}

func (s *imageService) GetImageByProduct(ctx context.Context, product uuid.UUID) ([]byte, error) {
	const op = "images.GetImageByProduct"

	cache, err := s.cache.Get(ctx, product)
	if err == nil {
		return s.GetImage(ctx, cache.ImageID)
	}
	
	img, err := s.clientPhoto.GetPhoto(ctx, cache.ImageID.String())
	if err != nil {
		return []byte{}, fmt.Errorf("%s: update in photo-service: %w", op, err)
	}

	return img, nil
}
