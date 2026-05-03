"""Middleware для проверки JWT через auth-service (gRPC)."""
from functools import wraps
from flask import request, jsonify
from app.clients import AuthGrpcClient

_auth_client = AuthGrpcClient()


def get_bearer_token() -> str | None:
    """Извлекает Bearer токен из заголовка Authorization."""
    auth = request.headers.get("Authorization")
    if not auth or not auth.startswith("Bearer "):
        return None
    return auth[7:].strip()


def require_auth(f):
    """Декоратор: требует валидный JWT. Возвращает 401 при отсутствии или невалидном токене."""

    @wraps(f)
    def decorated(*args, **kwargs):
        token = get_bearer_token()
        if not token:
            return jsonify({"message": "Missing or invalid Authorization header"}), 401
        try:
            user_id, app_id = _auth_client.validate(token)
            request.user_id = user_id
            request.app_id = app_id
            return f(*args, **kwargs)
        except Exception:
            return jsonify({"message": "Invalid or expired token"}), 401

    return decorated
