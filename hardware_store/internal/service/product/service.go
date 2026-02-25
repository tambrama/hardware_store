package product

import (
	"context"
	"errors"
	"hardware_store/internal/model/product"
	"hardware_store/internal/service/images"
	"hardware_store/internal/model/tx"

	"github.com/google/uuid"
)

var ErrAmountIsNegative = errors.New("amount must be positive")

type ProductRepository interface {
	Insert(ctx context.Context, product product.Product) error
	UpdateBalance(ctx context.Context, id uuid.UUID, col int) (product.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (product.Product, error)
	GetAll(ctx context.Context) ([]product.Product, error)
}

type productService struct {
	repo    ProductRepository
	img 	images.ImageService
	tx      tx.Manager
}

func NewProductService(repo ProductRepository, img images.ImageService, tx tx.Manager) *productService {
	return &productService{repo: repo, img: img, tx: tx}
}

func (s *productService) CreateProduct(ctx context.Context, product product.Product) error {

	return s.repo.Insert(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	getProduct, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if getProduct.ImageID != nil {
		return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
			err = s.repo.Delete(ctx, id)
			if err != nil {
				return err
			}
			return s.img.DeleteImage(ctx, *getProduct.ImageID)
		})
	} else {
		return s.repo.Delete(ctx, id)
	}
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, col int) (product.Product, error) {
	if col < 0 {
		return product.Product{}, ErrAmountIsNegative
	}

	return s.repo.UpdateBalance(ctx, id, col)
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (product.Product, error) {
	return s.repo.GetById(ctx, id)
}

func (s *productService) GetProducts(ctx context.Context) ([]product.Product, error) {
	return s.repo.GetAll(ctx)
}
