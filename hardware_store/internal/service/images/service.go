package images

import (
	"context"
	"fmt"
	"hardware_store/internal/logger"
	"hardware_store/internal/model/images"
	"hardware_store/internal/model/tx"
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
	repo   ImagesRepository
	tx     tx.Manager
	logger *slog.Logger
}

func NewImageService(repo ImagesRepository, tx tx.Manager, logger *slog.Logger) *imageService {
	return &imageService{
		repo:   repo,
		tx:     tx,
		logger: logger,
	}
}

func (s *imageService) CreateImage(ctx context.Context, image []byte, product uuid.UUID) (uuid.UUID, error) {
	imgID := uuid.New()
	s.logger.Info("Creating image",
		slog.String("image_id", imgID.String()),
		slog.String("product_id", product.String()),
		slog.Int("image_size", len(image)),
	)
	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.repo.Insert(ctx, images.Images{ImageID: imgID, Image: image}); err != nil {
			s.logger.Error("Failed to insert image",
				logger.Err(err),
				slog.String("image_id", imgID.String()),
			)
			return fmt.Errorf("failed to insert image: %w", err)
		}
		s.logger.Info("Image inserted successfully",
			slog.String("image_id", imgID.String()),
		)
		err := s.repo.Attach(ctx, imgID, product)
		if err != nil {
			fmt.Printf("Attach error: %v\n", err)
			return err
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return imgID, nil
}

func (s *imageService) UpdateImage(ctx context.Context, id uuid.UUID, image []byte) error {
	return s.repo.Update(ctx, images.Images{ImageID: id, Image: image})
}

func (s *imageService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *imageService) GetImage(ctx context.Context, id uuid.UUID) (images.Images, error) {
	return s.repo.GetById(ctx, id)
}

func (s *imageService) GetImageByProduct(ctx context.Context, product uuid.UUID) (images.Images, error) {
	return s.repo.GetByProduct(ctx, product)
}
