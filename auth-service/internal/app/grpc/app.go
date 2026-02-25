package grpcapp

import (
	"auth-service/internal/config"
	authgrpc "auth-service/internal/web/grpc/auth"
	"fmt"
	"log/slog"
	"net"

	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, authService authgrpc.Auth, cfg *config.Config, validate *validator.Validate) *App {
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()))
	authgrpc.Register(gRPCServer, authService, validate)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfg.GRPC.Port,
	}
}

func (a *App) Run() error {
	const component = "grpcapp.Run"
	log := a.log.With(slog.String("component", component), slog.Int("port", a.port))

	log.Info("Starting gRPC server")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error("Failed to listen", slog.Any("error", err))
		return err
	}

	log.Info("gRPC server is running", slog.String("address", l.Addr().String()))
	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("Failed to serve gRPC", slog.Any("error", err))
		return err
	}
	return nil
}

func (a *App) Stop() {
	a.log.Info("Stopping gRPC server")
	a.gRPCServer.GracefulStop()
	a.log.Info("gRPC server stopped")
}
