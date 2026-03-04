package auth

import (
	hashtoken "auth-service/internal/lib/hash_token"
	jwtMe "auth-service/internal/lib/jwt"
	"auth-service/internal/lib/logger/sl"
	"auth-service/internal/model"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	log         *slog.Logger
	storageU    UserStorage
	storageA    AppStorage
	storageJ    JWTStorage
	jwtProvider JWTProvider
}

type UserStorage interface {
	SaveUser(ctx context.Context, user model.Users) error
	GetUserByEmail(ctx context.Context, email string) (model.Users, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (model.Users, error)
	UpdateUserPassword(ctx context.Context, email string, newPassword []byte) error
}
type AppStorage interface {
	GetAppByID(ctx context.Context, appID string) (model.App, error)
}

type JWTStorage interface {
	SaveSession(ctx context.Context, session model.Session) error
	GetSession(ctx context.Context, tokenHash string) (model.Session, error)
	DeleteSession(ctx context.Context, userID uuid.UUID, appID string) error
	UpdateSession(ctx context.Context, oldToken, newToken string, duration time.Time) error
}

type JWTProvider interface {
	NewAccessToken(user, app uuid.UUID) (string, error)
	NewRefreshToken(user, app uuid.UUID) (string, error)
	ValidateToken(tokenString string) (*model.CustomClaims, error)
}

func NewAuth(log *slog.Logger, storageU UserStorage, storageA AppStorage, storageJ JWTStorage, jwtProvider JWTProvider) *auth {
	return &auth{
		log:         log,
		storageU:    storageU,
		storageA:    storageA,
		storageJ:    storageJ,
		jwtProvider: jwtProvider,
	}
}

func (a *auth) Login(ctx context.Context, email, password, appID string) (accessToken, refreshToken string, err error) {
	const op = "Auth.Login"
	log := a.log.With(slog.String("operation", op), slog.String("email", email), slog.String("appID", appID))

	log.Info("Attempting to log in")

	user, err := a.storageU.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			a.log.Error("user not found", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		a.log.Error("invalid credentials", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
	}
	app, err := a.storageA.GetAppByID(ctx, appID)
	if err != nil {
		if errors.Is(err, model.ErrAppNotFound) {
			return "", "", fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
		}
		a.log.Error("failed to get app", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err = a.jwtProvider.NewAccessToken(user.ID, app.ID)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	refreshToken, err = a.jwtProvider.NewRefreshToken(user.ID, app.ID)
	if err != nil {
		a.log.Error("failed to generate refresh token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	session := model.Session{
		SessionID: uuid.New(),
		UserID:    user.ID,
		AppID:     app.ID,
		HashToken: hashtoken.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(jwtMe.RefreshTokenTTL),
	}
	err = a.storageJ.SaveSession(ctx, session)
	if err != nil {
		a.log.Error("failed to save refresh token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}
func (a *auth) RegisterNewUser(ctx context.Context, email, password, name, surname, phoneNumber string) (userID string, err error) {
	const op = "Auth.RegisterNewUser"
	log := a.log.With(slog.String("operation", op), slog.String("email", email))
	log.Info("Attempting to register new user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", slog.Any("error", err))
		return "", err
	}
	user := model.Users{
		ID:          uuid.New(),
		Name:        name,
		Surname:     surname,
		Mail:        email,
		PhoneNumber: phoneNumber,
		Password:    passHash,
	}
	err = a.storageU.SaveUser(ctx, user)
	if err != nil {
		log.Error("Failed to save user", slog.Any("error", err))
		return "", err
	}
	return user.ID.String(), nil
}

func (a *auth) IsAdmin(ctx context.Context, userID uuid.UUID) (isAdmin bool, err error) {
	return false, nil
}

func (a *auth) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error {
	const op = "Auth.ChangePassword"
	log := a.log.With(slog.String("operation", op), slog.String("email", email))
	log.Info("Attempting to change password")

	user, err := a.storageU.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			a.log.Error("user not found", sl.Err(err))
			return fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(oldPassword)); err != nil {
		a.log.Error("invalid credentials", sl.Err(err))
		return fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", slog.Any("error", err))
		return err
	}
	err = a.storageU.UpdateUserPassword(ctx, email, passHash)
	if err != nil {
		log.Error("Failed to update password", slog.Any("error", err))
		return err
	}
	return nil
}

func (a *auth) RestorePassword(ctx context.Context, email string) error {
	return nil
}

func (a *auth) RefreshToken(ctx context.Context, refreshToken, appId string) (string, string, error) {
	const op = "Auth.RefreshToken"
	log := a.log.With(slog.String("operation", op), slog.String("appID", appId))
	log.Info("Attempting to refresh token")

	sessionHash, err := a.storageJ.GetSession(ctx, hashtoken.HashToken(refreshToken))
	if err != nil {
		if errors.Is(err, model.ErrRefreshTokenNotFound) {
			log.Info("Refresh token not found or expired", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
		}
		log.Error("Failed to get session", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if sessionHash.AppID.String() != appId {
		log.Info("Session does not belong to the specified app", slog.String("appID", appId))
		return "", "", fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
	}
	// Генерируем новый access и refresh токены
	newAccessToken, err := a.jwtProvider.NewAccessToken(sessionHash.UserID, sessionHash.AppID)
	if err != nil {
		log.Error("Failed to generate new access token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	newRefreshToken, err := a.jwtProvider.NewRefreshToken(sessionHash.UserID, sessionHash.AppID)
	if err != nil {
		log.Error("Failed to generate new refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.storageJ.UpdateSession(ctx, hashtoken.HashToken(refreshToken), hashtoken.HashToken(newRefreshToken), time.Now().Add(jwtMe.RefreshTokenTTL))
	if err != nil {
		if errors.Is(err, model.ErrRefreshTokenNotFound) {
			log.Info("Old refresh token not found or expired during update", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, model.ErrInvalidCredentials)
		}
		log.Error("Failed to update session", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return newAccessToken, newRefreshToken, nil
}

func (a *auth) Logout(ctx context.Context, token string) error {
	const op = "Auth.Logout"
	log := a.log.With(slog.String("operation", op))
	log.Info("Attempting to log out")
	claims, err := a.jwtProvider.ValidateToken(token)
	if err != nil {
		log.Error("Failed to validate token", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, model.ErrInvalidToken)
	}

	err = a.storageJ.DeleteSession(ctx, claims.UserID, claims.AppID.String())
	if err != nil {
		if errors.Is(err, model.ErrRefreshTokenNotFound) {
			log.Info("No active session found for user, nothing to delete", slog.Any("userID", claims.UserID), slog.Any("appID", claims.AppID))
			return model.ErrRefreshTokenNotFound
		}
		log.Error("Failed to delete session", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *auth) ValidateToken(ctx context.Context, token string) (*model.CustomClaims, error) {
	return a.jwtProvider.ValidateToken(token)
}