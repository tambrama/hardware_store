package client

import (
	"context"
	"fmt"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/client"
	"hardware_store/internal/storage"
	"hardware_store/internal/storage/postgres/dto"
	"hardware_store/internal/storage/postgres/mapper"
	"hardware_store/internal/storage/postgres/tx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type clientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *clientRepository {
	return &clientRepository{
		pool: db,
	}
}

func (r *clientRepository) Insert(ctx context.Context, client client.Client) error {
	query := `INSERT INTO client 
	(client_id, name, surname, birthday, gender, registration_date, address_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	ON CONFLICT (client_id) DO UPDATE SET
    name = $2,
    surname = $3,
    birthday = $4,
	gender = $5,
	registration_date = NOW(),
	address_id = $7`
	dto := mapper.ClientToDTO(client)
	_, err := r.pool.Exec(ctx, query, dto.ClientID, dto.Name, dto.Surname, dto.Birthday, dto.Gender, dto.RegistrationDate, dto.AddressID)
	if err != nil {
		return storage.ErrCreation
	}

	return nil
}

func (r *clientRepository) Delete(ctx context.Context, clientID uuid.UUID) error {
	query := `DELETE FROM client 
	WHERE client_id = $1`
	_, err := r.pool.Exec(ctx, query, clientID)
	if err != nil {
		return storage.ErrDelete
	}
	return nil
}

func (r *clientRepository) GetByName(ctx context.Context, name, surname string) (client.Client, error) {
	query := `SELECT client_id, name, surname, birthday, gender, registration_date, address_id FROM client 
	WHERE name = $1 AND surname = $2`

	var dto dto.ClientDTO

	err := r.pool.QueryRow(ctx, query, name, surname).Scan(&dto.ClientID, &dto.Name, &dto.Surname, &dto.Birthday, &dto.Gender, &dto.RegistrationDate, &dto.AddressID)
	if err != nil {
		return client.Client{}, storage.ErrClientNotFound
	}
	return mapper.ClientFromDTO(dto), nil
}

func (r *clientRepository) GetById(ctx context.Context, id uuid.UUID) (client.Client, error) {
	query := `SELECT client_id, name, surname, birthday, gender, registration_date, address_id FROM client 
	WHERE client_id = $1`

	var dto dto.ClientDTO

	err := r.pool.QueryRow(ctx, query, id).Scan(&dto.ClientID, &dto.Name, &dto.Surname, &dto.Birthday, &dto.Gender, &dto.RegistrationDate, &dto.AddressID)
	if err != nil {
		return client.Client{}, storage.ErrClientNotFound
	}
	return mapper.ClientFromDTO(dto), nil
}

func (r *clientRepository) GetAll(ctx context.Context, limit, offset int) ([]client.Client, error) {
	query := `SELECT client_id, name, surname, birthday, gender, registration_date, address_id FROM client`

	var row pgx.Rows
	var err error

	if limit > 0 {
		query += ` LIMIT $1 OFFSET $2`
		row, err = r.pool.Query(ctx, query, limit, offset)
	} else {
		row, err = r.pool.Query(ctx, query)
	}

	if err != nil {
		return []client.Client{}, storage.ErrClientNotFound
	}

	defer row.Close()
	var clients []client.Client
	for row.Next() {
		var dto dto.ClientDTO

		if err := row.Scan(&dto.ClientID, &dto.Name, &dto.Surname, &dto.Birthday, &dto.Gender, &dto.RegistrationDate, &dto.AddressID); err != nil {
			return []client.Client{}, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		clients = append(clients, mapper.ClientFromDTO(dto))
	}

	if err = row.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return clients, nil
}

func (r *clientRepository) UpdateAddress(ctx context.Context, clientUUID uuid.UUID, address address.Address) error {
	query := `UPDATE address 
	SET country = $2, city = $3, street = $4
	WHERE address_id = (SELECT address_id 
	FROM client
	WHERE client_id = $1)`

	dto := mapper.AddressToDTO(address)

	res, err := r.pool.Exec(ctx, query, clientUUID, dto.Country, dto.City, dto.Street)
	if err != nil {
		return storage.ErrUpdate
	}
	if res.RowsAffected() == 0 {
		return storage.ErrClientNotFound
	}

	return nil
}

func (r *clientRepository) UnsetAddress(ctx context.Context, address uuid.UUID) error {
	exec := tx.FromContext(ctx, r.pool)

	query := `UPDATE client 
	SET address_id = NULL
	WHERE address_id = $1`

	_, err := exec.Exec(ctx, query, address)
	return err
}
