from flask import Blueprint, request, jsonify
from app.models import Address
from app.repositories import AddressRepository
from app.dto import AddressSchema, AddressCreateSchema
from app.mappers import address_dao_to_dto, address_dto_to_dao

addresses_bp = Blueprint("addresses", __name__)
address_repo = AddressRepository()
address_create_schema = AddressCreateSchema()


@addresses_bp.route("", methods=["POST"])
def create_address():
    """
    Создание адреса (для использования в client/supplier)
    ---
    tags:
      - addresses
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - country
              - city
              - street
            properties:
              country:
                type: string
                example: "Россия"
              city:
                type: string
                example: "Москва"
              street:
                type: string
                example: "ул. Ленина, 1"
    responses:
      201:
        description: Адрес создан
      400:
        description: Ошибка валидации
      401:
        description: Не авторизован
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json()
    errors = address_create_schema.validate(data)
    if errors:
        return jsonify({"message": str(errors)}), 400
    try:
        loaded = address_create_schema.load(data)
    except Exception as e:
        return jsonify({"message": str(e)}), 400
    address = Address()
    address_dto_to_dao(loaded, address)
    address_repo.add(address)
    return jsonify(address_dao_to_dto(address)), 201
