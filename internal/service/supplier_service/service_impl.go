package supplierservice

import (
	"context"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/supplier"
	"hardware_store/internal/model/tx"

	"github.com/google/uuid"
)

type supplierService struct {
	repo     supplier.SupplierRepository
	repoAddr address.AddressRepository
	tx       tx.Manager
}

func NewSupplierService(repo supplier.SupplierRepository,
	repoAddr address.AddressRepository,
	tx tx.Manager) SupplierService {
	return &supplierService{repo: repo,
		repoAddr: repoAddr,
		tx:       tx}
}

func (s *supplierService) CreateSupplier(ctx context.Context, supplier supplier.Supplier, address address.Address) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.repoAddr.Insert(ctx, address); err != nil {
			return err
		}
		return s.repo.Insert(ctx, supplier)
	})
}
func (s *supplierService) DeleteSupplier(ctx context.Context, id uuid.UUID) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		getSup, err := s.repo.GetById(ctx, id)
		if err != nil {
			return err
		}
		if err = s.repo.Delete(ctx, id); err != nil {
			return err
		}

		return s.repoAddr.Delete(ctx, getSup.AddressID)
	})
}
func (s *supplierService) UpdateAddressSupplier(ctx context.Context, id uuid.UUID, address address.Address) error {
	return s.repo.UpdateAddress(ctx, id, address)
}
func (s *supplierService) GetSupplier(ctx context.Context, id uuid.UUID) (supplier.Supplier, error) {
	return s.repo.GetById(ctx, id)
}
func (s *supplierService) GetSuppliers(ctx context.Context) ([]supplier.Supplier, error) {
	return s.repo.GetAll(ctx)
}
