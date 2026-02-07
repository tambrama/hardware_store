package addressservice

import (
	"context"
	"fmt"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/client"
	"hardware_store/internal/model/supplier"
	"hardware_store/internal/model/tx"

	"github.com/google/uuid"
)

type addressService struct {
	repo         address.AddressRepository
	tx           tx.Manager
	clientRepo   client.ClientRepository
	supplierRepo supplier.SupplierRepository
}

func NewAddressService(tx tx.Manager, repo address.AddressRepository, clientRepo client.ClientRepository, supplierRepo supplier.SupplierRepository) AddressService {
	return &addressService{tx: tx, repo: repo, clientRepo: clientRepo, supplierRepo: supplierRepo}
}
func (s *addressService) CreateAddress(ctx context.Context, address address.Address) error {
	return s.repo.Insert(ctx, address)
}
func (s *addressService) DeleteAddress(ctx context.Context, id uuid.UUID) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.clientRepo.UnsetAddress(ctx, id); err != nil {
			return fmt.Errorf("failed to unset client address: %w", err)
		}
		if err := s.supplierRepo.UnsetAddress(ctx, id); err != nil {
			return fmt.Errorf("failed to unset supplier address: %w", err)
		}
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete address: %w", err)
		}

		return nil
	})
}
func (s *addressService) UpdateAddress(ctx context.Context, address address.Address) (address.Address, error) {
	return s.repo.Update(ctx, address)
}
func (s *addressService) GetAddress(ctx context.Context, id uuid.UUID) (address.Address, error) {
	return s.repo.Get(ctx, id)
}
