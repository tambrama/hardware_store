import uuid
from app.extensions import db


class Supplier(db.Model):
    """DAO: таблица поставщиков."""
    __tablename__ = "supplier"

    id = db.Column(db.Uuid(as_uuid=True), primary_key=True, default=uuid.uuid4)
    name = db.Column(db.String(200), nullable=False)
    address_id = db.Column(db.Uuid(as_uuid=True), db.ForeignKey("address.id"), nullable=False)
    phone_number = db.Column(db.String(50), nullable=False)

    address = db.relationship("Address", backref="suppliers", foreign_keys=[address_id])
