"""Pytest fixtures для тестов store_python."""
import os
import sys

# Добавляем src в путь
_src = os.path.join(os.path.dirname(__file__), "..", "src")
sys.path.insert(0, _src)

import pytest
from unittest.mock import patch, MagicMock


@pytest.fixture
def mock_auth_client():
    """Мок AuthGrpcClient — не вызывает реальный gRPC."""
    instance = MagicMock()
    instance.register.return_value = "user-uuid-123"
    instance.login.return_value = ("access-token-xyz", "refresh-token-abc")
    instance.validate.return_value = ("user-uuid-123", "app-id")
    instance.restore_password.return_value = None
    instance.change_password.return_value = None
    return instance


@pytest.fixture
def app(mock_auth_client):
    """Flask app с тестовой БД (SQLite в памяти) и моком auth."""
    os.environ["DATABASE_URI"] = "sqlite:///:memory:"
    os.environ["AUTH_GRPC_HOST"] = "localhost"
    _protos = os.path.normpath(
        os.path.join(os.path.dirname(__file__), "..", "..", "..", "protos", "gen", "python")
    )
    if os.path.exists(_protos):
        os.environ["PROTOS_PATH"] = _protos

    # Патчим app.clients.AuthGrpcClient (middleware) и app.api.auth._auth_client (auth endpoints)
    import app.clients
    import app.api.auth as auth_module
    mock_class = MagicMock(return_value=mock_auth_client)
    with patch.object(app.clients, "AuthGrpcClient", mock_class):
        from app import create_app
        app = create_app()
        auth_module._auth_client = mock_auth_client  # auth использует глобальный _auth_client
        app.config["TESTING"] = True
        yield app

    # Явно закрываем соединения с БД (SQLite in-memory), чтобы убрать ResourceWarning
    try:
        from app.extensions import db
        db.session.remove()
        # dispose() закрывает все соединения в пуле
        db.engine.dispose()
    except Exception:
        # В тестах лучше не падать из-за проблем с teardown
        pass


@pytest.fixture
def client(app):
    """Flask test client."""
    return app.test_client()


@pytest.fixture
def auth_headers():
    """Заголовок с валидным токеном для protected endpoints."""
    return {"Authorization": "Bearer valid-test-token"}
