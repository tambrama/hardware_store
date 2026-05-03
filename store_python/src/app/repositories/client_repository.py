import uuid
from app.extensions import db
from app.models import Client


class ClientRepository:
    """Репозиторий для работы с клиентами (DAO)."""

    def get_by_id(self, client_id: uuid.UUID) -> Client | None:
        return db.session.get(Client, client_id)

    def get_by_name_surname(self, name: str, surname: str) -> list[Client]:
        q = db.session.query(Client).filter(
            Client.client_name == name,
            Client.client_surname == surname,
        )
        return list(q.all())

    def get_all(self, limit: int | None = None, offset: int | None = None) -> list[Client]:
        q = db.session.query(Client).order_by(Client.registration_date.desc())
        if limit is not None:
            q = q.limit(limit)
        if offset is not None:
            q = q.offset(offset)
        return list(q.all())

    def add(self, client: Client) -> Client:
        db.session.add(client)
        db.session.commit()
        return client

    def update(self, client: Client) -> Client:
        db.session.commit()
        return client

    def delete(self, client: Client) -> None:
        db.session.delete(client)
        db.session.commit()
