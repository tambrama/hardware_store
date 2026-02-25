package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ProductID      uuid.UUID
	Name           string
	CategoryID     uuid.UUID
	Price          float64
	AvailableStock int
	LastUpdateDate time.Time
	SupplierID     uuid.UUID
	ImageID        *uuid.UUID
}
