import json
import queue
import threading
import time
from dataclasses import dataclass

from app.extensions import socketio


@dataclass(frozen=True)
class ProductUpdateEvent:
    product_id: str
    price: float
    available_stock: int
    ts: float

    def to_dict(self) -> dict:
        return {
            "product_id": self.product_id,
            "price": self.price,
            "available_stock": self.available_stock,
            "ts": self.ts,
        }


class RealtimeBroker:
    """In-process pubsub: SSE subscribers + Socket.IO broadcast."""

    def __init__(self):
        self._lock = threading.Lock()
        self._subscribers: set[queue.Queue] = set()

    def subscribe(self) -> queue.Queue:
        q: queue.Queue = queue.Queue(maxsize=100)
        with self._lock:
            self._subscribers.add(q)
        return q

    def unsubscribe(self, q: queue.Queue) -> None:
        with self._lock:
            self._subscribers.discard(q)

    def publish_product_update(self, event: ProductUpdateEvent) -> None:
        payload = event.to_dict()

        # Socket.IO broadcast (JS clients)
        socketio.emit("product_update", payload, namespace="/")

        # SSE broadcast
        with self._lock:
            subs = list(self._subscribers)
        for q in subs:
            try:
                q.put_nowait(payload)
            except queue.Full:
                # Drop slow subscriber updates
                pass


broker = RealtimeBroker()


def sse_stream(q: queue.Queue):
    """Generator for Flask Response with text/event-stream."""
    # Send first ping immediately so clients see activity
    yield f"event: ping\ndata: {int(time.time())}\n\n".encode("utf-8")
    try:
        while True:
            try:
                payload = q.get(timeout=3)
                data = json.dumps(payload, ensure_ascii=False)
                yield f"event: product_update\ndata: {data}\n\n".encode("utf-8")
            except queue.Empty:
                # keepalive ping
                yield f"event: ping\ndata: {int(time.time())}\n\n".encode("utf-8")
    finally:
        broker.unsubscribe(q)
