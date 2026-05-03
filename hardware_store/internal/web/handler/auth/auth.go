package auth

import (
	"context"
	"hardware_store/internal/config"
	"hardware_store/internal/web/dto"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	auth  AuthService
	user  UserService
	oauth GoogleOAuthService
	log   *slog.Logger
	cfg *config.Config
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, refreshTokenOut string, err error)
	Logout(ctx context.Context, userID, appID string) error
}

type UserService interface {
	RegisterNewUser(ctx context.Context, email, password, name, surname, phoneNumber string) (userID string, err error)
	ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error
	RestorePassword(ctx context.Context, email string) error
}

type TokenService interface {
	ValidateToken(ctx context.Context, token string) (string, string, error)
}

type GoogleOAuthService interface {
	GetGoogleAuthURL(ctx context.Context, state string) (authURL string, err error)
	GoogleCallback(ctx context.Context, code, state, appID string) (accessToken, refreshToken string, err error)
}

func NewAuthHandler(auth AuthService, user UserService, oauth GoogleOAuthService, log *slog.Logger, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		auth:  auth,
		user:  user,
		oauth: oauth,
		log:   log,
		cfg: cfg,
	}
}

func (a *AuthHandler) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", a.RegisterNewUser)
		auth.POST("", a.Login)
		auth.GET("/oauth/google", a.GoogleLoginRedirect)
		auth.GET("/google/callback", a.GoogleCallback)
		auth.POST("/refresh", a.Refresh)
		auth.PATCH("/reset", a.ChangePassword)
		auth.POST("/restore", a.RestorePassword)
	}
}

// GoogleLoginRedirect godoc
// @Summary      Начало авторизации через Google
// @Description  Возвращает URL для редиректа пользователя на сервер авторизации Google.
// @Tags         auth
// @Produce      json
// @Param        state query string false "CSRF state parameter"
// @Success      200  {object}  dto.GoogleAuthURLResponse  "URL для редиректа на Google"
// @Failure      400  {object}  dto.ErrorResponse  "Неверный запрос"
// @Failure      500  {object}  dto.InternalErrorResponse  "Ошибка сервера"
// @Router       /auth/oauth/google [get]
func (a *AuthHandler) GoogleLoginRedirect(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		state = "test_" + uuid.New().String()[:8] 
	}

	authURL, err := a.oauth.GetGoogleAuthURL(c.Request.Context(), state)
	if err != nil {
		a.log.Error("Failed to get Google auth URL", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to initiate Google auth"})
		return
	}

	// Вариант А: Вернуть URL в ответе (для Postman / SPA)
	c.JSON(http.StatusOK, dto.GoogleAuthURLResponse{AuthURL: authURL})

	// Сразу редиректить пользователя (для браузера)
	// c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback godoc
// @Summary      Завершение авторизации через Google
// @Description  Обрабатывает callback от Google, обменивает code на JWT-токены.
// @Tags         auth
// @Produce      json
// @Param        code  query  string  true  "Authorization code from Google"
// @Param        state query  string  true  "State parameter (must match)"
// @Success      200  {object}  dto.LoginResponse  "JWT tokens"
// @Failure      400  {object}  dto.ErrorResponse  "Missing code or state"
// @Failure      401  {object}  dto.ErrorResponse  "Invalid code or state"
// @Failure      500  {object}  dto.InternalErrorResponse  "Internal error"
// @Router       /auth/google/callback [get]
func (a *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "code and state are required"})
		return
	}
	appID := a.cfg.AppID

	accessToken, refreshToken, err := a.oauth.GoogleCallback(c.Request.Context(), code, state, appID)
	if err != nil {
		a.log.Error("Google callback failed", slog.Any("error", err))
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "authentication failed"})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	// Редирект на фронтенд с токеном в хеше (для браузера)
	// c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(
	//     "http://localhost:3000/auth/success#access_token=%s&refresh_token=%s",
	//     accessToken, refreshToken,
	// ))
}

// RegisterNewUser godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает новую учетную запись пользователя в системе. Возвращает ID созданного пользователя.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.RegisterRequest true "Данные для регистрации"
// @Success      201      {object} dto.RegisterResponse "Пользователь успешно зарегистрирован"
// @Failure      400      {object} dto.ErrorResponse "Неверный формат запроса или данные не прошли валидацию"
// @Failure      409      {object} dto.ErrorResponse "Пользователь с таким email уже существует"
// @Failure      500      {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (a *AuthHandler) RegisterNewUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}

	userID, err := a.user.RegisterNewUser(c.Request.Context(), req.Email, req.Password, req.Name, req.Surname, req.PhoneNumber)
	if err != nil {
		handleGrpcError(c, err, "registration failed")
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{UserID: userID})
}

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Проверяет учетные данные и возвращает пару токенов (Access и Refresh).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.LoginRequest true "Учетные данные (email и пароль)"
// @Success      200      {object} dto.LoginResponse "Успешный вход. Токены возвращены."
// @Failure      400      {object} dto.ErrorResponse "Неверный формат запроса"
// @Failure      401      {object} dto.ValidationErrorResponse "Неверный email или пароль"
// @Failure      500      {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router       /auth [post]
func (a *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}

	accessToken, refreshToken, err := a.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		handleGrpcError(c, err, "login failed")
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken})
}

// Refresh godoc
// @Summary      Обновление токенов
// @Description  Принимает refresh токен и возвращает новую пару токенов (Access и Refresh).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.RefreshRequest true "Refresh токен"
// @Success      200      {object} dto.LoginResponse "Новая пара токенов"
// @Failure      400      {object} dto.ErrorResponse "Неверный формат запроса"
// @Failure      401      {object} dto.ValidationErrorResponse "Неверный или протухший refresh токен"
// @Failure      500      {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router       /auth/refresh [post]
func (a *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}

	accessToken, refreshToken, err := a.auth.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleGrpcError(c, err, "refresh failed")
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary      Выход из системы (Logout)
// @Description  Завершает сеанс пользователя. Для выхода достаточно наличия валидного Access токена в заголовке Authorization. Сервис инвалидирует сессию для текущего пользователя и приложения.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]string "Успешный выход. Сессия завершена."
// @Failure      401 {object} map[string]string "Неверный или отсутствующий Access токен"
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /auth/logout [post]
func (a *AuthHandler) Logout(c *gin.Context) {

	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not in context"})
		return
	}
	appVal, ok := c.Get("app_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "app_id not in context"})
		return
	}
	userID, _ := uidVal.(string)
	appID, _ := appVal.(string)
	err := a.auth.Logout(c.Request.Context(), userID, appID)
	if err != nil {
		a.log.Warn("Logout warning", slog.Any("error", err))
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// ChangePassword godoc
// @Summary      Смена пароля
// @Description  Изменяет пароль пользователя. Требует указания старого пароля для подтверждения.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.ChangePasswordRequest true "Email, старый и новый пароль"
// @Security     BearerAuth
// @Success      200      {object} map[string]string "Пароль успешно изменен"
// @Failure      400      {object} dto.ErrorResponse "Неверный формат запроса"
// @Failure      401      {object} dto.ValidationErrorResponse "Неверный старый пароль"
// @Failure      404      {object} dto.ErrorResponse "Пользователь не найден"
// @Failure      500      {object} dto.InternalErrorResponse "Внутренняя ошибка сервера"
// @Router       /auth/reset [patch]
func (a *AuthHandler) ChangePassword(c *gin.Context) {

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}

	err := a.user.ChangePassword(c.Request.Context(), req.Email, req.OldPassword, req.NewPassword)
	if err != nil {
		a.log.Error("Change password failed", slog.Any("error", err), slog.String("email", req.Email))
		handleGrpcError(c, err, "failed to change password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// RestorePassword godoc
// @Summary      Восстановление пароля
// @Description  Генерирует новый пароль и отправляет его на email пользователя. Используется, если пользователь забыл свой пароль.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body     dto.RestorePasswordRequest true "Email пользователя для восстановления"
// @Success      200      {object} map[string]string "Инструкции по восстановлению отправлены на email"
// @Failure      400      {object} dto.ValidationErrorResponse "Неверный формат email"
// @Failure      404      {object} dto.ErrorResponse "Пользователь не найден"
// @Failure      500      {object} dto.InternalErrorResponse "Ошибка отправки письма или генерации пароля"
// @Router       /auth/restore [post]
func (a *AuthHandler) RestorePassword(c *gin.Context) {

	var req dto.RestorePasswordRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON format" + err.Error()})
		return
	}

	err := a.user.RestorePassword(c.Request.Context(), req.Email)
	if err != nil {
		a.log.Error("Restore password failed", slog.Any("error", err), slog.String("email", req.Email))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If an account with that email exists, a new password has been sent.",
	})
}

func handleGrpcError(c *gin.Context, err error, defaultMsg string) {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: st.Message()})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found"})
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid credentials"})
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "resource already exists"})
		case codes.PermissionDenied:
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: defaultMsg})
		}
		return
	}

	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: defaultMsg})
}
