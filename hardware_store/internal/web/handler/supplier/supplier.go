package supplier

import (
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/supplier"
	service "hardware_store/internal/service/supplier"
	"hardware_store/internal/web/dto"
	"hardware_store/internal/web/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SupplierHandler struct {
	validator *validator.Validate
	service   service.SupplierService
}

func NewSupplierHandler(validator *validator.Validate,
	service service.SupplierService) *SupplierHandler {
	return &SupplierHandler{
		validator: validator,
		service:   service,
	}
}

func (h *SupplierHandler) Register(c *gin.RouterGroup) {
	supplier := c.Group("suppliers")
	{
		supplier.POST("", h.Create)
		supplier.DELETE(":id", h.Delete)
		supplier.GET("/:id", h.Get)
		supplier.GET("", h.List)
		supplier.PUT("/:id", h.Update)
	}
}

// Create godoc
// @Summary Создать нового поставщика
// @Description Создаёт нового поставщика вместе с его адресом в системе
// @Tags suppliers
// @Accept json
// @Produce json
// @Param supplier body dto.SupplierRequest true "Данные поставщика для создания"
// @Success 201 {object} dto.SupplierResponse "Поставщик успешно создан"
// @Failure 400 {object} dto.ValidationErrorResponse "Ошибки валидации полей или некорректный формат запроса"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при сохранении поставщика"
// @Router /suppliers [post]
func (h *SupplierHandler) Create(c *gin.Context) {
	var req dto.SupplierRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	supplierID := uuid.New()
	addrID := uuid.New()
	sup := supplier.Supplier{
		SupplierID:  supplierID,
		Name:        req.Name,
		AddressID:   addrID,
		PhoneNumber: req.PhoneNumber,
	}
	addr := address.Address{
		AddressID: addrID,
		Country:   req.Address.Country,
		City:      req.Address.City,
		Street:    req.Address.Street}
	err := h.service.CreateSupplier(c.Request.Context(), sup, addr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, mapper.SupplierDomainToWeb(sup))
}

// Delete godoc
// @Summary Удалить поставщика
// @Description Удаляет поставщика по уникальному идентификатору UUID
// @Tags suppliers
// @Param id path string true "UUID поставщика" format(uuid)
// @Success 204 "Поставщик успешно удалён"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Поставщик не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при удалении"
// @Router /suppliers/{id} [delete]
func (h *SupplierHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	err = h.service.DeleteSupplier(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "supplier not found"})
	}
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary Получить поставщика
// @Description Возвращает полные данные поставщика по уникальному идентификатору
// @Tags suppliers
// @Produce json
// @Param id path string true "UUID поставщика" format(uuid)
// @Success 200 {object} dto.SupplierResponse "Поставщик успешно получен"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Поставщик не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении поставщика"
// @Router /suppliers/{id} [get]
func (h *SupplierHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}
	sup, err := h.service.GetSupplier(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "supplier not found"})
	}
	c.JSON(http.StatusOK, mapper.SupplierDomainToWeb(sup))
}

// List godoc
// @Summary Получить список поставщиков
// @Description Возвращает список всех поставщиков в системе
// @Tags suppliers
// @Produce json
// @Success 200 {array} dto.SupplierResponse "Список поставщиков успешно получен"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении списка"
// @Router /suppliers [get]
func (h *SupplierHandler) List(c *gin.Context) {
	suppliers, err := h.service.GetSuppliers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "supplier not found"})
	}
	var res []dto.SupplierResponse
	for _, sup := range suppliers {
		res = append(res, mapper.SupplierDomainToWeb(sup))
	}
	c.JSON(http.StatusOK, res)
}

// Update godoc
// @Summary Обновить адрес поставщика
// @Description Обновляет адрес поставщика по его уникальному идентификатору
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "UUID поставщика" format(uuid)
// @Param address body dto.AddressRequest true "Обновлённые данные адреса"
// @Success 200 "Адрес поставщика успешно обновлён"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID или ошибки валидации данных"
// @Failure 404 {object} dto.NotFoundErrorResponse "Поставщик не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при обновлении"
// @Router /suppliers/{id} [put]
func (h *SupplierHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}
	var req dto.AddressRequest
	if err = c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}
	addr := address.Address{
		Country: req.Country,
		City:    req.City,
		Street:  req.Street,
	}
	if err = h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	err = h.service.UpdateAddressSupplier(c.Request.Context(), id, addr)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "client not found"})
		return
	}
	c.Status(http.StatusOK)
}
