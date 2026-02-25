package address

import (
	"context"
	"errors"
	"fmt"
	"hardware_store/internal/model/address"
	"hardware_store/internal/storage"
	"hardware_store/internal/storage/postgres/dto"
	"hardware_store/internal/storage/postgres/mapper"
	"hardware_store/internal/storage/postgres/tx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type addressRepository struct {
	pool *pgxpool.Pool
}

func NewAddressRepository(db *pgxpool.Pool) *addressRepository {
	return &addressRepository{
		pool: db,
	}
}

func (r *addressRepository) Insert(ctx context.Context, address address.Address) error {
	const op = "storage.postgres.CreateAddress"
	dto := mapper.AddressToDTO(address)
	query := `INSERT INTO address 
	(address_id, country, city, street)
	VALUES ($1,$2,$3,$4)`
	_, err := r.pool.Exec(ctx, query, dto.AddressID, dto.Country, dto.City, dto.Street)
	if err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrCreation)
	}

	return nil
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	exec := tx.FromContext(ctx, r.pool)
	query := `DELETE FROM address 
	WHERE address_id = $1`
	_, err := exec.Exec(ctx, query, id)
	if err != nil {
		return storage.ErrDelete
	}

	return nil
}

func (r *addressRepository) Get(ctx context.Context, id uuid.UUID) (address.Address, error) {

	query := `SELECT * FROM address 
	WHERE address_id = $1`
	var dto dto.AddressDTO
	err := r.pool.QueryRow(ctx, query, id).Scan(&dto.AddressID, &dto.Country, &dto.City, &dto.Street)
	if err != nil {
		return address.Address{}, storage.ErrAddressNotFound
	}

	return mapper.AddressFromDTO(dto), nil
}

func (r *addressRepository) Update(ctx context.Context, addr address.Address) (address.Address, error) {
	query := `UPDATE address 
	SET country = $2, city = $3, street = $4
	WHERE address_id = $1
	RETURNING address_id, country, city, street`
	addrDto := mapper.AddressToDTO(addr)
	var dto dto.AddressDTO
	err := r.pool.QueryRow(ctx, query, addrDto.AddressID, addrDto.Country, addrDto.City, addrDto.Street).Scan(&dto.AddressID, &dto.Country, &dto.City, &dto.Street)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return address.Address{}, storage.ErrClientNotFound
		}
		return address.Address{}, storage.ErrUpdate
	}

	return mapper.AddressFromDTO(dto), nil
}
