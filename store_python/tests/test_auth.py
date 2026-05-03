"""Тесты эндпоинтов авторизации: /register, /auth, /reset."""
import pytest


def test_register_success(client, mock_auth_client):
    """Регистрация с валидными данными."""
    mock_auth_client.register.return_value = "user-uuid-123"
    resp = client.post(
        "/api/v1/auth/register",
        json={
            "email": "test@example.com",
            "password": "secret123",
            "name": "Иван",
            "surname": "Петров",
        },
        content_type="application/json",
    )
    assert resp.status_code == 201
    data = resp.get_json()
    assert data["user_id"] == "user-uuid-123"
    mock_auth_client.register.assert_called_once()
    call_kw = mock_auth_client.register.call_args
    assert call_kw[0][0] == "test@example.com"
    assert call_kw[0][1] == "secret123"
    assert call_kw[0][2] == "Иван"
    assert call_kw[0][3] == "Петров"


def test_register_missing_fields(client):
    """Регистрация без обязательных полей."""
    resp = client.post(
        "/api/v1/auth/register",
        json={"email": "test@example.com"},
        content_type="application/json",
    )
    assert resp.status_code == 400
    assert "Missing" in resp.get_json()["message"]


def test_register_duplicate_email(client, mock_auth_client):
    """Регистрация с уже занятым email."""
    mock_auth_client.register.side_effect = Exception("user already exists")
    mock_auth_client.register.return_value = None  # сброс return_value при side_effect
    resp = client.post(
        "/api/v1/auth/register",
        json={
            "email": "existing@example.com",
            "password": "secret",
            "name": "A",
            "surname": "B",
        },
        content_type="application/json",
    )
    assert resp.status_code == 409


def test_login_success(client, mock_auth_client):
    """Логин с валидными данными."""
    mock_auth_client.login.return_value = ("access-token-xyz", "refresh-token-abc")
    resp = client.post(
        "/api/v1/auth",
        json={"email": "test@example.com", "password": "secret"},
        content_type="application/json",
    )
    assert resp.status_code == 200
    data = resp.get_json()
    assert data["access_token"] == "access-token-xyz"
    assert data["refresh_token"] == "refresh-token-abc"


def test_login_invalid_credentials(client, mock_auth_client):
    """Логин с неверным паролем."""
    mock_auth_client.login.side_effect = Exception("invalid")
    mock_auth_client.login.return_value = None
    resp = client.post(
        "/api/v1/auth",
        json={"email": "test@example.com", "password": "wrong"},
        content_type="application/json",
    )
    assert resp.status_code == 401


def test_login_missing_fields(client):
    """Логин без email или password."""
    resp = client.post(
        "/api/v1/auth",
        json={},
        content_type="application/json",
    )
    assert resp.status_code == 400


def test_reset_password_success(client, mock_auth_client):
    """Восстановление пароля."""
    resp = client.post(
        "/api/v1/auth/reset",
        json={"email": "user@example.com"},
        content_type="application/json",
    )
    assert resp.status_code == 200
    # auth использует _auth_client из app.api.auth
    mock_auth_client.restore_password.assert_called()


def test_reset_password_missing_email(client):
    """Восстановление без email."""
    resp = client.post(
        "/api/v1/auth/reset",
        json={},
        content_type="application/json",
    )
    assert resp.status_code == 400
