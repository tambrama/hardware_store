package address

import (
	"context"

	"github.com/google/uuid"
)

type Address struct {
	AddressID uuid.UUID `validate:"required"`
	Country   string    `validate:"required,min=2,max=50"`
	City      string    `validate:"required,min=2,max=50"`
	Street    string
}

type AddressRepository interface {
	Insert(ctx context.Context, address Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (Address, error)
	Update(ctx context.Context, addr Address) (Address, error)
}
