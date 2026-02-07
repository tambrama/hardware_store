package supplier

import (
	"context"
	"hardware_store/internal/model/address"

	"github.com/google/uuid"
)

type Supplier struct {
	SupplierID  uuid.UUID
	Name        string
	AddressID   uuid.UUID
	PhoneNumber string
}

type SupplierRepository interface {
	Insert(ctx context.Context, supplier Supplier) error
	UpdateAddress(ctx context.Context, id uuid.UUID, addr address.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (Supplier, error)
	GetAll(ctx context.Context) ([]Supplier, error)
	UnsetAddress(ctx context.Context, addressId uuid.UUID) error
}
