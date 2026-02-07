package productservice

import (
	"context"
	"errors"
	"hardware_store/internal/model/images"
	"hardware_store/internal/model/product"
	"hardware_store/internal/model/tx"

	"github.com/google/uuid"
)

var ErrAmountIsNegative = errors.New("amount must be positive")

type productService struct {
	repo    product.ProductRepository
	repoImg images.ImagesRepository
	tx      tx.Manager
}

func NewProductService(repo product.ProductRepository, repoImg images.ImagesRepository, tx tx.Manager) ProductService {
	return &productService{repo: repo, repoImg: repoImg, tx: tx}
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
			return s.repoImg.Delete(ctx, *getProduct.ImageID)
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
