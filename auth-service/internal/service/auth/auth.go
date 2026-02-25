package auth

import (
	"auth-service/internal/lib/logger/sl"
	"auth-service/internal/model"
	jwtMe"auth-service/internal/lib/jwt"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	log         *slog.Logger
	storageU    UserStorage
	storageA    AppStorage
	storageJ    JWTStorage
	jwtProvider JWTProvider
}

type UserStorage interface {
	SaveUser(ctx context.Context, user model.Users) error
	GetUserByEmail(ctx context.Context, email string) (model.Users, error)
	UpdateUserPassword(ctx context.Context, email string, newPassword []byte) error
}
type AppStorage interface {
	GetAppByID(ctx context.Context, appID string) (model.App, error)
}

type JWTStorage interface {
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, appID string, refreshToken string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, userID uuid.UUID, appID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID uuid.UUID, appID string) error
}

type JWTProvider interface {
	NewAccessToken(user model.Users, app model.App) (string, error)
	NewRefreshToken(user model.Users, app model.App) (string, error)
	// RefreshToken(tokenString string, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewAuth(log *slog.Logger, storageU UserStorage, storageA AppStorage, storageJ JWTStorage, jwtProvider JWTProvider) *Auth {
	return &Auth{
		log:         log,
		storageU:    storageU,
		storageA:    storageA,
		storageJ:    storageJ,
		jwtProvider: jwtProvider,
	}
}

func (a *Auth) Login(ctx context.Context, email, password, appID string) (accessToken, refreshToken string, err error) {
	const op = "Auth.Login"
	log := a.log.With(slog.String("operation", op), slog.String("email", email), slog.String("appID", appID))

	log.Info("Attempting to log in")

	user, err := a.storageU.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			a.log.Error("user not found", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		a.log.Error("invalid credentials", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	app, err := a.storageA.GetAppByID(ctx, appID)
	if err != nil {
		if errors.Is(err, model.ErrAppNotFound) {
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get app", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err = a.jwtProvider.NewAccessToken(user, app)
	if err != nil {
		a.log.Error("failed to generate access token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	refreshToken, err = a.jwtProvider.NewRefreshToken(user, app)
	if err != nil {
		a.log.Error("failed to generate refresh token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	expiresAt := time.Now().Add(jwtMe.RefreshTokenTTL)
	err = a.storageJ.SaveRefreshToken(ctx, user.ID, appID, refreshToken, expiresAt)
	if err != nil {
		a.log.Error("failed to save refresh token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}
func (a *Auth) RegisterNewUser(ctx context.Context, email, password, name, surname, phoneNumber string) (userID string, err error) {
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

func (a *Auth) IsAdmin(ctx context.Context, userID uuid.UUID) (isAdmin bool, err error) {
	return false, nil
}

func (a *Auth) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error {
	const op = "Auth.ChangePassword"
	log := a.log.With(slog.String("operation", op), slog.String("email", email))
	log.Info("Attempting to change password")

	user, err := a.storageU.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			a.log.Error("user not found", sl.Err(err))
			return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(oldPassword)); err != nil {
		a.log.Error("invalid credentials", sl.Err(err))
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
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

func (a *Auth) RestorePassword(ctx context.Context, email string) error {
	return nil
}
