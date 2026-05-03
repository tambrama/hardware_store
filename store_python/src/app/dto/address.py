from marshmallow import Schema, fields, validate


class AddressSchema(Schema):
    """DTO: адрес (ответ)."""
    id = fields.UUID(dump_only=True)
    country = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    city = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    street = fields.Str(required=True, validate=validate.Length(min=1, max=255))


class AddressCreateSchema(Schema):
    """DTO: создание/обновление адреса (вход)."""
    country = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    city = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    street = fields.Str(required=True, validate=validate.Length(min=1, max=255))
