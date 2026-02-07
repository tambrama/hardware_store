package client

import (
	"context"
	"hardware_store/internal/model/address"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ClientID         uuid.UUID `validate:"required"`
	Name             string    `validate:"required,min=2,max=50"`
	Surname          string    `validate:"required,min=2,max=50"`
	Birthday         time.Time `validate:"datetime=2006-01-02"`
	Gender           string    `validate:"oneof=male female"`
	RegistrationDate time.Time
	AddressID        uuid.UUID
}

type ClientRepository interface {
	Insert(ctx context.Context, client Client) error
	Delete(ctx context.Context, clientID uuid.UUID) error
	GetByName(ctx context.Context, name, surname string) (Client, error)
	GetById(ctx context.Context, id uuid.UUID) (Client, error)
	GetAll(ctx context.Context, limit, offset int) ([]Client, error)
	UpdateAddress(ctx context.Context, clientUUID uuid.UUID, address address.Address) error
	UnsetAddress(ctx context.Context, addressId uuid.UUID) error
}
