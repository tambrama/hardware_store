from marshmallow import Schema, fields


class ImageResponseSchema(Schema):
    """DTO: метаданные изображения (ответ)."""
    id = fields.UUID(dump_only=True)
