from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()

# Socket.IO (WebSockets/long-poll) для пуша обновлений
import os

try:
    from flask_socketio import SocketIO

    # threading mode is the most stable for SSE + docker demo
    socketio = SocketIO(
        cors_allowed_origins="*",
        async_mode=os.getenv("SOCKETIO_ASYNC_MODE", "threading"),
        path="/api/v1/socket.io",
    )
except Exception:  # pragma: no cover
    # Allows running unit tests without optional realtime deps installed.
    class _DummySocketIO:
        def init_app(self, *args, **kwargs):
            return None

        def run(self, *args, **kwargs):
            raise RuntimeError("Flask-SocketIO is not installed")

    socketio = _DummySocketIO()


def get_redis_url() -> str:
    import os

    return os.getenv("REDIS_URL", "redis://redis:6379/0")


def get_redis_client():
    import redis

    return redis.Redis.from_url(get_redis_url())
