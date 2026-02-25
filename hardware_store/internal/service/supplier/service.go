package supplier

import (
	"context"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/supplier"
	"hardware_store/internal/model/tx"
	service "hardware_store/internal/service/address"

	"github.com/google/uuid"
)

type SupplierRepository interface {
	Insert(ctx context.Context, supplier supplier.Supplier) error
	UpdateAddress(ctx context.Context, id uuid.UUID, addr address.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (supplier.Supplier, error)
	GetAll(ctx context.Context) ([]supplier.Supplier, error)
	UnsetAddress(ctx context.Context, addressId uuid.UUID) error
}

type AddressRepository interface {
	Insert(ctx context.Context, address address.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type supplierService struct {
	repo  SupplierRepository
	addrr service.AddressService
	tx    tx.Manager
}

func NewSupplierService(repo SupplierRepository,
	addrr service.AddressService,
	tx tx.Manager) *supplierService {
	return &supplierService{repo: repo,
		addrr: addrr,
		tx:    tx}
}

func (s *supplierService) CreateSupplier(ctx context.Context, supplier supplier.Supplier, address address.Address) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.addrr.CreateAddress(ctx, address); err != nil {
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

		return s.addrr.DeleteAddress(ctx, getSup.AddressID)
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
