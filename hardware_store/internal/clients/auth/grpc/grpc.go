package grpc

import (
	"context"
	"hardware_store/internal/config"
	"log/slog"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	sso "github.com/tambrama/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	api sso.AuthClient
}

func NewClient(cfg *config.Config, log *slog.Logger) (*Client, error) { //колиечство повторов
	const op = "grpc.New"
	addr := cfg.Clients.SSO.Address
	timeout := cfg.Clients.SSO.Timeout
	retriesCount := cfg.Clients.SSO.RetriesCount
	// Опции для интерсептора grpcretry
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}
	cc, err := grpc.DialContext(context.Background(), addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, err
	}
	grpcClient := sso.NewAuthClient(cc)
	return &Client{
		api: grpcClient,
	}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (c *Client) GetGoogleAuthURL(ctx context.Context, state string) (*sso.GetGoogleAuthURLResponse, error) {
    return c.api.GetGoogleAuthURL(ctx, &sso.GetGoogleAuthURLRequest{
        State: state,
    })
}


func (c *Client) GoogleCallback(ctx context.Context, code, state, appID string) (*sso.LoginResponse, error) {
    return c.api.GoogleCallback(ctx, &sso.GoogleCallbackRequest{
        Code:  code,
        State: state,
        AppId: appID,
    })
}

func (c *Client) Register(ctx context.Context, email, password, name, surname, phoneNumber string) (*sso.RegisterResponse, error) {
	return c.api.Register(ctx, &sso.RegisterRequest{
		Email:       email,
		Password:    password,
		Name:        name,
		Surname:     surname,
		PhoneNumber: phoneNumber,
	})
}

func (c *Client) Login(ctx context.Context, email, password, appID string) (*sso.LoginResponse, error) {
	return c.api.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
}

func (c *Client) Refresh(ctx context.Context, refreshToken, appID string) (*sso.LoginResponse, error) {
	return c.api.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: refreshToken,
		AppId:        appID,
	})
}

func (c *Client) Validate(ctx context.Context, token string) (*sso.ValidateResponse, error) {
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return c.api.Validate(ctx, &sso.ValidateRequest{
		Token: token,
	})
}

func (c *Client) Logout(ctx context.Context, userID, appID string) (*sso.Empty, error) {
	return c.api.Logout(ctx, &sso.LogoutRequest{UserId: userID, AppId: appID})
}

func (c *Client) IsAdmin(ctx context.Context, userID string) (*sso.IsAdminResponse, error) {
	return c.api.IsAdmin(ctx, &sso.IsAdminRequest{UserId: userID})
}

func (c *Client) ChangePassword(ctx context.Context, email, oldPassword, newPassword string) (*sso.Empty, error) {
	return c.api.ChangePassword(ctx, &sso.ChangePasswordRequest{
		Email:       email,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})
}

func (c *Client) RestorePassword(ctx context.Context, email string) (*sso.Empty, error) {
	return c.api.RestorePassword(ctx, &sso.RestorePasswordRequest{Email: email})
}
