package address

import (
	"context"
	"hardware_store/internal/model/address"

	"github.com/google/uuid"
)

type AddressService interface {
	CreateAddress(ctx context.Context, address address.Address) error
	DeleteAddress(ctx context.Context, id uuid.UUID) error
	UpdateAddress(ctx context.Context, address address.Address) (address.Address, error)
	GetAddress(ctx context.Context, id uuid.UUID) (address.Address, error)
}
