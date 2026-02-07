package mapper

import (
	"hardware_store/internal/model/client"
	"hardware_store/internal/storage/postgres/dto"
)

func ClientToDTO(c client.Client) dto.ClientDTO {
	return dto.ClientDTO{
		ClientID:         c.ClientID,
		Name:             c.Name,
		Surname:          c.Surname,
		Birthday:         c.Birthday,
		Gender:           c.Gender,
		RegistrationDate: c.RegistrationDate,
		AddressID:        c.AddressID,
	}
}

func ClientFromDTO(d dto.ClientDTO) client.Client {
	return client.Client{
		ClientID:         d.ClientID,
		Name:             d.Name,
		Surname:          d.Surname,
		Birthday:         d.Birthday,
		Gender:           d.Gender,
		RegistrationDate: d.RegistrationDate,
		AddressID:        d.AddressID,
	}
}
