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

//для кафки
type ProductUpdate struct {
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	OldPrice     float64   `json:"old_price"`
	NewPrice     float64   `json:"new_price"`
	OldStock     int       `json:"old_stock"`
	NewStock     int       `json:"new_stock"`
	Timestamp    time.Time `json:"timestamp"`
}