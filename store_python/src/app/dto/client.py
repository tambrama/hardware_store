from marshmallow import Schema, fields, validate
from app.dto.address import AddressSchema


class ClientSchema(Schema):
    """DTO: клиент (ответ)."""
    id = fields.UUID(dump_only=True)
    client_name = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    client_surname = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    birthday = fields.Date(required=True)
    gender = fields.Str(required=True, validate=validate.Length(min=1, max=20))
    registration_date = fields.Date(required=True)
    address_id = fields.UUID(required=True)
    address = fields.Nested(AddressSchema, dump_only=True)


class ClientCreateSchema(Schema):
    """DTO: создание клиента (вход)."""
    client_name = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    client_surname = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    birthday = fields.Date(required=True)
    gender = fields.Str(required=True, validate=validate.Length(min=1, max=20))
    registration_date = fields.Date(required=True)
    address_id = fields.UUID(required=True)
