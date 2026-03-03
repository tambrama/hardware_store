package jwt

import (
	"auth-service/internal/config"
	"auth-service/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenTTL   = 60 * time.Minute   // короткий срок жизни
	RefreshTokenTTL  = 7 * 24 * time.Hour // 7 дней
	Issuer           = "auth-service"
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

func (p *iwtProvider) generateToken(user, app uuid.UUID, duration time.Duration) (string, error) {
	claims := model.CustomClaims{
		UserID: user,
		AppID:  app,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(p.secretKey)
}

func (p *iwtProvider) NewAccessToken(user, app uuid.UUID) (string, error) {
	return p.generateToken(user, app, AccessTokenTTL)
}

func (p *iwtProvider) NewRefreshToken(user, app uuid.UUID) (string, error) {
	return p.generateToken(user, app, RefreshTokenTTL)
}

func (p *iwtProvider) ValidateToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return p.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok || !token.Valid {
		return nil, model.ErrInvalidToken
	}
	return claims, nil
}
