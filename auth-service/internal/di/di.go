package di

import (
	"auth-service/internal/app"
	grpcapp "auth-service/internal/app/grpc"
	"auth-service/internal/config"
	"auth-service/internal/lib/jwt"
	"auth-service/internal/service/auth"
	"auth-service/internal/storage/sqlite"
	appstorage "auth-service/internal/storage/sqlite/app"
	jwtstorage "auth-service/internal/storage/sqlite/jwt"
	userstorage "auth-service/internal/storage/sqlite/user"
	web "auth-service/internal/web/grpc/auth"
	"database/sql"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

const (
	envLocal      = "local"
	envProduction = "production"
)

func NewValidator() *validator.Validate {
	return validator.New()
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger
	switch cfg.Env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		setupLogger,
		NewValidator,

		sqlite.NewSQLiteDB,
		fx.Annotate(appstorage.NewStorage, fx.As(new(auth.AppStorage))),
		fx.Annotate(userstorage.NewStorage, fx.As(new(auth.UserStorage))),
		fx.Annotate(jwtstorage.NewStorage, fx.As(new(auth.JWTStorage))),

		fx.Annotate(jwt.NewJWTProvider, fx.As(new(auth.JWTProvider))),

		fx.Annotate(auth.NewAuth, fx.As(new(web.Auth))),
		grpcapp.NewApp,
	),
	fx.Invoke(func(lc fx.Lifecycle, storage *sql.DB, log *slog.Logger) {
		sqlite.AddLifecycle(lc, storage, log)
	}),
	fx.Invoke(app.NewApp),
)
