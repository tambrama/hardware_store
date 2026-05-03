from marshmallow import Schema, fields, validate
from decimal import Decimal


class ProductSchema(Schema):
    """DTO: товар (ответ)."""
    id = fields.UUID(dump_only=True)
    name = fields.Str(required=True, validate=validate.Length(min=1, max=200))
    category = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    price = fields.Decimal(required=True, as_string=True, places=2)
    available_stock = fields.Int(required=True)
    last_update_date = fields.Date(allow_none=True)
    supplier_id = fields.UUID(required=True)
    image_id = fields.UUID(allow_none=True)


class ProductCreateSchema(Schema):
    """DTO: создание товара (вход)."""
    name = fields.Str(required=True, validate=validate.Length(min=1, max=200))
    category = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    price = fields.Decimal(required=True, as_string=True, places=2)
    available_stock = fields.Int(required=True, validate=validate.Range(min=0))
    last_update_date = fields.Date(allow_none=True)
    supplier_id = fields.UUID(required=True)
    image_id = fields.UUID(allow_none=True)


class ProductDecreaseSchema(Schema):
    """DTO: уменьшение количества товара (вход)."""
    amount = fields.Int(required=True, validate=validate.Range(min=1))
