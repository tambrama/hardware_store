from datetime import date
from flask import Blueprint, request, jsonify
from app.models import Product
from app.repositories import ProductRepository, SupplierRepository
from app.dto import ProductSchema, ProductCreateSchema, ProductDecreaseSchema
from app.mappers import product_dao_to_dto, product_dto_to_dao

products_bp = Blueprint("products", __name__)
product_repo = ProductRepository()
supplier_repo = SupplierRepository()
product_schema = ProductSchema()
product_create_schema = ProductCreateSchema()
product_decrease_schema = ProductDecreaseSchema()


@products_bp.route("", methods=["POST"])
def add_product():
    """
    Добавление товара
    ---
    tags:
      - products
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
              - category
              - price
              - available_stock
              - supplier_id
            properties:
              name:
                type: string
                example: "Холодильник"
              category:
                type: string
                example: "Бытовая техника"
              price:
                type: number
                example: 29999.99
              available_stock:
                type: integer
                example: 10
              last_update_date:
                type: string
                format: date
              supplier_id:
                type: string
                format: uuid
              image_id:
                type: string
                format: uuid
    responses:
      201:
        description: Товар создан
      400:
        description: Ошибка валидации
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json()
    errors = product_create_schema.validate(data)
    if errors:
        return jsonify({"message": str(errors)}), 400
    try:
        loaded = product_create_schema.load(data)
    except Exception as e:
        return jsonify({"message": str(e)}), 400
    if supplier_repo.get_by_id(loaded["supplier_id"]) is None:
        return jsonify({"message": "Supplier not found"}), 400
    product = Product()
    product_dto_to_dao(loaded, product)
    product_repo.add(product)
    return jsonify(product_dao_to_dto(product)), 201


@products_bp.route("/<uuid:product_id>/decrease", methods=["PATCH"])
def decrease_product_stock(product_id):
    """
    Уменьшение количества товара
    ---
    tags:
      - products
    security:
      - BearerAuth: []
    parameters:
      - name: product_id
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
              - amount
            properties:
              amount:
                type: integer
                minimum: 1
                example: 2
    responses:
      200:
        description: Количество уменьшено
      400:
        description: Ошибка валидации или недостаточно товара
      404:
        description: Товар не найден
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    product = product_repo.get_by_id(product_id)
    if product is None:
        return jsonify({"message": "Product not found"}), 404
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json()
    errors = product_decrease_schema.validate(data)
    if errors:
        return jsonify({"message": str(errors)}), 400
    try:
        loaded = product_decrease_schema.load(data)
    except Exception as e:
        return jsonify({"message": str(e)}), 400
    amount = loaded["amount"]
    if product.available_stock < amount:
        return jsonify({"message": "Insufficient stock"}), 400
    product.available_stock -= amount
    product.last_update_date = date.today()
    product_repo.update(product)
    return jsonify(product_dao_to_dto(product))


@products_bp.route("/<uuid:product_id>", methods=["GET"])
def get_product_by_id(product_id):
    """
    Получение товара по id
    ---
    tags:
      - products
    security:
      - BearerAuth: []
    parameters:
      - name: product_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      200:
        description: Товар найден
      404:
        description: Товар не найден
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    product = product_repo.get_by_id(product_id)
    if product is None:
        return jsonify({"message": "Product not found"}), 404
    return jsonify(product_dao_to_dto(product))


@products_bp.route("", methods=["GET"])
def get_all_products():
    """
    Получение всех доступных товаров
    ---
    tags:
      - products
    security:
      - BearerAuth: []
    responses:
      200:
        description: Список товаров (может быть пустым)
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    products = product_repo.get_all_available()
    return jsonify([product_dao_to_dto(p) for p in products])


@products_bp.route("/<uuid:product_id>", methods=["DELETE"])
def delete_product(product_id):
    """
    Удаление товара по id
    ---
    tags:
      - products
    security:
      - BearerAuth: []
    parameters:
      - name: product_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      204:
        description: Товар удалён
      404:
        description: Товар не найден
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    product = product_repo.get_by_id(product_id)
    if product is None:
        return jsonify({"message": "Product not found"}), 404
    product_repo.delete(product)
    return "", 204
