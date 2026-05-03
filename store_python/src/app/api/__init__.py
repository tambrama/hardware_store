from flask import Blueprint, request, jsonify, make_response


def _is_auth_exempt(path: str, method: str) -> bool:
    """Пути, не требующие JWT: register, login, reset, swagger."""
    if "/apidocs" in path or "/flasgger_static" in path:
        return True
    # realtime streams
    if "/realtime/stream/" in path:
        return True
    # socket.io transport endpoint
    if "/socket.io" in path:
        return True
    if path.endswith("/auth/register") or path.endswith("/auth/reset"):
        return True
    if path.endswith("/auth") and method == "POST":
        return True
    # Google OAuth endpoints are public entrypoints
    if path.endswith("/auth/oauth/google") and method == "GET":
        return True
    if path.endswith("/auth/google/callback") and method == "GET":
        return True
    return False


def register_blueprints(app):
    from app.api.auth import auth_bp
    from app.api.addresses import addresses_bp
    from app.api.clients import clients_bp
    from app.api.products import products_bp
    from app.api.suppliers import suppliers_bp
    from app.api.images import images_bp
    from app.api.realtime import realtime_bp

    api_v1 = Blueprint("api_v1", __name__, url_prefix="/api/v1")
    # Auth — без проверки токена
    api_v1.register_blueprint(auth_bp, url_prefix="/auth")
    # Остальные API — с проверкой токена (через before_request в blueprint)
    api_v1.register_blueprint(addresses_bp, url_prefix="/addresses")
    api_v1.register_blueprint(clients_bp, url_prefix="/clients")
    api_v1.register_blueprint(products_bp, url_prefix="/products")
    api_v1.register_blueprint(suppliers_bp, url_prefix="/suppliers")
    api_v1.register_blueprint(images_bp, url_prefix="/images")
    api_v1.register_blueprint(realtime_bp, url_prefix="/realtime")

    app.register_blueprint(api_v1)

    # Middleware: проверка JWT для всех /api/v1/* кроме auth (register, login, reset)
    @app.before_request
    def auth_middleware():
        path = request.path
        if not path.startswith("/api/v1"):
            return
        if _is_auth_exempt(path, request.method):
            return
        token = request.headers.get("Authorization")
        if not token or not token.startswith("Bearer "):
            return make_response(jsonify({"message": "Missing or invalid Authorization header"}), 401)
        token = token[7:].strip()
        try:
            from app.clients import AuthGrpcClient
            client = AuthGrpcClient()
            user_id, app_id = client.validate(token)
            request.user_id = user_id
            request.app_id = app_id
            return
        except Exception:
            return make_response(jsonify({"message": "Invalid or expired token"}), 401)

    from flasgger import Swagger
    import os
    import json
    _app_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    Swagger(app, template_file=os.path.join(_app_dir, "swagger_template.yaml"))

    # Убираем поле "swagger" из спецификации, чтобы не было конфликта с "openapi" в Swagger UI
    @app.after_request
    def strip_swagger_field(response):
        if response.content_type and "application/json" in response.content_type:
            try:
                data = response.get_json(silent=True)
                if isinstance(data, dict) and "openapi" in data:
                    # flasgger иногда отдаёт сразу и swagger+openapi — убираем swagger
                    data.pop("swagger", None)

                    # Переписываем пути для внешнего доступа через Nginx:
                    # внутреннее: /api/v1/... (Flask)
                    # внешнее:   /api/v2/... (Nginx)
                    paths = data.get("paths")
                    if isinstance(paths, dict):
                        new_paths = {}
                        for p, v in paths.items():
                            if isinstance(p, str) and p.startswith("/api/v1"):
                                new_paths["/api/v2" + p[len("/api/v1"):]] = v
                            else:
                                new_paths[p] = v
                        data["paths"] = new_paths

                    response.set_data(json.dumps(data))
            except Exception:
                pass
        return response
