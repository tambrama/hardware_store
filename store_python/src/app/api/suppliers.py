import uuid
from flask import Blueprint, request, jsonify
from app.models import Address, Supplier
from app.repositories import SupplierRepository, AddressRepository
from app.dto import SupplierSchema, SupplierCreateSchema, AddressCreateSchema
from app.mappers import supplier_dao_to_dto, supplier_dto_to_dao, address_dto_to_dao

suppliers_bp = Blueprint("suppliers", __name__)
supplier_repo = SupplierRepository()
address_repo = AddressRepository()
supplier_schema = SupplierSchema()
supplier_create_schema = SupplierCreateSchema()
address_create_schema = AddressCreateSchema()


@suppliers_bp.route("", methods=["POST"])
def add_supplier():
    """
    Добавление поставщика
    ---
    tags:
      - suppliers
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - name
              - address_id
              - phone_number
            properties:
              name:
                type: string
                example: "ООО Поставки"
              address_id:
                type: string
                format: uuid
              phone_number:
                type: string
                example: "+7 495 123-45-67"
    responses:
      201:
        description: Поставщик создан
      400:
        description: Ошибка валидации
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json()
    errors = supplier_create_schema.validate(data)
    if errors:
        return jsonify({"message": str(errors)}), 400
    try:
        loaded = supplier_create_schema.load(data)
    except Exception as e:
        return jsonify({"message": str(e)}), 400
    if address_repo.get_by_id(loaded["address_id"]) is None:
        return jsonify({"message": "Address not found"}), 400
    supplier = Supplier()
    supplier_dto_to_dao(loaded, supplier)
    supplier_repo.add(supplier)
    return jsonify(supplier_dao_to_dto(supplier)), 201


@suppliers_bp.route("/<uuid:supplier_id>/address", methods=["PATCH"])
def update_supplier_address(supplier_id):
    """
    Изменение адреса поставщика
    ---
    tags:
      - suppliers
    security:
      - BearerAuth: []
    parameters:
      - name: supplier_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
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
                example: "Санкт-Петербург"
              street:
                type: string
                example: "Невский пр., 1"
    responses:
      200:
        description: Адрес обновлён
      400:
        description: Ошибка валидации
      404:
        description: Поставщик не найден
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    supplier = supplier_repo.get_by_id(supplier_id)
    if supplier is None:
        return jsonify({"message": "Supplier not found"}), 404
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
    new_address = Address()
    address_dto_to_dao(loaded, new_address)
    address_repo.add(new_address)
    supplier.address_id = new_address.id
    supplier_repo.update(supplier)
    return jsonify(supplier_dao_to_dto(supplier))


@suppliers_bp.route("/<uuid:supplier_id>", methods=["DELETE"])
def delete_supplier(supplier_id):
    """
    Удаление поставщика по id
    ---
    tags:
      - suppliers
    security:
      - BearerAuth: []
    parameters:
      - name: supplier_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      204:
        description: Поставщик удалён
      404:
        description: Поставщик не найден
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    supplier = supplier_repo.get_by_id(supplier_id)
    if supplier is None:
        return jsonify({"message": "Supplier not found"}), 404
    supplier_repo.delete(supplier)
    return "", 204


@suppliers_bp.route("", methods=["GET"])
def get_all_suppliers():
    """
    Получение всех поставщиков
    ---
    tags:
      - suppliers
    security:
      - BearerAuth: []
    responses:
      200:
        description: Список поставщиков (может быть пустым)
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    suppliers = supplier_repo.get_all()
    return jsonify([supplier_dao_to_dto(s) for s in suppliers])


@suppliers_bp.route("/<uuid:supplier_id>", methods=["GET"])
def get_supplier_by_id(supplier_id):
    """
    Получение поставщика по id
    ---
    tags:
      - suppliers
    security:
      - BearerAuth: []
    parameters:
      - name: supplier_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      200:
        description: Поставщик найден
      404:
        description: Поставщик не найден
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    supplier = supplier_repo.get_by_id(supplier_id)
    if supplier is None:
        return jsonify({"message": "Supplier not found"}), 404
    return jsonify(supplier_dao_to_dto(supplier))
