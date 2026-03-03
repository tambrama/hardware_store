package suite

import (
	"auth-service/internal/config"
	"context"
	"net"
	"strconv"
	"testing"

	pb "github.com/tambrama/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient pb.AuthClient
}

const (
	localConfigPath = "../config/test_local.yaml"
	grpcHost        = "localhost"
	appSecret       = "your_secret_key_here"
)

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	// t.Parallel()

	cfg := config.ConfigByPath(localConfigPath)
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect to grpc server: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: pb.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
