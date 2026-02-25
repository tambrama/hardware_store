package jwt_test

import (
	"auth-service/internal/config"
	"auth-service/internal/lib/jwt"
	"auth-service/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJWTProvider(t *testing.T) {
	secretKey := "my_secret_key"
	cfg := config.Config{
		JWTSecretKey: secretKey,
	}
	provider := jwt.NewJWTProvider(&cfg)
	user := model.Users{
		ID:      uuid.New(),
		Name:    "Test User",
		Surname: "Test",
		Mail:    "test@mail.ru",
	}
	app := model.App{
		ID:   uuid.New(),
		Name: "Test App",
	}
	accessToken, err := provider.NewAccessToken(user, app)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
}
