package mapper

import (
	"hardware_store/internal/model/category"
	"hardware_store/internal/model/supplier"
	"time"

	"hardware_store/internal/model/address"
	"hardware_store/internal/model/client"
	"hardware_store/internal/model/images"
	"hardware_store/internal/model/product"
	"hardware_store/internal/web/dto"

	"github.com/google/uuid"
)

// === Address mappers ===

func AddressWebToDomain(req dto.AddressRequest, addrID uuid.UUID) address.Address {
	return address.Address{
		AddressID: addrID,
		Country:   req.Country,
		City:      req.City,
		Street:    req.Street,
	}
}

func AddressDomainToWeb(addr address.Address) dto.AddressResponse {
	return dto.AddressResponse{
		AddressID: addr.AddressID,
		Country:   addr.Country,
		City:      addr.City,
		Street:    addr.Street,
	}
}

// === Client mappers ===

func ClientWebToDomain(req dto.ClientRequest, clientID, addressID uuid.UUID, date time.Time) client.Client {
	return client.Client{
		ClientID: clientID,
		Name:     req.Name,
		Surname:  req.Surname,
		// Birthday:         req.Birthday,
		Gender:           req.Gender,
		RegistrationDate: date,
		AddressID:        addressID,
	}
}

func ClientDomainToWeb(client client.Client) dto.ClientResponse {
	return dto.ClientResponse{
		ClientID:         client.ClientID,
		Name:             client.Name,
		Surname:          client.Surname,
		Birthday:         client.Birthday,
		Gender:           client.Gender,
		RegistrationDate: client.RegistrationDate,
		AddressID:        client.AddressID,
	}
}

// === Product mappers ===
func ProductRequestToDomain(
	req dto.ProductRequest,
	productID uuid.UUID,
	lastUpdate time.Time,
) product.Product {
	return product.Product{
		ProductID:      productID,
		Name:           req.Name,
		CategoryID:     req.CategoryID,
		Price:          req.Price,
		AvailableStock: req.AvailableStock,
		LastUpdateDate: lastUpdate,
		SupplierID:     req.SupplierID,
	}
}
func ProductDomainToWeb(p product.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ProductID:      p.ProductID,
		Name:           p.Name,
		CategoryID:     p.CategoryID,
		Price:          p.Price,
		AvailableStock: p.AvailableStock,
		LastUpdateDate: p.LastUpdateDate,
		SupplierID:     p.SupplierID,
		ImageID:        p.ImageID,
	}
}

// === Image mappers ===
func ImageRequestToDomain(req dto.ImageRequest, imageID uuid.UUID) images.Images {
	return images.Images{
		ImageID: imageID,
		Image:   req.Image,
	}
}

func ImageDomainToWeb(img images.Images) dto.ImageResponse {
	return dto.ImageResponse{
		ImageID: img.ImageID,
	}
}

// === Category mappers ===
func CategoryRequestToDomain(req dto.CategoryRequest, id uuid.UUID) category.Category {
	return category.Category{
		CategoryID: id,
		Category:   req.Category,
	}
}

func CategoryDomainToWeb(category category.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		CategoryID: category.CategoryID,
		Category:   category.Category,
	}
}

// === Supplier mappers ===
func SupplierRequestToDomain(req dto.SupplierRequest, supplierID, addressID uuid.UUID) supplier.Supplier {
	return supplier.Supplier{
		SupplierID:  supplierID,
		Name:        req.Name,
		AddressID:   addressID,
		PhoneNumber: req.PhoneNumber,
	}
}

func SupplierDomainToWeb(s supplier.Supplier) dto.SupplierResponse {
	return dto.SupplierResponse{
		SupplierID:  s.SupplierID,
		Name:        s.Name,
		AddressID:   s.AddressID,
		PhoneNumber: s.PhoneNumber,
	}
}
