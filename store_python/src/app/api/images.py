import uuid
import os
import requests
from flask import Blueprint, request, Response, jsonify
from app.models import Product
from app.extensions import db
from app.repositories import ProductRepository
from app.extensions import get_redis_client

images_bp = Blueprint("images", __name__)
product_repo = ProductRepository()
_redis = None

PHOTO_SERVICE_URL = os.getenv("PHOTO_SERVICE_URL", "http://localhost:8000").rstrip("/")
PHOTO_TIMEOUT_S = float(os.getenv("PHOTO_SERVICE_TIMEOUT_S", "3.0"))

# Лимит размера изображения (5 MB)
MAX_IMAGE_SIZE = 5 * 1024 * 1024

# Магические байты форматов (начало файла)
IMAGE_SIGNATURES = (
    b"\x89PNG\r\n\x1a\n",           # PNG
    b"\xff\xd8\xff",                 # JPEG
    b"GIF87a",                       # GIF
    b"GIF89a",                       # GIF
    b"RIFF",                         # WebP (RIFF....WEBP)
)
WEBP_TAIL = b"WEBP"  # RIFF + 4 байта + WEBP


def _validate_image_data(data: bytes) -> tuple[bool, str]:
    """Проверка размера и формата. Возвращает (ok, message)."""
    if not data:
        return False, "Image data is required"
    if len(data) > MAX_IMAGE_SIZE:
        return False, f"Image size must not exceed {MAX_IMAGE_SIZE // (1024*1024)} MB"
    # Проверка на похожесть на изображение (магические байты)
    if data.startswith(b"RIFF") and len(data) >= 12 and data[8:12] == WEBP_TAIL:
        return True, ""
    for sig in IMAGE_SIGNATURES:
        if sig == b"RIFF":
            continue
        if data.startswith(sig):
            return True, ""
    return False, "Unsupported image format. Allowed: PNG, JPEG, GIF, WebP"


@images_bp.route("", methods=["POST"])
def add_image():
    """
    Добавление изображения (тело — byte array, product_id в query или body)
    ---
    tags:
      - images
    security:
      - BearerAuth: []
    parameters:
      - name: product_id
        in: query
        required: true
        schema:
          type: string
          format: uuid
    requestBody:
      required: true
      content:
        application/octet-stream:
          schema:
            type: string
            format: binary
    responses:
      201:
        description: Изображение создано, возвращается id
      400:
        description: Ошибка валидации
      404:
        description: Товар не найден
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    global _redis
    product_id = request.args.get("product_id")
    if not product_id:
        return jsonify({"message": "product_id is required"}), 400
    try:
        product_id = uuid.UUID(product_id)
    except (ValueError, TypeError):
        return jsonify({"message": "Invalid product_id"}), 400
    product = product_repo.get_by_id(product_id)
    if product is None:
        return jsonify({"message": "Product not found"}), 404
    data = request.get_data()
    ok, msg = _validate_image_data(data)
    if not ok:
        return jsonify({"message": msg}), 400

    try:
        r = requests.post(
            f"{PHOTO_SERVICE_URL}/images",
            data=data,
            headers={"Content-Type": "application/octet-stream"},
            timeout=PHOTO_TIMEOUT_S,
        )
    except Exception:
        return jsonify({"message": "Photo service is unavailable"}), 503

    if r.status_code != 200:
        return jsonify({"message": "Failed to store image"}), 502
    try:
        payload = r.json() or {}
        image_id = uuid.UUID(payload.get("id", ""))
        shard = payload.get("shard")
    except Exception:
        return jsonify({"message": "Invalid response from photo service"}), 502

    product.image_id = image_id
    product_repo.update(product)
    try:
        _redis = _redis or get_redis_client()
        _redis.delete(f"product_image:{product.id}")
    except Exception:
        pass
    resp = {"id": str(image_id)}
    if shard is not None:
        try:
            shard_i = int(shard)
            resp["shard"] = shard_i
            resp["database"] = f"photos{shard_i}"
        except Exception:
            pass
    return resp, 201


@images_bp.route("/<uuid:image_id>", methods=["PUT"])
def update_image(image_id):
    """
    Изменение изображения (новые байты в теле)
    ---
    tags:
      - images
    security:
      - BearerAuth: []
    parameters:
      - name: image_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    requestBody:
      required: true
      content:
        application/octet-stream:
          schema:
            type: string
            format: binary
    responses:
      200:
        description: Изображение обновлено
      400:
        description: Тело запроса пустое
      404:
        description: Изображение не найдено
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    global _redis
    data = request.get_data()
    ok, msg = _validate_image_data(data)
    if not ok:
        return jsonify({"message": msg}), 400

    try:
        r = requests.put(
            f"{PHOTO_SERVICE_URL}/images/{image_id}",
            data=data,
            headers={"Content-Type": "application/octet-stream"},
            timeout=PHOTO_TIMEOUT_S,
        )
    except Exception:
        return jsonify({"message": "Photo service is unavailable"}), 503

    if r.status_code == 404:
        return jsonify({"message": "Image not found"}), 404
    if r.status_code != 200:
        return jsonify({"message": "Failed to update image"}), 502

    try:
        _redis = _redis or get_redis_client()
        _redis.delete(f"image:{image_id}")
    except Exception:
        pass
    return jsonify({"id": str(image_id)}), 200


@images_bp.route("/<uuid:image_id>", methods=["DELETE"])
def delete_image(image_id):
    """
    Удаление изображения по id
    ---
    tags:
      - images
    security:
      - BearerAuth: []
    parameters:
      - name: image_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      204:
        description: Изображение удалено
      404:
        description: Изображение не найдено
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    global _redis
    try:
        r = requests.delete(f"{PHOTO_SERVICE_URL}/images/{image_id}", timeout=PHOTO_TIMEOUT_S)
    except Exception:
        return jsonify({"message": "Photo service is unavailable"}), 503
    if r.status_code not in {200, 204, 404}:
        return jsonify({"message": "Failed to delete image"}), 502

    for product in db.session.query(Product).filter(Product.image_id == image_id).all():
        product.image_id = None
        db.session.commit()
        try:
            _redis = _redis or get_redis_client()
            _redis.delete(f"product_image:{product.id}")
        except Exception:
            pass
    try:
        _redis = _redis or get_redis_client()
        _redis.delete(f"image:{image_id}")
    except Exception:
        pass
    return "", 204


@images_bp.route("/product/<uuid:product_id>", methods=["GET"])
def get_image_by_product_id(product_id):
    """
    Получение изображения по id товара
    ---
    tags:
      - images
    security:
      - BearerAuth: []
    parameters:
      - name: product_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      200:
        description: Изображение (application/octet-stream)
      404:
        description: Изображение или товар не найдены
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    global _redis
    try:
        _redis = _redis or get_redis_client()
        cached = _redis.get(f"product_image:{product_id}")
        if cached:
            return Response(
                cached,
                mimetype="application/octet-stream",
                headers={"Content-Disposition": "attachment", "X-Cache": "HIT"},
            )
    except Exception:
        pass

    product = product_repo.get_by_id(product_id)
    if product is None or not product.image_id:
        return jsonify({"message": "Image not found"}), 404
    image_id = product.image_id

    try:
        r = requests.get(f"{PHOTO_SERVICE_URL}/images/{image_id}", timeout=PHOTO_TIMEOUT_S)
    except Exception:
        return jsonify({"message": "Photo service is unavailable"}), 503

    if r.status_code == 404:
        return jsonify({"message": "Image not found"}), 404
    if r.status_code != 200:
        return jsonify({"message": "Failed to fetch image"}), 502

    data = r.content or b""
    try:
        _redis = _redis or get_redis_client()
        # 1 hour TTL
        _redis.setex(f"product_image:{product_id}", 3600, data)
    except Exception:
        pass
    return Response(data, mimetype="application/octet-stream", headers={"Content-Disposition": "attachment"})


@images_bp.route("/<uuid:image_id>", methods=["GET"])
def get_image_by_id(image_id):
    """
    Получение изображения по id изображения
    ---
    tags:
      - images
    security:
      - BearerAuth: []
    parameters:
      - name: image_id
        in: path
        required: true
        schema:
          type: string
          format: uuid
    responses:
      200:
        description: Изображение (application/octet-stream)
      404:
        description: Изображение не найдено
      401:
        description: Не авторизован — отсутствует или невалидный токен
    """
    global _redis
    try:
        _redis = _redis or get_redis_client()
        cached = _redis.get(f"image:{image_id}")
        if cached:
            return Response(
                cached,
                mimetype="application/octet-stream",
                headers={"Content-Disposition": "attachment", "X-Cache": "HIT"},
            )
    except Exception:
        pass

    try:
        r = requests.get(f"{PHOTO_SERVICE_URL}/images/{image_id}", timeout=PHOTO_TIMEOUT_S)
    except Exception:
        return jsonify({"message": "Photo service is unavailable"}), 503

    if r.status_code == 404:
        return jsonify({"message": "Image not found"}), 404
    if r.status_code != 200:
        return jsonify({"message": "Failed to fetch image"}), 502

    try:
        _redis = _redis or get_redis_client()
        _redis.setex(f"image:{image_id}", 3600, r.content or b"")
    except Exception:
        pass
    return Response(
        r.content or b"",
        mimetype="application/octet-stream",
        headers={"Content-Disposition": "attachment"},
    )
