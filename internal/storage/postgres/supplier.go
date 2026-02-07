package postgres

import (
	"context"
	"fmt"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/supplier"
	"hardware_store/internal/storage"
	"hardware_store/internal/storage/postgres/dto"
	"hardware_store/internal/storage/postgres/mapper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type supplierRepository struct {
	pool *pgxpool.Pool
}

func NewSupplierRepository(db *pgxpool.Pool) supplier.SupplierRepository {
	return &supplierRepository{
		pool: db,
	}
}

func (r *supplierRepository) Insert(ctx context.Context, supplier supplier.Supplier) error {
	query := `INSERT INTO supplier 
	(supplier_id, name, address_id, phone_number)
	VALUES ($1,$2,$3,$4)`
	dto := mapper.SupplierToDTO(supplier)
	_, err := r.pool.Exec(ctx, query, dto.SupplierID, dto.Name, dto.AddressID, dto.PhoneNumber)
	if err != nil {
		return storage.ErrCreation
	}

	return nil
}

func (r *supplierRepository) UpdateAddress(ctx context.Context, id uuid.UUID, addr address.Address) error {
	query := `UPDATE address 
	SET country = $2, city = $3, street = $4
	WHERE address_id = (SELECT address_id 
	FROM supplier
	WHERE supplier_id = $1)`
	addrDto := mapper.AddressToDTO(addr)
	_, err := r.pool.Exec(ctx, query, id, addrDto.Country, addrDto.City, addrDto.Street)
	if err != nil {
		return storage.ErrUpdate
	}
	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM supplier 
	WHERE supplier_id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return storage.ErrDelete
	}
	return nil
}

func (r *supplierRepository) GetById(ctx context.Context, id uuid.UUID) (supplier.Supplier, error) {
	query := `SELECT * FROM supplier 
	WHERE supplier_id = $1`

	var dto dto.SupplierDTO

	err := r.pool.QueryRow(ctx, query, id).Scan(&dto.SupplierID, &dto.Name, &dto.AddressID, &dto.PhoneNumber)
	if err != nil {
		return supplier.Supplier{}, storage.ErrSupplierNotFound
	}
	return mapper.SupplierFromDTO(dto), nil
}

func (r *supplierRepository) GetAll(ctx context.Context) ([]supplier.Supplier, error) {
	query := `SELECT * FROM supplier`

	row, err := r.pool.Query(ctx, query)
	if err != nil {
		return []supplier.Supplier{}, storage.ErrSupplierNotFound
	}
	defer row.Close()
	var suppliers []supplier.Supplier
	for row.Next() {
		var dto dto.SupplierDTO

		if err := row.Scan(&dto.SupplierID, &dto.Name, &dto.AddressID, &dto.PhoneNumber); err != nil {
			return []supplier.Supplier{}, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		suppliers = append(suppliers, mapper.SupplierFromDTO(dto))
	}

	if err = row.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return suppliers, nil
}

func (r *supplierRepository) UnsetAddress(ctx context.Context, address uuid.UUID) error {
	exec := fromContext(ctx, r.pool)

	query := `UPDATE supplier 
	SET address_id = NULL
	WHERE address_id = $1`

	_, err := exec.Exec(ctx, query, address)
	return err
}
