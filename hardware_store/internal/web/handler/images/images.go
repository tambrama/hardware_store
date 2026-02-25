package images

import (
	"errors"
	model "hardware_store/internal/model/error"
	"hardware_store/internal/service/images"

	"hardware_store/internal/web/dto"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ImageHandler struct {
	validator *validator.Validate
	service   images.ImageService
}

func NewImageHandler(validator *validator.Validate, service images.ImageService) *ImageHandler {
	return &ImageHandler{validator: validator, service: service}
}

func (h *ImageHandler) Register(r *gin.RouterGroup) {
	clients := r.Group("/images")
	{
		clients.DELETE("/:id", h.Delete)
		clients.GET("/:id", h.Get)
		clients.PUT("/:id", h.Update)
	}
	r.POST("/products/:id/image", h.Create)
	r.GET("/products/:id/image", h.GetImage)
}

// Create godoc
// @Summary Загрузить изображение для продукта
// @Description Загружает новое изображение для указанного продукта
// @Tags images
// @Accept application/octet-stream
// @Produce json
// @Param id path string true "UUID продукта" format(uuid)
// @Param image body string true "Бинарные данные изображения" binary
// @Success 201 {object} dto.ImageResponse "Изображение успешно загружено"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID продукта или некорректные данные изображения"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при сохранении изображения"
// @Router /products/{id}/image [post]
func (h *ImageHandler) Create(c *gin.Context) {
	prodyctId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product ID"})
		return
	}

	const maxImageSize = 5 << 20 //5Mb

	imageBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, maxImageSize))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "image too large or invalid"})
		return
	}
	defer c.Request.Body.Close()
	if len(imageBytes) == maxImageSize {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "image size too large"})
		return
	}
	if len(imageBytes) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "empty image"})
		return
	}

	imgID, err := h.service.CreateImage(c.Request.Context(), imageBytes, prodyctId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, dto.ImageResponse{ImageID: imgID})

}

// Delete godoc
// @Summary Удалить изображение
// @Description Удаляет изображение по уникальному идентификатору UUID
// @Tags images
// @Param id path string true "UUID изображения" format(uuid)
// @Success 204 "Изображение успешно удалено"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID изображения"
// @Failure 404 {object} dto.NotFoundErrorResponse "Изображение не найдено"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при удалении"
// @Router /images/{id} [delete]
func (h *ImageHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	err = h.service.DeleteImage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "image not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary Получить изображение по ID
// @Description Возвращает бинарные данные изображения по уникальному идентификатору
// @Tags images
// @Produce application/octet-stream
// @Param id path string true "UUID изображения" format(uuid)
// @Success 200 "Бинарные данные изображения"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID изображения"
// @Failure 404 {object} dto.NotFoundErrorResponse "Изображение не найдено"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении изображения"
// @Router /images/{id} [get]
func (h *ImageHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	img, err := h.service.GetImage(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrImageNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch image"})
		}
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", `attachment; filename="product_image.bin"`)

	c.Data(http.StatusOK, "", img.Image)
}

// GetImage godoc
// @Summary Получить изображение продукта
// @Description Возвращает изображение продукта по его уникальному идентификатору
// @Tags images
// @Produce application/octet-stream
// @Param id path string true "UUID продукта" format(uuid)
// @Success 200 "Бинарные данные изображения"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID продукта"
// @Failure 404 {object} dto.NotFoundErrorResponse "Изображение не найдено"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении изображения"
// @Router /products/{id}/image [get]
func (h *ImageHandler) GetImage(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}
	img, err := h.service.GetImageByProduct(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, model.ErrImageNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch image"})
		}
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", `attachment; filename="product_image.bin"`)

	c.Data(http.StatusOK, "", img.Image)
}

// Update godoc
// @Summary Обновить изображение
// @Description Обновляет существующее изображение новыми бинарными данными
// @Tags images
// @Accept application/octet-stream
// @Param id path string true "UUID изображения" format(uuid)
// @Param image body string true "Бинарные данные изображения" binary
// @Success 204 "Изображение успешно обновлено"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID или некорректные данные изображения"
// @Failure 404 {object} dto.NotFoundErrorResponse "Изображение не найдено"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при обновлении"
// @Router /images/{id} [put]
func (h *ImageHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product ID"})
		return
	}

	const maxImageSize = 5 << 20 //5Mb

	imageBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, maxImageSize))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "image too large or invalid"})
		return
	}
	defer c.Request.Body.Close()
	if len(imageBytes) == maxImageSize {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "image size too large"})
		return
	}
	if len(imageBytes) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "empty image"})
		return
	}

	err = h.service.UpdateImage(c.Request.Context(), id, imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create product"})
		return
	}
	c.Status(http.StatusNoContent)
}
