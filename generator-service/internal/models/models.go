package models

import "time"

// ProductUpdate - событие об изменении товара
type ProductUpdate struct {
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	OldPrice     float64   `json:"old_price"`
	NewPrice     float64   `json:"new_price"`
	OldStock     int       `json:"old_stock"`
	NewStock     int       `json:"new_stock"`
	Timestamp    time.Time `json:"timestamp"`
	Source       string    `json:"source"`
}

type Product struct {
	ID      string
	Name    string
	Price   float64
	Stock   int
}

type ProductSource string

const (
    SourceShopAPI    ProductSource = "shopapi"
    SourceHardware   ProductSource = "hardware"
)