package mapper

import (
	model "hardware_store/internal/model/images"
	"hardware_store/internal/storage/postgres/dto"
)

func ImageToDTO(i model.Images) dto.ImagesDTO {
	return dto.ImagesDTO{
		ImageID: i.ImageID,
		Image:   i.Image,
	}
}

func ImageFromDTO(d dto.ImagesDTO) model.Images {
	return model.Images{
		ImageID: d.ImageID,
		Image:   d.Image,
	}
}
