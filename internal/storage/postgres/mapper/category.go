package mapper

import (
	"hardware_store/internal/model/category"
	"hardware_store/internal/storage/postgres/dto"
)

func CategoryToDTO(i category.Category) dto.CategoryDTO {
	return dto.CategoryDTO{
		CategoryID: i.CategoryID,
		Category:   i.Category,
	}
}

func CategoryFromDTO(d dto.CategoryDTO) category.Category {
	return category.Category{
		CategoryID: d.CategoryID,
		Category:   d.Category,
	}
}
