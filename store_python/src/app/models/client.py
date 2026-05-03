import uuid
from datetime import date, datetime
from app.extensions import db


class Client(db.Model):
    """DAO: таблица клиентов."""
    __tablename__ = "client"

    id = db.Column(db.Uuid(as_uuid=True), primary_key=True, default=uuid.uuid4)
    client_name = db.Column(db.String(100), nullable=False)
    client_surname = db.Column(db.String(100), nullable=False)
    birthday = db.Column(db.Date, nullable=False)
    gender = db.Column(db.String(20), nullable=False)
    registration_date = db.Column(db.Date, nullable=False)
    address_id = db.Column(db.Uuid(as_uuid=True), db.ForeignKey("address.id"), nullable=False)

    address = db.relationship("Address", backref="clients", foreign_keys=[address_id])
