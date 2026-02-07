package supplierservice

import (
	"context"
	"hardware_store/internal/model/address"
	"hardware_store/internal/model/supplier"

	"github.com/google/uuid"
)

type SupplierService interface {
	CreateSupplier(ctx context.Context, supplier supplier.Supplier, address address.Address) error
	DeleteSupplier(ctx context.Context, id uuid.UUID) error
	UpdateAddressSupplier(ctx context.Context, id uuid.UUID, address address.Address) error
	GetSupplier(ctx context.Context, id uuid.UUID) (supplier.Supplier, error)
	GetSuppliers(ctx context.Context) ([]supplier.Supplier, error)
}
