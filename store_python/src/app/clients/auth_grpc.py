"""gRPC клиент для auth-service."""
import os
import sys

# Путь к protos: Docker — PROTOS_PATH=/app/protos_python, локально — относительно store_python
_protos_path = os.environ.get("PROTOS_PATH")
if not _protos_path:
    _base = os.path.dirname(os.path.abspath(__file__))
    # store_python/src/app/clients -> ../../../../protos/gen/python (до src/, где store_python)
    _protos_path = os.path.normpath(os.path.join(_base, "..", "..", "..", "..", "protos", "gen", "python"))
if os.path.exists(_protos_path) and _protos_path not in sys.path:
    sys.path.insert(0, _protos_path)

import grpc
from sso import sso_pb2, sso_pb2_grpc


class AuthGrpcClient:
    """Клиент для вызова auth-service по gRPC."""

    def __init__(self, host: str = None, port: int = None, app_id: str = None):
        self.host = host or os.getenv("AUTH_GRPC_HOST", "auth-service")
        self.port = port or int(os.getenv("AUTH_GRPC_PORT", "9090"))
        self.app_id = app_id or os.getenv("APP_ID", "a1b2c3d4-e5f6-7890-abcd-ef1234567890")
        self.rpc_timeout_s = float(os.getenv("AUTH_GRPC_TIMEOUT_S", "10"))
        self.oauth_timeout_s = float(os.getenv("AUTH_GRPC_OAUTH_TIMEOUT_S", "45"))

    def _channel(self):
        return grpc.insecure_channel(f"{self.host}:{self.port}")

    def register(self, email: str, password: str, name: str, surname: str, phone_number: str) -> str:
        """Регистрация пользователя. Возвращает user_id."""
        with self._channel() as ch:
            stub = sso_pb2_grpc.AuthStub(ch)
            req = sso_pb2.RegisterRequest(
                email=email,
                password=password,
                name=name,
                surname=surname,
                phone_number=phone_number or "",
            )
            resp = stub.Register(req, timeout=self.rpc_timeout_s)
            return resp.user_id

    def login(self, email: str, password: str) -> tuple[str, str]:
        """Логин. Возвращает (access_token, refresh_token)."""
        with self._channel() as ch:
            stub = sso_pb2_grpc.AuthStub(ch)
            req = sso_pb2.LoginRequest(email=email, password=password, app_id=self.app_id)
            resp = stub.Login(req, timeout=self.rpc_timeout_s)
            return resp.access_token, resp.refresh_token

    def validate(self, token: str) -> tuple[str, str]:
        """Проверка токена. Возвращает (user_id, app_id)."""
        with self._channel() as ch:
            stub = sso_pb2_grpc.AuthStub(ch)
            req = sso_pb2.ValidateRequest(token=token)
            # Validate защищён gRPC interceptor'ом в auth-service: нужен metadata authorization
            resp = stub.Validate(req, timeout=self.rpc_timeout_s, metadata=(("authorization", f"Bearer {token}"),))
            return resp.user_id, resp.app_id

    def change_password(self, email: str, old_password: str, new_password: str) -> None:
        """Смена пароля."""
        with self._channel() as ch:
            stub = sso_pb2_grpc.AuthStub(ch)
            req = sso_pb2.ChangePasswordRequest(
                email=email,
                old_password=old_password,
                new_password=new_password,
            )
            stub.ChangePassword(req, timeout=self.rpc_timeout_s)

    def restore_password(self, email: str) -> None:
        """Восстановление пароля (заглушка — вывод в консоль auth-service)."""
        with self._channel() as ch:
            stub = sso_pb2_grpc.AuthStub(ch)
            req = sso_pb2.RestorePasswordRequest(email=email)
            stub.RestorePassword(req, timeout=self.rpc_timeout_s)

    def get_google_auth_url(self, state: str = "") -> str:
        """Старт Google OAuth. Возвращает auth_url для редиректа на Google."""
        with self._channel() as ch:
            stub = sso_pb2_grpc.AuthStub(ch)
            req = sso_pb2.GetGoogleAuthURLRequest(state=state or "")
            resp = stub.GetGoogleAuthURL(req, timeout=self.rpc_timeout_s)
            return resp.auth_url

    def google_callback(self, code: str, state: str) -> tuple[str, str]:
        """Завершение Google OAuth. Возвращает (access_token, refresh_token)."""
        with self._channel() as ch:
            stub = sso_pb2_grpc.AuthStub(ch)
            req = sso_pb2.GoogleCallbackRequest(code=code, state=state, app_id=self.app_id)
            resp = stub.GoogleCallback(req, timeout=self.oauth_timeout_s)
            return resp.access_token, resp.refresh_token
