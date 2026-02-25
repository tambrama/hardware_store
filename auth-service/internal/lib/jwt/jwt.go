package jwt

import (
	"auth-service/internal/config"
	"auth-service/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL   = 60 * time.Minute   // короткий срок жизни
	RefreshTokenTTL  = 7 * 24 * time.Hour // 7 дней
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type iwtProvider struct {
	secretKey []byte
}

func NewJWTProvider(cfg *config.Config) *iwtProvider {
	return &iwtProvider{
		secretKey: []byte(cfg.JWTSecretKey),
	}
}

func (p *iwtProvider) generateToken(user model.Users, app model.App, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Mail,
		"exp":     time.Now().Add(duration),
		"app_id":  app.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(p.secretKey)
}

func (p *iwtProvider) NewAccessToken(user model.Users, app model.App) (string, error) {
	return p.generateToken(user, app, AccessTokenTTL)
}

func (p *iwtProvider) NewRefreshToken(user model.Users, app model.App) (string, error) {
	return p.generateToken(user, app, RefreshTokenTTL)
}

func (p *iwtProvider) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return p.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
