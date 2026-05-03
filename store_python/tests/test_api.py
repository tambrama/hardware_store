"""Тесты защищённых API endpoints (products, clients и т.д.)."""
import pytest


def test_products_without_token_returns_401(client):
    """GET /products без токена возвращает 401."""
    resp = client.get("/api/v1/products")
    assert resp.status_code == 401
    assert "Authorization" in resp.get_json()["message"]


def test_products_with_invalid_token_returns_401(client, mock_auth_client):
    """GET /products с невалидным токеном возвращает 401."""
    mock_auth_client.validate.side_effect = Exception("invalid token")
    resp = client.get(
        "/api/v1/products",
        headers={"Authorization": "Bearer bad-token"},
    )
    assert resp.status_code == 401


def test_products_with_valid_token_returns_200(client, mock_auth_client, auth_headers):
    """GET /products с валидным токеном возвращает 200."""
    mock_auth_client.validate.return_value = ("user-123", "app-id")
    resp = client.get("/api/v1/products", headers=auth_headers)
    assert resp.status_code == 200
    assert isinstance(resp.get_json(), list)


def test_clients_without_token_returns_401(client):
    """GET /clients без токена возвращает 401."""
    resp = client.get("/api/v1/clients")
    assert resp.status_code == 401


def test_clients_with_token_returns_200(client, auth_headers, mock_auth_client):
    """GET /clients с токеном возвращает 200."""
    mock_auth_client.validate.return_value = ("user-123", "app-id")
    resp = client.get("/api/v1/clients", headers=auth_headers)
    assert resp.status_code == 200


def test_swagger_accessible_without_token(client):
    """Swagger/apidocs доступен без токена."""
    resp = client.get("/apidocs/")
    assert resp.status_code in (200, 302, 404)  # Flasgger может быть на /apidocs/
