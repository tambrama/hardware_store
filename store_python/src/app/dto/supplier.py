from marshmallow import Schema, fields, validate
from app.dto.address import AddressSchema


class SupplierSchema(Schema):
    """DTO: поставщик (ответ)."""
    id = fields.UUID(dump_only=True)
    name = fields.Str(required=True, validate=validate.Length(min=1, max=200))
    address_id = fields.UUID(required=True)
    phone_number = fields.Str(required=True, validate=validate.Length(min=1, max=50))
    address = fields.Nested(AddressSchema, dump_only=True)


class SupplierCreateSchema(Schema):
    """DTO: создание поставщика (вход)."""
    name = fields.Str(required=True, validate=validate.Length(min=1, max=200))
    address_id = fields.UUID(required=True)
    phone_number = fields.Str(required=True, validate=validate.Length(min=1, max=50))
