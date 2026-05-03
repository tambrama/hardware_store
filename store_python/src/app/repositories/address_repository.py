import uuid
from app.extensions import db
from app.models import Address


class AddressRepository:
    """Репозиторий для работы с адресами (DAO)."""

    def get_by_id(self, address_id: uuid.UUID) -> Address | None:
        return db.session.get(Address, address_id)

    def add(self, address: Address) -> Address:
        db.session.add(address)
        db.session.commit()
        return address

    def update(self, address: Address) -> Address:
        db.session.commit()
        return address

    def delete(self, address: Address) -> None:
        db.session.delete(address)
        db.session.commit()
