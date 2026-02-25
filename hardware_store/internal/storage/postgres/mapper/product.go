package mapper

import (
	model "hardware_store/internal/model/product"
	"hardware_store/internal/storage/postgres/dto"
)

func ProductToDTO(p model.Product) dto.ProductDTO {
	return dto.ProductDTO{
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

func ProductFromDTO(d dto.ProductDTO) model.Product {
	return model.Product{
		ProductID:      d.ProductID,
		Name:           d.Name,
		CategoryID:     d.CategoryID,
		Price:          d.Price,
		AvailableStock: d.AvailableStock,
		LastUpdateDate: d.LastUpdateDate,
		SupplierID:     d.SupplierID,
		ImageID:        d.ImageID,
	}
}
