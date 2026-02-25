package client

import (
	"context"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/client"
	"hardware_store/internal/model/tx"
	service "hardware_store/internal/service/address"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Insert(ctx context.Context, client client.Client) error
	Delete(ctx context.Context, clientID uuid.UUID) error
	GetByName(ctx context.Context, name, surname string) (client.Client, error)
	GetById(ctx context.Context, id uuid.UUID) (client.Client, error)
	GetAll(ctx context.Context, limit, offset int) ([]client.Client, error)
	UpdateAddress(ctx context.Context, clientUUID uuid.UUID, address address.Address) error
	UnsetAddress(ctx context.Context, addressId uuid.UUID) error
}

type clientService struct {
	repo    ClientRepository
	address service.AddressService
	tx      tx.Manager
}

func NewClientService(repo ClientRepository, address service.AddressService, tx tx.Manager) *clientService {
	return &clientService{repo: repo, address: address, tx: tx}
}

func (s *clientService) CreateClient(ctx context.Context, client client.Client, address address.Address) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.address.CreateAddress(ctx, address); err != nil {
			return err
		}
		return s.repo.Insert(ctx, client)
	})
}
func (s *clientService) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		getCli, err := s.repo.GetById(ctx, id)
		if err != nil {
			return err
		}
		if err = s.repo.Delete(ctx, id); err != nil {
			return err
		}

		return s.address.DeleteAddress(ctx, getCli.AddressID)
	})
}
func (s *clientService) UpdateAddressClient(ctx context.Context, id uuid.UUID, address address.Address) error {
	return s.repo.UpdateAddress(ctx, id, address)

}
func (s *clientService) GetClient(ctx context.Context, name, surname string) (client.Client, error) {
	return s.repo.GetByName(ctx, name, surname)
}
func (s *clientService) GetClients(ctx context.Context, limit, offset int) ([]client.Client, error) {
	return s.repo.GetAll(ctx, limit, offset)
}
