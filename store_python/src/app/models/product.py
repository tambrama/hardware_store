import uuid
from datetime import date, datetime
from decimal import Decimal
from app.extensions import db


class Product(db.Model):
    """DAO: таблица товаров."""
    __tablename__ = "product"

    id = db.Column(db.Uuid(as_uuid=True), primary_key=True, default=uuid.uuid4)
    name = db.Column(db.String(200), nullable=False)
    category = db.Column(db.String(100), nullable=False)
    price = db.Column(db.Numeric(12, 2), nullable=False)
    available_stock = db.Column(db.Integer, nullable=False, default=0)
    last_update_date = db.Column(db.Date, nullable=True)
    supplier_id = db.Column(db.Uuid(as_uuid=True), db.ForeignKey("supplier.id"), nullable=False)
    # Image bytes are stored in external photo-service; we keep only image UUID here.
    image_id = db.Column(db.Uuid(as_uuid=True), nullable=True)

    supplier = db.relationship("Supplier", backref="products", foreign_keys=[supplier_id])
