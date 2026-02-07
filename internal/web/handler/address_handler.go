package handler

import (
	"hardware_store/internal/model/address"
	addressservice "hardware_store/internal/service/address_service"
	"hardware_store/internal/web/dto"
	"hardware_store/internal/web/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AddressHandler struct {
	validator *validator.Validate
	service   addressservice.AddressService
}

func NewAddressHandler(validator *validator.Validate, service addressservice.AddressService) *AddressHandler {
	return &AddressHandler{validator: validator, service: service}
}

func (h *AddressHandler) Register(r *gin.RouterGroup) {
	clients := r.Group("/addresses")
	{
		clients.POST("", h.Create)
		clients.DELETE("/:id", h.Delete)
		clients.GET("/:id", h.Get)
	}
}

// Create godoc
// @Summary Создать новый адрес
// @Description Создаёт новый адрес в системе на основе переданных данных
// @Tags addresses
// @Accept json
// @Produce json
// @Param address body dto.AddressRequest true "Данные адреса для создания"
// @Success 201 {object} dto.AddressResponse "Адрес успешно создан"
// @Failure 400 {object} dto.ValidationErrorResponse "Ошибки валидации полей или некорректный формат запроса"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при сохранении адреса"
// @Router /addresses [post]
func (h *AddressHandler) Create(c *gin.Context) {
	var req dto.AddressRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	addrID := uuid.New()

	addr := address.Address{
		AddressID: addrID,
		Country:   req.Country,
		City:      req.City,
		Street:    req.Street}

	err := h.service.CreateAddress(c.Request.Context(), addr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, mapper.AddressDomainToWeb(addr))

}

// Delete godoc
// @Summary Удалить адрес
// @Description Удаляет адрес по уникальному идентификатору UUID
// @Tags addresses
// @Param id path string true "UUID адреса" format(uuid)
// @Success 204 "Адрес успешно удалён"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Адрес не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при удалении"
// @Router /addresses/{id} [delete]
func (h *AddressHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	err = h.service.DeleteAddress(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "client not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary Получить адрес
// @Description Возвращает полные данные адреса по уникальному идентификатору
// @Tags addresses
// @Produce json
// @Param id path string true "UUID адреса" format(uuid)
// @Success 200 {object} dto.AddressResponse "Адрес успешно получен"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Адрес не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера при получении адреса"
// @Router /addresses/{id} [get]
func (h *AddressHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid UUID format"})
		return
	}

	addr, err := h.service.GetAddress(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "client not found"})
		return
	}

	c.JSON(http.StatusOK, mapper.AddressDomainToWeb(addr))
}

// func (h *AddressHandler) Update(c *gin.Context) {

// }

// func (h *AddressHandler) List(c *gin.Context) {

// }
