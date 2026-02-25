package address

import (
	"github.com/google/uuid"
)

type Address struct {
	AddressID uuid.UUID `validate:"required"`
	Country   string    `validate:"required,min=2,max=50"`
	City      string    `validate:"required,min=2,max=50"`
	Street    string
}
