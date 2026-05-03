import uuid
from app.extensions import db


class Address(db.Model):
    """DAO: таблица адресов."""
    __tablename__ = "address"

    id = db.Column(db.Uuid(as_uuid=True), primary_key=True, default=uuid.uuid4)
    country = db.Column(db.String(100), nullable=False)
    city = db.Column(db.String(100), nullable=False)
    street = db.Column(db.String(255), nullable=False)
