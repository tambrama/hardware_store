"""Запуск приложения: код в src/. Из корня: python run.py"""
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "src"))
os.chdir(os.path.join(os.path.dirname(__file__), "src"))

from app import create_app
from app.extensions import socketio

app = create_app()

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 5000))
    socketio.run(
        app,
        host="0.0.0.0",
        port=port,
        debug=True,
        use_reloader=False,
        allow_unsafe_werkzeug=True,
    )
