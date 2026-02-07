package dto

import (
	"time"

	"github.com/google/uuid"
)

// ProductRequest запрос на создание товара
// @Description Запрос на создание нового товара с категорией, ценой и информацией о поставщике
// swagger:model ProductRequest
type ClientRequest struct {
	Name     string         `json:"name" validate:"required,min=2,max=50" example:"Иван"`
	Surname  string         `json:"surname" validate:"required,min=2,max=50" example:"Иванов"`
	Birthday string         `json:"birthday" validate:"datetime=2006-01-02" example:"1999-01-01"`
	Gender   string         `json:"gender" validate:"oneof=male female" example:"male"`
	Address  AddressRequest `json:"address"`
}

// ClientResponse ответ с информацией о клиенте
// @Description Данные клиента включая дату регистрации и ссылку на адрес
// swagger:model ClientResponse
type ClientResponse struct {
	ClientID         uuid.UUID `json:"client_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name             string    `json:"name"`
	Surname          string    `json:"surname"`
	Birthday         time.Time `json:"birthday"`
	Gender           string    `json:"gender"`
	RegistrationDate time.Time `json:"registration_date"`
	AddressID        uuid.UUID `json:"address_uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// UpdateStockCountRequest запрос на обновление остатков
// @Description Запрос на обновление количества товара на складе
// swagger:model UpdateStockCountRequest
type ProductRequest struct {
	Name           string    `json:"name" validate:"required,min=2,max=100" example:"Холодильник Samsung RB38A7861B1"`
	CategoryID     uuid.UUID `json:"category_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Price          float64   `json:"price" validate:"required,gt=0" example:"75990.00"`
	AvailableStock int       `json:"available_stock" validate:"required,gte=0" example:"15"`
	SupplierID     uuid.UUID `json:"supplier_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdateStockCountRequest struct {
	Amount int `json:"amount" validate:"required,gt=0" example:"5"`
}

// ProductResponse ответ с информацией о товаре
// @Description Полная информация о товаре включая цену, остатки и информацию о поставщике
// swagger:model ProductResponse
type ProductResponse struct {
	ProductID      uuid.UUID  `json:"product_id" validate:"required" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Name           string     `json:"name" validate:"required,min=2,max=100"`
	CategoryID     uuid.UUID  `json:"category" example:"550e8400-e29b-41d4-a716-446655440000"`
	Price          float64    `json:"price" validate:"required,gt=0"`
	AvailableStock int        `json:"available_stock" validate:"required,gte=0"`
	LastUpdateDate time.Time  `json:"last_update_date" validate:"required"`
	SupplierID     uuid.UUID  `json:"supplier" example:"550e8400-e29b-41d4-a716-446655440000"`
	ImageID        *uuid.UUID `json:"image" example:"b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22"`
}

// SupplierRequest запрос на создание поставщика
// @Description Запрос на создание нового поставщика с контактной информацией и адресом
// swagger:model SupplierRequest
type SupplierRequest struct {
	Name        string         `json:"name" validate:"required,min=2,max=100" example:"ООО 'ТехноСнаб'"`
	Address     AddressRequest `json:"address"`
	PhoneNumber string         `json:"phone_number" validate:"required,e164" example:"+79561234567"`
}

// SupplierResponse ответ с информацией о поставщике
// @Description Данные поставщика включая контактную информацию и ссылку на адрес
// swagger:model SupplierResponse
type SupplierResponse struct {
	SupplierID  uuid.UUID `json:"supplier_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"ООО 'ТехноСнаб'"`
	AddressID   uuid.UUID `json:"address_uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	PhoneNumber string    `json:"phone_number" example:"+79561234567"`
}

// ImageRequest запрос на загрузку изображения
// @Description Запрос на загрузку файла изображения для товара
// swagger:model ImageRequest
type ImageRequest struct {
	Image []byte `json:"image" validate:"required" example:"base64-encoded-image-data"`
}

// ImageResponse ответ с информацией об изображении
// @Description Данные изображения включая уникальный идентификатор
// swagger:model ImageResponse
type ImageResponse struct {
	ImageID uuid.UUID `json:"image_id" example:"b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22"`
}

// AddressRequest запрос на создание адреса
// @Description Запрос на создание нового адреса с информацией о стране, городе и улице
// swagger:model AddressRequest
type AddressRequest struct {
	Country string `json:"country" validate:"required,min=2,max=50" example:"Россия"`
	City    string `json:"city" validate:"required,min=2,max=50" example:"Москва"`
	Street  string `json:"street" validate:"required,min=2,max=100" example:"Технопарк, 15"`
}

// AddressResponse ответ с информацией об адресе
// @Description Данные адреса включая информацию о местоположении
// swagger:model AddressResponse
type AddressResponse struct {
	AddressID uuid.UUID `json:"address_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Country   string    `json:"country" example:"Россия"`
	City      string    `json:"city" example:"Москва"`
	Street    string    `json:"street" example:"Технопарк, 15"`
}

// CategoryRequest запрос на создание категории
// @Description Запрос на создание новой категории товаров
// swagger:model CategoryRequest
type CategoryRequest struct {
	Category string `json:"category" validate:"required,min=2,max=50" example:"Холодильники"`
}

// CategoryResponse ответ с информацией о категории
// @Description Данные категории включая уникальный идентификатор и название
// swagger:model CategoryResponse
type CategoryResponse struct {
	CategoryID uuid.UUID `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Category   string    `json:"category" example:"Холодильники"`
}
