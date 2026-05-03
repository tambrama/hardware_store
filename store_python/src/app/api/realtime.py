from flask import Blueprint, Response, stream_with_context

from app.realtime.broker import broker, sse_stream


realtime_bp = Blueprint("realtime", __name__)


@realtime_bp.route("/stream/products", methods=["GET"])
def stream_products():
    """
    SSE stream обновлений товара (price/available_stock).
    ---
    tags:
      - realtime
    responses:
      200:
        description: text/event-stream
    """
    q = broker.subscribe()
    return Response(
        stream_with_context(sse_stream(q)),
        mimetype="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "X-Accel-Buffering": "no",
            "Connection": "keep-alive",
        },
        direct_passthrough=True,
    )
