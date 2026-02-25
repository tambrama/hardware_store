package product

import (
	"errors"
	"hardware_store/internal/logger"
	model "hardware_store/internal/model/error"
	"hardware_store/internal/model/product"
	service "hardware_store/internal/service/product"

	"hardware_store/internal/web/dto"
	"hardware_store/internal/web/mapper"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProductHandler struct {
	validator *validator.Validate
	service   service.ProductService
	logger    *slog.Logger
}

func NewProductHandler(validator *validator.Validate, service service.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{validator: validator, service: service, logger: logger}
}

func (h *ProductHandler) Register(r *gin.RouterGroup) {
	clients := r.Group("/products")
	{
		clients.POST("", h.Create)
		clients.DELETE("/:id", h.Delete)
		clients.GET("/:id", h.Get)
		clients.PUT("/:id/stock", h.Update)
		clients.GET("", h.List)
	}
}

// Create godoc
// @Summary Создать новый продукт
// @Description Создаёт новый продукт в системе на основе переданных данных
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.ProductRequest true "Данные продукта для создания"
// @Success 201 {object} dto.ProductResponse "Продукт успешно создан"
// @Failure 400 {object} dto.ValidationErrorResponse "Ошибки валидации полей или некорректный формат запроса"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при сохранении продукта"
// @Router /products [post]
func (h *ProductHandler) Create(c *gin.Context) {

	var req dto.ProductRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Error("Invalid JSON format", logger.Err(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed",
			logger.Err(err),
			slog.Any("request", req),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	productID := uuid.New()
	//imgID := uuid.New()
	updDate := time.Now()
	product := product.Product{
		ProductID:      productID,
		Name:           req.Name,
		CategoryID:     req.CategoryID,
		Price:          req.Price,
		AvailableStock: req.AvailableStock,
		LastUpdateDate: updDate,
		SupplierID:     req.SupplierID,
	}

	err := h.service.CreateProduct(c.Request.Context(), product)
	if err != nil {
		h.logger.Error("Failed to create product",
			logger.Err(err),
			slog.String("product_id", productID.String()),
			slog.String("name", req.Name),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create product"})
		return
	}
	h.logger.Info("Product created successfully",
		slog.String("product_id", productID.String()),
		slog.String("name", req.Name),
	)

	c.JSON(http.StatusCreated, mapper.ProductDomainToWeb(product))

}

// Delete godoc
// @Summary Удалить продукт
// @Description Удаляет продукт по уникальному идентификатору UUID
// @Tags products
// @Param id path string true "UUID продукта" format(uuid)
// @Success 204 "Продукт успешно удалён"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Продукт не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при удалении"
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	err = h.service.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "client not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary Получить продукт
// @Description Возвращает полные данные продукта по уникальному идентификатору
// @Tags products
// @Produce json
// @Param id path string true "UUID продукта" format(uuid)
// @Success 200 {object} dto.ProductResponse "Продукт успешно получен"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Продукт не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении продукта"
// @Router /products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product not found"})
		return
	}

	c.JSON(http.StatusOK, mapper.ProductDomainToWeb(product))
}

// Update godoc
// @Summary Обновить количество товара на складе
// @Description Обновляет доступное количество товара на складе для указанного продукта
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "UUID продукта" format(uuid)
// @Param stock body dto.UpdateStockCountRequest true "Новое количество товара"
// @Success 200 {object} dto.ProductResponse "Количество товара успешно обновлено"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат запроса, отрицательное количество или недостаточный остаток"
// @Failure 404 {object} dto.NotFoundErrorResponse "Продукт не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при обновлении"
// @Router /products/{id}/stock [put]
func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}
	var req dto.UpdateStockCountRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Error("Invalid JSON format", logger.Err(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}
	product, err := h.service.UpdateProduct(c.Request.Context(), id, req.Amount)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product not found"})
			return
		}
		if errors.Is(err, model.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, model.ErrAmountIsNegative) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "amount must be positive"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update product stock"})
		return
	}
	c.JSON(http.StatusOK, mapper.ProductDomainToWeb(product))
}

// List godoc
// @Summary Получить список продуктов
// @Description Возвращает список всех продуктов в системе
// @Tags products
// @Produce json
// @Success 200 {array} dto.ProductResponse "Список продуктов успешно получен"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении списка"
// @Router /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var res []dto.ProductResponse
	products, err := h.service.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch products"})
		return
	}

	for _, product := range products {
		res = append(res, mapper.ProductDomainToWeb(product))
	}
	c.JSON(http.StatusOK, res)
}
