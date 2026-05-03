"""Эндпоинты авторизации: /register, /auth, /reset."""
import os
from flask import Blueprint, request, jsonify
from app.clients import AuthGrpcClient

auth_bp = Blueprint("auth", __name__)
_auth_client = AuthGrpcClient()


@auth_bp.route("/register", methods=["POST"])
def register():
    """
    Регистрация пользователя
    ---
    tags:
      - auth
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - email
              - password
              - name
              - surname
            properties:
              email:
                type: string
                example: user@example.com
              password:
                type: string
                example: secret123
              name:
                type: string
                example: Иван
              surname:
                type: string
                example: Петров
              phone_number:
                type: string
                example: "+79001234567"
    responses:
      201:
        description: Пользователь создан
      400:
        description: Ошибка валидации
      409:
        description: Email уже занят
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json() or {}
    email = data.get("email")
    password = data.get("password")
    name = data.get("name")
    surname = data.get("surname")
    phone_number = data.get("phone_number", "")
    if not all([email, password, name, surname]):
        return jsonify({"message": "Missing required fields: email, password, name, surname"}), 400
    try:
        user_id = _auth_client.register(email, password, name, surname, phone_number)
        return jsonify({"user_id": user_id}), 201
    except Exception as e:
        err_msg = str(e)
        if "already exists" in err_msg.lower() or "409" in err_msg:
            return jsonify({"message": "User with this email already exists"}), 409
        return jsonify({"message": err_msg}), 400


@auth_bp.route("", methods=["POST"])
def login():
    """
    Логин (получение токенов)
    ---
    tags:
      - auth
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - email
              - password
            properties:
              email:
                type: string
                example: user@example.com
              password:
                type: string
                example: secret123
    responses:
      200:
        description: Токены получены
      401:
        description: Неверные учётные данные
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json() or {}
    email = data.get("email")
    password = data.get("password")
    if not email or not password:
        return jsonify({"message": "Missing email or password"}), 400
    try:
        access_token, refresh_token = _auth_client.login(email, password)
        return jsonify({
            "access_token": access_token,
            "refresh_token": refresh_token,
        }), 200
    except Exception as e:
        return jsonify({"message": "Invalid credentials"}), 401


@auth_bp.route("/oauth/google", methods=["GET"])
def google_oauth_url():
    """
    Начало Google OAuth: вернуть auth_url
    ---
    tags:
      - auth
    parameters:
      - name: state
        in: query
        required: false
        schema:
          type: string
    responses:
      200:
        description: URL для редиректа на Google
      500:
        description: Ошибка сервера
    """
    state = (request.args.get("state") or "").strip()
    if not state:
        # Keep behavior aligned with Go shop: auto-generate state when absent.
        import uuid
        state = "test_" + uuid.uuid4().hex[:8]
    try:
        url = _auth_client.get_google_auth_url(state=state)
        return jsonify({"auth_url": url}), 200
    except Exception as e:
        return jsonify({"message": str(e)}), 500


@auth_bp.route("/google/callback", methods=["GET"])
def google_callback():
    """
    Callback Google OAuth: обмен code/state на JWT токены
    ---
    tags:
      - auth
    parameters:
      - name: code
        in: query
        required: true
        schema:
          type: string
      - name: state
        in: query
        required: true
        schema:
          type: string
    responses:
      200:
        description: JWT токены
      400:
        description: Нет code/state
      401:
        description: Авторизация не удалась
    """
    code = (request.args.get("code") or "").strip()
    state = (request.args.get("state") or "").strip()
    if not code or not state:
        return jsonify({"message": "code and state are required"}), 400
    try:
        access_token, refresh_token = _auth_client.google_callback(code=code, state=state)
        return jsonify({"access_token": access_token, "refresh_token": refresh_token}), 200
    except Exception:
        return jsonify({"message": "authentication failed"}), 401


@auth_bp.route("/reset", methods=["POST"])
def reset_password():
    """
    Восстановление пароля (заглушка — пароль выводится в консоль auth-service)
    ---
    tags:
      - auth
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - email
            properties:
              email:
                type: string
    responses:
      200:
        description: Пароль отправлен (заглушка)
      400:
        description: Email не найден
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json() or {}
    email = data.get("email")
    if not email:
        return jsonify({"message": "Missing email"}), 400
    try:
        _auth_client.restore_password(email)
        return jsonify({"message": "Password reset link sent (stub: check auth-service console)"}), 200
    except Exception as e:
        return jsonify({"message": str(e)}), 400


@auth_bp.route("/change-password", methods=["PATCH"])
def change_password():
    """
    Смена пароля (требует токен)
    ---
    tags:
      - auth
    security:
      - BearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required:
              - email
              - old_password
              - new_password
            properties:
              email:
                type: string
              old_password:
                type: string
              new_password:
                type: string
    responses:
      200:
        description: Пароль изменён
      401:
        description: Неверный старый пароль
    """
    if not request.is_json:
        return jsonify({"message": "Request must be JSON"}), 400
    data = request.get_json() or {}
    email = data.get("email")
    old_password = data.get("old_password")
    new_password = data.get("new_password")
    if not all([email, old_password, new_password]):
        return jsonify({"message": "Missing email, old_password or new_password"}), 400
    try:
        _auth_client.change_password(email, old_password, new_password)
        return jsonify({"message": "Password changed"}), 200
    except Exception as e:
        return jsonify({"message": str(e)}), 400
