package dto

import (
	"time"

	"github.com/google/uuid"
)

type ClientDTO struct {
	ClientID         uuid.UUID `db:"client_id"`
	Name             string    `db:"name"`
	Surname          string    `db:"surname"`
	Birthday         time.Time `db:"birthday"`
	Gender           string    `db:"gender"`
	RegistrationDate time.Time `db:"registration_date"`
	AddressID        uuid.UUID `db:"address_id"`
}

type ProductDTO struct {
	ProductID      uuid.UUID  `db:"product_id"`
	Name           string     `db:"name"`
	CategoryID     uuid.UUID  `db:"category_id"`
	Price          float64    `db:"price"`
	AvailableStock int        `db:"available_stock"`
	LastUpdateDate time.Time  `db:"last_update_date"`
	SupplierID     uuid.UUID  `db:"supplier_id"`
	ImageID        *uuid.UUID `db:"image_uuid"`
}

type SupplierDTO struct {
	SupplierID  uuid.UUID `db:"supplier_id"`
	Name        string    `db:"name"`
	AddressID   uuid.UUID `db:"address_id"`
	PhoneNumber string    `db:"phone_number"`
}

type ImagesDTO struct {
	ImageID uuid.UUID `db:"images_id"`
	Image   []byte    `db:"image"`
}

type AddressDTO struct {
	AddressID uuid.UUID `db:"address_id"`
	Country   string    `db:"country"`
	City      string    `db:"city"`
	Street    string    `db:"street"`
}

type CategoryDTO struct {
	CategoryID uuid.UUID `db:"category_id"`
	Category   string    `db:"category"`
}
