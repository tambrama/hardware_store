package app

import (
	grpcapp "auth-service/internal/app/grpc"
	"context"
	"log/slog"

	"go.uber.org/fx"
)

func NewApp(lc fx.Lifecycle, log *slog.Logger, server *grpcapp.App) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting application lifecycle")
			go func() {
				if err := server.Run(); err != nil {
					log.Error("gRPC server error", slog.Any("error", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("shutting down server")
			server.Stop()
			return nil
		},
	})
}
