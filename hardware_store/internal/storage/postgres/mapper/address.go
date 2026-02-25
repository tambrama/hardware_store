package mapper

import (
	model "hardware_store/internal/model/address"
	"hardware_store/internal/storage/postgres/dto"
)

func AddressToDTO(a model.Address) dto.AddressDTO {
	return dto.AddressDTO{
		AddressID: a.AddressID,
		Country:   a.Country,
		City:      a.City,
		Street:    a.Street,
	}
}

func AddressFromDTO(d dto.AddressDTO) model.Address {
	return model.Address{
		AddressID: d.AddressID,
		Country:   d.Country,
		City:      d.City,
		Street:    d.Street,
	}
}
