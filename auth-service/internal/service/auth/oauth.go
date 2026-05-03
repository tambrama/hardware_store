package auth

import (
	"auth-service/internal/config"
	"auth-service/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthService struct {
	log           *slog.Logger
	config        *oauth2.Config
	jwtSecret     []byte
	jwtExpiration time.Duration
}

func NewOAuthService(log *slog.Logger, cfg *config.Config) *OAuthService {
	return &OAuthService{
		log: log,
		config: &oauth2.Config{
			ClientID:     cfg.OAuth.ClientID,
			ClientSecret: cfg.OAuth.ClientSecret,
			RedirectURL:  cfg.OAuth.RedirectURL,
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		},
		jwtSecret:     []byte(cfg.JWTSecretKey),
		jwtExpiration: cfg.OAuth.JwtExpiration,
	}
}

func (s *OAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
}

func (s *OAuthService) ExchangeCode(ctx context.Context, code string) (*model.GoogleUserInfo, *oauth2.Token, error) {
	requestID := uuid.New().String()[:8]
	const op = "OAuth.ExchangeCode"

	// 👇 Логируем входные данные (без секретов!)
	s.log.Debug("Exchanging code for token",
		slog.String("code_prefix", code[:10]+"..."),
		slog.String("redirect_uri", s.config.RedirectURL),
		slog.Bool("client_secret_empty", s.config.ClientSecret == ""),
	)
	fmt.Printf("🔥 [%s] Calling Exchange...\n", requestID)
	token, err := s.config.Exchange(ctx, code)
	fmt.Printf("🔥 [%s] Exchange returned: err=%v\n", requestID, err)
	if err != nil {
		// 👇 Логируем ПОЛНУЮ ошибку от Google
		s.log.Error("Google token exchange failed",
			slog.String("operation", op),
			slog.Any("error", err),
			slog.String("error_type", fmt.Sprintf("%T", err)),
			slog.String("error_text", err.Error()),
		)
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	fmt.Printf("🔥 [%s] Got token, creating client...\n", requestID)
	
	client := s.config.Client(ctx, token)
	
	fmt.Printf("🔥 [%s] Client created, fetching userinfo...\n", requestID)
	
	 resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	  fmt.Printf("🔥 [%s] UserInfo GET returned, err=%v, status=%v\n", 
        requestID, err, 
        func() string { if resp != nil { return resp.Status }; return "nil" }(),
    )

	if err != nil {
		   fmt.Printf("🔥 [%s] UserInfo request failed: %v\n", requestID, err)
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo model.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		 fmt.Printf("🔥 [%s] JSON decode failed: %v\n", requestID, err)
		return nil, nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	  fmt.Printf("🔥 [%s] COMPLETE! google_id=%s, email=%s\n", 
        requestID, userInfo.GoogleID, userInfo.Mail)
		
	return &userInfo, token, nil
}
