package middleware

import (
	"auth-service/internal/model"
	"context"
	"errors"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtProvider JWTProvider
	service     AuthService
}

type JWTProvider interface {
	ValidateToken(tokenString string) (*model.CustomClaims, error)
}

type AuthService interface {
	RefreshToken(ctx context.Context, refreshToken, appId string) (string, string, error)
}

func NewAuthInterceptor(jwtProvider JWTProvider, service AuthService) *AuthInterceptor {
	return &AuthInterceptor{
		jwtProvider: jwtProvider,
		service:     service,
	}
}

func (i *AuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Here you can add logic to extract and validate JWT from the incoming request metadata
		// For example, you can look for an "Authorization" header and validate the token
		// If the token is valid, you can proceed to call the handler
		// If the token is invalid or missing, you can return an error
		if i.isPublicEndpoint(info.FullMethod) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		accessToken := md.Get("authorization")
		refreshToken := md.Get("x-refresh-token")
		if len(accessToken) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing access token")
		}
		tokenString := strings.TrimPrefix(accessToken[0], "Bearer ")
		token, err := i.jwtProvider.ValidateToken(tokenString)
		if err != nil {
			if !errors.Is(err, jwt.ErrTokenExpired) {
				return nil, status.Error(codes.Unauthenticated, "invalid access token")
			}
			if len(refreshToken) == 0 {
				return nil, status.Error(codes.Unauthenticated, "missing refresh token")
			}
			newAccessToken, newRefreshToken, err := i.service.RefreshToken(ctx, refreshToken[0], token.AppID.String())
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, "failed to refresh session")
			}
			header := metadata.Pairs(
				"x-new-access-token", newAccessToken,
				"x-new-refresh-token", newRefreshToken,
			)

			if err := grpc.SetHeader(ctx, header); err != nil {
				log.Printf("failed to set header: %v", err)
			}
			return handler(ctx, req)
		}
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) isPublicEndpoint(method string) bool {
	publicEndpoints := map[string]bool{
		"/auth.AuthService/Login":           true,
		"/auth.AuthService/Register":        true,
		"/auth.AuthService/RestorePassword": true,
		"/auth.AuthService/ChangePassword":  true,
	}
	for endpoint := range publicEndpoints {
		if endpoint == method {
			return true
		}
	}
	return false
}
