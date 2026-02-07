package clientservice

import (
	"context"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/client"

	"github.com/google/uuid"
)

type ClientService interface {
	CreateClient(context.Context, client.Client, address.Address) error
	DeleteClient(ctx context.Context, id uuid.UUID) error
	UpdateAddressClient(ctx context.Context, id uuid.UUID, address address.Address) error
	GetClient(ctx context.Context, name, surname string) (client.Client, error)
	GetClients(ctx context.Context, limit, offset int) ([]client.Client, error)
}
