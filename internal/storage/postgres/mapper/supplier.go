package mapper

import (
	model "hardware_store/internal/model/supplier"
	"hardware_store/internal/storage/postgres/dto"
)

func SupplierToDTO(s model.Supplier) dto.SupplierDTO {
	return dto.SupplierDTO{
		SupplierID:  s.SupplierID,
		Name:        s.Name,
		AddressID:   s.AddressID,
		PhoneNumber: s.PhoneNumber,
	}
}

func SupplierFromDTO(d dto.SupplierDTO) model.Supplier {
	return model.Supplier{
		SupplierID:  d.SupplierID,
		Name:        d.Name,
		AddressID:   d.AddressID,
		PhoneNumber: d.PhoneNumber,
	}
}
