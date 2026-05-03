import uuid
from app.extensions import db
from app.models import Supplier


class SupplierRepository:
    """Репозиторий для работы с поставщиками (DAO)."""

    def get_by_id(self, supplier_id: uuid.UUID) -> Supplier | None:
        return db.session.get(Supplier, supplier_id)

    def get_all(self) -> list[Supplier]:
        return list(db.session.query(Supplier).order_by(Supplier.name).all())

    def add(self, supplier: Supplier) -> Supplier:
        db.session.add(supplier)
        db.session.commit()
        return supplier

    def update(self, supplier: Supplier) -> Supplier:
        db.session.commit()
        return supplier

    def delete(self, supplier: Supplier) -> None:
        db.session.delete(supplier)
        db.session.commit()
