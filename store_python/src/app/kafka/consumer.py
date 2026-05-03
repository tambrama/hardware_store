import json
import os
import threading
import time
from datetime import date
import uuid

from kafka import KafkaConsumer

from app.extensions import db
from app.models import Product
from app.realtime.broker import ProductUpdateEvent, broker


def _env_bool(name: str, default: str = "0") -> bool:
    return os.getenv(name, default).lower() in {"1", "true", "yes", "y", "on"}


def start_product_updates_consumer(app) -> None:
    """Start background Kafka consumer (only once per process)."""
    if not _env_bool("KAFKA_CONSUMER_ENABLED", "0"):
        return

    if getattr(app, "_kafka_consumer_started", False):
        return
    app._kafka_consumer_started = True

    t = threading.Thread(target=_run_consumer, args=(app,), daemon=True)
    t.start()


def _run_consumer(app) -> None:
    bootstrap = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")
    topic = os.getenv("KAFKA_TOPIC", "product-updates")
    group_id = os.getenv("KAFKA_GROUP_ID", "store_python_gateway")

    consumer = KafkaConsumer(
        topic,
        bootstrap_servers=bootstrap.split(","),
        group_id=group_id,
        enable_auto_commit=True,
        auto_offset_reset=os.getenv("KAFKA_AUTO_OFFSET_RESET", "latest"),
        value_deserializer=lambda v: v,
        consumer_timeout_ms=1000,
    )

    while True:
        try:
            for msg in consumer:
                _handle_message(app, msg.value)
        except Exception:
            # transient errors: sleep and retry loop
            time.sleep(2)


def _handle_message(app, raw_value: bytes) -> None:
    try:
        payload = json.loads(raw_value.decode("utf-8"))
    except Exception:
        return

    product_id_raw = payload.get("product_id") or payload.get("id") or ""
    if not product_id_raw:
        return
    try:
        product_id = uuid.UUID(str(product_id_raw))
    except Exception:
        return

    # accept common field names
    price = payload.get("price", payload.get("cost"))
    stock = payload.get("available_stock", payload.get("remaining_quantity", payload.get("quantity")))
    ts = float(payload.get("ts") or payload.get("timestamp") or time.time())

    try:
        price_f = float(price)
        stock_i = int(stock)
    except Exception:
        return

    with app.app_context():
        product = db.session.get(Product, product_id)
        if product is None:
            return
        product.price = price_f
        product.available_stock = stock_i
        product.last_update_date = date.today()
        db.session.commit()

    broker.publish_product_update(
        ProductUpdateEvent(product_id=str(product_id), price=price_f, available_stock=stock_i, ts=ts)
    )
