package address

import (
	"context"
	"hardware_store/internal/model/address"

	"github.com/google/uuid"
)

type AddressRepository interface {
	Insert(ctx context.Context, address address.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (address.Address, error)
	Update(ctx context.Context, addr address.Address) (address.Address, error)
}

type addressService struct {
	repo AddressRepository
}

func NewAddressService(repo AddressRepository) *addressService {
	return &addressService{repo: repo}
}
func (s *addressService) CreateAddress(ctx context.Context, address address.Address) error {
	return s.repo.Insert(ctx, address)
}
func (s *addressService) DeleteAddress(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
func (s *addressService) UpdateAddress(ctx context.Context, address address.Address) (address.Address, error) {
	return s.repo.Update(ctx, address)
}
func (s *addressService) GetAddress(ctx context.Context, id uuid.UUID) (address.Address, error) {
	return s.repo.Get(ctx, id)
}
