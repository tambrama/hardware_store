package supplier

import (
	"github.com/google/uuid"
)

type Supplier struct {
	SupplierID  uuid.UUID
	Name        string
	AddressID   uuid.UUID
	PhoneNumber string
}
