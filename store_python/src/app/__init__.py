import os
import re
from flask import Flask
from dotenv import load_dotenv

load_dotenv()


def create_app():
    app = Flask(__name__)
    app.config["SQLALCHEMY_DATABASE_URI"] = os.getenv(
        "DATABASE_URI",
        "postgresql://postgres:postgres@localhost:5432/shopapi",
    )
    app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
    app.config["SWAGGER"] = {"title": "ShopAPI", "uiversion": 3}

    from app.extensions import db
    db.init_app(app)

    from app.extensions import socketio
    socketio.init_app(app)

    from app.api import register_blueprints
    register_blueprints(app)

    from flask import request

    @app.after_request
    def fix_swagger_none_in_js(response):
        """Flasgger: None → null в JS; убираем запрос к Google Fonts (таймаут за стеной)."""
        try:
            if "apidocs" in request.path and response.content_type and "text/html" in response.content_type:
                data = response.get_data()
                if b"None" in data:
                    data = data.replace(b"= None", b"= null").replace(b", None", b", null").replace(b"(None)", b"(null)")
                if b"fonts.googleapis.com" in data:
                    data = re.sub(rb"<link[^>]*fonts\.googleapis\.com[^>]*>", b"", data)
                response.set_data(data)
        except Exception:
            pass
        return response

    from flask import redirect
    @app.route("/swagger/")
    @app.route("/swagger/index.html")
    def swagger_ui():
        return redirect("/apidocs/")

    with app.app_context():
        db.create_all()

    # Kafka consumer (enabled only on write instance via env)
    try:
        from app.kafka.consumer import start_product_updates_consumer

        start_product_updates_consumer(app)
    except Exception:
        pass

    return app
