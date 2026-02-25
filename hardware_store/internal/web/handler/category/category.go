package category

import (
	service "hardware_store/internal/service/category"
	"hardware_store/internal/web/dto"
	"hardware_store/internal/web/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	validator *validator.Validate
	service   service.CategoryService
}

func NewCategoryHandler(validator *validator.Validate,
	service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		validator: validator,
		service:   service,
	}
}
func (h *CategoryHandler) Register(r *gin.RouterGroup) {
	category := r.Group("/categories")
	{
		category.POST("", h.Create)
		category.DELETE("/:id", h.Delete)
		category.GET("/:id", h.Get)
		category.GET("", h.List)
		category.POST("/:id", h.Update)

	}
}

// Create godoc
// @Summary Создать новую категорию
// @Description Создаёт новую категорию в системе на основе переданных данных
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dto.CategoryRequest true "Данные категории для создания"
// @Success 200 {object} dto.CategoryResponse "Категория успешно создана"
// @Failure 400 {object} dto.ValidationErrorResponse "Ошибки валидации полей или некорректный формат запроса"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при сохранении категории"
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	newID := uuid.New()
	categor := mapper.CategoryRequestToDomain(req, newID)
	err := h.service.CreateCategory(c.Request.Context(), categor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create category"})
		return
	}

	c.JSON(http.StatusOK, mapper.CategoryDomainToWeb(categor))
}

// Delete godoc
// @Summary Удалить категорию
// @Description Удаляет категорию по уникальному идентификатору UUID
// @Tags categories
// @Param id path string true "UUID категории" format(uuid)
// @Success 204 "Категория успешно удалена"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Категория не найдена"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при удалении"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	err = h.service.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary Получить категорию
// @Description Возвращает полные данные категории по уникальному идентификатору
// @Tags categories
// @Produce json
// @Param id path string true "UUID категории" format(uuid)
// @Success 200 {object} dto.CategoryResponse "Категория успешно получена"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Категория не найдена"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении категории"
// @Router /categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	cat, err := h.service.GetCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
		return
	}
	c.JSON(http.StatusOK, mapper.CategoryDomainToWeb(cat))
}

// List godoc
// @Summary Получить список категорий
// @Description Возвращает список всех категорий в системе
// @Tags categories
// @Produce json
// @Success 200 {array} dto.CategoryResponse "Список категорий успешно получен"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении списка"
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	cat, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
		return
	}
	var categories []dto.CategoryResponse
	for _, category := range cat {
		categories = append(categories, mapper.CategoryDomainToWeb(category))
	}
	c.JSON(http.StatusOK, categories)
}

// Update godoc
// @Summary Обновить категорию
// @Description Обновляет данные категории по уникальному идентификатору
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "UUID категории" format(uuid)
// @Param category body dto.CategoryRequest true "Обновлённые данные категории"
// @Success 200 {object} dto.CategoryResponse "Категория успешно обновлена"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный UUID или ошибки валидации данных"
// @Failure 404 {object} dto.NotFoundErrorResponse "Категория не найдена"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при обновлении"
// @Router /categories/{id} [post]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}
	var req dto.CategoryRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	cat, err := h.service.UpdateCategory(c.Request.Context(), mapper.CategoryRequestToDomain(req, id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "category not found"})
		return
	}
	c.JSON(http.StatusOK, mapper.CategoryDomainToWeb(cat))
}
