package handler

import (
	"fmt"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/client"
	addressservice "hardware_store/internal/service/address_service"
	clientservice "hardware_store/internal/service/client_service"
	"hardware_store/internal/web/dto"
	"hardware_store/internal/web/mapper"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ClientHandler struct {
	validator *validator.Validate
	service   clientservice.ClientService
	address   addressservice.AddressService
}

func NewClientHandler(validator *validator.Validate, service clientservice.ClientService, address addressservice.AddressService) *ClientHandler {
	return &ClientHandler{validator: validator, service: service, address: address}
}

func (h *ClientHandler) Register(r *gin.RouterGroup) {
	clients := r.Group("/clients")
	{
		clients.POST("", h.Create)
		clients.DELETE("/:id", h.Delete)
		clients.GET("/search", h.Get)
		clients.PUT("/:id", h.Update)
		clients.GET("", h.List)
	}
}

// Create godoc
// @Summary Создание клиента
// @Description Создаёт нового клиента вместе с адресом
// @Tags clients
// @Accept json
// @Produce json
// @Param client body dto.ClientRequest true "Данные клиента"
// @Success 201 {object} dto.ClientResponse "Клиент успешно создан"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный запрос"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /clients [post]
func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.ClientRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid JSON format: " + err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "validation error: " + err.Error()})
		return
	}
	clientID := uuid.New()
	addrID := uuid.New()
	dateRegistration := time.Now()
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "birthday must be in format YYYY-MM-DD",
		})
		return
	}
	cl := client.Client{
		ClientID:         clientID,
		Name:             req.Name,
		Surname:          req.Surname,
		Birthday:         birthday,
		Gender:           req.Gender,
		RegistrationDate: dateRegistration,
		AddressID:        addrID,
	}
	addr := address.Address{
		AddressID: addrID,
		Country:   req.Address.Country,
		City:      req.Address.City,
		Street:    req.Address.Street}

	err = h.service.CreateClient(c.Request.Context(), cl, addr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "failed to create client",
		})
		return
	}

	c.JSON(http.StatusCreated, mapper.ClientDomainToWeb(cl))
}

// Delete godoc
// @Summary Удаление клиента
// @Description Удаляет клиента по его уникальному идентификатору
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "UUID клиента" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Клиент не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /clients/{id} [delete]
func (h *ClientHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid UUID format"})
		return
	}
	err = h.service.DeleteClient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "client not found"})
	}
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary Получить клиента
// @Description Возвращает информацию о клиенте по имени и фамилии
// @Tags clients
// @Accept json
// @Produce json
// @Param name query string true "Имя клиента"
// @Param surname query string true "Фамилия клиента"
// @Success 200 {object} dto.ClientResponse "Успешно"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Клиент не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /clients/search [get]
func (h *ClientHandler) Get(c *gin.Context) {
	name := c.Query("name")
	surname := c.Query("surname")

	if name == "" || surname == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "name and surname are required"})
		return
	}
	getClient, err := h.service.GetClient(c.Request.Context(), name, surname)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "client not found"})
		return
	}

	c.JSON(http.StatusOK, mapper.ClientDomainToWeb(getClient))
}

// Update godoc
// @Summary Обновить клиента
// @Description Обновляет адрес клиента по его уникальному идентификатору
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "UUID клиента" format(uuid)
// @Param address body dto.AddressRequest true "Данные адреса для обновления"
// @Success 200 "Адрес успешно обновлён"
// @Failure 400 {object} dto.ValidationErrorResponse "Невалидный формат UUID"
// @Failure 404 {object} dto.NotFoundErrorResponse "Клиент не найден"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /clients/{id} [put]
func (h *ClientHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid UUID format"})
		return
	}

	var req dto.AddressRequest
	if err = c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}

	if err = h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "validation error: " + err.Error()})
		return
	}
	addr := address.Address{
		Country: req.Country,
		City:    req.City,
		Street:  req.Street,
	}
	err = h.service.UpdateAddressClient(c.Request.Context(), id, addr)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "client not found"})
		return
	}
	c.Status(http.StatusOK)
}

// List godoc
// @Summary Получить список всех клиентов
// @Description Возвращает список всех клиентов в системе
// @Tags clients
// @Accept json
// @Produce json
// @Success 200 {array} dto.ClientResponse "Список клиентов"
// @Failure 500 {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router /clients [get]
func (h *ClientHandler) List(c *gin.Context) {
	limit, err := parsePaginationParam(c.Query("limit"), "limit", 100)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	offset, err := parsePaginationParam(c.Query("offset"), "offset", 100)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var res []dto.ClientResponse
	clients, err := h.service.GetClients(c.Request.Context(), int(limit), offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch clients"})
		return
	}
	for _, client := range clients {
		res = append(res, mapper.ClientDomainToWeb(client))
	}
	c.JSON(http.StatusOK, res)
}

func parsePaginationParam(s string, paramName string, max int) (int, error) {
	if s == "" {
		return 0, nil
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: must be an integer", paramName)
	}

	if n < 0 {
		return 0, fmt.Errorf("%s must be >= 0", paramName)
	}

	if max > 0 && n > max {
		return 0, fmt.Errorf("%s must be <= %d", paramName, max)
	}

	return n, nil
}
