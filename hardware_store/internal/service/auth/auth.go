package auth

import (
	"context"
	"fmt"
	"hardware_store/internal/clients/auth/grpc"
	"hardware_store/internal/config"
	"log/slog"
)

type Auth struct {
	client *grpc.Client
	log    *slog.Logger
	appID  string
}

func NewAuth(client *grpc.Client, log *slog.Logger, cfg *config.Config) *Auth {
	return &Auth{
		client: client,
		log:    log,
		appID:  cfg.AppID,
	}
}

func (a *Auth) GetGoogleAuthURL(ctx context.Context, state string) (string, error) {
	const op = "auth.Service.GetGoogleAuthURL"

	resp, err := a.client.GetGoogleAuthURL(ctx, state)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.AuthUrl, nil
}

// GoogleCallback проксирует callback в auth-service через gRPC
func (a *Auth) GoogleCallback(ctx context.Context, code, state, appID string) (string, string, error) {
	const op = "auth.Service.GoogleCallback"

	resp, err := a.client.GoogleCallback(ctx, code, state, appID)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.AccessToken, resp.RefreshToken, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password, name, surname, phoneNumber string) (string, error) {
	const op = "auth.Register.Store"
	req, err := a.client.Register(ctx, email, password, name, surname, phoneNumber)
	if err != nil {
		a.log.Error("Register call failed",
			slog.String("op", op),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return req.UserId, nil
}

func (a *Auth) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	const op = "auth.Login.Store"
	req, err := a.client.Login(ctx, email, password, a.appID)
	if err != nil {
		a.log.Error("Login call failed",
			slog.String("op", op),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return req.AccessToken, req.RefreshToken, nil
}

func (a *Auth) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) error {
	const op = "auth.ChangePassword.Store"
	_, err := a.client.ChangePassword(ctx, email, oldPassword, newPassword)
	if err != nil {
		a.log.Error("ChangePassword call failed",
			slog.String("op", op),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *Auth) RestorePassword(ctx context.Context, email string) error {
	const op = "auth.RestorePassword.Store"
	_, err := a.client.RestorePassword(ctx, email)
	if err != nil {
		a.log.Error("RestorePassword call failed",
			slog.String("op", op),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID string) (isAdmin bool, err error) {
	const op = "auth.IsAdmin.Store"
	req, err := a.client.IsAdmin(ctx, userID)
	if err != nil {
		a.log.Error("IsAdmin call failed",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return req.IsAdmin, nil
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "auth.RefreshToken.Store"
	req, err := a.client.Refresh(ctx, refreshToken, a.appID)
	if err != nil {
		a.log.Error("RefreshToken call failed",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return req.AccessToken, req.RefreshToken, nil
}

func (a *Auth) Logout(ctx context.Context, userID, appID string) error {
	const op = "auth.Logout.Store"
	_, err := a.client.Logout(ctx, userID, appID)
	if err != nil {
		a.log.Error("Logout call failed",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Auth) ValidateToken(ctx context.Context, token string) (string, string, error) {
	const op = "auth.ValidateToken.Store"
	req, err := a.client.Validate(ctx, token)
	if err != nil {
		a.log.Error("ValidateToken call failed",
			slog.String("op", op),
			slog.Any("error", err),
		)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return req.AppId, req.UserId, nil
}
