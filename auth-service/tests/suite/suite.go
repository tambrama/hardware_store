package suite

import (
	"auth-service/internal/config"
	"context"
	"testing"

	pb "github.com/tambrama/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type Suite struct {
	*testing.T
	Cfg *config.Config
	AuthClient pb.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.NewConfig()	
	ctx, cancelCtx := grpc.WithContextDialer(context.Background(), )
}