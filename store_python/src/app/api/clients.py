from flask import Blueprint, request, jsonify

from app.models import Client, Address
from app.repositories import ClientRepository, AddressRepository
from app.dto import ClientSchema, ClientCreateSchema, AddressCreateSchema
from app.mappers import client_dao_to_dto, client_dto_to_dao, address_dto_to_dao

clients_bp = Blueprint("clients", __name__)
client_repo = ClientRepository()
address_repo = AddressRepository()
client_schema = ClientSchema()
client_create_schema = ClientCreateSchema()
address_create_schema = AddressCreateSchema()


@clients_bp.route("", methods=["POST"])
def add_client():
    """
    Добавление клиента
    ---
    tags:
      - clients
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - client_name
              - client_surname
              - birthday
              - gender
              - registration_date
              - address_id
            properties:
              client_name:
                type: string
                example: "Иван"
              client_surname:
                type: string
                example: "Петров"
              birthday:
                type: string
                format: date
                example: "1990-05-15"
              gender:
                type: string
                example: "male"
              registration_date:
                type: string
                format: date
                example: "2024-01-10"
              address_id:
                type: string
                format: uuid
    responses:
      201:
        description: Клиент создан
      400:
        description: Ошибка валидации
      401:
        description: Не авторизован
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json()
    errors = client_create_schema.validate(data)
    if errors:
        return jsonify({"message": str(errors)}), 400
    try:
        loaded = client_create_schema.load(data)
    except Exception as e:
        return jsonify({"message": str(e)}), 400
    if address_repo.get_by_id(loaded["address_id"]) is None:
        return jsonify({"message": "Address not found"}), 400
    client = Client()
    client_dto_to_dao(loaded, client)
    client_repo.add(client)
    return jsonify(client_dao_to_dto(client)), 201


@clients_bp.route("/<uuid:client_id>", methods=["DELETE"])
def delete_client(client_id):
    """
    Удаление клиента по id
    ---
    tags:
      - clients
    security:
      - BearerAuth: []
    parameters:
      - name: client_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      204:
        description: Клиент удалён
      404:
        description: Клиент не найден
      401:
        description: Не авторизован
    """
    client = client_repo.get_by_id(client_id)
    if client is None:
        return jsonify({"message": "Client not found"}), 404
    client_repo.delete(client)
    return "", 204


@clients_bp.route("", methods=["GET"])
def get_clients():
    """
    Получение клиентов (все или по имени/фамилии, с пагинацией)
    ---
    tags:
      - clients
    security:
      - BearerAuth: []
    parameters:
      - name: name
        in: query
        schema:
          type: string
      - name: surname
        in: query
        schema:
          type: string
      - name: limit
        in: query
        schema:
          type: integer
      - name: offset
        in: query
        schema:
          type: integer
    responses:
      200:
        description: Список клиентов (может быть пустым)
      401:
        description: Не авторизован
    """
    name = request.args.get("name")
    surname = request.args.get("surname")
    limit = request.args.get("limit", type=int)
    offset = request.args.get("offset", type=int)

    if name is not None and surname is not None:
        clients = client_repo.get_by_name_surname(name, surname)
    else:
        clients = client_repo.get_all(limit=limit, offset=offset)

    return jsonify([client_dao_to_dto(c) for c in clients])


@clients_bp.route("/<uuid:client_id>", methods=["GET"])
def get_client_by_id(client_id):
    """
    Получение клиента по id
    ---
    tags:
      - clients
    security:
      - BearerAuth: []
    parameters:
      - name: client_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      200:
        description: Клиент найден
      404:
        description: Клиент не найден
      401:
        description: Не авторизован
    """
    client = client_repo.get_by_id(client_id)
    if client is None:
        return jsonify({"message": "Client not found"}), 404
    return jsonify(client_dao_to_dto(client))


@clients_bp.route("/<uuid:client_id>/address", methods=["PATCH"])
def update_client_address(client_id):
    """
    Изменение адреса клиента
    ---
    tags:
      - clients
    security:
      - BearerAuth: []
    parameters:
      - name: client_id
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
                example: "Москва"
              street:
                type: string
                example: "ул. Ленина, 1"
    responses:
      200:
        description: Адрес обновлён
      400:
        description: Ошибка валидации
      404:
        description: Клиент не найден
      401:
        description: Не авторизован
    """
    client = client_repo.get_by_id(client_id)
    if client is None:
        return jsonify({"message": "Client not found"}), 404
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
    client.address_id = new_address.id
    client_repo.update(client)
    return jsonify(client_dao_to_dto(client))
